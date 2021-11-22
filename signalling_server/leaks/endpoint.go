package leaks

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mayankkumar2/SecurumExireSignalServer/globals"
	"log"
	"net/http"
	"net/url"
	"time"
)

func BlockEndpoint(endpointHash string) {
	urlEndpoint, _ := url.Parse(fmt.Sprintf("http://%s/block_endpoint", globals.DeploymentConfig.LeakServerAddress))
	var req = &http.Request{
		URL: urlEndpoint,
		Header: map[string][]string{
			"endpoint": {endpointHash},
		},
		Method: http.MethodGet,
	}
	cl := http.Client{
		Timeout: time.Minute * 5,
	}
	_, err := cl.Do(req)
	if err != nil {
		color.Set(color.FgRed)
		log.Println("error: couldn't report leak :", err)
		color.Unset()
	}
}
