package api

import (
	"testing"
)

func Test_encryption(t *testing.T) {
	type test struct {
		input   string
		pass    string
		wantErr bool
	}

	tests := []test{
		{
			input:   "test",
			pass:    "password",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		encrypted, err := encrypt(tc.input, tc.pass)
		if (err != nil) != tc.wantErr {
			t.Errorf("encrypt() error = %v, wantErr %v", err, tc.wantErr)

			return
		}

		if encrypted == "" {
			t.Errorf("encrypt() encrypted = %v, want not empty", encrypted)
		}

		decrypted, err := decrypt(encrypted, tc.pass)
		if (err != nil) != tc.wantErr {
			t.Errorf("decrypt() error = %v, wantErr %v", err, tc.wantErr)

			return
		}

		if decrypted != tc.input {
			t.Errorf("decrypt() decrypted = %v, want %v", decrypted, tc.input)
		}
		t.Logf("input: %v, decrypted: %v", tc.input, decrypted)
	}
}

func Test_base64(t *testing.T) {
	type test struct {
		input   string
		wantErr bool
	}

	tests := []test{
		{
			input:   "test",
			wantErr: false,
		},
	}

	for _, tc := range tests {
		encoded := toBase64([]byte(tc.input))
		decoded, err := fromBase64(encoded)
		if (err != nil) != tc.wantErr {
			t.Errorf("base64Decode() error = %v, wantErr %v", err, tc.wantErr)

			return
		}

		if string(decoded) != tc.input {
			t.Errorf("base64Decode() decoded = %v, want %v", string(decoded), tc.input)
		}
		t.Logf("input: %v, decoded: %v", tc.input, string(decoded))
	}
}
