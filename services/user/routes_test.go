package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/API/types"
	"github.com/gin-gonic/gin"
)

type mockUserStore struct {
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) GetUserById(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(user types.User) error {
	return nil
}

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)
	t.Run("Should fail if user payload is invalid", func(t *testing.T) {
		payload := types.SignInPayload{
			FirstName: "testname",
			LastName:  "testsurname",
			Email:     "testgmail.com",
			Password:  "123",
		}
		marshaledPayload, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/signin", handler.handleSignIn)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d ,got %d", http.StatusBadRequest, rr.Code)
		}

	})

	t.Run("Should pass if user payload is correct", func(t *testing.T) {
		payload := types.SignInPayload{
			FirstName: "testName",
			LastName:  "testSurname",
			Email:     "example@blabla.com",
			Password:  "1231232",
		}
		marshaledPayload, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/signin", handler.handleSignIn)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d ,got %d", http.StatusCreated, rr.Code)
		}
	})

}
