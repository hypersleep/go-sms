package main

import(
	"os"
	"log"
	"errors"
	"strings"
	"time"
	"net/http"

	consul "github.com/hashicorp/consul/api"
)

type ServerConfig struct {
	Port        string
	OpenApiKey  string
	CloseApiKey string
	DeliveryID  string
	Receivers   []string
	AuthToken   string
	AppEnv      string
}

var defaults map[string]string = map[string]string{
	"port": "8080",
}

func fetchConsul(kv *consul.KV, env string, key string) (value string, conf_err error) {
	pair, _, _ := kv.Get("go-sms/" + env + "/" + key, nil)
	if pair == nil {
		conf_err = errors.New("Failed to fetch \"" + key + "\" key for you environment. Leave default: " + defaults[key])
		value = defaults[key]
		return
	}

	value = string(pair.Value)
	return
}

func fetchConfig(continuousFetch bool) {
	var duration time.Duration

	if continuousFetch {
		duration = time.Duration(5 * time.Second)
	}

	FetchLoop:
		for {
			<- time.After(duration)

			consulConfig := &consul.Config{
				Address:    "consul.service.consul:8500",
				Scheme:     "http",
				HttpClient: http.DefaultClient,
			}

			client, err := consul.NewClient(consulConfig)
			if err != nil {
				log.Println("Failed to connect consul", err)
				return
			}

			kv := client.KV()

			serverConfig.Port, err = fetchConsul(kv, serverConfig.AppEnv, "port")
			if err != nil {
				log.Println(err)
			}

			serverConfig.OpenApiKey, err = fetchConsul(kv, serverConfig.AppEnv, "open_api_key")
			if err != nil {
				log.Println(err)
			}

			serverConfig.CloseApiKey, err = fetchConsul(kv, serverConfig.AppEnv, "close_api_key")
			if err != nil {
				log.Println(err)
			}

			serverConfig.DeliveryID, err = fetchConsul(kv, serverConfig.AppEnv, "delivery_id")
			if err != nil {
				log.Println(err)
			}

			receivers, err := fetchConsul(kv, serverConfig.AppEnv, "receivers")
			if err != nil {
				log.Println(err)
			}

			serverConfig.Receivers = strings.Split(receivers, ",")

			serverConfig.AuthToken, err = fetchConsul(kv, serverConfig.AppEnv, "auth_token")
			if err != nil {
				log.Println(err)
			}


			if duration == time.Duration(0) { break FetchLoop }
		}
}

func (serverConfig *ServerConfig) loadConfig() {
	serverConfig.AppEnv = "development"

	env := os.Getenv("_ENV")
	if env != "" { serverConfig.AppEnv = env }

	log.Println("Running in", serverConfig.AppEnv, "environment")

	fetchConfig(false)
	go fetchConfig(true)
}
