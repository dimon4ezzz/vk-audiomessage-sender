package sender

import (
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

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

func sendMessage(document Document, token string, auth Auth) Message {
	uri := getMessagesSendURI(document, token, auth)
	println(uri)
	println(document.String())

	resp, err := client.Get(uri)
	checkErr(err)
	defer resp.Body.Close()
	resp.Write(os.Stdout)

	messageResponse := &MessageResponse{}
	getStructFromJSON(resp.Body, messageResponse)

	message := Message{
		User: auth.Recipient,
		ID:   messageResponse.Response,
	}

	return message
}

func getMessagesSendURI(document Document, token string, auth Auth) string {
	rawQuery := getMessageSendRawQuery(document, token, auth)

	uri := url.URL{
		Scheme:   "https",
		Host:     "api.vk.com",
		Path:     "method/messages.send",
		RawQuery: rawQuery,
	}

	return uri.String()
}

func getMessageSendRawQuery(document Document, token string, auth Auth) string {
	randomStr := getRandomString()
	query := url.Values{}
	query.Set("access_token", token)
	query.Add("user_id", strconv.Itoa(auth.Recipient))
	query.Add("attachment", document.String())
	query.Add("v", version)
	query.Add("random_id", randomStr)

	return query.Encode()
}

func getRandomString() string {
	randomSource := rand.NewSource(time.Now().UnixNano())
	r := rand.New(randomSource)
	// get 5-digit number
	random := r.Intn(90000) + 9999
	return strconv.Itoa(random)
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
