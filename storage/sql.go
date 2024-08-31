package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/CTRL-Impact-Team4/khair-backend/core"
	_ "github.com/mattn/go-sqlite3" // Import for SQLite3
)

func SetupInMemoryDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Execute table creation statements here
	statements := []string{
		`CREATE TABLE organizations (id TEXT PRIMARY KEY, name TEXT, phone TEXT, latitude REAL, longitude REAL)`,
		`CREATE TABLE services (id TEXT PRIMARY KEY, name TEXT)`,
		`CREATE TABLE organization_services (organization_id TEXT, service_id TEXT, PRIMARY KEY (organization_id, service_id), FOREIGN KEY (organization_id) REFERENCES organizations(id), FOREIGN KEY (service_id) REFERENCES services(id))`,
	}

	for _, stmt := range statements {
		_, err := db.Exec(stmt)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func InsertPredefinedServices(db *sql.DB, services []core.Service) error {
	stmt, err := db.Prepare("INSERT INTO services (id, name) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, service := range services {
		_, err := stmt.Exec(service.ID, service.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateOrganization(db *sql.DB, org core.Organization) error {
	_, err := db.Exec("INSERT INTO organizations (id, name, phone, latitude, longitude) VALUES (?, ?, ?, ?, ?)", org.ID, org.Name, org.Phone, org.Location.Latitude, org.Location.Longitude)
	return err
}

func AddServicesToOrganization(db *sql.DB, orgID string, serviceIDs []string) error {
	stmt, err := db.Prepare("INSERT INTO organization_services (organization_id, service_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, serviceID := range serviceIDs {
		_, err := stmt.Exec(orgID, serviceID)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetOrganizationsByServices(db *sql.DB, serviceIDs []string) ([]core.Organization, error) {
	// Create placeholders for the IN clause
	placeholders := make([]string, len(serviceIDs))
	args := make([]interface{}, len(serviceIDs))
	for i, id := range serviceIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	// Query to find organizations that offer all specified services
	query := fmt.Sprintf(`
		SELECT o.id, o.name, o.phone, o.latitude, o.longitude
		FROM organizations o
		JOIN organization_services os ON o.id = os.organization_id
		WHERE os.service_id IN (%s)
		GROUP BY o.id, o.name, o.phone, o.latitude, o.longitude
		HAVING COUNT(DISTINCT os.service_id) = ?
	`, strings.Join(placeholders, ","))

	// Add the count of service IDs to the arguments list
	args = append(args, len(serviceIDs))

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []core.Organization
	for rows.Next() {
		var org core.Organization
		if err := rows.Scan(&org.ID, &org.Name, &org.Phone, &org.Location.Latitude, &org.Location.Longitude); err != nil {
			return nil, err
		}
		organizations = append(organizations, org)
	}

	return organizations, nil
}

func GetOrganizationByID(db *sql.DB, orgID string) (core.Organization, error) {
	var org core.Organization
	err := db.QueryRow(`
        SELECT id, name, phone, latitude, longitude
        FROM organizations
        WHERE id = ?
    `, orgID).Scan(&org.ID, &org.Name, &org.Phone, &org.Location.Latitude, &org.Location.Longitude)
	if err != nil {
		return core.Organization{}, err
	}
	return org, nil
}

func DeleteOrganizationByID(db *sql.DB, orgID string) error {
	result, err := db.Exec("DELETE FROM organizations WHERE id = ?", orgID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func GetPredefinedServices(db *sql.DB) ([]core.Service, error) {
	var services []core.Service
	rows, err := db.Query("SELECT id, name FROM services")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var svc core.Service
		if err := rows.Scan(&svc.ID, &svc.Name); err != nil {
			return nil, err
		}
		services = append(services, svc)
	}

	return services, nil
}

func GetServicesByOrganizationID(db *sql.DB, orgID string) ([]core.Service, error) {
	rows, err := db.Query(`
		SELECT s.id, s.name
		FROM services s
		JOIN organization_services os ON s.id = os.service_id
		WHERE os.organization_id = ?
	`, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []core.Service
	for rows.Next() {
		var svc core.Service
		if err := rows.Scan(&svc.ID, &svc.Name); err != nil {
			return nil, err
		}
		services = append(services, svc)
	}
	return services, nil
}

func GetServicesByID(db *sql.DB, serviceIDs []string) ([]core.Service, error) {
	placeholders := make([]string, len(serviceIDs))
	args := make([]interface{}, len(serviceIDs))
	for i, id := range serviceIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("SELECT id, name FROM services WHERE id IN (%s)", strings.Join(placeholders, ","))
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []core.Service
	for rows.Next() {
		var svc core.Service
		if err := rows.Scan(&svc.ID, &svc.Name); err != nil {
			return nil, err
		}
		services = append(services, svc)
	}

	// Check if all service IDs were found
	if len(services) != len(serviceIDs) {
		return nil, fmt.Errorf("one or more services do not exist")
	}

	return services, nil
}
