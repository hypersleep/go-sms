package main

import(
	"log"
	"time"
	"bytes"
	"net/url"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"crypto/sha256"
	"encoding/base64"
)

type BillingradMessage struct {
	Did     string `json:"did"`
	To      string `json:"to"`
	Text    string `json:"text"`
	Planned string `json:"planned"`
}

func billingradSign(close, data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(close + data))
	escaped := url.QueryEscape(string(base64.URLEncoding.EncodeToString(hasher.Sum(nil))))
	escaped = strings.Replace(escaped, "-", "%2B", -1)
	escaped = strings.Replace(escaped, "_", "%2F", -1)
	return escaped
}

func sendBillingrad(message *BillingradMessage) error {
	var err error

	url := "http://my.billingrad.com/api/delivery/createMessage?_open=" + serverConfig.OpenApiKey

	b, err := json.Marshal(message)
	if err != nil {
		log.Println("Failed to marshal JSON:", err)
		return err
	}

	data := string(b)
	log.Println(data)

	url += "&_key=" + billingradSign(serverConfig.CloseApiKey, data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	return nil
}

func billingradSender() {
	var message *BillingradMessage
	var err error

	for {
		message = <- messagesChannel
		err = sendBillingrad(message)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			messagesChannel <- message
		}
	}
}
