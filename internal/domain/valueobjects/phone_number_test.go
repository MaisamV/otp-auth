package valueobjects

import (
	"testing"
)

func TestPhoneNumber_NewPhoneNumber_LocalToInternational(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "valid local number should be normalized to international",
			input:    "09123456789",
			expected: "+989123456789",
			wantErr:  false,
		},
		{
			name:     "valid international number should remain unchanged",
			input:    "+989123456789",
			expected: "+989123456789",
			wantErr:  false,
		},
		{
			name:     "invalid local number should return error",
			input:    "08123456789",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "empty phone number should return error",
			input:    "",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, err := NewPhoneNumber(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(phone) != tt.expected {
				t.Errorf("NewPhoneNumber() = %v, want %v", string(phone), tt.expected)
			}
		})
	}
}

func TestPhoneNumber_String_AlwaysInternational(t *testing.T) {
	tests := []struct {
		name     string
		phone    PhoneNumber
		expected string
	}{
		{
			name:     "international number should return as is",
			phone:    PhoneNumber("+989123456789"),
			expected: "+989123456789",
		},
		{
			name:     "local number should be converted to international",
			phone:    PhoneNumber("09123456789"),
			expected: "+989123456789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.phone.String()
			if result != tt.expected {
				t.Errorf("PhoneNumber.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}