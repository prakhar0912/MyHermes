package leaks

import (
	"bytes"
	"encoding/json"
	"github.com/fatih/color"
	"github.com/mayankkumar2/SecurumExireSignalServer/constants"
	"github.com/mayankkumar2/SecurumExireSignalServer/globals"
	"log"
	"net/http"
)

func ReportLeak(endpoint, secrets string, endpointHash string) {
	var requestBody = map[string]string{
		"endpoint":      endpoint,
		"secret":        globals.BotToken,
		"secret_name":   secrets,
		"endpoint_hash": endpointHash,
	}
	var b bytes.Buffer
	_ = json.NewEncoder(&b).Encode(&requestBody)
	_, err := http.Post(constants.SecurumExireBotURL+"/report/leak", "application/json", &b)
	if err != nil {
		color.Set(color.FgRed)
		log.Println("error: couldn't report leak :", err)
		color.Unset()
	}
}
