package auth

import "testing"

func TestCreateJWT(t *testing.T) {
	secret := []byte("secret")

	tokenString, err := CreateJWT(secret, 100)
	if err != nil {
		t.Errorf("failed to create jwt token: %v", err)
	}
	if tokenString == "" {
		t.Errorf("token is empty")
	}

}
