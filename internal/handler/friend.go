package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/jwtauth/v5"
	"github.com/vandenbill/social-media-10k-rps/internal/dto"
	"github.com/vandenbill/social-media-10k-rps/internal/ierr"
	"github.com/vandenbill/social-media-10k-rps/internal/service"
	response "github.com/vandenbill/social-media-10k-rps/pkg/resp"
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

func (h *friendHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var param dto.ParamGetFriends

	param.Limit, _ = strconv.Atoi(queryParams.Get("limit"))
	param.Offset, _ = strconv.Atoi(queryParams.Get("offset"))
	param.SortBy = queryParams.Get("sortBy")
	param.OrderBy = queryParams.Get("orderBy")
	param.Search = queryParams.Get("search")
	param.OnlyFriend, _ = strconv.ParseBool(queryParams.Get("onlyFriend"))

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	res, meta, err := h.friendSvc.GetFriends(r.Context(), param, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessPageReponse{}
	successRes.Message = "Get friends successfully"
	successRes.Data = res
	successRes.Meta = meta

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
