package main

import (
	"fmt"
	"github.com/s0vunia/password-manager/internal/app"
	"github.com/s0vunia/password-manager/internal/config"
	appRepo "github.com/s0vunia/password-manager/internal/repositories/app"
	itemRepo "github.com/s0vunia/password-manager/internal/repositories/item"
	loginItemRepo "github.com/s0vunia/password-manager/internal/repositories/item/loginItem"
	"github.com/s0vunia/password-manager/internal/repositories/user"
	"github.com/s0vunia/password-manager/internal/services/auth"
	"github.com/s0vunia/password-manager/internal/services/manager/item"
	"github.com/s0vunia/password-manager/internal/services/manager/loginItem"
	log "github.com/sirupsen/logrus"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Add this line for logging filename and line number!
	log.SetReportCaller(true)

	// Only log the debug severity or above.
	log.SetLevel(log.DebugLevel)
}

// Start инициализирует и запускает оркестратор
func Start() {
	cfg := config.MustLoad()
	dataSourceName := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DbName, cfg.Postgres.User, cfg.Postgres.Password)
	userRepository, err := user.NewPostgresRepository(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to init user repo: %v", err)
	}
	appRepository, err := appRepo.NewPostgresRepository(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to init app repo: %v", err)
	}
	itemRepository, err := itemRepo.NewPostgresRepository(dataSourceName)
	if err != nil {
		log.Fatalf("Failed to init item repo: %v", err)
	}
	loginItemRepository, err := loginItemRepo.NewPostgresRepository(dataSourceName, itemRepository, itemRepository)
	if err != nil {
		log.Fatalf("Failed to init item repo: %v", err)
	}

	logSlog := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	newItem := item.New(logSlog, itemRepository)
	newLoginItem := loginItem.New(logSlog, loginItemRepository, loginItemRepository)
	newAuth := auth.New(logSlog, userRepository, userRepository, appRepository, cfg.TokenTTL)

	// Регистрация хендлеров
	application := app.New(logSlog, newItem, newLoginItem, appRepository, newAuth, cfg.GRPC.Port)
	go func() {
		application.GRPCServer.MustRun()
	}()
	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info("Gracefully stopped")

}

func main() {
	Start()
}
