package services

import (
    "context"
    "fs-backend/models"
    "fs-backend/repository"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService interface {
    CreatePost(ctx context.Context, title, content string) (*models.Post, error)
    GetAllPosts(ctx context.Context) ([]models.Post, error)
}

type postService struct {
    repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
    return &postService{repo: repo}
}

func (s *postService) CreatePost(ctx context.Context, title, content string) (*models.Post, error) {
    post := &models.Post{
        ID:      primitive.NewObjectID(),
        Title:   title,
        Content: content,
    }
    if err := s.repo.Create(ctx, post); err != nil {
        return nil, err
    }
    return post, nil
}

func (s *postService) GetAllPosts(ctx context.Context) ([]models.Post, error) {
    return s.repo.FindAll(ctx)
}
