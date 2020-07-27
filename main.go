package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	clientID     = 0
	clientSecret = ""
	grantType    = "password"
	username     = ""
	password     = ""
	version      = "5.120"
)

const act = "authcheck_code"
const audioMessageType = "audio_message"
const filename = ""
const recipient = 0

const (
	authcheckRegexp = `authcheck_code&hash=([^\"]+)`
	tokenRegexp     = `access_token=([0-9a-f]+)`
)

type OauthResponse struct {
	URI string `json:"redirect_uri,omitempty"`
}

type UploadURLResponse struct {
	Response struct {
		URI string `json:"upload_url,omitempty"`
	} `json:"response,omitempty"`
}

type FileResponse struct {
	File string `json:"file,omitempty"`
}

type DocumentResponse struct {
	Response struct {
		Message struct {
			Id    int `json:"id,omitempty"`
			Owner int `json:"owner_id,omitempty"`
		} `json:"audio_message,omitempty"`
	} `json:"response,omitempty"`
}

type MessageResponse struct {
	Response int `json:"response,omitempty"`
}

type Document struct {
	Owner int
	ID    int
}

type Message struct {
	User int
	ID   int
}

var client http.Client

var (
	errConstantNotSet = errors.New("some constant did not set")
	errResponseStatus = errors.New("response is not 200")
	err2faCode        = errors.New("your 2fa code is not a number")
	errFragmentToken  = errors.New("cannot find access token")
	errHashNotFound   = errors.New("not found match for hash")
)

func main() {
	checkVariables()
	setupClient()

	redirectURI := getRedirectURI()
	hash := getAuthHash(redirectURI)
	token := getAccessToken(hash)
	uploadURI := getUploadServer(token)
	fileData := uploadFileAndGetData(uploadURI)
	document := getDocument(fileData, token)
	message := sendMessage(document, token)

	println("look at " + message.String())
}

func checkVariables() {
	if clientID == 0 ||
		clientSecret == "" ||
		username == "" ||
		password == "" ||
		filename == "" ||
		recipient == 0 {
		panic(errConstantNotSet)
	}
}

func setupClient() {
	cookieJar, err := cookiejar.New(nil)
	checkErr(err)
	client = http.Client{
		Jar: cookieJar,
	}
}

func getRedirectURI() string {
	uri := getOauthURI()

	resp, err := client.Get(uri)
	checkErr(err)
	defer resp.Body.Close()

	oauthResp := &OauthResponse{}
	getStructFromJSON(resp.Body, oauthResp)

	return oauthResp.URI
}

func getOauthURI() string {
	rawQuery := getOauthRawQuery()

	uri := &url.URL{
		Scheme:   "https",
		Host:     "oauth.vk.com",
		Path:     "token",
		RawQuery: rawQuery,
	}

	return uri.String()
}

func getOauthRawQuery() string {
	query := url.Values{}
	query.Set("client_id", strconv.Itoa(clientID))
	query.Add("client_secret", clientSecret)
	query.Add("grant_type", grantType)
	query.Add("password", password)
	query.Add("username", username)
	query.Add("version", version)
	query.Add("2fa_supported", "1")

	return query.Encode()
}

func getAuthHash(oauthRedirectURI string) string {
	resp, err := client.Get(oauthRedirectURI)
	checkErr(err)
	defer resp.Body.Close()

	re := regexp.MustCompile(authcheckRegexp)
	hash := getMatchFromReader(resp.Body, re)
	return hash
}

func getMatchFromReader(body io.ReadCloser, re *regexp.Regexp) string {
	sc := bufio.NewScanner(body)

	var line string
	var matches []string
	for sc.Scan() {
		line = sc.Text()
		matches = re.FindStringSubmatch(line)
		if len(matches) != 0 {
			// [0] is full string
			// [1] is the first substring
			return matches[1]
		}
	}

	panic(errHashNotFound)
}

func getAccessToken(hash string) string {
	code := get2faCode()
	loginURI := getLoginURI(hash, code)

	resp, err := client.Get(loginURI)
	checkErr(err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(errResponseStatus)
	}

	fragment := resp.Request.URL.Fragment
	token := getTokenFromFragment(fragment)
	return token
}

func get2faCode() string {
	reader := bufio.NewReader(os.Stdin)
	println("enter 2fa code:")
	text, err := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	checkErr(err)
	ok, err := regexp.MatchString(`\d+`, text)
	if !ok || err != nil {
		panic(err2faCode)
	}
	return text
}

func getLoginURI(hash string, code string) string {
	rawQuery := getLoginRawQuery(hash, code)

	uri := url.URL{
		Scheme:   "https",
		Host:     "m.vk.com",
		Path:     "login",
		RawQuery: rawQuery,
	}

	return uri.String()
}

