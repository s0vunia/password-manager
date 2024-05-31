package grpcapp

import (
	"bytes"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	authgrpc "github.com/s0vunia/password-manager/internal/grpc/auth"
	managergrpc "github.com/s0vunia/password-manager/internal/grpc/manager"
	"github.com/s0vunia/password-manager/internal/repositories/app"
	authService "github.com/s0vunia/password-manager/internal/services/auth"
	"github.com/s0vunia/password-manager/internal/services/manager/item"
	"github.com/s0vunia/password-manager/internal/services/manager/loginItem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"log/slog"
	"net"
	"runtime/pprof"
)

var (
	listOfRoutesJWTMiddleware = []string{
		"/manager.Manager/CreateLoginItem",
		"/manager.Manager/GetItem",
		"/manager.Manager/GetItems",
		"/manager.Manager/GetLoginItem",
		"/manager.Manager/GetLoginItems",
		"/manager.Manager/GetItemsByFolder",
		"/manager.Manager/DeleteLoginItem",
	}
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	authService authService.IOAuth,
	itemService item.IItemService,
	loginItemService loginItem.ILoginItemService,
	appRepo app.Repository,
	port int,

) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			//logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			stackTrace := getStackTrace() // Получаем трассировку стека перед записью в лог
			log.Error("Recovered from panic",
				slog.Any("panic", p),
				slog.Any("stacktrace", stackTrace))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}
	gRPCServer := grpc.NewServer(grpc.MaxRecvMsgSize(1024*1024*600), grpc.MaxSendMsgSize(1024*1024*600),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
			logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
			selector.UnaryServerInterceptor(authgrpc.JWTMiddleware(appRepo), selector.MatchFunc(checkGrpcNameForJWT)),
		))
	authgrpc.Register(gRPCServer, authService)
	managergrpc.Register(gRPCServer, itemService, loginItemService)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}

func checkGrpcNameForJWT(ctx context.Context, callMeta interceptors.CallMeta) bool {
	fullMethName := callMeta.FullMethod()
	log.Printf(fullMethName)
	for _, name := range listOfRoutesJWTMiddleware {
		if name == fullMethName {
			return true
		}
	}
	return false
}

func getStackTrace() string {
	buf := bytes.NewBuffer(nil)
	if err := pprof.Lookup("goroutine").WriteTo(buf, 1); err != nil {
		log.Printf("Failed to capture goroutine stack trace: %v", err)
		return ""
	}
	stackTrace := buf.String()
	return stackTrace
}
