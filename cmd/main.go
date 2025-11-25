package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/database"
	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/handlers"
	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/repository"
	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/service"
)

func main() {

	db, err := database.NewPostgresDB()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	mux := http.NewServeMux()
	teamHandler := handlers.NewTeamHandler(service.NewTeamService(repository.NewTeamRepo(db)))

	mux.HandleFunc("POST /team/add", teamHandler.AddTeam)
	mux.HandleFunc("GET /team/get", teamHandler.GetTeam)

	res, err := db.Exec("INSERT INTO teams VALUES ($1), ('Name 2'), ('Surprise Name 3')", "Name 1")
	if err != nil {
		slog.Error(err.Error())
	}
	rows, _ := res.RowsAffected()
	slog.Info(fmt.Sprintf("Insert %v rows", rows))

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		slog.Info(r.RequestURI)
		w.Write([]byte("Hello"))
	})
	http.ListenAndServe(":8080", mux)
}
