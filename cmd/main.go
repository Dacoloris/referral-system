package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"referral-system-test/configs"
	"referral-system-test/internal/cache"
	"referral-system-test/internal/handler"
	"referral-system-test/internal/repository"
	"referral-system-test/internal/service"
	"referral-system-test/pkg/db"
)

func main() {
	cfg := configs.LoadConfig()

	db := db.NewDb(cfg.Dsn)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.JwtSecret)

	inMemoryCache := cache.NewReferralCache()
	router := gin.Default()

	handler.NewHandler(router, userService, inMemoryCache)

	log.Printf("Server is running on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
