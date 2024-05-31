package managergrpc

import (
	"context"
	"errors"
	"github.com/google/uuid"
	mngv1 "github.com/s0vunia/password-manager-protos/gen/go/manager"
	"github.com/s0vunia/password-manager/internal/domain"
	"github.com/s0vunia/password-manager/internal/repositories"
	"github.com/s0vunia/password-manager/internal/services/manager/item"
	"github.com/s0vunia/password-manager/internal/services/manager/loginItem"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverApi struct {
	mngv1.UnimplementedManagerServer
	itemService      item.IItemService
	loginItemService loginItem.ILoginItemService
}

func Register(gRPCServer *grpc.Server, itemService item.IItemService, loginItemService loginItem.ILoginItemService) {
	mngv1.RegisterManagerServer(gRPCServer, &serverApi{itemService: itemService, loginItemService: loginItemService})
}

func (s serverApi) CreateLoginItem(ctx context.Context, request *mngv1.CreateLoginItemRequest) (*mngv1.CreateLoginItemResponse, error) {
	if request.Item.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "item.name is required")
	}
	if request.Item.FolderId == nil {
		return nil, status.Error(codes.InvalidArgument, "item.folder_id is required")
	}
	if request.Item.UserId == nil {
		var err error
		_, err = uuid.Parse(ctx.Value("userID").(string))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "user_id is required")
		}
	} else {
		_, _ = uuid.Parse(request.Item.UserId.Value)
	}

	if request.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	if request.EncryptPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "encrypt password is required")
	}

	logItem := RequestToLoginItemModel(request, ctx)
	id, err := s.loginItemService.CreateLoginItem(ctx, logItem)
	if err != nil {
		if errors.Is(err, repositories.ErrItemExists) {
			return nil, status.Error(codes.NotFound, "item exists")
		}
		return nil, status.Error(codes.Internal, "failed to create login item")
	}
	return &mngv1.CreateLoginItemResponse{
		Item: &mngv1.CreateItemResponse{
			Id: &mngv1.UUID{Value: id.String()},
		},
	}, nil
}

func RequestToLoginItemModel(request *mngv1.CreateLoginItemRequest, ctx context.Context) domain.LoginItem {
	folderId, _ := uuid.Parse(request.Item.FolderId.Value)
	var userId uuid.UUID
	if request.Item.UserId == nil {
		userId, _ = uuid.Parse(ctx.Value("userID").(string))
	} else {
		userId, _ = uuid.Parse(request.Item.UserId.Value)
	}
	return domain.LoginItem{
		Item: domain.Item{
			Type:       domain.ItemType(request.Item.Type),
			Name:       request.Item.Name,
			FolderId:   folderId,
			UserId:     userId,
			IsFavorite: false,
		},
		Login:           request.Login,
		EncryptPassword: request.EncryptPassword,
	}
}

func (s serverApi) GetItem(ctx context.Context, request *mngv1.GetItemRequest) (*mngv1.GetItemResponse, error) {
	if request.Id == nil {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	var userId uuid.UUID
	if request.UserId == nil {
		var err error
		userId, err = uuid.Parse(ctx.Value("userID").(string))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "user_id is required")
		}
	} else {
		userId, _ = uuid.Parse(request.UserId.Value)
	}
	id, _ := uuid.Parse(request.Id.Value)
	item, err := s.itemService.GetItem(ctx, id, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get item")
	}
	return s.GetItemModelToResponse(*item), nil
}

func (s serverApi) GetItemModelToResponse(model domain.Item) *mngv1.GetItemResponse {
	return &mngv1.GetItemResponse{
		Id:         &mngv1.UUID{Value: model.ID.String()},
		Name:       model.Name,
		Type:       mngv1.ItemType(model.Type),
		FolderId:   &mngv1.UUID{Value: model.FolderId.String()},
		UserId:     &mngv1.UUID{Value: model.UserId.String()},
		IsFavorite: model.IsFavorite,
	}
}

