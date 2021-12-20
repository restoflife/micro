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
	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
	"github.com/restoflife/micro/gateway/internal/protocol"
	"github.com/sony/sonyflake"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	store = base64Captcha.DefaultMemStore
	sf    = sonyflake.NewSonyflake(sonyflake.Settings{
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

// VerifyCode 验证验证码
func VerifyCode(uuid, code string) bool {
	return store.Verify(uuid, code, true)
}

// DriverDigitFunc 生成验证码
func DriverDigitFunc() (id, b64s string, err error) {
	e := protocol.ConfigJsonBody{}
	e.Id = uuid.New().String()
	e.DriverDigit = base64Captcha.DefaultDriverDigit
	driver := e.DriverDigit
	captcha := base64Captcha.NewCaptcha(driver, store)
	return captcha.Generate()
}
