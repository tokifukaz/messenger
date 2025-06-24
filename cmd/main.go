package main

import (
    "log"

    "messenger/internal/api"
    _ "messenger/docs"
)

// @title Insider Message API
// @version 1.0
// @description This is an automatic message sender system.
// @host localhost:8080
// @BasePath /api

func main() {
    router := api.SetupRouter()
    log.Println("Server started on :8080")
    router.Run(":8080")
}