package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func booksTestHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", "https://laravel-books-db.herokuapp.com/api/books?page=1&category=historical-fiction&language=en", nil)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	req.Header.Set("user-agent", "golang application")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer 39|6leo0PRGqqn1nyo47ykqpTGunoEJnAAHss9OxIXs")
	response, err := client.Do(req)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	var template interface{}

	err = json.Unmarshal(b, &template)
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	fmt.Println(template.(map[string]interface{}))

	json.NewEncoder(w).Encode(template)

}
