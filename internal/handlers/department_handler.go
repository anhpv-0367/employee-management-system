package handlers

import (
	"log"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		writeError(w, http.StatusBadRequest, "method not allowed")
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		writeError(w, http.StatusBadRequest, "name is required")
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

func (h *DepartmentHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	q := r.URL.Query()
	limit := 100
	offset := 0
	if l := q.Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := q.Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	depts, total, err := h.service.FindAll(r.Context(), limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type respDept struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}

	var out []respDept
	for _, d := range depts {
		out = append(out, respDept{ID: d.ID, Name: d.Name, CreatedAt: d.CreatedAt.Format(time.RFC3339), UpdatedAt: d.UpdatedAt.Format(time.RFC3339)})
	}

	resp := struct {
		TotalCount int64     `json:"totalCount"`
		Departments []respDept `json:"departments"`
	}{TotalCount: total, Departments: out}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
