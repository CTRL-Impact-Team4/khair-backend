package api

import (
	"database/sql"
	"encoding/json"
	"math"
	"net/http"

	core "github.com/CTRL-Impact-Team4/khair-backend/core"
	"github.com/CTRL-Impact-Team4/khair-backend/storage"
	"github.com/go-chi/chi/v5"
)

func GetServices(store *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services, err := storage.GetPredefinedServices(store)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		jsonResponse, err := json.Marshal(services)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

func PostServicesByOrgIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := chi.URLParam(r, "org_id")

		_, err := storage.GetOrganizationByID(db, orgID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		var serviceIDs []string
		if err := json.NewDecoder(r.Body).Decode(&serviceIDs); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Validate that services exist in the predefined list
		services, err := storage.GetServicesByID(db, serviceIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Associate the services with the organization
		err = storage.AddServicesToOrganization(db, orgID, serviceIDs)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Return the services that were associated with the organization
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(services)
	}
}

func GetServicesByOrgIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := chi.URLParam(r, "org_id")

		_, err := storage.GetOrganizationByID(db, orgID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		services, err := storage.GetServicesByOrganizationID(db, orgID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		if len(services) == 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(services)
	}
}

// Haversine function to calculate distance between two coordinates
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers
	lat1Rad, lon1Rad := lat1*math.Pi/180, lon1*math.Pi/180
	lat2Rad, lon2Rad := lat2*math.Pi/180, lon2*math.Pi/180

	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func GetNearestOrganizationHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Services  []string `json:"services"`
			Latitude  float64  `json:"latitude"`
			Longitude float64  `json:"longitude"`
		}

		type organizationWithDistance struct {
			core.Organization
			Distance float64 `json:"distance"`
		}

		// Parse the request body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// Validate the services
		services, err := storage.GetServicesByID(db, req.Services)
		if err != nil {
			http.Error(w, "One or more services do not exist", http.StatusBadRequest)
			return
		}

		// Get organizations offering all specified services
		orgs, err := storage.GetOrganizationsByServices(db, req.Services)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		// Find the closest organization
		var closestOrg *organizationWithDistance
		minDistance := math.MaxFloat64

		for _, org := range orgs {
			distance := haversine(req.Latitude, req.Longitude, org.Location.Latitude, org.Location.Longitude)
			if distance < minDistance {
				minDistance = distance
				closestOrg = &organizationWithDistance{
					Organization: org,
					Distance:     distance,
				}
			}
		}

		// Check if no organization was found
		if closestOrg == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		closestOrg.Services = services

		// Return the closest organization
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(closestOrg)
	}
}
