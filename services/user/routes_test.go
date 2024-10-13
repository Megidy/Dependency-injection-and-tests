package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/API/services/auth"
	"github.com/API/types"
	"github.com/gin-gonic/gin"
)

type mockUserStoreSignIn struct {
}

func (m *mockUserStoreSignIn) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStoreSignIn) GetUserById(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStoreSignIn) CreateUser(user types.User) error {
	return nil
}

func TestUserServiceHandlerSignIn(t *testing.T) {
	userStore := &mockUserStoreSignIn{}
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

type mockUserStoreLogIn struct {
}

func (m *mockUserStoreLogIn) CreateUser(user types.User) error {
	return nil
}
func (m *mockUserStoreLogIn) GetUserByEmail(email string) (*types.User, error) {
	hashedPassword, _ := auth.HashPassword("test")
	return &types.User{
		Id:        100,
		FirstName: "testName",
		LastName:  "testSurname",
		Email:     "test@gmail.com",
		Password:  hashedPassword,
	}, nil
}
func (m *mockUserStoreLogIn) GetUserById(id int) (*types.User, error) {
	return nil, nil
}

func TestUserServiceHandlerLogIn(t *testing.T) {
	NewStore := &mockUserStoreLogIn{}
	handler := NewHandler(NewStore)

	t.Run("Should Pass if user payload is correct", func(t *testing.T) {

		payload := types.LogInPayload{
			Email:    "test@gmail.com",
			Password: "test",
		}
		marshaledPayload, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/login", handler.handleLogIn)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected %d , got %d", http.StatusOK, rr.Code)
		}
	})
	t.Run("Should fail if user email is incorrect", func(t *testing.T) {
		payload := types.LogInPayload{
			Email:    "fail.com",
			Password: "123",
		}
		marshaledPayload, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshaledPayload))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := gin.Default()
		router.POST("/login", handler.handleLogIn)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, rr.Code)
		}

	})
}
