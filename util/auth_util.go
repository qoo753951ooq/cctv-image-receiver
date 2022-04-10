package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

//驗證CCTV
func BasicAuth(user, pwg string) string {
	auth := user + ":" + pwg
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func GetDigestParts(resp *http.Response) map[string]string {
	result := map[string]string{}
	if len(resp.Header["Www-Authenticate"]) > 0 {
		wantedHeaders := []string{"nonce", "realm", "qop", "algorithm"}
		responseHeaders := strings.Split(resp.Header["Www-Authenticate"][0], ",")
		for _, r := range responseHeaders {
			for _, w := range wantedHeaders {
				if strings.Contains(r, w) {
					result[w] = strings.Split(r, `"`)[1]
				}
			}
		}

	}
	return result
}

func GetDigestAuthrization(digestParts map[string]string) string {
	var ha1, ha2, response string
	d := digestParts

	getMD5 := func(text string) string {
		hasher := md5.New()
		hasher.Write([]byte(text))
		return hex.EncodeToString(hasher.Sum(nil))
	}

	getCnonce := func() string {
		b := make([]byte, 8)
		_, _ = io.ReadFull(rand.Reader, b)
		return fmt.Sprintf("%x", b)[:16]
	}
	cnonce := getCnonce()
	ha1 = getMD5(d["username"] + ":" + d["realm"] + ":" + d["password"])
	if strings.Compare(d["algorithm"], "MD5-sess") == 0 {
		ha1 = getMD5(ha1 + ":" + d["nonce"] + ":" + cnonce)
	}

	if strings.Compare(d["qop"], "auth-int") != 0 {
		ha2 = getMD5(d["method"] + ":" + d["uri"])
	}
	nonceCount := 00000001
	if len(d["qop"]) == 0 {
		response = getMD5(fmt.Sprintf("%s:%v:%s", ha1, nonceCount, ha2))
	} else {
		response = getMD5(fmt.Sprintf("%s:%s:%v:%s:%s:%s", ha1, d["nonce"], nonceCount, cnonce, d["qop"], ha2))
	}

	authorization := fmt.Sprintf(`Digest username="%s", realm="%s", nonce="%s", uri="%s", cnonce="%s", nc="%v", qop="%s", response="%s"`,
		d["username"], d["realm"], d["nonce"], d["uri"], cnonce, nonceCount, d["qop"], response)
	return authorization
}
