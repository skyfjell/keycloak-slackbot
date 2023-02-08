package api

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"keycloakslackbot/logs"
	"net/http"
	"net/url"
	"strings"
)

type User struct {
	CreatedTimeStamp int64  `json:"createdTimestamp"`
	Email            string `json:"email"`
}

type KeyCloak struct {
	Realm    string
	Host     string
	user     string
	password string
}

func NewKeyCloak(host string, realm string, user string, password string) KeyCloak {
	return KeyCloak{
		Realm:    realm,
		Host:     host,
		user:     user,
		password: password,
	}
}

// Generates token url
func (k *KeyCloak) tokenURL() string {
	u, _ := url.Parse(k.Host)
	u.Path = fmt.Sprintf("auth/realms/%s/protocol/openid-connect/token", k.Realm)
	return u.String()
}

// Generates token url
func (k *KeyCloak) userURL(query map[string]string) string {
	u, _ := url.Parse(k.Host)
	u.Path = fmt.Sprintf("auth/admin/realms/%s/users", k.Realm)
	if query != nil {
		q := u.Query()
		for key, elm := range query {
			q.Add(key, elm)
		}
		u.RawQuery = q.Encode()
	}
	return u.String()
}

// Utility for setting proper header auth on get token
func (k *KeyCloak) setAuthHeader(header *http.Header, content string, auth ...string) {
	var ab string
	if len(auth) == 0 {
		ab = "Basic " + b64.StdEncoding.EncodeToString([]byte(k.user+":"+k.password))
	} else {
		ab = "Bearer " + auth[0]
	}

	header.Set("Authorization", ab)
	header.Set("Content-Type", content)
}

func (k *KeyCloak) getToken() (string, error) {

	data := url.Values{"grant_type": []string{"client_credentials"}}
	URL := k.tokenURL()

	logs.Logger.Info(fmt.Sprintf("Posting to %s", URL))
	req, _ := http.NewRequest("POST", URL, strings.NewReader(data.Encode()))
	k.setAuthHeader(&req.Header, "application/x-www-form-urlencoded;charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Logger.Error("Error getting token POST")
		return "", err
	}
	defer resp.Body.Close()

	var body struct {
		Token string `json:"access_token"`
	}
	if err := readResponse(resp.Body, &body); err != nil {
		logs.Logger.Error("Error reading token body")
		return "", err
	}

	return body.Token, nil
}

func (k *KeyCloak) ListUsers() ([]User, error) {
	token, err := k.getToken()
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", k.userURL(map[string]string{"max": "1000"}), nil)
	k.setAuthHeader(&req.Header, "application/json", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logs.Logger.Error("Error getting users")
		return nil, err
	}

	defer resp.Body.Close()

	var results []User
	if err := readResponse(resp.Body, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func readResponse(r io.Reader, i interface{}) error {
	rawbody, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(rawbody, &i); err != nil {
		return err
	}
	return nil
}
