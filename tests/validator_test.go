package tests

import (
	"qasynda/pkg/validator"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "valid email",
			email: "test@example.com",
			want:  true,
		},
		{
			name:  "invalid email - no @",
			email: "testexample.com",
			want:  false,
		},
		{
			name:  "invalid email - no domain",
			email: "test@",
			want:  false,
		},
		{
			name:  "valid email with subdomain",
			email: "test@mail.example.com",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validator.ValidateEmail(tt.email); got != tt.want {
				t.Errorf("ValidateEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateRequired(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "non-empty string",
			value: "test",
			want:  true,
		},
		{
			name:  "empty string",
			value: "",
			want:  false,
		},
		{
			name:  "whitespace only",
			value: "   ",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validator.ValidateRequired(tt.value); got != tt.want {
				t.Errorf("ValidateRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateMinLength(t *testing.T) {
	tests := []struct {
		name  string
		value string
		min   int
		want  bool
	}{
		{
			name:  "valid length",
			value: "test",
			min:   3,
			want:  true,
		},
		{
			name:  "exact length",
			value: "test",
			min:   4,
			want:  true,
		},
		{
			name:  "too short",
			value: "te",
			min:   3,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validator.ValidateMinLength(tt.value, tt.min); got != tt.want {
				t.Errorf("ValidateMinLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

