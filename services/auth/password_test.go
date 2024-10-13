package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hashedPassword, err := HashPassword("testpass")
	if err != nil {
		t.Errorf("error when hashing password , %v", err)
	}
	if hashedPassword == "" {
		t.Errorf("didnt mean password to be empty")
	}
	if hashedPassword == "testpass" {
		t.Errorf("didnt hash password")

	}

}
func TestComparePasswords(t *testing.T) {
	hashedPassword, err := HashPassword("testpass")
	if err != nil {
		t.Errorf("error when hashing password , %v", err)
	}
	ok := ComparePasswords(hashedPassword, []byte("testpass"))
	if !ok {
		t.Errorf("password didnt match")
	}

	if ComparePasswords(hashedPassword, []byte("BaDpAsSwOrD")) {
		t.Errorf("password matched but didnt have to ")
	}
}
