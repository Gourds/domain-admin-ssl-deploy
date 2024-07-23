package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port   int               `yaml:"port"`
	Level  string            `yaml:"level"`
	Cmds   map[string]string `yaml:"cmds"`
	Secret string            `yaml:"secret"`
}

type CertificateRequest struct {
	Domains           []string `json:"domains"`
	SslCertificate    string   `json:"ssl_certificate"`
	SslCertificateKey string   `json:"ssl_certificate_key"`
	StartTime         string   `json:"start_time"`
	ExpireTime        string   `json:"expire_time"`
}

var config Config

func loadConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &config)
}

func issueCertificateHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Token")
	if token != config.Secret {
		log.Println("Invalid Token attempt")
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return
	}

	keySavePath := r.Header.Get("Key-Save-Path")
	deployCmd := r.Header.Get("Deploy-Cmd")

	var certReq CertificateRequest
	err := json.NewDecoder(r.Body).Decode(&certReq)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, domain := range certReq.Domains {
		certPath := fmt.Sprintf("%s/%s.pem", keySavePath, domain)
		keyPath := fmt.Sprintf("%s/%s.key", keySavePath, domain)

		err := os.WriteFile(certPath, []byte(certReq.SslCertificate), 0644)
		if err != nil {
			log.Printf("Error writing certificate for domain %s: %v", domain, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = os.WriteFile(keyPath, []byte(certReq.SslCertificateKey), 0644)
		if err != nil {
			log.Printf("Error writing key for domain %s: %v", domain, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	cmd, ok := config.Cmds[deployCmd]
	if !ok {
		log.Printf("Invalid Deploy-Cmd: %s", deployCmd)
		http.Error(w, "Invalid Deploy-Cmd", http.StatusBadRequest)
		return
	}

	// parts := strings.Fields(cmd)
	// out, err := exec.Command(parts[0], parts[1:]...).Output()
	out, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		log.Printf("Command execution failed: %s\nError: %v", cmd, err)
		http.Error(w, fmt.Sprintf("Command execution failed: %s", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Certificate issued and command executed successfully for domains: %s. Command: %s", strings.Join(certReq.Domains, ", "), cmd)
	log.Printf("cmd run result %s", out)
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func main() {
	err := loadConfig("ops-cdsm.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/issueCertificate", issueCertificateHandler).Methods("POST")

	log.Printf("Starting server on port %d", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), router))
}
