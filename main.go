package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"Flow_gym_go_project/database"
	"Flow_gym_go_project/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not loaded")
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	fmt.Println("Database connected successfully")

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Could not load homepage", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/recommendation", handlers.RecommendationHandler(db))
	http.HandleFunc("/machines/update-availability", handlers.UpdateMachineAvailabilityHandler(db))
	
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}