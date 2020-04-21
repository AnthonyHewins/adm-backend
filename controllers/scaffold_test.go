package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

const (
	testPolyreg = "/polyreg"
	testFeatureEngineering = "/feature-engineering"
	testRegistration = "/registration"
	testConfirmation = "/confirmation"
)

func buildRouter() *gin.Engine {
	return Router(&Routes{
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
