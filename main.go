package main

import (
	"fmt"
	"log"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

func main() {
	auth, err := ioutil.ReadFile("ServiceAccount.json")
	if err != nil {
		log.Fatal(err)
	}

	cred, err := google.CredentialsFromJSON(context.Background(), auth, "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		log.Fatalln(err)
	}

	token, err := cred.TokenSource.Token()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(token.AccessToken)

	sendNotification(token.AccessToken)
}

func sendNotification(bearer string) {
	url := "https://fcm.googleapis.com/v1/projects/[PROJECT_ID]/messages:send"

	notification := map[string]interface{}{
		"message": map[string]interface{}{
			"token": "[USER_TOKEN]",
			"notification": map[string]interface{}{
				"title": "Notification Title",
				"body": "Notification Body Text",
			},
			"data": map[string]interface{}{
				"data_1": "Some Data",
				"data_2": "More Data",
				"data_3": "And some more",
			},
			"apns": map[string]interface{}{
				"payload": map[string]interface{}{
					"aps" : map[string]interface{}{
						"category": "push-category",
						"mutable-content": 1,
						"content-available": 1,
					},
				},
			},
		},
	}

	data, err := json.Marshal(notification)
	if err != nil {
		log.Fatal(err)
	}

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearer))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("Request error: %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("--------------------------------")
		fmt.Println(string(data))
		fmt.Println("--------------------------------")
	}
}

