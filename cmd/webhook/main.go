package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	DeployScriptPath = "/DATA/AppData/backend_codingin/scripts/deploy.sh"
	Port             = ":9090"
)

func main() {
	http.HandleFunc("/webhook", webhookHandler)

	log.Printf("GitHub Webhook service listening on %s", Port)
	if err := http.ListenAndServe(Port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	secret := os.Getenv("DEPLOY_SECRET")
	if secret == "" {
		log.Println("Error: DEPLOY_SECRET environment variable not set")
		http.Error(w, "Server misconfiguration", http.StatusInternalServerError)
		return
	}

	// 1. Read Payload
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 2. Verify Signature (X-Hub-Signature-256)
	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		log.Printf("Missing signature header")
		http.Error(w, "Missing signature", http.StatusUnauthorized)
		return
	}

	if !verifySignature(secret, signature, payload) {
		log.Printf("Invalid signature from %s", r.RemoteAddr)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	log.Println("Received valid GitHub push event. Executing script...")

	// 3. Execute Script
	cmd := exec.Command("/bin/bash", DeployScriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Deployment failed: %v\nOutput: %s", err, string(output))
		http.Error(w, fmt.Sprintf("Deployment failed: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Deployment successful:\n%s", string(output))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deployment triggered successfully"))
}

// verifySignature checks the HMAC-SHA256 signature
func verifySignature(secret, headerSignature string, payload []byte) bool {
	// Header format: sha256=...
	parts := strings.SplitN(headerSignature, "=", 2)
	if len(parts) != 2 || parts[0] != "sha256" {
		return false
	}
	
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)
	expectedSignature := hex.EncodeToString(expectedMAC)
	
	return hmac.Equal([]byte(parts[1]), []byte(expectedSignature))
}
