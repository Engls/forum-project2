package http

import (
	"forum-project/forum_service/internal/entity"
	"forum-project/forum_service/internal/repository"
	"forum-project/forum_service/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type PostHandler struct {
	postUsecase usecase.PostUsecase
	postRepo    repository.PostRepository
}

func NewPostHandler(postUsecase usecase.PostUsecase, postRepo repository.PostRepository) *PostHandler {
	return &PostHandler{postUsecase: postUsecase, postRepo: postRepo}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	if tokenString == authHeader {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	userID, err := h.postRepo.GetUserIDByToken(c.Request.Context(), tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var post entity.Post
	if err := c.BindJSON(&post); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.AuthorId = userID

	createdPost, err := h.postUsecase.CreatePost(c.Request.Context(), post)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPost)
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := h.postRepo.GetPosts(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}
