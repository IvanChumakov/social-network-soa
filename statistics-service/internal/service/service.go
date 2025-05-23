package service

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "social-network/protos"
	"social-network/statistics-service/internal/logger"
	"social-network/statistics-service/internal/repository/statistics"
)

type Transactor interface {
	WithTransaction(context.Context, func(context.Context) error) error
	WithTransactionValue(ctx context.Context, txFunc func(ctx context.Context) (any, error)) (any, error)
}

type StatisticsService struct {
	repo       *statistics.StatisticsRepository
	transactor Transactor
}

func NewStatisticsService(repo *statistics.StatisticsRepository, transactor Transactor) *StatisticsService {
	return &StatisticsService{
		repo:       repo,
		transactor: transactor,
	}
}

func (ss *StatisticsService) GetLikesViewsComments(ctx context.Context, postId *pb.PostId) (*pb.LikesViewsComments, error) {
	likes, err := ss.transactor.WithTransactionValue(ctx, func(ctx context.Context) (any, error) {
		return ss.repo.GetLikesCount(ctx, postId.PostId)
	})
	if err != nil {
		return nil, err
	}

	views, err := ss.transactor.WithTransactionValue(ctx, func(ctx context.Context) (any, error) {
		return ss.repo.GetViewsCount(ctx, postId.PostId)
	})
	if err != nil {
		return nil, err
	}

	comments, err := ss.transactor.WithTransactionValue(ctx, func(ctx context.Context) (any, error) {
		return ss.repo.GetCommentsCount(ctx, postId.PostId)
	})
	if err != nil {
		return nil, err
	}

	return &pb.LikesViewsComments{
		Comments: comments.(int32),
		Views:    views.(int32),
		Likes:    likes.(int32),
	}, nil
}

func (ss *StatisticsService) GetViewsDynamic(ctx context.Context, postId *pb.PostId) (*pb.ViewsDynamicArr, error) {
	dynamic, err := ss.transactor.WithTransactionValue(ctx, func(ctx context.Context) (any, error) {
		return ss.repo.GetDynamic(ctx, postId.PostId, "views")
	})
	if err != nil {
		logger.Error("")
	}

	var viewsDynamicProto pb.ViewsDynamicArr

	viewsDynamicProto.Views = make([]*pb.ViewsDynamic, 0)
	dynamicCasted := dynamic.([]statistics.Dynamic)
	for _, v := range dynamicCasted {
		viewsDynamicProto.Views = append(viewsDynamicProto.Views, &pb.ViewsDynamic{
			Views: v.Count,
			Date:  timestamppb.New(v.Time),
		})
	}

	return &viewsDynamicProto, nil
}

func (ss *StatisticsService) GetLikesDynamic(ctx context.Context, postId *pb.PostId) (*pb.LikesDynamicArr, error) {
	dynamic, err := ss.transactor.WithTransactionValue(ctx, func(ctx context.Context) (any, error) {
		return ss.repo.GetDynamic(ctx, postId.PostId, "likes")
	})
	if err != nil {
		logger.Error("")
	}

	var likesDynamicProto pb.LikesDynamicArr

	likesDynamicProto.Likes = make([]*pb.LikesDynamic, 0)
	dynamicCasted := dynamic.([]statistics.Dynamic)
	for _, v := range dynamicCasted {
		likesDynamicProto.Likes = append(likesDynamicProto.Likes, &pb.LikesDynamic{
			Likes: v.Count,
			Date:  timestamppb.New(v.Time),
		})
	}

	return &likesDynamicProto, nil
}

func (ss *StatisticsService) GetCommentsDynamic(ctx context.Context, postId *pb.PostId) (*pb.CommentsDynamicArr, error) {
	dynamic, err := ss.transactor.WithTransactionValue(ctx, func(ctx context.Context) (any, error) {
		return ss.repo.GetDynamic(ctx, postId.PostId, "comments")
	})
	if err != nil {
		logger.Error("")
	}

	var commentsDynamicProto pb.CommentsDynamicArr

	commentsDynamicProto.Comments = make([]*pb.CommentsDynamic, 0)
	dynamicCasted := dynamic.([]statistics.Dynamic)
	for _, v := range dynamicCasted {
		commentsDynamicProto.Comments = append(commentsDynamicProto.Comments, &pb.CommentsDynamic{
			Comments: v.Count,
			Date:     timestamppb.New(v.Time),
		})
	}

	return &commentsDynamicProto, nil
}

func (ss *StatisticsService) GetPostsTop(ctx context.Context, eventType *pb.Type) (*pb.PostsTop, error) {
	ids, err := ss.transactor.WithTransactionValue(ctx, func(ctx context.Context) (any, error) {
		return ss.repo.GetTop(ctx, eventType.Type, "PostId")
	})
	if err != nil {
		return nil, err
	}

	idsCasted := ids.([]int32)

	var postsTop pb.PostsTop
	postsTop.PostId = idsCasted

	return &postsTop, nil
}

func (ss *StatisticsService) GetUsersTop(ctx context.Context, eventType *pb.Type) (*pb.UsersTop, error) {
	ids, err := ss.transactor.WithTransactionValue(ctx, func(ctx context.Context) (any, error) {
		return ss.repo.GetTop(ctx, eventType.Type, "UserId")
	})
	if err != nil {
		return nil, err
	}

	idsCasted := ids.([]int32)
	var usersTop pb.UsersTop
	usersTop.UserId = idsCasted

	return &usersTop, nil
}
