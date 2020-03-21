package main

import (
	"testing"
	"bytes"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/AnthonyHewins/adm/endpoints"
)

func TestFeatureEngineering(t *testing.T) {
	router := routerSetup()

	test := func(d *tools.FeatureEngineeringData, eCode int, err string) {
		resp, code := requestFeatEng(router, d)
		assert.Equal(t, eCode, code)
		assert.Equal(t, err, resp["error"])
	}

	// Test missing params
	z := "zscore"
	test(nil, 400, tools.ERR_PARAM)
	test(&tools.FeatureEngineeringData{X: &[][]float64{}}, 400, tools.ERR_PARAM)
	test(&tools.FeatureEngineeringData{Mode: &z},          400, tools.ERR_PARAM)

	one := []float64{1}

	// Test length of array
	test(&tools.FeatureEngineeringData{X: &[][]float64{[]float64{}}, Mode: &z}, 204, "")

	tooMany := make([][]float64, tools.MAX_ELEMENTS + 1)
	for i, _ := range tooMany { tooMany[i] = one }
	test(&tools.FeatureEngineeringData{X: &tooMany, Mode: &z}, 422, tools.ERR_LENGTH)

	// Test non-rectangular matrix
	test(
		&tools.FeatureEngineeringData{X: &[][]float64{ one, []float64{1,2} }, Mode: &z},
		422,
		tools.ERR_LENGTH_MISMATCH,
	)

	// Test invalid mode
	bs := "asd"
	test(
		&tools.FeatureEngineeringData{X: &[][]float64{one, one}, Mode: &bs},
		422,
		tools.ERR_CMD,
	)
}

func TestPolynomialRegression(t *testing.T) {
	router := routerSetup()

	test := func(d *tools.PolyRegData, eCode int, err string)  {
		resp, code := requestPolyReg(router, d)
		assert.Equal(t, eCode, code)
		assert.Equal(t, err, resp["error"])
	}

	// Test missing params
	test(nil, 400, tools.ERR_PARAM)
	test(&tools.PolyRegData{X: &[]float64{}}, 400, tools.ERR_PARAM)
	test(&tools.PolyRegData{Y: &[]float64{}}, 400, tools.ERR_PARAM)

	// test varying MaxDeg
	for i := -1; i <= 0; i++ {
		test(
			&tools.PolyRegData{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &i},
			422,
			tools.ERR_DEGREE,
		)
	}

	deg := 1
	test(
		&tools.PolyRegData{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &deg},
		422,
		tools.ERR_LENGTH,
	)

	deg = 6
	test(
		&tools.PolyRegData{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &deg},
		422,
		tools.ERR_DEGREE,
	)

	// test length of arrays
	deg = 1
	arr := make([]float64, 101)
	test(
		&tools.PolyRegData{X: &arr, Y: &[]float64{}, MaxDeg: &deg},
		422,
		tools.ERR_LENGTH,
	)

	test(
		&tools.PolyRegData{Y: &arr, X: &[]float64{}, MaxDeg: &deg},
		422,
		tools.ERR_LENGTH,
	)

	// test matrix that is guaranteed to be singular
	test(
		&tools.PolyRegData{X: &[]float64{1,1}, Y: &[]float64{1,1}, MaxDeg: &deg},
		500,
		tools.ERR_GENERAL,
	)
}

func requestFeatEng(r *gin.Engine, body *tools.FeatureEngineeringData) (map[string]string, int) {
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

func requestPolyReg(r *gin.Engine, body *tools.PolyRegData) (map[string]string, int) {
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
