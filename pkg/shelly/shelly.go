package shelly

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mongodb-forks/digest"
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

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("request to %s returned non 200 code: %d", url.String(), resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, err
}

func DigestAuthedRequest(url *url.URL, auth *Auth, params map[string]string) ([]byte, error) {
	t := digest.NewTransport(auth.User, auth.Password)
	client, err := t.Client()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request to %s returned non 200 code: %d", url.String(), resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, err
}
