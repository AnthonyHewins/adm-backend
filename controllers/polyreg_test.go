package controllers

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

const polyreg = "/polyreg"

func TestPolynomialRegression(t *testing.T) {
	_, router := buildRouterAndDB(t)

	test := func(d *PolyRegData, eCode int, err string)  {
		resp, code := buildRequestFn(router, "POST", "/polyreg", d)
		assert.Equal(t, eCode, code)
		assert.Equal(t, err, resp["error"])
	}

	// Test missing params
	test(nil, 400, ERR_PARAM)
	test(&PolyRegData{X: &[]float64{}}, 400, ERR_PARAM)
	test(&PolyRegData{Y: &[]float64{}}, 400, ERR_PARAM)

	// test varying MaxDeg
	for i := -1; i <= 0; i++ {
		test(
			&PolyRegData{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &i},
			422,
			ERR_DEGREE,
		)
	}

	deg := 1
	test(
		&PolyRegData{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &deg},
		422,
		ERR_LENGTH,
	)

	deg = 6
	test(
		&PolyRegData{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &deg},
		422,
		ERR_DEGREE,
	)

	// test length of arrays
	deg = 1
	arr := make([]float64, 101)
	test(
		&PolyRegData{X: &arr, Y: &[]float64{}, MaxDeg: &deg},
		422,
		ERR_LENGTH,
	)

	test(
		&PolyRegData{Y: &arr, X: &[]float64{}, MaxDeg: &deg},
		422,
		ERR_LENGTH,
	)

	// test matrix that is guaranteed to be singular
	test(
		&PolyRegData{X: &[]float64{1,1}, Y: &[]float64{1,1}, MaxDeg: &deg},
		500,
		ERR_GENERAL,
	)
}
