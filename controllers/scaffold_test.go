package controllers

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AnthonyHewins/adm-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	testPolyreg            = "/polyreg"
	testFeatureEngineering = "/feature-engineering"
	testRegistration       = "/registration"
	testConfirmation       = "/confirmation"
	testTicker             = "/ticker"
	testLogin              = "/auth/login"
	testPwReset            = "/pw-reset"
	testConfirmPwReset     = "/confirm-pw-reset"

	testRsaPrivkey = "../fixtures/testkey"
	testRsaPubkey  = "../fixtures/testkey.pub"
)

func buildRouterAndDB(t *testing.T) (*gorm.DB, *gin.Engine) {
	models.DBSetup(&models.DB{
		Host:     "localhost",
		Port:     5432,
		Name:     "admtest",
		User:     "test",
		Password: "test",
	})

	db, err := models.Connect()

	if err != nil {
		t.Fatalf("db opening failed: %v", err)
	}

	return db, buildRouter()
}

func buildRouter() *gin.Engine {
	return Router(
		Routes{
			Base:               "",
			Polyreg:            testPolyreg,
			FeatureEngineering: testFeatureEngineering,
			Registration:       testRegistration,
			AcctConfirmation:   testConfirmation,
			DcfValuation:       testTicker,
			ConfirmPwReset:     testConfirmPwReset,
			PasswordReset: testPwReset,
		},
		testRsaPrivkey,
		testRsaPubkey,
	)
}

func buildRequestFn(r *gin.Engine, method, endpoint string, body interface{}) (map[string]string, int) {
	w := httptest.NewRecorder()

	var req *http.Request
	if body == nil {
		req, _ = http.NewRequest(method, endpoint, nil)
	} else {
		marshalled, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, endpoint, bytes.NewBuffer(marshalled))
	}

	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	resp := make(map[string]string)
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp, w.Code
}

func buildProtectedRequestFn(r *gin.Engine, method, endpoint string, jwt []byte, body interface{}) (map[string]string, int) {
	w := httptest.NewRecorder()

	var req *http.Request
	if body == nil {
		req, _ = http.NewRequest(method, endpoint, nil)
	} else {
		marshalled, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, endpoint, bytes.NewBuffer(marshalled))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))
	r.ServeHTTP(w, req)

	resp := make(map[string]string)
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp, w.Code
}
