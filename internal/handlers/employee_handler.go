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

func (h *EmployeeHandler) ListEmployees(w http.ResponseWriter, r *http.Request) {
	log.Println("ListEmployees handler called")

	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// parse query params
	q := r.URL.Query()
	limit := 10
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

	var deptID *int64
	if d := q.Get("departmentId"); d != "" {
		if v, err := strconv.ParseInt(d, 10, 64); err == nil {
			deptID = &v
		} else {
			writeError(w, http.StatusBadRequest, "invalid departmentId")
			return
		}
	}

	var keyword *string
	if k := q.Get("keyword"); k != "" {
		keyword = &k
	}

	employees, total, err := h.service.List(r.Context(), limit, offset, deptID, keyword)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var out []EmployeeResponse
	for _, e := range employees {
		out = append(out, EmployeeResponse{
			ID:           e.ID,
			Name:         e.Name,
			Age:          derefInt(e.Age),
			Position:     derefString(e.Position),
			DepartmentID: e.DepartmentID,
			Salary:       derefFloat(e.Salary),
			CreatedAt:    e.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    e.UpdatedAt.Format(time.RFC3339),
		})
	}

	resp := struct {
		TotalCount int64             `json:"totalCount"`
		Employees  []EmployeeResponse `json:"employees"`
	}{
		TotalCount: total,
		Employees:  out,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
  log.Println("CreateEmployee handler called")

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
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
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	employee := &models.Employee{
		Name:         req.Name,
		Email:        req.Email,
		DepartmentID: req.DepartmentID,
	}

	if err := h.service.CreateEmployee(r.Context(), employee); err != nil {
    writeError(w, http.StatusBadRequest, err.Error())
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
		writeError(w, http.StatusBadRequest, "missing id")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	employee, err := h.service.GetByID(r.Context(), id)
	if err != nil {
    writeError(w, http.StatusNotFound, "employee not found")
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

