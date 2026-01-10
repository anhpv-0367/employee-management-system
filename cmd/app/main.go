package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"app/internal/config"
	"app/internal/handlers"
	"app/internal/repositories"
	"app/internal/services"
)

func main() {
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Database connection established successfully")

	repo := repositories.NewEmployeeRepository(db)

	deptRepo := repositories.NewDepartmentRepository(db)
	deptService := services.NewDepartmentService(deptRepo)
	deptHandler := handlers.NewDepartmentHandler(deptService)

	employeeService := services.NewEmployeeService(repo, deptRepo)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)


	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// /departments: GET=list, POST=create
	mux.HandleFunc("/departments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			deptHandler.ListDepartments(w, r)
			return
		}
		if r.Method == http.MethodPost {
			deptHandler.CreateDepartment(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	// /employees: GET=list, POST=create
	mux.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			employeeHandler.ListEmployees(w, r)
			return
		}
		if r.Method == http.MethodPost {
			employeeHandler.CreateEmployee(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	// /employees/{id}: GET, PUT, DELETE
	mux.HandleFunc("/employees/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			employeeHandler.GetByID(w, r)
		case http.MethodPut:
			employeeHandler.UpdateEmployee(w, r)
		case http.MethodDelete:
			employeeHandler.DeleteEmployee(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Employee Management System running at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
