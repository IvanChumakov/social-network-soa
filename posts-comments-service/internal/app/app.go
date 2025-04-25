package app

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"social-network/posts-comments-service/internal/logger"
	"social-network/posts-comments-service/internal/models"
	eventsservice "social-network/posts-comments-service/internal/service/events-service"
	pb "social-network/protos"
	"strconv"
	"time"
)

type Service interface {
	AddPost(post *pb.PostEssential, userId int32) error
	DeletePost(postId int32) error
	UpdatePost(post *pb.PostWithNoUser) error
	GetPostById(postId int32) (*pb.Post, error)
	GetAllPosts(pagination *pb.Pagination, userId int32) (*pb.AllPosts, error)
	AddComment(comment *pb.Comment, userId int32) error
	GetAllComments(pagination *pb.Pagination, postId int32) (*pb.AllComments, error)
}

type Server struct {
	pb.UnimplementedPostsServiceServer
	service Service
	broker  *eventsservice.KafkaEvents
}

func NewServer(service Service, broker *eventsservice.KafkaEvents) *Server {
	return &Server{
		service: service,
		broker:  broker,
	}
}

func (s *Server) extractFromCtx(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata")
	}

	data := md.Get(key)
	if len(data) == 0 {
		return "", fmt.Errorf("no data")
	}

	return data[0], nil
}

func (s *Server) AddPost(ctx context.Context, post *pb.PostEssential) (*emptypb.Empty, error) {
	logger.Info("add post called")
	userIds, err := s.extractFromCtx(ctx, "user_id")
	if err != nil {
		logger.Error(fmt.Sprintf("error with extracting data from context: %v", err))
		return nil, err
	}

	userId, _ := strconv.Atoi(userIds)
	return &emptypb.Empty{}, s.service.AddPost(post, int32(userId))
}

func (s *Server) DeletePost(_ context.Context, post *pb.PostId) (*emptypb.Empty, error) {
	logger.Info("delete post called")
	return &emptypb.Empty{}, s.service.DeletePost(post.PostId)
}

func (s *Server) UpdatePost(_ context.Context, post *pb.PostWithNoUser) (*emptypb.Empty, error) {
	logger.Info("update post called")
	return &emptypb.Empty{}, s.service.UpdatePost(post)
}

func (s *Server) GetPostById(_ context.Context, id *pb.PostId) (*pb.Post, error) {
	logger.Info("get post by id called")
	return s.service.GetPostById(id.PostId)
}

func (s *Server) GetAllPostsPaginated(ctx context.Context, pagination *pb.Pagination) (*pb.AllPosts, error) {
	logger.Info("get all posts paginated")
	userIds, err := s.extractFromCtx(ctx, "user_id")
	if err != nil {
		logger.Error(fmt.Sprintf("error with extracting data from context: %v", err))
		return nil, err
	}

	userId, _ := strconv.Atoi(userIds)

	return s.service.GetAllPosts(pagination, int32(userId))
}

func (s *Server) LikeEvent(ctx context.Context, like *pb.Like) (*emptypb.Empty, error) {
	logger.Info("like event called")

	userIds, err := s.extractFromCtx(ctx, "user_id")
	if err != nil {
		logger.Error(fmt.Sprintf("error with extracting data from context: %v", err))
		return nil, err
	}

	userIDCasted, _ := strconv.Atoi(userIds)

	return &emptypb.Empty{}, s.broker.SendEvent(models.Like{
		PostId: like.PostId,
		UserId: int32(userIDCasted),
		Time:   time.Now(),
	}, "likes-topic")
}

func (s *Server) AddComment(ctx context.Context, comment *pb.Comment) (*emptypb.Empty, error) {
	logger.Info("add comment called")

	userIds, err := s.extractFromCtx(ctx, "user_id")
	if err != nil {
		logger.Error(fmt.Sprintf("error with extracting data from context: %v", err))
		return nil, err
	}

	userIDCasted, _ := strconv.Atoi(userIds)
	err = s.service.AddComment(comment, int32(userIDCasted))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, s.broker.SendEvent(models.Comment{
		UserId: int32(userIDCasted),
		PostId: comment.PostId,
		Time:   time.Now(),
	}, "comments-topic")
}

func (s *Server) ViewEvent(ctx context.Context, view *pb.View) (*emptypb.Empty, error) {
	logger.Info("view event called")

	userIds, err := s.extractFromCtx(ctx, "user_id")
	if err != nil {
		logger.Error(fmt.Sprintf("error with extracting data from context: %v", err))
		return nil, err
	}

	userIDCasted, _ := strconv.Atoi(userIds)

	return &emptypb.Empty{}, s.broker.SendEvent(models.View{
		PostId: view.PostId,
		UserId: int32(userIDCasted),
		Time:   time.Now(),
	}, "views-topic")
}

func (s *Server) GetAllCommentsPaginated(ctx context.Context, pagination *pb.Pagination) (*pb.AllComments, error) {
	logger.Info("get all comments paginated called")

	postId, err := s.extractFromCtx(ctx, "post_id")
	if err != nil {
		logger.Error(fmt.Sprintf("error with extracting post_id data from context: %v", err))
		return nil, err
	}

	postIDCasted, _ := strconv.Atoi(postId)

	return s.service.GetAllComments(pagination, int32(postIDCasted))
}
