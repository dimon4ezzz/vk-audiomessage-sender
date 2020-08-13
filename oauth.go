package sender

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const act = "authcheck_code"

const (
	authcheckRegexp = `authcheck_code&hash=([^\"]+)`
	tokenRegexp     = `access_token=([0-9a-f]+)`
)

const (
	grantType = "password"
	version   = "5.120"
)

type OauthResponse struct {
	URI   string `json:"redirect_uri,omitempty"`
	Token string `json:"access_token,omitempty"`
}

var (
	errResponseStatus = errors.New("response is not 200")
	err2faCode        = errors.New("your 2fa code is not a number")
	errFragmentToken  = errors.New("cannot find access token")
	errHashNotFound   = errors.New("not found match for hash")
)

func getAccessToken(auth Auth) string {
	response := getOauthResponse(auth)

	if response.Token != "" {
		return response.Token
	}

	// get hash for special oauth link
	hash := getAuthHash(response.URI)

	code := get2faCode()
	loginURI := getLoginURI(hash, code)

	resp, err := client.Get(loginURI)
	checkErr(err)
	defer resp.Body.Close()

	// it can be 301 when wrong
	if resp.StatusCode != 200 {
		log.Fatal(errResponseStatus)
	}

	fragment := resp.Request.URL.Fragment
	token := getTokenFromFragment(fragment)
	return token
}

func getOauthResponse(auth Auth) *OauthResponse {
	uri := getOauthURI(auth)

	resp, err := client.Get(uri)
	checkErr(err)
	defer resp.Body.Close()

	oauthResp := &OauthResponse{}
	getStructFromJSON(resp.Body, oauthResp)

	return oauthResp
}

func getOauthURI(auth Auth) string {
	rawQuery := getOauthRawQuery(auth)

	uri := &url.URL{
		Scheme:   "https",
		Host:     "oauth.vk.com",
		Path:     "token",
		RawQuery: rawQuery,
	}

	return uri.String()
}

func getOauthRawQuery(auth Auth) string {
	query := url.Values{}
	query.Set("client_id", strconv.Itoa(auth.ClientID))
	query.Add("client_secret", auth.ClientSecret)
	query.Add("grant_type", grantType)
	query.Add("password", auth.Password)
	query.Add("username", auth.Username)
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

func get2faCode() string {
	reader := bufio.NewReader(os.Stdin)
	println("enter 2fa code:")
	text, err := reader.ReadString('\n')
	checkErr(err)
	text = strings.TrimSpace(text)
	ok, err := regexp.MatchString(`\d+`, text)
	if !ok || err != nil {
		log.Fatal(err2faCode)
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
		// [0] is full string
		// [1] is the first substring
		return matches[1]
	}

	panic(errFragmentToken)
}
