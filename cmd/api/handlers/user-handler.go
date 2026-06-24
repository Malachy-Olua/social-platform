package handlers

import (
	"net/http"

	"github.com/Malachy-Olua/social-platform/helpers"
	"github.com/Malachy-Olua/social-platform/internal/store"
	// "github.com/google/uuid"
)

type UserHandler struct {
	Store store.Storage
}

func NewUserHandler(store store.Storage) *UserHandler {
	return &UserHandler{Store: store}
}

type CreateUserPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
	}

	ctx := r.Context()

	if err := h.Store.Users.Create(ctx, user); err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := helpers.WriteJSON(w, http.StatusCreated, user); err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
