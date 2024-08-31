package main

import (
	"log"
	"net/http"

	api "github.com/CTRL-Impact-Team4/khair-backend/api"
	mw "github.com/CTRL-Impact-Team4/khair-backend/api/middleware"
	"github.com/CTRL-Impact-Team4/khair-backend/core"
	"github.com/CTRL-Impact-Team4/khair-backend/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/joho/godotenv/autoload"
)

const addr = "localhost:8080"

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	store, _ := storage.SetupInMemoryDatabase()
	storage.InsertPredefinedServices(store, []core.Service{
		{ID: "1", Name: "Bed"},
		{ID: "2", Name: "Food"},
	})

	authenticationMiddleware := r.With(
		mw.ValidateApiKey(),
		middleware.AllowContentType("application/json"),
	)

	authenticationMiddleware.Post("/orgs", api.PostOrgsHandler(store))
	authenticationMiddleware.Get("/orgs/{org_id}", api.GetOrgByID(store))
	authenticationMiddleware.Delete("/orgs/{org_id}", api.DeleteOrgByID(store))
	authenticationMiddleware.Get("/services", api.GetServices(store))
	authenticationMiddleware.Post("/orgs/{org_id}/services", api.PostServicesByOrgIDHandler(store))
	authenticationMiddleware.Get("/orgs/{org_id}/services", api.GetServicesByOrgIDHandler(store))
	authenticationMiddleware.Get("/services/nearest", api.GetNearestOrganizationHandler(store))

	log.Printf("serving http://%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
