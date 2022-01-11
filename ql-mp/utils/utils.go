/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-23 16:15
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-23 16:15
 * @FilePath: ql-mp/utils/utils.go
 */

package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/restoflife/micro/mp/internal/component/log"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	sf = sonyflake.NewSonyflake(sonyflake.Settings{
		MachineID: func() (uint16, error) {
			return 2, nil
		},
		StartTime: time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local),
	})
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

// EncryptionPassword 加密密码
func EncryptionPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CompareHashAndPassword 比较密码 e 密码 p 需验证密码
func CompareHashAndPassword(e, p string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(e), []byte(p))
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetUUID 生成uuid
func GetUUID() (uint64, error) {
	nextID, err := sf.NextID()
	if err != nil {
		return 0, err
	}
	return nextID, nil
}

// Get HTTP GET request
func Get(url string) (int, []byte) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Error(zap.Error(err))
		return http.StatusInternalServerError, nil
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Error(zap.Error(err))
		}
	}(resp.Body)
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error(zap.Any("http get ", fmt.Sprintf("error decoding response from GET request, url: %s, %q", url, err)))
		}
	}
	return resp.StatusCode, result.Bytes()
}

func PageIndex(page, pageSize int) (limit, offset int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset = (page - 1) * pageSize
	limit = pageSize
	return
}
