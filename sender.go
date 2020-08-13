package sender

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
)

// Auth is data for send message
type Auth struct {
	ClientID     int    // VK Application Client ID
	ClientSecret string // VK Application Client secret
	Username     string // Sender username, e.g. email
	Password     string // Sender password in plain text
	Filename     string // File to send; must be .ogg (mono)
	Recipient    int    // Recipient VK ID
}

var client http.Client

var errConstantNotSet = errors.New("some Auth field did not set")

// Send is a main function to send audio message
func Send(auth Auth) {
	// check all fields in Auth were set
	checkVariables(auth)
	// create http client
	setupClient()

	// get vk access token
	token := getAccessToken(auth)
	// get upload uri
	uploadURI := getUploadServer(token)
	// upload file and get vk internal data about file
	fileData := uploadFileAndGetData(uploadURI, auth)
	// save file into personal documents and get link
	document := getDocument(fileData, token)
	// send message with audiomessage document
	message := sendMessage(document, token, auth)

	println("look at " + message.String())
}

func checkVariables(auth Auth) {
	if auth.ClientID == 0 ||
		auth.ClientSecret == "" ||
		auth.Username == "" ||
		auth.Password == "" ||
		auth.Filename == "" ||
		auth.Recipient == 0 {
		log.Fatal(errConstantNotSet)
	}
}

func setupClient() {
	cookieJar, err := cookiejar.New(nil)
	checkErr(err)
	client = http.Client{
		Jar: cookieJar,
	}
}

func getStructFromJSON(body io.ReadCloser, obj interface{}) {
	err := json.NewDecoder(body).Decode(obj)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
