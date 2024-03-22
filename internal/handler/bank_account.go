package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/vandenbill/marketplace-10k-rps/internal/dto"
	"github.com/vandenbill/marketplace-10k-rps/internal/ierr"
	"github.com/vandenbill/marketplace-10k-rps/internal/service"
	response "github.com/vandenbill/marketplace-10k-rps/pkg/resp"
	"github.com/vandenbill/marketplace-10k-rps/pkg/validator"
)

type bankAccountHandler struct {
	bankAccountSvc *service.BankAccountService
}

func newBankAccountHandler(bankAccountSvc *service.BankAccountService) *bankAccountHandler {
	return &bankAccountHandler{bankAccountSvc}
}

func (h *bankAccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqCreateBankAccount

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

	err = h.bankAccountSvc.Create(r.Context(), req, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *bankAccountHandler) Delete(w http.ResponseWriter, r *http.Request) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	bankID := chi.URLParam(r, "bank_account_id")

	err = h.bankAccountSvc.Delete(r.Context(), bankID, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *bankAccountHandler) Update(w http.ResponseWriter, r *http.Request) {
	bankID := chi.URLParam(r, "bank_account_id")
	if bankID == "" {
		http.Error(w, "blank bank id", http.StatusNotFound)
		return
	}

	if !validator.ValidateUUID(bankID) {
		http.Error(w, "blank bank id", http.StatusNotFound)
		return
	}

	var req dto.ReqCreateBankAccount

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

	err = h.bankAccountSvc.Update(r.Context(), req, bankID, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *bankAccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	bankAccs, err := h.bankAccountSvc.Get(r.Context(), token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	res := response.SuccessReponse{}
	res.Message = "success"
	res.Data = bankAccs

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
