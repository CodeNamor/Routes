package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	//create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(HealthHandler)
	handler.ServeHTTP(rr, req)

	//checks
	if rr.Code != http.StatusOK {
		t.Errorf("Handler returned status %v", rr.Code)
	}

	want := "OK"

	if rr.Body.String() != want {
		t.Errorf("Expected %v, got %v in response body", want, rr.Body.String())
	}
}
