package util

import (
	"encoding/base64"
)

//驗證CCTV
func BasicAuth(user, pwg string) string {
	auth := user + ":" + pwg
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
