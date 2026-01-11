package main

import (
    "fs-backend/config"
    "fs-backend/connections"
    "fs-backend/repository"
    "fs-backend/routes"
    "fs-backend/services"
    "log"

    "github.com/gin-gonic/gin"
)

func main() {
    // 1. Initialize Configuration
    config.Init()
    mongoURI := config.GetString("mongo.uri")
    mongoDB := config.GetString("mongo.database")
    port := config.GetString("server.port")

    // 2. Initialize Database
    db := connections.ConnectMongo(mongoURI, mongoDB)

    // 3. Initialize Dependency Injection (Manual)
    postRepo := repository.NewPostRepository(db)
    postService := services.NewPostService(postRepo)

    // 4. Initialize Router
    r := gin.Default()

    // 5. Register Routes
    routes.RegisterRoutes(r, postService)

    // 6. Start Server
    log.Println("Server starting on " + port)
    if err := r.Run(port); err != nil {
        log.Fatal(err)
    }
}
