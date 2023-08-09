package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/evermos/boilerplate-go/internal/domain/users"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/go-chi/chi"
)

type UserHandler struct {
	UserService    users.UserService
	Authentication *middleware.Authentication
}

func ProvideUserHandler(service users.UserService, auth *middleware.Authentication) UserHandler {
	return UserHandler{
		UserService:    service,
		Authentication: auth,
	}
}

// Router sets up the router for this domain.
func (h *UserHandler) Router(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Use(h.Authentication.VerifyJWT)
		r.Get("/profiles", h.GetProfile)
		r.Group(func(r chi.Router) {
			r.Use(h.Authentication.IsAdmin)
			r.Get("/users", h.ReadUser)
		})
	})
}

func (h *UserHandler) ReadUser(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	q := r.URL.Query()
	name := q.Get("name")
	city := q.Get("city")
	province := q.Get("province")
	jobRole := q.Get("jobRole")
	status := q.Get("status")
	page, _ := strconv.Atoi(q.Get("page"))
	size, _ := strconv.Atoi(q.Get("size"))

	// Set default values for page and size if they are not provided or invalid
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 5
	}

	// Call the service to fetch the user with the specified filters and sorting
	response, err := h.UserService.ReadUser(users.UserFilter{
		Name:     name,
		City:     city,
		Province: province,
		JobRole:  jobRole,
		Status:   status,
	},
		page,
		size)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	uuid := r.Context().Value("user_id").(string)
	profile, err := h.UserService.GetProfile(uuid)
	if err != nil {
		http.Error(w, "Failed to fetch profile", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}
