package main

import (
	"fmt"
	"golang-task-manager/config"
	"golang-task-manager/internal/handler"
	"golang-task-manager/internal/logger"
	"golang-task-manager/internal/middleware"
	"golang-task-manager/internal/migrations"
	"golang-task-manager/internal/repository"
	"golang-task-manager/internal/service"
	"golang-task-manager/internal/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	ServerConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		os.Exit(1)
	}

	fmt.Println("Config loaded successfully")
	fmt.Println(ServerConfig)

	db, err := sqlx.Connect("postgres", ServerConfig.DB_URL)
	if err != nil {
		fmt.Println("failed to connect to db:", err)
		os.Exit(1)
	}
	defer db.Close()

	migrations.RunMigrations(db.DB)

	fmt.Println("Database ready")

	rateLimit := middleware.RateLimiter(ServerConfig.REDIS_URL,ServerConfig.REDIS_PASSWORD,0, ServerConfig.RATE_LIMIT)

	logger, err := logger.BuildLogger(ServerConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger.Info("Logger initialised")
	
	utils.InitJwt(ServerConfig.JWT_SECRET)
	
	userRepo := repository.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepo, ServerConfig.JWT_EXPIRATION_HOURS)
	userHandler := handler.NewUserHandler(userService)
	
	taskRepo := repository.NewPostgresTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)
	
	auth := middleware.Auth()
	
	r := gin.Default()

	r.Use(rateLimit)

	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		tasks := api.Group("/tasks")
		{
			tasks.POST("",auth, taskHandler.Create)
			tasks.GET("",auth, taskHandler.List)
			tasks.GET("/:id",auth, taskHandler.Get)
			tasks.PUT("/:id",auth, taskHandler.Update)
			tasks.DELETE("/:id",auth, taskHandler.Delete)
		}
		api.GET("/public/:slug", taskHandler.GetPublic)
	}

	r.Run(":" + ServerConfig.PORT)

}
