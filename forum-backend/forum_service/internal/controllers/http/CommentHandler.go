package http

import (
	utils "github.com/Engls/EnglsJwt"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

type CommentHandler struct {
	commentUsecase usecase.CommentsUsecases
	jwtUtil        *utils.JWTUtil
	logger         *zap.Logger
}

func NewCommentHandler(commentUsecase usecase.CommentsUsecases, jwtUtil *utils.JWTUtil, logger *zap.Logger) *CommentHandler {
	return &CommentHandler{commentUsecase: commentUsecase, jwtUtil: jwtUtil, logger: logger}
}

// CreateComment godoc
// @Summary Создать новый комментарий
// @Description Создает новый комментарий к указанному посту
// @Tags Комментарии
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID поста"
// @Param comment body entity.Comment true "Данные комментария"
// @Success 201 {object} entity.Comment
// @Failure 400 {object} entity.ErrorResponse
// @Failure 401 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /posts/{id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	postIDStr := c.Param("id") // Получаем ID поста из URL
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.logger.Error("Invalid post ID", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var comment entity.Comment
	if err := c.BindJSON(&comment); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.logger.Error("Authorization header required")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	if tokenString == authHeader {
		h.logger.Error("Invalid Authorization header format")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	// Получаем ID пользователя из токена
	userID, err := h.jwtUtil.GetUserIDFromToken(tokenString)
	if err != nil {
		h.logger.Error("Invalid token or user ID", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or user ID"})
		return
	}

	comment.PostId = postID
	comment.AuthorId = userID // Устанавливаем AuthorId из userID

	createdComment, err := h.commentUsecase.CreateComment(c.Request.Context(), comment)
	if err != nil {
		h.logger.Error("Failed to create comment", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Comment created successfully", zap.Int("postID", postID), zap.Int("userID", userID))
	c.JSON(http.StatusCreated, createdComment)
}

// GetCommentsByPostID godoc
// @Summary Получить комментарии поста
// @Description Возвращает все комментарии для указанного поста
// @Tags Комментарии
// @Accept json
// @Produce json
// @Param id path int true "ID поста"
// @Success 200 {array} entity.Comment
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /posts/{id}/comments [get]
func (h *CommentHandler) GetCommentsByPostID(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.logger.Error("Invalid post ID", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	comments, err := h.commentUsecase.GetCommentByPostID(c.Request.Context(), postID)
	if err != nil {
		h.logger.Error("Failed to get comments", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Comments retrieved successfully", zap.Int("postID", postID))
	c.JSON(http.StatusOK, comments)
}
