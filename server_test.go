package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var s *server

func TestMain(m *testing.M) {

	s = newServer()

	code := m.Run()
	os.Exit(code)
}

func TestHealthy(t *testing.T) {

	request, err := http.NewRequest("GET", "/healthy", nil)
	if err != nil {
		t.Fatal(err)
	}

	response := executeRequest(request)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestPOST(t *testing.T) {

	var testPost = []struct {
		name        string
		path        string
		description string
		content     string
		code        int
	}{
		{
			"send post request",
			"/",
			"my description",
			"my content",
			http.StatusOK,
		},
	}

	for _, testData := range testPost {

		id := postRequest(t, testData.path, `{"description":"`+testData.description+`","content":"`+testData.content+`"}`, testData.name)

		text, ok := s.entries[id]
		if !ok {
			t.Errorf("%v not found", id)
		}

		if text.Description != testData.description {
			t.Errorf("\nDescription expected: %v\n Got: %v\n", testData.description, text.Description)
		}
		if text.Content != testData.content {
			t.Errorf("\nContent expected: %v\n Got: %v\n", testData.content, text.Content)
		}
	}
}

func TestGET(t *testing.T) {

	id := "any-id"
	s.entries = make(map[string]textData)
	s.entries[id] = textData{
		Description: "any-description",
		Content:     "any-content",
	}

	var testGet = []struct {
		name string
		path string
		code int
		body string
	}{
		{
			"read existing entry",
			"/paste/" + id,
			http.StatusOK,
			`{"id":"` + id + `","description":"any-description","content":"any-content"}`,
		},
		{
			"read non-existing entry",
			"/paste/bad-id",
			http.StatusNotFound,
			``,
		},
	}

	for _, testData := range testGet {

		request, err := http.NewRequest("GET", testData.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		response := executeRequest(request)
		checkResponseCode(t, testData.code, response.Code, testData.name)
		checkResponseBody(t, testData.body, response.Body.String(), testData.name)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	return rec
}

func checkResponseCode(t *testing.T, expected, actual int, testName ...string) {
	if expected != actual {
		t.Errorf("Error in test%s: expected response code %d. Got %d\n", testName, expected, actual)
	}
}

func checkResponseBody(t *testing.T, expected, actual string, testName ...string) {
	if expected != actual {
		t.Errorf("Error in test%s: expected response body %v. Got %v\n", testName, expected, actual)
	}
}

func postRequest(t *testing.T, path, jsonStr string, testName ...string) string {
	request, err := http.NewRequest("POST", path, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")

	response := executeRequest(request)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var entryID textID
	if err := json.Unmarshal(response.Body.Bytes(), &entryID); err != nil {
		t.Fatal(err)
	}
	return entryID.ID
}
