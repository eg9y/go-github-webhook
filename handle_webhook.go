package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func ValidateSignature(payload []byte, signature string, secret string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	sig := strings.TrimPrefix(signature, "sha256=")

	expectedSig, err := hex.DecodeString(sig)

	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	actualSignature := mac.Sum(nil)

	return hmac.Equal(expectedSig, actualSignature)
}

func prettifyJSON(data []byte) (string, error) {
	var jsonData interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return "", errors.New("Error unmarshalling data")
	}
	prettyJSON, err := json.MarshalIndent(jsonData, "", "	")
	return string(prettyJSON), nil
}

func (a *ApiConfig) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error reading request body"))
		return
	}
	defer r.Body.Close()

	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		signature = r.Header.Get("X-Hub-Signature")
		if signature == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing signature header"))
			return
		}
	}

	isValid := ValidateSignature(body, signature, a.ActionsSecret)
	if !isValid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Signature does not match"))
		return
	}

	prettyPayload, err := prettifyJSON(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong prettifying payload data"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Webhook received and validated success! \n\npayload: %v", prettyPayload)))
}
