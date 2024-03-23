package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/vandenbill/social-media-10k-rps/internal/cfg"
	"github.com/vandenbill/social-media-10k-rps/internal/service"
)

// var (
// 	requestsTotal = prometheus.NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: "http_requests_total",
// 			Help: "Total number of HTTP requests.",
// 		},
// 		[]string{"method", "path", "status"},
// 	)
// 	requestDuration = prometheus.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "http_request_duration_seconds",
// 			Help:    "Histogram of request duration in seconds.",
// 			Buckets: prometheus.DefBuckets,
// 		},
// 		[]string{"method", "path", "status"},
// 	)
// )

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Status() int {
	return rw.statusCode
}

type Handler struct {
	router  *chi.Mux
	service *service.Service
	cfg     *cfg.Cfg
}

func NewHandler(router *chi.Mux, service *service.Service, cfg *cfg.Cfg) *Handler {
	handler := &Handler{router, service, cfg}
	handler.registRoute()

	return handler
}

func (h *Handler) registRoute() {
	// prometheus.MustRegister(requestsTotal, requestDuration)

	r := h.router
	var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(h.cfg.JWTSecret), nil, jwt.WithAcceptableSkew(30*time.Second))

	userH := newUserHandler(h.service.User)
	fileH := newFileHandler(h.cfg)
	friendH := newFriendHandler(h.service.Friend)
	postH := newPostHandler(h.service.Post)

	// r.Use(middleware.RedirectSlashes)
	// r.Use(prometheusMiddleware)

	// r.Get("/metrics", func(h http.Handler) http.HandlerFunc {
	// 	return func(w http.ResponseWriter, r *http.Request) {
	// 		h.ServeHTTP(w, r)
	// 	}
	// }(promhttp.Handler()))

	r.Post("/v1/user/register", userH.Register)
	r.Post("/v1/user/login", userH.Login)

	// protected route
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/v1/user/link", userH.LinkEmail)
		r.Post("/v1/user/link/phone", userH.LinkPhone)
		r.Patch("/v1/user", userH.UpdateAccount)

		r.Get("/v1/friend", friendH.GetFriends)
		r.Post("/v1/friend", friendH.AddFriend)
		r.Delete("/v1/friend", friendH.DeleteFriend)

		r.Post("/v1/post", postH.AddPost)

		r.Post("/v1/post/comment", postH.AddComment)

		r.Post("/v1/image", fileH.Upload)
	})
}

// func prometheusMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		startTime := time.Now()
// 		rw := newResponseWriter(w)
// 		defer func() {
// 			status := rw.Status()
// 			duration := time.Since(startTime).Seconds()
// 			requestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(status)).Inc()
// 			requestDuration.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(status)).Observe(duration)
// 		}()
// 		next.ServeHTTP(rw, r)
// 	})
// }
