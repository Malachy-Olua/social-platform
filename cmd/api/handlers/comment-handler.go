package handlers

import (
	"net/http"

	"github.com/Malachy-Olua/social-platform/helpers"
	"github.com/Malachy-Olua/social-platform/internal/store"
	"github.com/google/uuid"
)

type CommentsHandler struct {
	Store store.Storage
}

func NewCommnetsHandler(store store.Storage) *CommentsHandler {
	return &CommentsHandler{Store: store}
}

type CreateCommentsPayload struct {
	PostID  uuid.UUID `json:"post_id"`
	UserID  uuid.UUID `json:"user_id"`
	Content string    `json:"content"`
}

func (h *CommentsHandler) CreateCommentsHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentsPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	comment := &store.Comment{
		PostID:  payload.PostID,
		Content: payload.Content,
		UserID:  payload.UserID,
	}

	ctx := r.Context()

	if err := h.Store.Comments.Create(ctx, comment); err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := helpers.WriteJSON(w, http.StatusCreated, comment); err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
