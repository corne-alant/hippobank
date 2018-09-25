package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestPaymentHandler(t *testing.T) {
	srv := server{router: mux.NewRouter()}
	srv.router.HandleFunc("/payments", srv.paymentHandler())

	req, _ := http.NewRequest("GET", "/payments", nil)
	res := httptest.NewRecorder()
	srv.router.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	result := PaymentResponse{}
	err := json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		t.Fatal("Could not parse request", err)
	}

	assert.Equal(t, "Account 1", result.From)
	assert.Equal(t, "Account 2", result.To)
	assert.Equal(t, 200, result.Amount)

}
