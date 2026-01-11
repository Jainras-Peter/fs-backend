package routes

import (
    "fs-backend/http/controllers"
    "fs-backend/services"
    "github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, postService services.PostService) {
    postController := controllers.NewPostController(postService)

    api := router.Group("/api")
    {
        api.POST("/posts", postController.Create)
        api.GET("/posts", postController.GetAll)
    }
}
