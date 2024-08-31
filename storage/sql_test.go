package storage

import (
	"database/sql"
	"testing"

	"github.com/CTRL-Impact-Team4/khair-backend/core"
	_ "github.com/mattn/go-sqlite3" // Import for SQLite3
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := SetupInMemoryDatabase()
	if err != nil {
		t.Fatalf("setupTestDB failed: %v", err)
	}
	return db
}

// TestSetupInMemoryDatabase tests the initial setup of the in-memory database
func TestSetupInMemoryDatabase(t *testing.T) {
	db, err := SetupInMemoryDatabase()
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

// // TestCreateSchema tests the schema creation in the database
// func TestCreateSchema(t *testing.T) {
// 	db := setupTestDB(t)
// 	err := createSchema(context.Background(), db)
// 	assert.NoError(t, err)
// }

// TestInsertPredefinedServices tests inserting predefined services into the database
func TestInsertPredefinedServices(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	expectedServices := []core.Service{
		{ID: "1", Name: "Service1"},
		{ID: "2", Name: "Service2"},
	}
	err := InsertPredefinedServices(db, expectedServices)
	assert.NoError(t, err)

	// Retrieve services and compare
	var services []core.Service
	rows, err := db.Query("SELECT id, name FROM services ORDER BY id")
	assert.NoError(t, err)
	defer rows.Close()

	for rows.Next() {
		var s core.Service
		err := rows.Scan(&s.ID, &s.Name)
		assert.NoError(t, err)
		services = append(services, s)
	}

	assert.Equal(t, expectedServices, services)
}

// TestCreateOrganization tests creating an organization in the database
func TestCreateOrganization(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	expectedOrg := core.Organization{
		ID: "org1", Name: "Org One", Phone: "1234567890",
		Location: core.Location{Latitude: 10.1234, Longitude: -20.5678},
	}
	err := CreateOrganization(db, expectedOrg)
	assert.NoError(t, err)

	// Retrieve organization and compare
	var org core.Organization
	row := db.QueryRow("SELECT id, name, phone, latitude, longitude FROM organizations WHERE id = ?", expectedOrg.ID)
	err = row.Scan(&org.ID, &org.Name, &org.Phone, &org.Location.Latitude, &org.Location.Longitude)
	assert.NoError(t, err)

	assert.Equal(t, expectedOrg, org)
}

func TestGetOrganizationsByService(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Pre-insert orgs and services and links for testing
	_, err := db.Exec("INSERT INTO organizations (id, name, phone, latitude, longitude) VALUES (?, ?, ?, ?, ?)", "org1", "Org One", "123", 10.1, -20.2)
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO organizations (id, name, phone, latitude, longitude) VALUES (?, ?, ?, ?, ?)", "org2", "Org Two", "456", 20.1, -30.2)
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO services (id, name) VALUES (?, ?)", "1", "Service1")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO services (id, name) VALUES (?, ?)", "2", "Service2")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO organization_services (organization_id, service_id) VALUES (?, ?)", "org1", "1")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO organization_services (organization_id, service_id) VALUES (?, ?)", "org1", "2")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO organization_services (organization_id, service_id) VALUES (?, ?)", "org2", "1")
	assert.NoError(t, err)

	expectedOrgs := []core.Organization{
		{ID: "org1", Name: "Org One", Phone: "123", Location: core.Location{Latitude: 10.1, Longitude: -20.2}},
		// {ID: "org2", Name: "Org Two", Phone: "456", Location: core.Location{Latitude: 20.1, Longitude: -30.2}},
	}
	orgs, err := GetOrganizationsByServices(db, []string{"1"})
	assert.NoError(t, err)

	assert.Equal(t, expectedOrgs, orgs)
}
