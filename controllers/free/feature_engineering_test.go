package free

import (
	"reflect"
	"testing"

	"github.com/AnthonyHewins/adm-backend/controllers/api"
)

func Test_featureEngineering(t *testing.T) {
	z := "zscore"
	garbageMode := "kasda"

	one := []float64{1}
	blank := &featureEngineeringReq{X: &[][]float64{[]float64{}}, Mode: &z}
	invalidTensor := &featureEngineeringReq{X: &[][]float64{one, []float64{1, 2}}, Mode: &z}
	invalidMode := &featureEngineeringReq{X: &[][]float64{one, one}, Mode: &garbageMode}

	tests := []struct {
		name  string
		args  *featureEngineeringReq
		want  api.Payload
		want1 *api.Error
	}{
		{
			"Zero length",
			blank,
			blank,
			nil,
		},
		{
			"Length mismatch",
			invalidTensor,
			nil,
			&api.Error{Http: 422, Code: ErrLength, Msg: "need a rectangular matrix, but row 1 had length 2, expected length 1"},
		},
		{
			"Invalid mode",
			invalidMode,
			nil,
			&api.Error{Http: 422, Code: ErrCmd, Msg: "don't understand mode kasda"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := featureEngineering(tt.args)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("featureEngineering() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("featureEngineering() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
