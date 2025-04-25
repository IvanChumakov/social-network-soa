package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	customerror "social-network/posts-comments-service/internal/errors"
	"social-network/posts-comments-service/internal/logger"

	"github.com/uptrace/bun"
)

type PostRepository struct {
	db *bun.DB
}

func NewPostRepository(db *bun.DB) *PostRepository {
	return &PostRepository{db}
}

func (pr *PostRepository) AddPost(post Post) error {
	_, err := pr.db.NewInsert().
		Model(&post).
		Exec(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("error adding post: %v", err))
		return err
	}

	return nil
}

func (pr *PostRepository) DeletePost(id int32) error {
	_, err := pr.GetPostById(id)
	if err != nil {
		if errors.Is(err, &customerror.NotFoundError{}) {
			return &customerror.NotFoundError{}
		}
	}

	_, err = pr.db.NewDelete().
		Model((*Post)(nil)).
		Where("id = ?", id).
		Exec(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Info("not fount")
			return &customerror.NotFoundError{}
		}
		logger.Error(fmt.Sprintf("error deleting post: %v", err))
		return err
	}

	return nil
}

func (pr *PostRepository) UpdatePost(post Post) error {
	_, err := pr.GetPostById(post.Id)
	if err != nil {
		if errors.Is(err, &customerror.NotFoundError{}) {
			logger.Info("not found")
			return &customerror.NotFoundError{}
		}
	}

	_, err = pr.db.NewUpdate().
		Model(&post).
		Where("id = ?", post.Id).
		OmitZero().
		Exec(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("error updating post: %v", err))
		return err
	}

	return nil
}

func (pr *PostRepository) GetPostById(id int32) (Post, error) {
	var post Post
	err := pr.db.NewSelect().
		Model(&post).
		Where("id = ?", id).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Info("not found")
			return Post{}, &customerror.NotFoundError{}
		}
		logger.Error(fmt.Sprintf("error getting post: %v", err))
		return Post{}, err
	}

	return post, nil
}

func (pr *PostRepository) GetAllPosts(limit int32, offset int32, userId int32) ([]Post, error) {
	var posts []Post
	err := pr.db.NewSelect().
		Model(&posts).
		Where("is_private = ?", false).
		WhereOr("creator_id = ?", int(userId)).
		Order("created_at DESC").
		Limit(int(limit)).
		Offset(int(offset)).
		Scan(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("error getting all posts: %v", err))
		return nil, err
	}

	return posts, nil
}

func (pr *PostRepository) AddComment(comment Comment) error {
	_, err := pr.db.NewInsert().
		Model(&comment).
		Exec(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("error adding comment: %v", err))
		return err
	}

	return nil
}

func (pr *PostRepository) GetAllComments(limit int32, offset int32, postId int32) ([]Comment, error) {
	var comments []Comment

	err := pr.db.NewSelect().
		Model(&comments).
		Where("post_id = ?", postId).
		Order("created_at DESC").
		Limit(int(limit)).
		Offset(int(offset)).
		Scan(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("error getting all comments: %v", err))
	}

	return comments, nil
}
