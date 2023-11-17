package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AfroSMS struct {
	log   *log.Logger
	Name  string
	Token string
	URL   string
}

func New(log *log.Logger) AfroSMS {
	return AfroSMS{
		log:   log,
		Name:  "Afro SMS",
		Token: "eyJhbGciOiJIUzI1NiJ9.eyJpZGVudGlmaWVyIjoiZ3cwWXN0cXlTZjNOTmlDdlpiTGhkSlVCREY3R1A4ZGsiLCJleHAiOjE4NDUzNTE4NjEsImlhdCI6MTY4NzQ5OTA2MSwianRpIjoiODQxZGYyZGQtOWZlMy00ODc4LTkxN2QtM2VkMjY4Nzk2ZjAxIn0.AAkhdnl7ma0boH7wKbsu76bQce4Xel5Kf9JzqzMnGJ0",
		URL:   "https://api.afromessage.com/api/send",
	}
}

func (sms AfroSMS) SendSMS(phone, message string) error {
	sms.log.Println("Send SMS")
	type SMSBody struct {
		From string `json:"from"`
		// Sender   string
		To      string `json:"to"`
		Message string `json:"message"`
		// Callback string
	}
	sms.log.Println(phone)
	sms.log.Println(message)
	body := SMSBody{
		From: "e80ad9d8-adf3-463f-80f4-7c4b39f7f164",
		// Sender:  "LakiPay_",
		To:      phone,
		Message: message,
	}
	serBody, err := json.Marshal(body)
	if err != nil {
		sms.log.Println(err)
		return err
	}
	req, err := http.NewRequest(http.MethodPost, sms.URL, bytes.NewReader(serBody))
	if err != nil {
		sms.log.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+sms.Token)
	// var i interface{}
	// err = json.NewDecoder(req.Body).Decode(&i)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s", i)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		sms.log.Println(err)
		return err
	}

	defer res.Body.Close()

	// sms.log.Println(res.StatusCode)
	// sms.log.Println(res.Body)
	// sms.log.Println(res)

	var j interface{}
	err = json.NewDecoder(res.Body).Decode(&j)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", j)

	if res.StatusCode != 200 {
		sms.log.Println(res)
		msg, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(msg))
	}
	return nil
}
