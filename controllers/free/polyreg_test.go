package free

import (
	"reflect"
	"testing"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
)

const polyreg = "/polyreg"

func Test_polynomialRegression(t *testing.T) {
	invalidNegative1 := -1
	invalidZero := 0
	validOne := 1
	invalidSix := 6

	tooLong := make([]float64, maxElements+1)

	tests := []struct {
		name  string
		args  *polyRegReq
		want  api.Payload
		want1 *api.Error
	}{
		{
			"Degree can't be -1",
			&polyRegReq{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &invalidNegative1},
			nil,
			&api.Error{Http: 422, Code: ErrDegree, Msg: "maxDeg must satisfy 0 <= maxDeg <= 5"},
		},
		{
			"Degree can't be 0",
			&polyRegReq{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &invalidZero},
			nil,
			&api.Error{Http: 422, Code: ErrDegree, Msg: "maxDeg must satisfy 0 <= maxDeg <= 5"},
		},
		{
			"Degree can't be greater than 5",
			&polyRegReq{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &invalidSix},
			nil,
			&api.Error{Http: 422, Code: ErrDegree, Msg: "maxDeg must satisfy 0 <= maxDeg <= 5"},
		},
		{
			"Length can't be nothing",
			&polyRegReq{X: &[]float64{}, Y: &[]float64{}, MaxDeg: &validOne},
			nil,
			&api.Error{Http: 422, Code: ErrLength, Msg: "must have len(x) == len(y) && maxDeg < len(x, y) <= 100, got len(x) = 0 and len(y) = 0"},
		},
		{
			"Lengths must match",
			&polyRegReq{X: &[]float64{1}, Y: &[]float64{}, MaxDeg: &validOne},
			nil,
			&api.Error{Http: 422, Code: ErrLength, Msg: "must have len(x) == len(y) && maxDeg < len(x, y) <= 100, got len(x) = 1 and len(y) = 0"},
		},
		{
			"Lengths cannot exceed maxLength",
			&polyRegReq{X: &tooLong, Y: &[]float64{}, MaxDeg: &validOne},
			nil,
			&api.Error{Http: 422, Code: ErrLength, Msg: "must have len(x) == len(y) && maxDeg < len(x, y) <= 100, got len(x) = 101 and len(y) = 0"},
		},
		{
			"Matrix cannot be singular",
			&polyRegReq{X: &[]float64{1, 1}, Y: &[]float64{1, 1}, MaxDeg: &validOne},
			nil,
			&api.Error{Http: 500, Code: ErrSingular, Msg: "matrix singular or near-singular with condition number +Inf"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := polynomialRegression(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("polynomialRegression() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("polynomialRegression() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
