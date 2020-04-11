package main

import (
	"testing"
	"bytes"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/AnthonyHewins/adm-backend/controllers"
)

func TestFeatureEngineering(t *testing.T) {
	router := routerSetup()

	test := func(d *controllers.FeatureEngineeringData, eCode int, err string) {
		resp, code := requestFeatEng(router, d)
		assert.Equal(t, eCode, code)
		assert.Equal(t, err, resp["error"])
	}

	// Test missing params
	z := "zscore"
	test(nil, 400, controllers.ERR_PARAM)
	test(&controllers.FeatureEngineeringData{X: &[][]float64{}}, 400, controllers.ERR_PARAM)
	test(&controllers.FeatureEngineeringData{Mode: &z},          400, controllers.ERR_PARAM)

	one := []float64{1}

	// Test length of array
	test(&controllers.FeatureEngineeringData{X: &[][]float64{[]float64{}}, Mode: &z}, 204, "")

	tooMany := make([][]float64, controllers.MAX_ELEMENTS + 1)
	for i, _ := range tooMany { tooMany[i] = one }
	test(&controllers.FeatureEngineeringData{X: &tooMany, Mode: &z}, 422, controllers.ERR_LENGTH)

	// Test non-rectangular matrix
	test(
		&controllers.FeatureEngineeringData{X: &[][]float64{ one, []float64{1,2} }, Mode: &z},
		422,
		controllers.ERR_LENGTH_MISMATCH,
	)

	// Test invalid mode
	bs := "asd"
	test(
		&controllers.FeatureEngineeringData{X: &[][]float64{one, one}, Mode: &bs},
		422,
		controllers.ERR_CMD,
	)
}

func TestPolynomialRegression(t *testing.T) {
	router := routerSetup()

	test := func(d *controllers.PolyRegData, eCode int, err string)  {
		resp, code := requestPolyReg(router, d)
		assert.Equal(t, eCode, code)
		assert.Equal(t, err, resp["error"])
	}

	// Test missing params
	test(nil, 400, controllers.ERR_PARAM)
	test(&controllers.PolyRegData{X: &[]float64{}}, 400, controllers.ERR_PARAM)
	test(&controllers.PolyRegData{Y: &[]float64{}}, 400, controllers.ERR_PARAM)

	// test varying MaxDeg
	for i := -1; i <= 0; i++ {
		test(
			&controllers.PolyRegData{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &i},
			422,
			controllers.ERR_DEGREE,
		)
	}

	deg := 1
	test(
		&controllers.PolyRegData{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &deg},
		422,
		controllers.ERR_LENGTH,
	)

	deg = 6
	test(
		&controllers.PolyRegData{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &deg},
		422,
		controllers.ERR_DEGREE,
	)

	// test length of arrays
	deg = 1
	arr := make([]float64, 101)
	test(
		&controllers.PolyRegData{X: &arr, Y: &[]float64{}, MaxDeg: &deg},
		422,
		controllers.ERR_LENGTH,
	)

	test(
		&controllers.PolyRegData{Y: &arr, X: &[]float64{}, MaxDeg: &deg},
		422,
		controllers.ERR_LENGTH,
	)

	// test matrix that is guaranteed to be singular
	test(
		&controllers.PolyRegData{X: &[]float64{1,1}, Y: &[]float64{1,1}, MaxDeg: &deg},
		500,
		controllers.ERR_GENERAL,
	)
}

func requestFeatEng(r *gin.Engine, body *controllers.FeatureEngineeringData) (map[string]string, int) {
	w := httptest.NewRecorder()

	var req *http.Request
	if body == nil {
		req, _ = http.NewRequest("POST", FEATURE_ENGINEERING, nil)
	} else {
		marshalled, _ := json.Marshal(*body)
		req, _ = http.NewRequest("POST", FEATURE_ENGINEERING, bytes.NewBuffer(marshalled))
	}

	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	resp := make(map[string]string)
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp, w.Code
}

func requestPolyReg(r *gin.Engine, body *controllers.PolyRegData) (map[string]string, int) {
	w := httptest.NewRecorder()

	var req *http.Request
	if body == nil {
		req, _ = http.NewRequest("POST", POLYREG, nil)
	} else {
		marshalled, _ := json.Marshal(*body)
		req, _ = http.NewRequest("POST", POLYREG, bytes.NewBuffer(marshalled))
	}

	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	resp := make(map[string]string)
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp, w.Code
}
