package auth

import "encoding/base64"

type Auth struct {
	User     string
	Password string
}

func (a Auth) BasicAuth() string {
	auth := a.User + ":" + a.Password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
