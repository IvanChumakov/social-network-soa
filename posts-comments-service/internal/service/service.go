package service

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"social-network/posts-comments-service/internal/repository"
	pb "social-network/protos"
	"time"
)

type Repository interface {
	AddPost(post repository.Post) error
	DeletePost(id int32) error
	UpdatePost(post repository.Post) error
	GetPostById(id int32) (repository.Post, error)
	GetAllPosts(limit int32, offset int32, userId int32) ([]repository.Post, error)
}

type PostService struct {
	repository Repository
}

func NewPostService(repo Repository) *PostService {
	return &PostService{
		repo,
	}
}

func (ps *PostService) AddPost(post *pb.PostEssential, userId int32) error {
	dbPost := repository.Post{
		Name:        post.Name,
		Description: post.Description,
		CreatorId:   userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsPrivate:   post.IsPrivate,
		Tags:        post.Tags,
	}

	return ps.repository.AddPost(dbPost)
}

func (ps *PostService) DeletePost(postId int32) error {
	return ps.repository.DeletePost(postId)
}

func (ps *PostService) UpdatePost(post *pb.PostWithNoUser) error {
	dbPost := repository.Post{
		Id:          post.Id,
		Name:        post.Name,
		Description: post.Description,
		UpdatedAt:   time.Now(),
		IsPrivate:   post.IsPrivate,
		Tags:        post.Tags,
	}

	return ps.repository.UpdatePost(dbPost)
}

func (ps *PostService) GetPostById(id int32) (*pb.Post, error) {
	post, err := ps.repository.GetPostById(id)
	if err != nil {
		return nil, err
	}
	grpcPost := &pb.Post{
		Name:        post.Name,
		Description: post.Description,
		CreatedAd:   timestamppb.New(post.CreatedAt),
		UpdatedAt:   timestamppb.New(post.UpdatedAt),
		IsPrivate:   post.IsPrivate,
		Tags:        post.Tags,
		Id:          post.Id,
		UserId:      post.CreatorId,
	}

	return grpcPost, nil
}

func (ps *PostService) GetAllPosts(pagination *pb.Pagination, userId int32) (*pb.AllPosts, error) {
	posts, err := ps.repository.GetAllPosts(pagination.PageSize, pagination.PageIndex, userId)
	if err != nil {
		return nil, err
	}

	var allPorts pb.AllPosts
	allPorts.Posts = make([]*pb.Post, 0)
	for _, post := range posts {
		allPorts.Posts = append(allPorts.Posts, &pb.Post{
			Name:        post.Name,
			Id:          post.Id,
			Description: post.Description,
			UserId:      post.CreatorId,
			CreatedAd:   timestamppb.New(post.CreatedAt),
			UpdatedAt:   timestamppb.New(post.UpdatedAt),
			IsPrivate:   post.IsPrivate,
			Tags:        post.Tags,
		})
	}

	return &allPorts, nil
}
