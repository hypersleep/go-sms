package main

import(
	"log"
	"errors"
	"net/http"
	"io/ioutil"
)

func auth(password string) error {
	if password != serverConfig.AuthToken {
		return errors.New("Token not match")
	}
	return nil
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "POST":

			err := auth(r.Header.Get("X-AUTH-TOKEN"))
			if err != nil {
				log.Println(err)
				http.Error(w, "Go away", http.StatusForbidden)
				return
			}

			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("Failed to read request body:", err)
				return
			}

			for i := range serverConfig.Receivers {
				message := &BillingradMessage{
					Did: serverConfig.DeliveryID,
					To: serverConfig.Receivers[i],
					Text: string(body),
					Planned: "1",
				}

				messagesChannel <- message
			}
	}
}

var serverConfig *ServerConfig = &ServerConfig{}
var messagesChannel chan *BillingradMessage = make(chan *BillingradMessage, 10)

func main() {
	serverConfig.loadConfig()
	go billingradSender()
	http.HandleFunc("/v1/send", sendHandler)
	log.Println("Server running on port", serverConfig.Port)
	http.ListenAndServe(":" + serverConfig.Port, nil)
}