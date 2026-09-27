package service_test

import (
	"testing"

	"api-students/app/model"
	"api-students/helper"
)

func TestValidateRegister(t *testing.T) {
	tests := []struct {
		name        string
		req         model.RegisterRequest
		expectedErr bool
		errorField  string
	}{
		{
			name: "Registrasi valid",
			req: model.RegisterRequest{
				Username: "mahasiswa1",
				Email:    "mahasiswa1@example.com",
				Password: "Rahasia123!",
			},
			expectedErr: false,
		},
		{
			name: "Username terlalu pendek",
			req: model.RegisterRequest{
				Username: "ab",
				Email:    "ab@example.com",
				Password: "Rahasia123!",
			},
			expectedErr: true,
			errorField:  "username",
		},
		{
			name: "Format email tidak valid",
			req: model.RegisterRequest{
				Username: "user_valid",
				Email:    "bukan-email",
				Password: "Rahasia123!",
			},
			expectedErr: true,
			errorField:  "email",
		},
		{
			name: "Password lemah",
			req: model.RegisterRequest{
				Username: "user_valid",
				Email:    "valid@example.com",
				Password: "admin123",
			},
			expectedErr: true,
			errorField:  "password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := helper.ValidateStruct(tt.req)
			if tt.expectedErr && len(errs) == 0 {
				t.Errorf("ValidateStruct() expected error on field %s, got none", tt.errorField)
			}
			if !tt.expectedErr && len(errs) > 0 {
				t.Errorf("ValidateStruct() expected no errors, got: %v", errs)
			}
		})
	}
}

func TestValidateLogin(t *testing.T) {
	tests := []struct {
		name        string
		req         model.LoginRequest
		expectedErr bool
	}{
		{
			name: "Login valid",
			req: model.LoginRequest{
				Username: "mahasiswa1",
				Password: "bebas-password",
			},
			expectedErr: false,
		},
		{
			name: "Username kosong",
			req: model.LoginRequest{
				Username: "",
				Password: "secret123",
			},
			expectedErr: true,
		},
		{
			name: "Password kosong",
			req: model.LoginRequest{
				Username: "mahasiswa1",
				Password: "",
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := helper.ValidateStruct(tt.req)
			if tt.expectedErr && len(errs) == 0 {
				t.Errorf("ValidateStruct() expected error, got none")
			}
			if !tt.expectedErr && len(errs) > 0 {
				t.Errorf("ValidateStruct() expected no error, got: %v", errs)
			}
		})
	}
}
