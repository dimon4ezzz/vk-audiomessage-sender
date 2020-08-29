package sender

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const encodingStringLength = 32

func encode(password string, token string, file *os.File) {
	password = normalizeEncodingString(password)

	text := []byte(token)
	key := []byte(password)

	ci, err := aes.NewCipher(key)
	checkErr(err)

	gcm, err := cipher.NewGCM(ci)
	checkErr(err)

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	seal := gcm.Seal(nonce, nonce, text, nil)

	_, err = file.Write(seal)
	checkErr(err)
}

func decode(password string, file *os.File) string {
	password = normalizeEncodingString(password)
	key := []byte(password)

	ci, err := aes.NewCipher(key)
	checkErr(err)

	gcm, err := cipher.NewGCM(ci)
	checkErr(err)

	data, err := ioutil.ReadAll(file)
	checkErr(err)

	size := gcm.NonceSize()
	if size > len(data) {
		log.Fatal("bad nonce size")
	}

	nonce, data := data[:size], data[size:]
	token, err := gcm.Open(nil, nonce, data, nil)
	checkErr(err)

	return string(token)
}

func normalizeEncodingString(str string) string {
	length := len(str)
	if length < encodingStringLength {
		strIndex := 0
		for i := length; i < encodingStringLength; i++ {
			str += string(str[strIndex])
			if strIndex >= length {
				strIndex = 0
			} else {
				strIndex++
			}
		}
		return str
	}

	return str[:encodingStringLength]
}
