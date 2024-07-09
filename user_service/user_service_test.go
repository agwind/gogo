package user_service

import (
	"encoding/json"
	"io"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"

	"github.com/google/go-cmp/cmp"
)

func TestGetUserBadID(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/user/pear", nil)
	w := httptest.NewRecorder()

	getUser(w, req)

	res := w.Result()

	if res.StatusCode != 400 {
		t.Errorf("Expected status code 400 when id is not a number.  Actual status code %d", res.StatusCode)
	}

}

func TestGetUserNotFound(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/user/999999", nil)
	w := httptest.NewRecorder()

	getUser(w, req)

	res := w.Result()

	if res.StatusCode != 404 {
		t.Errorf("Expected status code 404 when the user is not in the database.  Actual status code %d", res.StatusCode)
	}

}

func TestGetUserExportsJSON(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	w := httptest.NewRecorder()

	getUser(w, req)

	res := w.Result()

	if res.StatusCode != 200 {
		t.Errorf("Expected status code 200 when the user is found.  Actual status code %d", res.StatusCode)
	}

	expected_user := user{Id: 1, First_name: "first", Last_name: "last", Email: "f.last@example.com"}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading body: %v", err)
	}

	var retrieved_user user
	json.Unmarshal(data, &retrieved_user)

	if !cmp.Equal(retrieved_user, expected_user) {
		t.Errorf("Retrieved user does not equal expected user: %v", retrieved_user)
	}

}

func TestMain(m *testing.M) {
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("config_test.yaml")
	if err != nil {
		panic(err)
	}

	setup_db()
	defer shutdown_db()

}
