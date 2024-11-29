package shelly

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
)

type Auth struct {
	User     string
	Password string
}

func (a Auth) basicAuth() string {
	auth := a.User + ":" + a.Password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func Request(url *url.URL, auth *Auth) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	if auth != nil {
		req.Header.Add("Authorization", "Basic "+auth.basicAuth())
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, err
}
