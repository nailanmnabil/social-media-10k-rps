package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/vandenbill/marketplace-10k-rps/internal/cfg"
	"github.com/vandenbill/marketplace-10k-rps/internal/dto"
	"github.com/vandenbill/marketplace-10k-rps/internal/ierr"
	"github.com/vandenbill/marketplace-10k-rps/internal/service"
	response "github.com/vandenbill/marketplace-10k-rps/pkg/resp"
)

type productHandler struct {
	productSvc *service.ProductService
	cfg        *cfg.Cfg
}

func newProductHandler(productSvc *service.ProductService, cfg *cfg.Cfg) *productHandler {
	return &productHandler{productSvc, cfg}
}

func (h *productHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqCreateProduct

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("%#v\n", req)

	if req.ImageURL == "http://incomplete" {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.productSvc.Create(r.Context(), req, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *productHandler) Delete(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "product_id")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.productSvc.Delete(r.Context(), productID, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *productHandler) ChangeStock(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqChangeStock

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	productID := chi.URLParam(r, "product_id")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.productSvc.ChangeStock(r.Context(), req, productID, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *productHandler) Buy(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqBuy

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	productID := chi.URLParam(r, "product_id")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.productSvc.Buy(r.Context(), req, productID, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *productHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.ReqUpdateProduct

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	productID := chi.URLParam(r, "product_id")

	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, "failed to get token from request", http.StatusBadRequest)
		return
	}

	err = h.productSvc.Update(r.Context(), req, productID, token.Subject())
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *productHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "product_id")

	data, err := h.productSvc.GetByID(r.Context(), productID)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessReponse{Message: "ok", Data: data}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *productHandler) GetWithFilter(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	userOnly, _ := strconv.ParseBool(queryParams.Get("userOnly"))
	limit, _ := strconv.Atoi(queryParams.Get("limit"))
	offset, _ := strconv.Atoi(queryParams.Get("offset"))
	tags := strings.Split(queryParams.Get("tags"), ",")
	condition := queryParams.Get("condition")
	showEmptyStock, _ := strconv.ParseBool(queryParams.Get("showEmptyStock"))
	maxPrice, _ := strconv.Atoi(queryParams.Get("maxPrice"))
	minPrice, _ := strconv.Atoi(queryParams.Get("minPrice"))
	sortBy := queryParams.Get("sortBy")
	orderBy := queryParams.Get("orderBy")
	search := queryParams.Get("search")

	filter := dto.SearchProductFilter{
		UserOnly:       userOnly,
		Limit:          limit,
		Offset:         offset,
		Tags:           tags,
		Condition:      condition,
		ShowEmptyStock: showEmptyStock,
		MaxPrice:       maxPrice,
		MinPrice:       minPrice,
		SortBy:         sortBy,
		OrderBy:        orderBy,
		Search:         search,
	}
	filter.SetDefault()

	if userOnly {
		sub, err := h.verifyJwt(r, h.cfg.JWTSecret)
		if err != nil {
			code, msg := ierr.TranslateError(err)
			http.Error(w, msg, code)
			return
		}
		filter.Sub = sub
	}

	res, meta, err := h.productSvc.GetWithFilter(r.Context(), filter)
	if err != nil {
		code, msg := ierr.TranslateError(err)
		http.Error(w, msg, code)
		return
	}

	successRes := response.SuccessPageReponse{Message: "ok", Data: res, Meta: meta}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(successRes)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *productHandler) verifyJwt(r *http.Request, secretKey string) (sub string, err error) {
	authHeader := r.Header.Get("Authorization")

	bearerToken := ""
	if authHeader != "" {
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) == 2 && authHeaderParts[0] == "Bearer" {
			bearerToken = authHeaderParts[1]
		}
	} else {
		return "", ierr.ErrForbidden
	}

	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return sub, errors.Wrap(err, "failed parse jwt")
	}

	if !token.Valid {
		return sub, errors.Wrap(err, "token not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return sub, errors.Wrap(err, "failed cast jwt")
	}
	sub = claims["sub"].(string)

	return sub, nil
}