func (s serverApi) GetItems(ctx context.Context, request *mngv1.GetItemsRequest) (*mngv1.GetItemsResponse, error) {
	var userId uuid.UUID
	if request.UserId == nil {
		var err error
		userId, err = uuid.Parse(ctx.Value("userID").(string))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "user_id is required")
		}
	} else {
		userId, _ = uuid.Parse(request.UserId.Value)
	}
	items, err := s.itemService.GetItems(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get items")
	}
	var listOfItems []*mngv1.GetItemResponse
	for _, expression := range items {
		listOfItems = append(listOfItems, s.GetItemModelToResponse(*expression))
	}
	return &mngv1.GetItemsResponse{ListOfItems: listOfItems}, nil
}

func (s serverApi) GetLoginItem(ctx context.Context, request *mngv1.GetLoginItemRequest) (*mngv1.GetLoginItemResponse, error) {
	log.Info(request)
	if request.Item.Id == nil {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	var userId uuid.UUID
	if request.Item.UserId == nil {
		var err error
		userId, err = uuid.Parse(ctx.Value("userID").(string))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "user_id is required")
		}
	} else {
		userId, _ = uuid.Parse(request.Item.UserId.Value)
	}

	id, _ := uuid.Parse(request.Item.Id.Value)
	item, err := s.loginItemService.GetLoginItem(ctx, id, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get login item")
	}
	return s.GetLoginItemModelToResponse(*item), nil
}

func (s serverApi) GetLoginItemModelToResponse(model domain.LoginItem) *mngv1.GetLoginItemResponse {
	return &mngv1.GetLoginItemResponse{
		Id:              &mngv1.UUID{Value: model.ID.String()},
		Item:            s.GetItemModelToResponse(model.Item),
		Login:           model.Login,
		EncryptPassword: model.EncryptPassword,
	}
}
func (s serverApi) GetLoginItems(ctx context.Context, request *mngv1.GetLoginItemsRequest) (*mngv1.GetLoginItemsResponse, error) {
	var userId uuid.UUID
	if request.Items.UserId == nil {
		var err error
		userId, err = uuid.Parse(ctx.Value("userID").(string))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "user_id is required")
		}
	} else {
		userId, _ = uuid.Parse(request.Items.UserId.Value)
	}
	items, err := s.loginItemService.GetLoginItems(ctx, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get login items")
	}
	var listOfItems []*mngv1.GetLoginItemResponse
	for _, expression := range items {
		listOfItems = append(listOfItems, s.GetLoginItemModelToResponse(*expression))
	}
	return &mngv1.GetLoginItemsResponse{ListOfItems: listOfItems}, nil
}

func (s serverApi) GetItemsByFolder(ctx context.Context, request *mngv1.GetItemsByFolderRequest) (*mngv1.GetItemsByFolderRequest, error) {
	if request.FolderId == nil {
		return nil, status.Error(codes.InvalidArgument, "folder_id is required")
	}
	var userId uuid.UUID
	if request.UserId == nil {
		var err error
		userId, err = uuid.Parse(ctx.Value("userID").(string))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "user_id is required")
		}
	} else {
		userId, _ = uuid.Parse(request.UserId.Value)
	}
	log.Info(userId)
	return nil, status.Error(codes.Canceled, "not supported")
}

func (s serverApi) DeleteLoginItem(ctx context.Context, request *mngv1.DeleteLoginItemRequest) (*mngv1.DeleteLoginItemResponse, error) {
	if request.ItemId == nil {
		return nil, status.Error(codes.InvalidArgument, "item id is required")
	}
	var userId uuid.UUID
	if request.UserId == nil {
		var err error
		userId, err = uuid.Parse(ctx.Value("userID").(string))
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "user_id is required")
		}
	} else {
		userId, _ = uuid.Parse(request.UserId.Value)
	}
	itemId, err := uuid.Parse(request.ItemId.Value)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get login items")
	}
	err = s.loginItemService.DeleteLoginItem(ctx, userId, itemId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get login items")
	}
	return &mngv1.DeleteLoginItemResponse{}, nil
}
