package controllers

import (
	"fmt"
	"time"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	_, router := buildRouterAndDB(t)

	test := func(eCode int, err string, body interface{}) {
		resp, code := buildRequestFn(router, "POST", testRegistration, body)
		assert.Equal(t, eCode, code)
		assert.Equal(t, err,   resp["error"])
	}

	// failed binding
	test(400, ERR_PARAM, nil)

	// invalid email (fails validation)
	test(422, ERR_EMAIL, &credentials{Email: "sdf", Password: "sndfiush923"})

	// invalid password (fails validation in length)
	test(422, ERR_PASSWORD, &credentials{Email: "sdf@jsaiod.com", Password: "sadhj"})

	now := time.Now().UnixNano()

	// valid
	test(200, "", &credentials{
		Email: fmt.Sprintf("sd%vf@jsaiod.com", now),
		Password: "sasdasdhj",
	})

	// already taken
	test(422, ERR_ALREADY_EXISTS, &credentials{
		Email: fmt.Sprintf("sd%vf@jsaiod.com", now),
		Password: "sasdasdhj",
	})
}
