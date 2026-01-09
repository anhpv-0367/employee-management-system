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
	employeeService := services.NewEmployeeService(repo)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

  deptRepo := repositories.NewDepartmentRepository(db)
  deptService := services.NewDepartmentService(deptRepo)
  deptHandler := handlers.NewDepartmentHandler(deptService)


	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// POST /departments
	mux.HandleFunc("/departments", deptHandler.CreateDepartment)

	// POST /employees
	mux.HandleFunc("/employees", employeeHandler.CreateEmployee)

	// GET /employees/{id}
	mux.HandleFunc("/employees/", employeeHandler.GetByID)

	log.Println("Employee Management System running at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
