package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/vandenbill/social-media-10k-rps/internal/dto"
	"github.com/vandenbill/social-media-10k-rps/internal/ierr"
	"github.com/vandenbill/social-media-10k-rps/internal/service"
)

type friendHandler struct {
	friendSvc *service.FriendService
}

func newFriendHandler(friendSvc *service.FriendService) *friendHandler {
	return &friendHandler{friendSvc}
}

func (h *friendHandler) AddFriend(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqAddFriend

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.friendSvc.AddFriend(r.Context(), req, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *friendHandler) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqDeleteFriend

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.friendSvc.DeleteFriend(r.Context(), req, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}
