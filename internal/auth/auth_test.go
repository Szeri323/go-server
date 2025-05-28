package auth

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "123456"

	token, err := HashPassword(password1)
	if err != nil {
		t.Errorf(`Couldn't hash password`)
	}

	err = CheckPasswordHash(token, password1)
	if err != nil {
		t.Errorf(`Hash check failed`)
	}

}
