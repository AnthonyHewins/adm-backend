package controllers

import (
	"bytes"
	"testing"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"

	"github.com/AnthonyHewins/adm-backend/models"
)

const (
	testPolyreg = "/polyreg"
	testFeatureEngineering = "/feature-engineering"
	testRegistration = "/registration"
	testConfirmation = "/confirmation"
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

	if err != nil { t.Fatalf("db opening failed: %v", err) }

	return db, Router(&Routes{
		Polyreg: testPolyreg,
		FeatureEngineering: testFeatureEngineering,
		Registration: testRegistration,
		AcctConfirmation: testConfirmation,
	})
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
