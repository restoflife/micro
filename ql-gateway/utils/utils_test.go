/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2022-08-08 15:48
 * @LastEditors: Administrator
 * @LastEditTime: 2022-08-08 15:48
 * @FilePath: ql-gateway/utils/utils_test.go
 */

package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestAesDecrypt(t *testing.T) {
	key := "aFs25NHp2Kt8KYw6HxNZmET0"
	iv := "t19yZlW85z59RBTF"
	data := "M0U0O56KRZ9h/65HzlKw4ZqiNwe2/Lyi6Rt67iBNeDho5A8qadmnrXdgpobdBWXnxwYpxTkK75SlJ0pX4qYbGEZx/c+RV17y9eGw2v52iOtS7pm+1RjJam0bATMC5W5G1g9peqnfA5AbP7h1qA2WdQ=="
	dataArray := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	// base64解码
	n, err := base64.StdEncoding.Decode(dataArray, []byte(data))
	if err != nil {
		fmt.Println(err, "----------err")
		return
	}
	dataArray = dataArray[:n]
	block, _ := aes.NewCipher([]byte(key))
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))

	// 创建解密后的byte数组的内存空间，其与密文长度一致
	resultArray := make([]byte, len(dataArray))

	blockMode.CryptBlocks(resultArray, dataArray)

	r := pkcs5UnPadding(resultArray)
	fmt.Println(string(r))
}
func TestAesEncrypt(t *testing.T) {
	key := "aFs25NHp2Kt8KYw6HxNZmET0"
	iv := "t19yZlW85z59RBTF"
	data := []byte(`{"account":"yinggdfsfsd","password":"123456","username":"胡2磊","avatar":"http://dummyimage.com/100x100"}`)
	block, _ := aes.NewCipher([]byte(key))
	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	// PKCS5Padding填充
	data = pkcs5Padding(data, block.BlockSize())

	resultArray := make([]byte, len(data))

	blockMode.CryptBlocks(resultArray, data)

	r := make([]byte, base64.StdEncoding.EncodedLen(len(resultArray)))
	base64.StdEncoding.Encode(r, resultArray)
	fmt.Println(string(r))
}
