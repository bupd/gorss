package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bupd/gorss/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errResponse{
		Error: msg,
	})
}

func main() {
	godotenv.Load()
	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("Port is not found in the .env")
	}

	DbUrl := os.Getenv("DB_URL")
	if DbUrl == "" {
		log.Fatal("DbUrl is not found in the .env")
	}

	conn, err := sql.Open("postgres", DbUrl)
	if err != nil {
		log.Fatal("Can't connect to the database.")
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	fmt.Println("Port:", portString)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handleErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
