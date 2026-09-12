package service_test

import (
	"testing"

	"api-students/app/model"
	"api-students/app/service"
)

func TestCheckPasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected string
	}{
		{
			name:     "Password terlalu pendek (< 8 karakter)",
			password: "pass1",
			expected: "minimal 8 karakter",
		},
		{
			name:     "Password hanya berisi huruf tanpa angka",
			password: "passwordrahasia",
			expected: "harus memuat huruf dan angka",
		},
		{
			name:     "Password hanya berisi angka tanpa huruf",
			password: "1234567890",
			expected: "harus memuat huruf dan angka",
		},
		{
			name:     "Password terdaftar dalam kamus password lemah / umum",
			password: "password1",
			expected: "password terlalu umum",
		},
		{
			name:     "Password kuat dan valid",
			password: "rahasia123",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CheckPasswordStrength(tt.password)
			if result != tt.expected {
				t.Errorf("CheckPasswordStrength(%q) = %q; expected %q", tt.password, result, tt.expected)
			}
		})
	}
}

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
				Password: "rahasia123",
			},
			expectedErr: false,
		},
		{
			name: "Username terlalu pendek",
			req: model.RegisterRequest{
				Username: "ab",
				Email:    "ab@example.com",
				Password: "rahasia123",
			},
			expectedErr: true,
			errorField:  "username",
		},
		{
			name: "Format email tidak valid",
			req: model.RegisterRequest{
				Username: "user_valid",
				Email:    "bukan-email",
				Password: "rahasia123",
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
			errs := service.ValidateRegister(tt.req)
			if tt.expectedErr && len(errs) == 0 {
				t.Errorf("ValidateRegister() expected error on field %s, got none", tt.errorField)
			}
			if !tt.expectedErr && len(errs) > 0 {
				t.Errorf("ValidateRegister() expected no errors, got: %v", errs)
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
			errs := service.ValidateLogin(tt.req)
			if tt.expectedErr && len(errs) == 0 {
				t.Errorf("ValidateLogin() expected error, got none")
			}
			if !tt.expectedErr && len(errs) > 0 {
				t.Errorf("ValidateLogin() expected no error, got: %v", errs)
			}
		})
	}
}
