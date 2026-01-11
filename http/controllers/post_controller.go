package controllers

import (
    "fs-backend/services"
    "net/http"

    "github.com/gin-gonic/gin"
)

type PostController struct {
    service services.PostService
}

func NewPostController(service services.PostService) *PostController {
    return &PostController{service: service}
}

func (c *PostController) Create(ctx *gin.Context) {
    var req struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    post, err := c.service.CreatePost(ctx, req.Title, req.Content)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusCreated, post)
}

func (c *PostController) GetAll(ctx *gin.Context) {
    posts, err := c.service.GetAllPosts(ctx)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, posts)
}
