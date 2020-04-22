package controllers

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestFeatureEngineering(t *testing.T) {
	_, router := buildRouterAndDB(t)

	test := func(d *FeatureEngineeringData, eCode int, err string) {
		resp, code := buildRequestFn(router, "POST", "/feature-engineering", d)
		assert.Equal(t, eCode, code)
		assert.Equal(t, err, resp["error"])
	}

	// Test missing params
	z := "zscore"
	test(nil, 400, ERR_PARAM)
	test(&FeatureEngineeringData{X: &[][]float64{}}, 400, ERR_PARAM)
	test(&FeatureEngineeringData{Mode: &z},          400, ERR_PARAM)

	one := []float64{1}

	// Test length of array
	test(&FeatureEngineeringData{X: &[][]float64{[]float64{}}, Mode: &z}, 204, "")

	tooMany := make([][]float64, MAX_ELEMENTS + 1)
	for i, _ := range tooMany { tooMany[i] = one }
	test(&FeatureEngineeringData{X: &tooMany, Mode: &z}, 422, ERR_LENGTH)

	// Test non-rectangular matrix
	test(
		&FeatureEngineeringData{X: &[][]float64{ one, []float64{1,2} }, Mode: &z},
		422,
		ERR_LENGTH_MISMATCH,
	)

	// Test invalid mode
	bs := "asd"
	test(
		&FeatureEngineeringData{X: &[][]float64{one, one}, Mode: &bs},
		422,
		ERR_CMD,
	)
}
