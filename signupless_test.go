package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func Test_generateSalt(t *testing.T) {
	tests := []struct {
		length int
	}{
		{length: 8},
		{length: 16},
		{length: 32},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("length=%d", tt.length), func(t *testing.T) {
			salt, err := generateSalt(tt.length)
			if err != nil {
				t.Errorf("generateSalt() error = %v", err)
				return
			}
			if len(salt) != tt.length {
				t.Errorf("generateSalt() = %v, want %v", len(salt), tt.length)
			}
		})
	}
}

func Test_generateToken(t *testing.T) {
	tests := []struct {
		salt string
		time int64
		hash string
	}{
		{
			salt: "174baa8db2b92718a0099507b18702f7a1c1a3c3397e5b0b",
			time: 1678622811691240000,
			hash: "0ee8a6b3e60c97ed6bf3edd09ca69f30af92b7167b4eb70d1026065665c466fc",
		},
		{
			salt: "174baa8e393d2848c0bd386ad8d0dc30a0edc9ab86e9514c",
			time: 1678622813948039000,
			hash: "7ff9e17f12dd39105ea7e9f0acf2531c90227ed53f1c230351701763b2736074",
		},
		{
			salt: "174baa8e94c159781a3d27b38883bb108ead6b577fe819b4",
			time: 1678622815483429000,
			hash: "dc4fc37622aef1783da4256b38a48e59b6a74c82dcff71bc5be5cf6005a4bd91",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("salt=%s,time=%d", tt.salt, tt.time), func(t *testing.T) {
			got := generateToken(tt.salt, tt.time)
			if len(got) != 64 {
				t.Errorf("generateToken() length = %d, want 64", len(got))
			}
			match, _ := regexp.MatchString("^[a-f0-9]+$", got)
			if !match {
				t.Errorf("generateToken() = %v, want match for pattern '^[a-f0-9]+$'", got)
			}
			if got != tt.hash {
				t.Errorf("generateToken() = %v, want %v", got, tt.hash)
			}
		})
	}
}
func TestHandleAddSession(t *testing.T) {
	reqBody := []byte(`{"id": "1234567890"}`)
	req := httptest.NewRequest("POST", "/session/add", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	handleAddSession(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	resBody := rr.Body.Bytes()
	var res session
	err := json.Unmarshal(resBody, &res)
	if err != nil {
		t.Errorf("error unmarshaling response body: %v", err)
	}

	if len(res.Token) != 64 {
		t.Errorf("token length is incorrect: got %v want %v", len(res.Token), 64)
	}

	if !validateTokenFormat(res.Token) {
		t.Errorf("token has invalid format: %v", res.Token)
	}
}
func TestValidateTokenFormat_ValidToken(t *testing.T) {
	token := "a3a3f8c5b78fd9135c5e5b5e14b1d9b1fd11b7a3d3a8f5cdd0f72c492fd05019"
	isValid := validateTokenFormat(token)
	if !isValid {
		t.Errorf("Expected token %s to be valid, but got invalid", token)
	}
}

func TestValidateTokenFormat_InvalidToken(t *testing.T) {
	token := "a3a3f8c5b78fd9135c5e5b5e14b1d9b1fd11b7a3d3a8f5cdd0f72c492fd05019#"
	isValid := validateTokenFormat(token)
	if isValid {
		t.Errorf("Expected token %s to be invalid, but got valid", token)
	}
}

func TestValidateTokenFormat_TooShortToken(t *testing.T) {
	token := "a3a3f8c5b78fd9135c5e5b5e14b1d9b1fd11b7a3d3a8f5cdd0f72c492fd0501"
	isValid := validateTokenFormat(token)
	if isValid {
		t.Errorf("Expected token %s to be invalid, but got valid", token)
	}
}

func TestValidateTokenFormat_TooLongToken(t *testing.T) {
	token := "a3a3f8c5b78fd9135c5e5b5e14b1d9b1fd11b7a3d3a8f5cdd0f72c492fd050199"
	isValid := validateTokenFormat(token)
	if isValid {
		t.Errorf("Expected token %s to be invalid, but got valid", token)
	}
}

func TestValidateTokenFormat_EmptyToken(t *testing.T) {
	token := ""
	isValid := validateTokenFormat(token)
	if isValid {
		t.Errorf("Expected empty token to be invalid, but got valid")
	}
}
