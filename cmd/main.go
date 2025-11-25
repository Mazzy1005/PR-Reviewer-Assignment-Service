package main

import (
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

	prHandler := handlers.NewPullRequestHandler(
		service.NewPullRequestService(
			repository.NewPullRequestRepo(db),
			repository.NewUserRepo(db),
			repository.NewTeamRepo(db),
		),
	)

	userHandler := handlers.NewUserHandler(
		service.NewUserService(
			repository.NewUserRepo(db),
			repository.NewPullRequestRepo(db),
		),
	)

	mux.HandleFunc("POST /pullRequest/create", prHandler.Create)
	mux.HandleFunc("POST /pullRequest/reassign", prHandler.Reassign)
	mux.HandleFunc("POST /pullRequest/merge", prHandler.Merge)
	mux.HandleFunc("POST /users/setIsActive", userHandler.SetIsActive)
	mux.HandleFunc("GET /users/getReview", userHandler.GetReviews)
	mux.HandleFunc("POST /team/add", teamHandler.AddTeam)
	mux.HandleFunc("GET /team/get", teamHandler.GetTeam)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		slog.Info(r.RequestURI)
		w.Write([]byte("Hello"))
	})
	http.ListenAndServe(":8080", mux)
}
