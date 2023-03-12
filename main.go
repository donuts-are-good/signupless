package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/crypto/sha3"
)

// session tokens and ids
var sessions = make(map[string]string)

// globalSalt will hold the salt used for generating tokens
var globalSalt string

// session is a single record
type session struct {
	ID    string `json:"id,omitempty"`
	Token string `json:"token,omitempty"`
}

func init() {

	// generate 16 byte salt
	salt, err := generateSalt(16)

	if err != nil {

		// log error and exit if salt generation failed
		log.Fatalf("failed to generate salt: %s", err)

	}

	// concatenate current UnixNano time with salt and encode as hex string
	globalSalt := fmt.Sprintf("%x", time.Now().UnixNano()) + hex.EncodeToString(salt)

	// generate a new token using the globalSalt and current UnixNano time
	generateToken(globalSalt, time.Now().UnixNano())

}

func main() {

	// add a session
	http.HandleFunc("/session/add", handleAddSession)

	// check a session
	http.HandleFunc("/session/check", handleCheckSession)

	// start HTTP server listening on port 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {

		// log fatal error and exit if server fails to start
		log.Fatalf("Error starting server: %s", err)

	}

}

// generateToken function generates a token from a given salt and Unix timestamp.
func generateToken(salt string, timestamp int64) string {

	// convert Unix timestamp to hex string
	thisTime := fmt.Sprintf("%x", timestamp)

	// concatenate salt and timestamp, convert to byte slice
	data := []byte(salt + thisTime)

	// hash the byte slice with sha3-256 and get the hash as a byte slice
	hash := sha3.Sum256(data)

	// return the token as a hex string
	return hex.EncodeToString(hash[:])

}

// generateSalt returns a new salt of specified length
func generateSalt(length int) ([]byte, error) {

	// create byte slice of specified length for salt
	salt := make([]byte, length)

	// generate random bytes for salt
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	return salt, nil

}

// function will handle adding sessions
func handleAddSession(w http.ResponseWriter, r *http.Request) {

	// defer function to handle recovery from panic
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Panic: %v", err)
		}

		// close the request body after function execution
		err := r.Body.Close()
		if err != nil {
			log.Printf("Failed to close request body: %s", err)
		}
	}()

	// struct to hold the request
	var req struct {
		ID string `json:"id"`
	}

	// decode the JSON request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Failed to decode request: %s", err)
		return
	}

	// generate a new token and store the session
	token := generateToken(globalSalt, time.Now().UnixNano())

	// set this token's new value
	sessions[token] = req.ID

	// create a response with the session ID and token
	response := session{ID: req.ID, Token: token}

	// encode the response as JSON and send it
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Failed to encode response: %s", err)
		return
	}

}

// function will handle checking a session
func handleCheckSession(w http.ResponseWriter, r *http.Request) {

	// defer function to handle recovery from panic
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Panic: %v", err)
		}

		// close the request body after execution
		err := r.Body.Close()
		if err != nil {
			log.Printf("Failed to close request body: %s", err)
		}

	}()

	// get the session token from the request header
	token := r.Header.Get("session-token")

	// validate the format of the token
	if !validateTokenFormat(token) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Invalid token format: %s", token)
		return
	}

	// check if the token exists in the session map
	id, exists := sessions[token]
	if !exists {
		http.Error(w, "Forbidden", http.StatusForbidden)
		log.Printf("Token not found: %s", token)
		return
	}

	// generate a new token and update the session map
	newToken := generateToken(globalSalt, time.Now().UnixNano())
	sessions[newToken] = id

	// create a response with the session ID and new token
	response := session{ID: id, Token: newToken}

	// encode the response as JSON and send it
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Failed to encode response: %s", err)
		return
	}
}

// validateTokenFormat checks our token to see if it looks like one we'd use
func validateTokenFormat(token string) bool {

	// regular expression pattern to match the token format
	pattern := "^[a-f0-9]{64}$"

	// use the MatchString function to check if the token matches the pattern
	matched, err := regexp.MatchString(pattern, token)
	if err != nil {
		log.Printf("Failed to validate token format: %s", err)
		return false
	}

	// return whether the token matches the pattern
	return matched

}
