package free

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAddRoutes(t *testing.T) {
	type args struct {
		r       *gin.Engine
		apiBase string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddRoutes(tt.args.r, tt.args.apiBase)
		})
	}
}
