package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/mayankkumar2/SecurumExireSignalServer/globals"
	"github.com/mayankkumar2/SecurumExireSignalServer/leaks"
	"github.com/mayankkumar2/SecurumExireSignalServer/utils"
	"net/http"
)

func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	go utils.RegisterAfterHeartbeat()
	w.WriteHeader(http.StatusOK)
}

func BlockEndpointHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		AuthSecret   string `json:"authSecret"`
		EndpointHash string `json:"endpoint_hash"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := jwt.Parse(requestBody.AuthSecret, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error: Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(globals.Secret), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		leaks.BlockEndpoint(requestBody.EndpointHash)
	}
}

func ReportLeakHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Endpoint          string   `json:"endpoint"`
		LeakedCredentials []string `json:"leaked_credentials"`
		EndpointHash      string   `json:"endpoint_hash"`
	}
	secretHeader := r.Header.Get("SECRET")
	token, err := jwt.Parse(secretHeader, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error: Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(globals.Secret), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		b, _ := json.Marshal(requestBody.LeakedCredentials)
		leaks.ReportLeak(requestBody.Endpoint, string(b), requestBody.EndpointHash)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
