package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"Flow_gym_go_project/database"
	"Flow_gym_go_project/handlers"

	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from the .env file when running locally.
	// In production (Render), environment variables are usually provided directly.
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not loaded")
	}

	// Establish connection to the PostgreSQL database.
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	fmt.Println("Database connected successfully")

	// Serve static assets such as CSS, JavaScript, and images.
	http.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	// =====================================================
	// PAGE ROUTES
	// =====================================================

	// Landing page route.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Could not load homepage", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// Main application page.
	// Requires a valid user session cookie.
	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		tmpl, err := template.ParseFiles("templates/app.html")
		if err != nil {
			http.Error(w, "Could not load app page", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)
	})

	// Login page route.
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Could not load login page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// Registration page route.
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/register.html")
		if err != nil {
			http.Error(w, "Could not load register page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// =====================================================
	// API ROUTES
	// =====================================================

	// Health check endpoint used for monitoring.
	http.HandleFunc("/health", handlers.HealthHandler)

	// Recommendation engine endpoint.
	http.HandleFunc("/recommendation", handlers.RecommendationHandler(db))

	// Machine management endpoints.
	http.HandleFunc("/machines/update-availability", handlers.UpdateMachineAvailabilityHandler(db))
	http.HandleFunc("/machines", handlers.GetMachinesHandler(db))
	http.HandleFunc("/machines/available", handlers.GetAvailableMachinesHandler(db))
	http.HandleFunc("/machines/update-availability-post", handlers.UpdateMachineAvailabilityPostHandler(db))
	http.HandleFunc("/machine", handlers.GetMachineByIDHandler(db))

	// Exercise endpoints.
	http.HandleFunc("/exercises", handlers.GetExercisesHandler(db))
	http.HandleFunc("/exercise", handlers.GetExerciseByNameHandler(db))

	// Authentication endpoints.
	http.HandleFunc("/api/register", handlers.RegisterHandler(db))
	http.HandleFunc("/api/login", handlers.LoginHandler(db))
	http.HandleFunc("/api/logout", handlers.LogoutHandler())
	http.HandleFunc("/api/me", handlers.MeHandler(db))

	// Routine management endpoints.
	http.HandleFunc("/routines", handlers.GetRoutinesByUserIDHandler(db))
	http.HandleFunc("/routines/create", handlers.CreateRoutineHandler(db))
	http.HandleFunc("/routines/delete", handlers.DeleteRoutineHandler(db))
	http.HandleFunc("/routines/update", handlers.UpdateRoutineHandler(db))

	// Read the server port from environment variables.
	// Use 8080 as the default when no port is provided.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port " + port)

	// Start the HTTP server.
	log.Fatal(http.ListenAndServe(":"+port, nil))

	// =====================================================
	// DEBUG ROUTE
	// =====================================================
	// NOTE:
	// This code is currently unreachable because it is placed
	// after ListenAndServe(). It should be moved before
	// ListenAndServe() if you want it to be available.
	http.HandleFunc("/debug-db", func(w http.ResponseWriter, r *http.Request) {

		var databaseName string
		var currentUser string

		err := db.QueryRow(
			"SELECT current_database(), current_user",
		).Scan(&databaseName, &currentUser)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(
			w,
			"database=%s user=%s",
			databaseName,
			currentUser,
		)
	})
}