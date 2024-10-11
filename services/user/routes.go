package user

import (
	"net/http"
	"strings"

	"github.com/API/services/auth"
	"github.com/API/types"
	"github.com/API/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.POST("/signin", h.handleSignIn)
	router.POST("/login", h.handleLogIn)

}
func (h *Handler) handleSignIn(c *gin.Context) {
	var payload types.SignInPayload

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		utils.HandleError(c, err, "failed to read body", http.StatusBadRequest)
		return

	}
	if !strings.Contains(payload.Email, "@") {
		utils.HandleError(c, nil, "failed to read body", http.StatusBadRequest)
		return

	}
	_, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.HandleError(c, err, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.HandleError(c, err, "failed to hash password", http.StatusInternalServerError)
	}

	if err := h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	}); err != nil {
		utils.HandleError(c, err, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": "created account",
	})

}
func (h *Handler) handleLogIn(c *gin.Context) {

}