func getLoginRawQuery(hash string, code string) string {
	query := url.Values{}
	query.Set("act", act)
	query.Add("hash", hash)
	query.Add("code", code)

	return query.Encode()
}

func getTokenFromFragment(fragment string) string {
	re := regexp.MustCompile(tokenRegexp)
	matches := re.FindStringSubmatch(fragment)
	if len(matches) != 0 {
		return matches[1]
	}

	panic(errFragmentToken)
}

func getUploadServer(token string) string {
	uri := getMessagesUploadServerURI(token)

	resp, err := client.Get(uri)
	checkErr(err)
	defer resp.Body.Close()

	uploadURLResponse := &UploadURLResponse{}
	getStructFromJSON(resp.Body, uploadURLResponse)

	return uploadURLResponse.Response.URI
}

func getMessagesUploadServerURI(token string) string {
	rawQuery := getMessagesUploadServerRawQuery(token)

	uri := url.URL{
		Scheme:   "https",
		Host:     "api.vk.com",
		Path:     "method/docs.getMessagesUploadServer",
		RawQuery: rawQuery,
	}

	return uri.String()
}

func getMessagesUploadServerRawQuery(token string) string {
	query := url.Values{}
	query.Set("access_token", token)
	query.Add("type", audioMessageType)
	query.Add("v", version)

	return query.Encode()
}

func uploadFileAndGetData(uploadURI string) string {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	file, err := os.Open(filename)
	checkErr(err)
	defer file.Close()

	part, err := writer.CreateFormFile("file", filename)
	checkErr(err)

	_, err = io.Copy(part, file)
	checkErr(err)

	err = writer.Close()
	checkErr(err)

	req, err := http.NewRequest("POST", uploadURI, buf)
	checkErr(err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	checkErr(err)
	defer resp.Body.Close()

	fileResponse := &FileResponse{}
	getStructFromJSON(resp.Body, fileResponse)

	return fileResponse.File
}

func getDocument(fileData string, token string) Document {
	uri := getSaveURI(fileData, token)

	resp, err := client.Get(uri)
	checkErr(err)
	defer resp.Body.Close()

	document := &DocumentResponse{}
	getStructFromJSON(resp.Body, document)

	return Document{
		ID:    document.Response.Message.Id,
		Owner: document.Response.Message.Owner,
	}
}

func getSaveURI(fileData string, token string) string {
	rawQuery := getSaveRawQuery(fileData, token)

	uri := url.URL{
		Scheme:   "https",
		Host:     "api.vk.com",
		Path:     "method/docs.save",
		RawQuery: rawQuery,
	}

	return uri.String()
}

func getSaveRawQuery(fileData string, token string) string {
	query := url.Values{}
	query.Set("access_token", token)
	query.Add("file", fileData)
	query.Add("v", version)

	return query.Encode()
}

func sendMessage(document Document, token string) Message {
	uri := getMessagesSendURI(document, token)

	resp, err := client.Get(uri)
	checkErr(err)
	defer resp.Body.Close()

	messageResponse := &MessageResponse{}
	getStructFromJSON(resp.Body, messageResponse)

	message := Message{
		User: recipient,
		ID:   messageResponse.Response,
	}

	return message
}

func getMessagesSendURI(document Document, token string) string {
	rawQuery := getMessageSendRawQuery(document, token)

	uri := url.URL{
		Scheme:   "https",
		Host:     "api.vk.com",
		Path:     "method/messages.send",
		RawQuery: rawQuery,
	}

	return uri.String()
}

func getMessageSendRawQuery(document Document, token string) string {
	// get 5-digit number
	random := rand.Intn(90000) + 9999
	randomStr := strconv.Itoa(random)

	query := url.Values{}
	query.Set("access_token", token)
	query.Add("user_id", strconv.Itoa(recipient))
	query.Add("attachment", document.String())
	query.Add("v", version)
	query.Add("random_id", randomStr)

	return query.Encode()
}

func (doc Document) String() string {
	builder := strings.Builder{}
	builder.WriteString("doc")
	builder.WriteString(strconv.Itoa(doc.Owner))
	builder.WriteString("_")
	builder.WriteString(strconv.Itoa(doc.ID))

	return builder.String()
}

func (mes Message) String() string {
	builder := strings.Builder{}
	builder.WriteString("https://vk.com/im?sel=")
	builder.WriteString(strconv.Itoa(mes.User))
	builder.WriteString("&msgid=")
	builder.WriteString(strconv.Itoa(mes.ID))

	return builder.String()
}

func getStructFromJSON(body io.ReadCloser, obj interface{}) {
	err := json.NewDecoder(body).Decode(obj)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
