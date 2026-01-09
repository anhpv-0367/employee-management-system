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

type EmployeeHandler struct {
	service *services.EmployeeService
}

type EmployeeResponse struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Age          int     `json:"age"`
	Position     string  `json:"position"`
	DepartmentID int64   `json:"departmentId"`
	Salary       float64 `json:"salary"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

func derefInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func derefFloat(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func NewEmployeeHandler(service *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service: service,
	}
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
  log.Println("CreateEmployee handler called")

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
    Name         string   `json:"name"`
    Email        *string  `json:"email"`
    DepartmentID int64    `json:"departmentId"`
    Age          *int     `json:"age"`
    Position     *string  `json:"position"`
    Salary       *float64 `json:"salary"`
  }

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	employee := &models.Employee{
		Name:         req.Name,
		Email:        req.Email,
		DepartmentID: req.DepartmentID,
	}

	if err := h.service.CreateEmployee(r.Context(), employee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(employee)
}


func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    return
  }

	const prefix = "/employees/"

	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.NotFound(w, r)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, prefix)
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	employee, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "employee not found", http.StatusNotFound)
		return
	}

	resp := EmployeeResponse{
    ID:           employee.ID,
    Name:         employee.Name,
    Age:          derefInt(employee.Age),
    Position:     derefString(employee.Position),
    DepartmentID: employee.DepartmentID,
    Salary:       derefFloat(employee.Salary),
    CreatedAt:    employee.CreatedAt.Format(time.RFC3339),
    UpdatedAt:    employee.UpdatedAt.Format(time.RFC3339),
  }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

