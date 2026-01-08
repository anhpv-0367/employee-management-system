package main

import (
	"fmt"
	"log"
	"net/http"
	_ "github.com/lib/pq"
	"app/internal/config"
)

func main() {
	// http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintln(w, "OK")
	// })

	// log.Println("Employee Management System running at :8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
	db, err := config.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

  log.Println("Database connection established successfully")

  mux := http.NewServeMux()

  mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintln(w, "OK")
  })

  log.Println("Employee Management System running at :8080")
  log.Fatal(http.ListenAndServe(":8080", mux))
}
