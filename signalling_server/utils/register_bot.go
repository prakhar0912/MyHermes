package utils

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mayankkumar2/SecurumExireSignalServer/globals"
	"log"
	"net/http"
	"net/url"
	"time"
)

func RegisterSignalServer(token string) (*http.Response, error) {
	urlEndpoint, _ := url.Parse(fmt.Sprintf("http://%s/register_signal_server", globals.DeploymentConfig.LeakServerAddress))
	var req = &http.Request{
		URL: urlEndpoint,
		Header: map[string][]string{
			"secret": {token},
		},
		Method: http.MethodGet,
	}
	cl := http.Client{
		Timeout: time.Minute * 5,
	}
	return cl.Do(req)
}

func RegisterAfterHeartbeat() {
	color.Set(color.FgRed)
	log.Println("retry: responding to heartbeat, re-registering...")
	color.Unset()
	response, err := RegisterSignalServer(token)
	now := time.Now()
	finish := now.Add(time.Minute * 10)
	for err != nil && time.Now().Unix() <= finish.Unix() {
		color.Set(color.FgRed)
		log.Println("retry: responding to heartbeat, re-registering...")
		color.Unset()
		response, err = RegisterSignalServer(token)
		time.Sleep(time.Second * 5)
	}
	if response != nil {
		if response.StatusCode != http.StatusOK {
			color.Set(color.FgRed)
			log.Println("error: can't register the server")
			color.Unset()
		}
	}

}
