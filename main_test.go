package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	 "github.com/JoddieNgatz/go-postgres-art-records-api"
)

var a main.App

func TestMain(m *testing.M) {
    a.Initialize(
        os.Getenv("DB_USERNAME"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"))

    ensureTableExists()
    code := m.Run()
    clearTable()
    os.Exit(code)
}

func ensureTableExists() {
    if _, err := a.DB.Exec(tableCreationQuery); err != nil {
        log.Fatal(err)
    }
}

func clearTable() {
    a.DB.Exec("DELETE FROM art")
    a.DB.Exec("ALTER SEQUENCE arts _id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    a.Router.ServeHTTP(rr, req)

    return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}

func TestEmptyTable(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/arts", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    if body := response.Body.String(); body != "[]" {
        t.Errorf("Expected an empty array. Got %s", body)
    }
}

func TestGetNonExistentArt(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/art/11", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusNotFound, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "Art not found" {
        t.Errorf("Expected the 'error' key of the response to be set to 'Art not found'. Got '%s'", m["error"])
    }
}

func TestCreateArt(t *testing.T) {

    clearTable()

    var jsonStr = []byte(`{"name":"test art", "price": 11.22}`)
    req, _ := http.NewRequest("POST", "/art", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    response := executeRequest(req)
    checkResponseCode(t, http.StatusCreated, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    if m["name"] != "test art" {
        t.Errorf("Expected art name to be 'test art'. Got '%v'", m["name"])
    }

    if m["price"] != 11.22 {
        t.Errorf("Expected art price to be '11.22'. Got '%v'", m["price"])
    }

    // the id is compared to 1.0 because JSON unmarshaling converts numbers to
    // floats, when the target is a map[string]interface{}
    if m["id"] != 1.0 {
        t.Errorf("Expected art ID to be '1'. Got '%v'", m["id"])
    }
}

func TestGetArt(t *testing.T) {
    clearTable()
    addArt(1)

    req, _ := http.NewRequest("GET", "/art/1", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)
}

func addArt(count int) {
    if count < 1 {
        count = 1
    }

    for i := 0; i < count; i++ {
        a.DB.Exec("INSERT INTO art(name, price) VALUES($1, $2)", "art "+strconv.Itoa(i), (i+1.0)*10)
    }
}



func TestUpdateArt(t *testing.T) {

    clearTable()
    addArt(1)

    req, _ := http.NewRequest("GET", "/art/1", nil)
    response := executeRequest(req)
    var originalart map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &originalart)

    var jsonStr = []byte(`{"name":"test art - updated name", "price": 11.22}`)
    req, _ = http.NewRequest("PUT", "/art/1", bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    if m["id"] != originalart["id"] {
        t.Errorf("Expected the id to remain the same (%v). Got %v", originalart["id"], m["id"])
    }

    if m["name"] == originalart["name"] {
        t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalart["name"], m["name"], m["name"])
    }

    if m["price"] == originalart["price"] {
        t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalart["price"], m["price"], m["price"])
    }
}


func TestDeleteArt(t *testing.T) {
    clearTable()
    addArt(1)

    req, _ := http.NewRequest("GET", "/art/1", nil)
    response := executeRequest(req)
    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("DELETE", "/art/1", nil)
    response = executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    req, _ = http.NewRequest("GET", "/art/1", nil)
    response = executeRequest(req)
    checkResponseCode(t, http.StatusNotFound, response.Code)
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS arts
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT arts_pkey PRIMARY KEY (id)
)`
