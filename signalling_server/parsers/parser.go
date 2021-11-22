package parsers

import (
	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/mayankkumar2/SecurumExireSignalServer/globals"
	"github.com/mayankkumar2/SecurumExireSignalServer/models"
	"log"
	"os"
	"regexp"
	"strings"
)

func ValidateDeploymentConfig(d *models.Deployment) bool {
	urlRegex, _ := regexp.Compile("^https{0,1}://[a-zA-Z0-9\\-\\.]+:{0,1}\\d{0,5}$")
	listenAddress, _ := regexp.Compile("^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}:[0-9]{2,5}$")
	signalServerAddress, _ := regexp.Compile("^[a-zA-Z0-9\\-\\.]+:{0,1}\\d{0,5}$")
	if !urlRegex.MatchString(d.SelfWebhook) {
		color.Set(color.FgRed)
		log.Println("error: self_webhook url invalid")
		color.Unset()
		os.Exit(1)
	}
	if !listenAddress.MatchString(d.ListenAt) {
		color.Set(color.FgRed)
		log.Println("error: listen_address invalid")
		color.Unset()
		os.Exit(1)
	}
	if !signalServerAddress.MatchString(d.LeakServerAddress) {
		color.Set(color.FgRed)
		log.Println("error: leak_server_address invalid")
		color.Unset()
		os.Exit(1)
	}
	if d.BotSecret == "" {
		color.Set(color.FgRed)
		log.Println("error: bot_secret invalid")
		color.Unset()
		os.Exit(1)
	}

	if d.BotUID == "" {
		color.Set(color.FgRed)
		log.Println("error: bot_uid invalid")
		color.Unset()
		os.Exit(1)
	}

	return true
}

func LoadConf(deploy string) {
	homePath := os.Getenv("HOME")
	if homePath == "" {
		color.Set(color.FgRed)
		log.Println("error: HOME variable not found")
		color.Unset()
	}

	appDir := homePath + "/.securum_exire"
	_, err := os.Stat(appDir)
	if os.IsNotExist(err) {
		color.Set(color.FgRed)
		log.Println("error: .securum_exire (app directory) doesn't exit")
		color.Unset()
		err = os.MkdirAll(appDir, os.ModePerm)
		if err != nil {
			color.Set(color.FgRed)
			log.Println("error: unable to create a directory")
			color.Unset()
			os.Exit(1)
		}
	}


	configFilePath := appDir + "/signal_server.config.toml"
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		color.Set(color.FgRed)
		log.Println("error: signal_server.config.toml doesn't exit")
		color.Unset()
		_, err = os.Create(configFilePath)
		if err != nil {
			color.Set(color.FgRed)
			log.Println("error: unable to create signal_server.config.toml")
			color.Unset()
			os.Exit(1)
		}
	}

	var conf models.Conf
	if _, err := toml.DecodeFile(configFilePath, &conf); err != nil {
		color.Set(color.FgRed)
		log.Println("error: unable to parse signal_server.config.toml ")
		color.Unset()
		os.Exit(1)
	}
	if strings.ToUpper(deploy) == "PRODUCTION" {
		if conf.Production == nil {
			color.Set(color.FgRed)
			log.Println("error: Production config not defined")
			color.Unset()
			os.Exit(1)
		}
		ValidateDeploymentConfig(conf.Production)
		globals.DeploymentConfig = conf.Production
	} else if strings.ToUpper(deploy) == "DEVELOPMENT" {
		if conf.Development == nil {
			color.Set(color.FgRed)
			log.Println("error: Development config not defined")
			color.Unset()
			os.Exit(1)
		}
		ValidateDeploymentConfig(conf.Development)
		globals.DeploymentConfig = conf.Development
	} else if strings.ToUpper(deploy) == "STAGING" {
		if conf.Staging == nil {
			color.Set(color.FgRed)
			log.Println("error: Development config not defined")
			color.Unset()
			os.Exit(1)
		}
		ValidateDeploymentConfig(conf.Staging)
		globals.DeploymentConfig = conf.Staging
	}
}

