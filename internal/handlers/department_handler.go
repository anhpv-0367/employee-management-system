package handlers

import (
  "log"
  "encoding/json"
  "net/http"
  "strings"

	"app/internal/models"
	"app/internal/services"
)

type DepartmentHandler struct {
	service *services.DepartmentService
}

func NewDepartmentHandler(service *services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{
		service: service,
	}
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateDepartment handler called")

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	dept := &models.Department{
		Name: req.Name,
	}

	if err := h.service.Create(r.Context(), dept); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dept)
}
