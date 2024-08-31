package core

import (
	_ "github.com/mattn/go-sqlite3"
)

type Organization struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	Location Location  `json:"location"`
	Services []Service `json:"services"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Service struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OrganizationService is the join table between Organizations and Services
type OrganizationService struct {
	ID           string        `json:"id"`
	Organization *Organization `json:"organization"`
	ServiceID    string        `json:"service_id"`
	Service      *Service      `json:"service"`
}
