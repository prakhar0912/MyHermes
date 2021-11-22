package main

import (
	"flag"
	"github.com/fatih/color"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	h "github.com/mayankkumar2/SecurumExireSignalServer/api/handlers"
	"github.com/mayankkumar2/SecurumExireSignalServer/globals"
	"github.com/mayankkumar2/SecurumExireSignalServer/parsers"
	"github.com/mayankkumar2/SecurumExireSignalServer/utils"
	"log"
	"net/http"
	"os"
	"strings"
)


func main() {
	deploy := flag.String("deploy", "<production/development/staging>", "describe the deployment mode for signal server")
	flag.Parse()
	if strings.ToUpper(*deploy) != "PRODUCTION" && strings.ToUpper(*deploy) != "DEVELOPMENT" && strings.ToUpper(*deploy) != "STAGING" {
		color.Set(color.FgRed)
		log.Println("error: flag deploy is required [usage: -deploy <production/development/staging>]")
		color.Unset()
		os.Exit(1)
	}
	parsers.LoadConf(*deploy)
	utils.Register()

	router := mux.NewRouter()
	router.HandleFunc("/report/leak", h.ReportLeakHandler)
	router.HandleFunc("/block/endpoint", h.BlockEndpointHandler)
	router.HandleFunc("/heartbeat", h.HeartbeatHandler)

	color.Set(color.FgGreen)
	log.Println("Binding to:", globals.DeploymentConfig.ListenAt)
	color.Unset()

	err := http.ListenAndServe(globals.DeploymentConfig.ListenAt, handlers.LoggingHandler(os.Stderr, router))
	if err != nil {
		color.Set(color.FgRed)
		log.Println("error: failed to listen at:", globals.DeploymentConfig.ListenAt)
		color.Unset()
		os.Exit(1)
	}
}

/*
http.HandleFunc("/respond/user", func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Token string `json:"token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token, err := jwt.Parse(requestBody.Token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("error: Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims)
		} else {
			fmt.Println(err)
		}
	})
*/
