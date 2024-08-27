package http

import (
	"github.com/gin-gonic/gin"
	"github.com/yrss1/todo/internal/domain/user"
	"github.com/yrss1/todo/internal/service/auth"
	"github.com/yrss1/todo/pkg/server/response"
	"net/http"
	"strings"
)

type AuthHandler struct {
	authService *auth.Service
	JWTKey      []byte
}

func NewAuthHandler(s *auth.Service, jwtKey []byte) *AuthHandler {
	return &AuthHandler{
		authService: s,
		JWTKey:      jwtKey,
	}
}

func (h *AuthHandler) Routes(r *gin.RouterGroup) {
	api := r.Group("/auth")
	{
		api.POST("/register", h.register)
		api.POST("/login", h.login)
	}
}

// register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body user.Request true "User registration data"
// @Success 200 {object} user.Response "User registered successfully"
// @Failure 400 {object} response.Object "Bad Request"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /auth/register [post]
func (h *AuthHandler) register(c *gin.Context) {
	req := user.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}
	if err := req.Validate(); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.authService.Register(c, req)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, res)
}

// login godoc
// @Summary Login a user
// @Description Login a user and receive JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body user.Request true "User login data"
// @Success 200 {object} map[string]string "JWT token"
// @Failure 400 {object} response.Object "Bad Request"
// @Failure 500 {object} response.Object "Internal Server Error"
// @Router /auth/login [post]
func (h *AuthHandler) login(c *gin.Context) {
	req := user.Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err, req)
		return
	}

	id, err := h.authService.ValidateUser(c, req)
	if err != nil {
		response.BadRequest(c, err, req)
		return
	}

	res, err := h.authService.GenerateJWT(c, id, h.JWTKey)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.OK(c, map[string]string{"token": "Bearer " + res})
}

func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := h.authService.ValidateJWT(c, tokenString, h.JWTKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
