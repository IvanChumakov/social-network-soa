package app

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"social-network/posts-comments-service/internal/logger"
	pb "social-network/protos"
	"strconv"
)

type Service interface {
	AddPost(post *pb.PostEssential, userId int32) error
	DeletePost(postId int32) error
	UpdatePost(post *pb.PostWithNoUser) error
	GetPostById(postId int32) (*pb.Post, error)
	GetAllPosts(pagination *pb.Pagination, userId int32) (*pb.AllPosts, error)
}

type Server struct {
	pb.UnimplementedPostsServiceServer
	service Service
}

func NewServer(service Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) AddPost(ctx context.Context, post *pb.PostEssential) (*emptypb.Empty, error) {
	logger.Info("add post called")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata")
	}
	userIDs := md.Get("user_id")
	if len(userIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "no user_id")
	}

	userId, _ := strconv.Atoi(userIDs[0])
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no metadata")
	}
	userIDs := md.Get("user_id")
	if len(userIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "no user_id")
	}

	userId, _ := strconv.Atoi(userIDs[0])

	return s.service.GetAllPosts(pagination, int32(userId))
}
