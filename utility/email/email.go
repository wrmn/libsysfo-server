package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func (content Content) SendEmail(receiver ToData) (err error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	values := data{
		Sender: senderData{
			Name:  "Libsysfo Admin",
			Email: "moawarman@gmail.com",
		},
		To:          []ToData{receiver},
		Subject:     content.Subject,
		HtmlContent: content.HtmlContent,
	}
	json_data, err := json.Marshal(values)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", "https://api.sendinblue.com/v3/smtp/email", bytes.NewBuffer(json_data))
	if err != nil {
		return
	}
	req.Header.Set("user-agent", "golang application")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("api-key", os.Getenv("SENDINBLUE_TOKEN"))

	response, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	defer response.Body.Close()
	fmt.Printf("Code: %d\n", response.StatusCode)
	fmt.Printf("Body: %s\n", body)
	return nil
}
