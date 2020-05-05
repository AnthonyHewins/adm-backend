package free

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func test(t *testing.T, router *gin.Engine) func(int, interface{}) {
	r := gofight.New()

	return func(httpcode int, expected interface{}) {
		r.POST("/").
			SetDebug(true).
			Run(router, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, "Hello World", r.Body.String())
				assert.Equal(t, http.StatusOK, r.Code)
			})
	}
}
