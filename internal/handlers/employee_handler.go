package handlers

import (
	"log"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
	"database/sql"
	"os"
	"encoding/csv"
	"sync"
	"fmt"
	"path/filepath"
	"bytes"

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

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	const prefix = "/employees/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.NotFound(w, r)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, prefix)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}


	var req struct {
		Name         *string  `json:"name"`
		Email        *string  `json:"email"`
		DepartmentID *int64   `json:"departmentId"`
		Age          *int     `json:"age"`
		Position     *string  `json:"position"`
		Salary       *float64 `json:"salary"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	existing, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "employee not found")
		return
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Email != nil {
		existing.Email = req.Email
	}
	if req.DepartmentID != nil {
		existing.DepartmentID = *req.DepartmentID
	}
	if req.Age != nil {
		existing.Age = req.Age
	}
	if req.Position != nil {
		existing.Position = req.Position
	}
	if req.Salary != nil {
		existing.Salary = req.Salary
	}

	if err := h.service.Update(r.Context(), existing); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existing)
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	const prefix = "/employees/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.NotFound(w, r)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, prefix)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "employee not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func employeeToCSVRow(e *models.Employee) []string {
	age := ""
	if e.Age != nil {
		age = strconv.Itoa(*e.Age)
	}
	position := ""
	if e.Position != nil {
		position = *e.Position
	}
	email := ""
	if e.Email != nil {
		email = *e.Email
	}
	salary := ""
	if e.Salary != nil {
		salary = fmt.Sprintf("%v", *e.Salary)
	}
	return []string{
		fmt.Sprintf("%d", e.ID),
		e.Name,
		email,
		fmt.Sprintf("%d", e.DepartmentID),
		age,
		position,
		salary,
		e.CreatedAt.Format(time.RFC3339),
		e.UpdatedAt.Format(time.RFC3339),
	}
}

func writeCSV(wtr *csv.Writer, employees []*models.Employee) error {
	header := []string{"id", "name", "email", "departmentId", "age", "position", "salary", "createdAt", "updatedAt"}
	if err := wtr.Write(header); err != nil {
		return err
	}
	for _, e := range employees {
		if err := wtr.Write(employeeToCSVRow(e)); err != nil {
			return err
		}
	}
	wtr.Flush()
	return nil
}

func (h *EmployeeHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	log.Println("ExportCSV handler called")

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

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

	employees, _, err := h.service.List(r.Context(), limit, offset, deptID, keyword)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	download := q.Get("download") == "true"
	format := q.Get("format")
	if format == "" {
		format = "csv"
	}
	ts := time.Now().Unix()
	jsonFile := fmt.Sprintf("employees_%d.json", ts)
	csvFile := fmt.Sprintf("employees_%d.csv", ts)
	exportDir := os.Getenv("EXPORT_DIR")
	if exportDir == "" {
		exportDir = "."
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var csvBuf *bytes.Buffer
	var jsonBuf *bytes.Buffer

	wg.Add(2)

	go func() {
		defer wg.Done()
		data, err := json.MarshalIndent(employees, "", "  ")
		if err != nil {
			log.Println(err)
			return
		}
		if download && format == "json" {
			mu.Lock()
			jsonBuf = bytes.NewBuffer(data)
			mu.Unlock()
			return
		}
		mu.Lock()
		defer mu.Unlock()
		f, err := os.Create(filepath.Join(exportDir, jsonFile))
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		f.Write(data)
	}()

	go func() {
		defer wg.Done()
		if download && format == "csv" {
			buf := &bytes.Buffer{}
			wtr := csv.NewWriter(buf)
			if err := writeCSV(wtr, employees); err != nil {
				log.Println(err)
				return
			}
			mu.Lock()
			csvBuf = buf
			mu.Unlock()
			return
		}
		mu.Lock()
		defer mu.Unlock()
		f, err := os.Create(filepath.Join(exportDir, csvFile))
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		writeCSV(csv.NewWriter(f), employees)
	}()

	wg.Wait()

	if download {
		if format == "json" {
			if jsonBuf == nil {
				writeError(w, http.StatusInternalServerError, "json generation failed")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Disposition", "attachment; filename=\""+jsonFile+"\"")
			http.ServeContent(w, r, jsonFile, time.Now(), bytes.NewReader(jsonBuf.Bytes()))
			return
		}
		if csvBuf == nil {
			writeError(w, http.StatusInternalServerError, "csv generation failed")
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+csvFile+"\"")
		http.ServeContent(w, r, csvFile, time.Now(), bytes.NewReader(csvBuf.Bytes()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"jsonFile":  jsonFile,
		"csvFile":   csvFile,
		"exportDir": exportDir,
	})
}

