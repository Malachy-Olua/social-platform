package handlers

import (
	"net/http"

	"github.com/Malachy-Olua/social-platform/helpers"
	"github.com/Malachy-Olua/social-platform/internal/store"
	"github.com/go-chi/chi/v5"

	// "github.com/google/uuid"
	"strconv"
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

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "no id provided")
		return
	}

	ctx := r.Context()

	user, err := h.Store.Users.GetUserById(ctx, userId)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := helpers.WriteJSON(w, http.StatusOK, user); err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (h *UserHandler) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	followeeId, err := strconv.ParseInt(chi.URLParam(r, "followeeId"), 10, 64)

	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "no id provided")
		return
	}

	ctx := r.Context()

	user, err := h.Store.Users.GetUserById(ctx, userId)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Store.Followers.FollowUser(ctx, user.ID.String(), strconv.FormatInt(followeeId, 10))
	if err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := helpers.WriteJSON(w, http.StatusOK, user); err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
