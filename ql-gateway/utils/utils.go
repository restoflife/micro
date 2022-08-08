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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
	"github.com/restoflife/micro/gateway/conf"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/protocol"
	"github.com/sony/sonyflake"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"net/url"
	"reflect"
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
// that the reverse proxy (nginx or haproxy) can work properly.
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
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // 加密处理
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

// DeleteSlice 删除切片index
func DeleteSlice(slice interface{}, index int) (interface{}, error) {
	sliceValue := reflect.ValueOf(slice)
	length := sliceValue.Len()
	if slice == nil || length == 0 || (length-1) < index {
		return nil, fmt.Errorf("error deleting slice")
	}
	if length-1 == index {
		return sliceValue.Slice(0, index).Interface(), nil
	}

	return reflect.AppendSlice(sliceValue.Slice(0, index), sliceValue.Slice(index+1, length)).Interface(), nil
}

func ConvertToString(v interface{}) (str string) {
	if v == nil {
		return
	}
	var (
		bs  []byte
		err error
	)
	if bs, err = json.Marshal(v); err != nil {
		return
	}
	str = string(bs)
	return
}

// ConvertStrSliceMap 将字符串 slice 转为 map[string]struct{}。
func ConvertStrSliceMap(sl []string) map[string]struct{} {
	set := make(map[string]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// InMap 判断字符串是否在 map 中。
func InMap(m map[string]struct{}, s string) bool {
	_, ok := m[s]
	return ok
}

// ToMapSetStrictE converts a slice or array to map set with error strictly.
// The result of map key type is equal to the element type of input.
func ToMapSetStrictE(i interface{}) (interface{}, error) {
	// check param
	if i == nil {
		return nil, fmt.Errorf("unable to converts %#v of type %T to map[interface{}]struct{}", i, i)
	}
	t := reflect.TypeOf(i)
	kind := t.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, fmt.Errorf("the input %#v of type %T isn't a slice or array", i, i)
	}
	// execute the convert
	v := reflect.ValueOf(i)
	mT := reflect.MapOf(t.Elem(), reflect.TypeOf(struct{}{}))
	mV := reflect.MakeMapWithSize(mT, v.Len())
	for j := 0; j < v.Len(); j++ {
		mV.SetMapIndex(v.Index(j), reflect.ValueOf(struct{}{}))
	}
	return mV.Interface(), nil
	// 	var sl = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	//	m, _ := ToMapSetStrictE(sl)
	//	mSet = m.(map[string]struct{})
	//	if _, ok := m["m"]; ok {
	//		fmt.Println("in")
	//	}
	//	if _, ok := m["mm"]; !ok {
	//		fmt.Println("not in")
	//	}

}

// AesEncrypt 加密
func AesEncrypt(data []byte) (r []byte) {
	block, _ := aes.NewCipher([]byte(conf.C.Encryption.Key))
	blockMode := cipher.NewCBCEncrypter(block, []byte(conf.C.Encryption.Iv))
	// PKCS5Padding填充
	data = pkcs5Padding(data, block.BlockSize())

	// salt(data)

	resultArray := make([]byte, len(data))

	blockMode.CryptBlocks(resultArray, data)

	r = make([]byte, base64.StdEncoding.EncodedLen(len(resultArray)))
	base64.StdEncoding.Encode(r, resultArray)
	return
}
func pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// AesDecrypt 解密
func AesDecrypt(data []byte) (r []byte, err error) {
	dataArray := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	// base64解码
	n, err := base64.StdEncoding.Decode(dataArray, data)
	if err != nil {
		return
	}
	dataArray = dataArray[:n]
	log.Debug(zap.String("key", conf.C.Encryption.Key))
	log.Debug(zap.String("iv", conf.C.Encryption.Iv))
	block, _ := aes.NewCipher([]byte(conf.C.Encryption.Key))
	blockMode := cipher.NewCBCDecrypter(block, []byte(conf.C.Encryption.Iv))

	// 创建解密后的byte数组的内存空间，其与密文长度一致
	resultArray := make([]byte, len(dataArray))

	blockMode.CryptBlocks(resultArray, dataArray)

	r = pkcs5UnPadding(resultArray)
	return
}
func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:(length - unPadding)]
}

func salt(origin []byte) {
	l := len(origin)
	for i := 0; i < l; i += 2 {
		origin[i] ^= 0x9
	}

	for j := 1; j < l; j += 2 {
		origin[j] ^= 0x7
	}
}

func MakeBlocksOrigin(src []byte) []byte {

	// 1. 获取src长度
	length := len(src)

	// 2. 得到最后一个字符
	lastChar := src[length-1] // '4'

	// 3. 将字符转换为数字
	number := int(lastChar) // 4

	// 4. 截取需要的长度
	return src[:length-number]
}

func MakeBlocksFull(src []byte, blockSize int) []byte {

	// 1. 获取src的长度， blockSize对于des是8
	length := len(src)

	// 2. 对blockSize进行取余数， 4
	remains := length % blockSize

	// 3. 获取要填的数量 = blockSize - 余数
	paddingNumber := blockSize - remains // 4

	// 4. 将填充的数字转换成字符， 4， '4'， 创建了只有一个字符的切片
	// s1 = []byte{'4'}
	s1 := []byte{byte(paddingNumber)}

	// 5. 创造一个有4个'4'的切片
	// s2 = []byte{'4', '4', '4', '4'}
	s2 := bytes.Repeat(s1, paddingNumber)

	// 6. 将填充的切片追加到src后面
	s3 := append(src, s2...)

	return s3
}
