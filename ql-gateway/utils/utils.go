/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 17:04
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 17:04
 * @FilePath: ql-gateway/utils/utils.go
 */

package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func GetUrls(u string) (addrs []string, err error) {
	addr, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	values := strings.Split(addr.Host, ",")
	for _, value := range values {
		if strings.Contains(addr.String(), "https") {
			addrs = append(addrs, fmt.Sprintf("https://%s", value))
		} else {
			addrs = append(addrs, fmt.Sprintf("http://%s", value))
		}

	}
	return addrs, err
}

// ClientIp Resolve x-real-ip and x-forwarded-for so
//that the reverse proxy (nginx or haproxy) can work properly.
func ClientIp(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return ip
	}
	return ""
}

// MD5String md5 digest in string
func MD5String(plain string) string {
	cipher := MD5([]byte(plain))
	return hex.EncodeToString(cipher)
}

// MD5 md5 digest
func MD5(plain []byte) []byte {
	md5Ctx := md5.New()
	md5Ctx.Write(plain)
	cipher := md5Ctx.Sum(nil)
	return cipher
}
