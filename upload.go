package sender

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

const audioMessageType = "audio_message"

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
			ID    int `json:"id,omitempty"`
			Owner int `json:"owner_id,omitempty"`
		} `json:"audio_message,omitempty"`
	} `json:"response,omitempty"`
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

func uploadFileAndGetData(uploadURI string, auth Auth) string {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	file, err := os.Open(auth.Filename)
	checkErr(err)
	defer file.Close()

	part, err := writer.CreateFormFile("file", auth.Filename)
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
		ID:    document.Response.Message.ID,
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
