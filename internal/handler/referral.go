package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"referral-system-test/internal/cache"
	"referral-system-test/internal/payload"
	"referral-system-test/internal/service"
	"referral-system-test/pkg/middleware"
	"strconv"
)

type UserHandler struct {
	userService   service.UserService
	referralCache *cache.ReferralCache
}

func NewHandler(router *gin.Engine, userService service.UserService, referralCache *cache.ReferralCache) {
	handler := &UserHandler{
		userService:   userService,
		referralCache: referralCache,
	}

	router.Use(middleware.LoggerMiddleware())

	router.POST("/login", handler.Login)
	router.POST("/register", handler.Register)
	router.POST("/referral-code", middleware.WrapAuth(handler.CreateReferralCode))
	router.DELETE("/referral-code", middleware.WrapAuth(handler.DeleteReferralCode))
	router.GET("/referral-code/:email", middleware.WrapAuth(handler.GetReferralCode))
	router.GET("/referrals/:user_id", middleware.WrapAuth(handler.GetReferrals))
}

func (h *UserHandler) Register(c *gin.Context) {
	var req payload.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := h.userService.RegisterUser(req.Email, req.Password, req.ReferrerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req payload.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	token, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong email or password"})
		return
	}
	c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) CreateReferralCode(c *gin.Context) {
	var req payload.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	code := h.userService.GenerateReferralCode(req.UserID, req.ExpiryDate)
	h.referralCache.Set(req.Email, code)

	c.JSON(http.StatusOK, gin.H{"referral_code": code})
}

func (h *UserHandler) GetReferralCode(c *gin.Context) {
	userEmail := c.Param("email")

	if cachedCode, exists := h.referralCache.Get(userEmail); exists {
		c.JSON(http.StatusOK, gin.H{"referral_code": cachedCode})
		return
	}

	code, err := h.userService.GetReferralCode(userEmail)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Referral code not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"referral_code": code.Code})
}

func (h *UserHandler) GetReferrals(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong user_id"})
	}
	referrals, err := h.userService.GetReferrals(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch referrals"})
		return
	}

	c.JSON(http.StatusOK, referrals)
}

func (h *UserHandler) DeleteReferralCode(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong user_id"})
	}

	if err := h.userService.DeleteReferralCode(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete referral code"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
