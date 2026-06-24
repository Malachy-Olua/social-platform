package handlers

import (
	"net/http"

	"errors"

	"github.com/Malachy-Olua/social-platform/helpers"
	"github.com/Malachy-Olua/social-platform/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PostHandler struct {
	Store store.Storage
}

type CreatePostPayload struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Tags    []string  `json:"tags"`
	UserID  uuid.UUID `json:"user_id"`
}

func NewPostHandler(store store.Storage) *PostHandler {
	return &PostHandler{Store: store}
}

func (h *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	err := helpers.ReadJSON(w, r, &payload)
	if err != nil {
		helpers.WriteJSONError(w, http.StatusBadRequest, "invalid request payload: "+err.Error())
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		UserID:  payload.UserID,
		Tags:    payload.Tags,
	}

	ctx := r.Context()

	if err := h.Store.Posts.Create(ctx, post); err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := helpers.WriteJSON(w, http.StatusCreated, post); err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *PostHandler) GetPostHandler(w http.ResponseWriter, r *http.Request) {

	postId := chi.URLParam(r, "id")

	if postId == "" {
		helpers.WriteJSONError(w, http.StatusBadRequest, "no id provided")
		return
	}

	ctx := r.Context()

	post, err := h.Store.Posts.GetPostById(ctx, postId)
	if err != nil {

		switch {
		case errors.Is(err, store.ErrNotFound):
			helpers.WriteJSONError(w, http.StatusNotFound, err.Error())
		default:
			helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	comments := []store.Comment{}
	if comments, err = h.Store.Comments.GetCommentsByPostId(ctx, postId); err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	post.Comments = comments

	if err := helpers.WriteJSON(w, http.StatusOK, post); err != nil {
		helpers.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *PostHandler) ListPostsHandler(w http.ResponseWriter, r *http.Request)  {}
func (h *PostHandler) ShowPostHandler(w http.ResponseWriter, r *http.Request)   {}
func (h *PostHandler) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {}
func (h *PostHandler) DeletePostHandler(w http.ResponseWriter, r *http.Request) {}
