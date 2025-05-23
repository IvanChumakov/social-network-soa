package app

import (
	"context"
	pb "social-network/protos"
)

type Service interface {
	GetLikesViewsComments(context.Context, *pb.PostId) (*pb.LikesViewsComments, error)
	GetViewsDynamic(context.Context, *pb.PostId) (*pb.ViewsDynamicArr, error)
	GetLikesDynamic(context.Context, *pb.PostId) (*pb.LikesDynamicArr, error)
	GetCommentsDynamic(context.Context, *pb.PostId) (*pb.CommentsDynamicArr, error)
	GetPostsTop(context.Context, *pb.Type) (*pb.PostsTop, error)
	GetUsersTop(context.Context, *pb.Type) (*pb.UsersTop, error)
}

type App struct {
	pb.UnimplementedStatisticsServiceServer
	service Service
}

func NewApp(service Service) *App {
	return &App{
		service: service,
	}
}

func (s *App) GetLikesViewsComments(ctx context.Context, postId *pb.PostId) (*pb.LikesViewsComments, error) {
	return s.service.GetLikesViewsComments(ctx, postId)
}

func (s *App) GetViewsDynamic(ctx context.Context, postId *pb.PostId) (*pb.ViewsDynamicArr, error) {
	return s.service.GetViewsDynamic(ctx, postId)
}

func (s *App) GetLikesDynamic(ctx context.Context, postId *pb.PostId) (*pb.LikesDynamicArr, error) {
	return s.service.GetLikesDynamic(ctx, postId)
}

func (s *App) GetCommentsDynamic(ctx context.Context, postId *pb.PostId) (*pb.CommentsDynamicArr, error) {
	return s.service.GetCommentsDynamic(ctx, postId)
}

func (s *App) GetPostsTop(ctx context.Context, eventType *pb.Type) (*pb.PostsTop, error) {
	return s.service.GetPostsTop(ctx, eventType)
}

func (s *App) GetUsersTop(ctx context.Context, eventType *pb.Type) (*pb.UsersTop, error) {
	return s.service.GetUsersTop(ctx, eventType)
}
