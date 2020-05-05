package free

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFeatureEngineering(t *testing.T) {
	router := gin.Default()
	router.POST("/feature-engineering")

	test := func(d *featureEngineeringData, eCode int, err string) {
		resp, code := buildRequestFn(router, "POST", "/feature-engineering", d)
		assert.Equal(t, eCode, code)
		assert.Equal(t, err, resp["error"])
	}

	// Test missing params
	z := "zscore"
	test(nil, 400, api.ErrIncorrectParams)
	test(&featureEngineeringData{X: &[][]float64{}}, 400, api.ErrIncorrectParams)
	test(&featureEngineeringData{Mode: &z},          400, api.ErrIncorrectParams)

	one := []float64{1}

	// Test length of array
	test(&featureEngineeringData{X: &[][]float64{[]float64{}}, Mode: &z}, 204, "")

	tooMany := make([][]float64, maxElements + 1)
	for i, _ := range tooMany { tooMany[i] = one }
	test(&featureEngineeringData{X: &tooMany, Mode: &z}, 422, errLength)

	// Test non-rectangular matrix
	test(
		&featureEngineeringData{X: &[][]float64{ one, []float64{1,2} }, Mode: &z},
		422,
		errLength,
	)

	// Test invalid mode
	bs := "asd"
	test(
		&featureEngineeringData{X: &[][]float64{one, one}, Mode: &bs},
		422,
		api.ErrInvalidParam,
	)
}
