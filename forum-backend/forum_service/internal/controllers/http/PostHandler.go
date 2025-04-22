package http

import (
	utils "github.com/Engls/EnglsJwt"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/internal/repository"
	"github.com/Engls/forum-project2/forum_service/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

type PostHandler struct {
	postUsecase usecase.PostUsecase
	postRepo    repository.PostRepository
	jwtUtil     *utils.JWTUtil
	logger      *zap.Logger
}

func NewPostHandler(
	postUsecase usecase.PostUsecase,
	postRepo repository.PostRepository,
	jwtUtil *utils.JWTUtil,
	logger *zap.Logger,
) *PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
		postRepo:    postRepo,
		jwtUtil:     jwtUtil,
		logger:      logger,
	}
}

// CreatePost godoc
// @Summary Создать новый пост
// @Description Создает новый пост в системе
// @Tags Посты
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param post body entity.Post true "Данные поста"
// @Success 201 {object} entity.Post
// @Failure 400 {object} entity.ErrorResponse
// @Failure 401 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.logger.Warn("Authorization header required")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	if tokenString == authHeader {
		h.logger.Warn("Invalid Authorization header format", zap.String("header", authHeader))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	userID, err := h.postRepo.GetUserIDByToken(c.Request.Context(), tokenString)
	if err != nil {
		h.logger.Warn("Invalid token", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var post entity.Post
	if err := c.BindJSON(&post); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.AuthorId = userID

	h.logger.Info("Creating post", zap.Any("post", post))
	createdPost, err := h.postUsecase.CreatePost(c.Request.Context(), post)
	if err != nil {
		h.logger.Error("Failed to create post", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Post created successfully", zap.Any("createdPost", createdPost))
	c.JSON(http.StatusCreated, createdPost)
}

// GetPosts godoc
// @Summary Получить все посты
// @Description Возвращает список всех постов в системе
// @Tags Посты
// @Accept json
// @Produce json
// @Success 200 {array} entity.Post
// @Failure 500 {object} entity.ErrorResponse
// @Router /posts [get]
func (h *PostHandler) GetPosts(c *gin.Context) {
	h.logger.Info("Getting all posts")
	posts, err := h.postRepo.GetPosts(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get posts", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Successfully retrieved posts", zap.Int("count", len(posts)))
	c.JSON(http.StatusOK, posts)
}

// DeletePost godoc
// @Summary Удалить пост
// @Description Удаляет пост по ID (доступно автору или администратору)
// @Tags Посты
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID поста"
// @Success 204 "No Content"
// @Failure 400 {object} entity.ErrorResponse
// @Failure 401 {object} entity.ErrorResponse
// @Failure 403 {object} entity.ErrorResponse
// @Failure 404 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.logger.Warn("Authorization header required")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	if tokenString == authHeader {
		h.logger.Warn("Invalid Authorization header format", zap.String("header", authHeader))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.logger.Warn("Invalid post ID", zap.String("postID", postIDStr), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	userID, err := h.jwtUtil.GetUserIDFromToken(tokenString)
	if err != nil {
		h.logger.Warn("Invalid token or user ID", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or user ID"})
		return
	}

	userRole, err := h.jwtUtil.GetRoleFromToken(tokenString)
	if err != nil {
		h.logger.Warn("Invalid token or user role", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or user role"})
		return
	}

	if userRole != "admin" {
		post, err := h.postRepo.GetPostByID(c.Request.Context(), postID)
		if err != nil {
			h.logger.Error("Failed to get post", zap.Int("postID", postID), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get post"})
			return
		}

		if post.AuthorId != userID {
			h.logger.Warn("Unauthorized attempt to delete post",
				zap.Int("userID", userID),
				zap.Int("postAuthorID", post.AuthorId),
				zap.Int("postID", postID))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this post"})
			return
		}
	}

	h.logger.Info("Deleting post", zap.Int("postID", postID))
	err = h.postUsecase.DeletePost(c.Request.Context(), postID)
	if err != nil {
		h.logger.Error("Failed to delete post", zap.Int("postID", postID), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	h.logger.Info("Post deleted successfully", zap.Int("postID", postID))
	c.Status(http.StatusNoContent)
}
