package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"
	"github.com/mayankkumar2/SecurumExireSignalServer/constants"
	"github.com/mayankkumar2/SecurumExireSignalServer/globals"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var token string

func Register() {
	globals.Secret = RandStringRunes(32)
	Url := fmt.Sprintf("%s/login", constants.SecurumExireBotURL)
	webhook := globals.DeploymentConfig.SelfWebhook
	uid := globals.DeploymentConfig.BotUID
	botSecret := globals.DeploymentConfig.BotSecret
	tokenizer := jwt.New(jwt.SigningMethodHS256)
	var err error
	token, err = tokenizer.SignedString([]byte(globals.Secret))
	if err != nil {
		color.Set(color.FgRed)
		log.Println("error:", err.Error())
		color.Unset()
		os.Exit(1)
	}

	var request = map[string]string{
		"identity_string": token,
		"webhook":         webhook,
		"secret":          botSecret,
		"uid":             uid,
	}

	var b bytes.Buffer

	_ = json.NewEncoder(&b).Encode(request)
	response, err := http.Post(Url, "application/json", &b)
	if err != nil {
		color.Set(color.FgRed)
		log.Println("error:", err.Error())
		color.Unset()
		os.Exit(1)
	}

	if response.StatusCode != http.StatusOK {
		switch response.StatusCode {
		case http.StatusUnauthorized:
			color.Set(color.FgRed)
			log.Println("error: Credentials are incorrect")
			color.Unset()
			os.Exit(1)
		case http.StatusNotFound:
			color.Set(color.FgRed)
			log.Println("error: UID not found")
			color.Unset()
			os.Exit(1)
		}
	}
	r, err := ioutil.ReadAll(response.Body)
	if err != nil {
		color.Set(color.FgRed)
		log.Println("error:", err.Error())
		color.Unset()
		os.Exit(1)
	}
	globals.BotToken = string(r)
	log.Println("Login successful!")
	response, err = RegisterSignalServer(token)
	if err != nil {
		color.Set(color.FgRed)
		log.Println("error:", err.Error())
		color.Unset()
		os.Exit(1)
	}
	if response.StatusCode != http.StatusOK {
		color.Set(color.FgRed)
		log.Println("error: failed to register to domain server")
		color.Unset()
		os.Exit(1)
	}
}
