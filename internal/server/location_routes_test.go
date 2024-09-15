package server

import (
	_ "bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time-management/internal/database"

	_ "github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDbService struct {
	mock.Mock
}

func (m *MockDbService) GetLocations() ([]database.Location, error) {
	args := m.Called()
	return args.Get(0).([]database.Location), args.Error(1)
}

func (m *MockDbService) GetLocationById(id string) (*database.Location, error) {
	args := m.Called(id)
	return args.Get(0).(*database.Location), args.Error(1)
}

func (m *MockDbService) CreateLocation(name string) (*database.Location, error) {
	args := m.Called(name)
	return args.Get(0).(*database.Location), args.Error(1)
}

func (m *MockDbService) UpdateLocation(id, name string) (*database.Location, error) {
	args := m.Called(id, name)
	return args.Get(0).(*database.Location), args.Error(1)
}

func (m *MockDbService) DeleteLocationById(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDbService) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestLocationsHandler(t *testing.T) {
	mockDb := new(MockDbService)
	locations := []database.Location{
		{Id: "1", Name: "Location1", CreatedAt: 1633024800},
		{Id: "2", Name: "Location2", CreatedAt: 1633024900},
	}
	mockDb.On("GetLocations").Return(locations, nil)

	server := &Server{db: mockDb}

	req := httptest.NewRequest(http.MethodGet, "/locations", nil)
	r := chi.NewRouter()
	r.Get("/locations", MakeHttpHandleFunc(server.locationsHandler))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	var response []database.Location
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.ElementsMatch(t, locations, response)

	mockDb.AssertExpectations(t)
}

func TestLocationHandler(t *testing.T) {
	mockDb := new(MockDbService)
	id := "1"
	location := &database.Location{Id: id, Name: "New Location", CreatedAt: 1633025000}
	mockDb.On("GetLocationById", id).Return(location, nil)

	server := &Server{db: mockDb}

	req := httptest.NewRequest(http.MethodGet, "/locations/1", nil)
	r := chi.NewRouter()
	r.Get("/locations/{id}", MakeHttpHandleFunc(server.locationHandler))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	var response database.Location
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, *location, response)
	mockDb.AssertExpectations(t)
}

func TestCreateLocationHandler(t *testing.T) {
	mockDb := new(MockDbService)
	loc := &database.Location{Id: "1", Name: "New Location", CreatedAt: 1633025000}
	mockDb.On("CreateLocation", "New Location").Return(loc, nil)

	server := &Server{db: mockDb}

	reqBody := `{"name": "New Location"}`
	req := httptest.NewRequest(http.MethodPost, "/locations", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r := chi.NewRouter()
	r.Post("/locations", MakeHttpHandleFunc(server.createLocationHandler))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	var response database.Location
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	assert.Equal(t, loc, &response)

	mockDb.AssertExpectations(t)
}

func TestDeleteLocationHandler(t *testing.T) {
	mockDb := new(MockDbService)
	id := "1"
	mockDb.On("DeleteLocationById", id).Return(nil)

	server := &Server{db: mockDb}

	req := httptest.NewRequest(http.MethodDelete, "/locations/"+id, nil)
	r := chi.NewRouter()
	r.Delete("/locations/{id}", MakeHttpHandleFunc(server.deleteLocationHandler))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)

	mockDb.AssertExpectations(t)
}

func TestUpdateLocationHandler(t *testing.T) {
	mockDb := new(MockDbService)
	id := "1"
	updatedLocation := &database.Location{Id: id, Name: "Updated Location", CreatedAt: 1633025100}
	mockDb.On("UpdateLocation", id, "Updated Location").Return(updatedLocation, nil)

	server := &Server{db: mockDb}

	reqBody := `{"name": "Updated Location"}`
	req := httptest.NewRequest(http.MethodPut, "/locations/"+id, strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r := chi.NewRouter()
	r.Put("/locations/{id}", MakeHttpHandleFunc(server.updateLocationHandler))

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	var response database.Location
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Equal(t, updatedLocation, &response)

	mockDb.AssertExpectations(t)
}
