/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2022-08-08 10:25
 * @LastEditors: Administrator
 * @LastEditTime: 2022-08-08 10:25
 * @FilePath: ql-gateway/internal/middleware/crypto.go
 */

package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/constant"
	"github.com/restoflife/micro/gateway/internal/encoding"
	"github.com/restoflife/micro/gateway/internal/errutil"
	"github.com/restoflife/micro/gateway/internal/protocol"
	"github.com/restoflife/micro/gateway/utils"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type aesWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *aesWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *aesWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func AesDecrypt() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
		} else {
			handleAes(c)
		}
	}
}
func handleAes(c *gin.Context) {
	contentType := c.Request.Header.Get("Content-Type")
	isJson := strings.Contains(contentType, constant.Types[constant.TypeJSON])
	isFile := strings.Contains(contentType, constant.Types[constant.TypeMultipartFormData])
	isFormUrl := strings.Contains(contentType, constant.Types[constant.TypeForm])
	if c.Request.Method == http.MethodGet {
		err := parseQuery(c)
		if err != nil {
			encoding.Error(c, errutil.ErrIllegalParameter)
			return
		}
	} else if isJson {
		err := parseJson(c)
		if err != nil {
			encoding.Error(c, errutil.ErrIllegalParameter)
			return
		}
	} else if isFormUrl {
		err := parseForm(c)
		if err != nil {
			encoding.Error(c, errutil.ErrIllegalParameter)
			return
		}
	} else if isFile {
		err := parseFile(c)
		if err != nil {
			encoding.Error(c, errutil.ErrIllegalParameter)
			return
		}
	}

	// 截取 response body
	oldWriter := c.Writer
	blw := &aesWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw

	c.Next()

	// 获取返回数据
	responseByte := blw.body.Bytes()

	c.Writer = oldWriter
	// 如果返回的不是json格式 那么直接返回,应为文件下载之类的不应该加密
	if !isJsonResponse(c) {
		_, _ = c.Writer.Write(responseByte)
		return
	}

	// 加密
	// encryptStr := utils.AesEncrypt(responseByte)
	encryptStr := responseByte

	_, _ = c.Writer.WriteString(string(encryptStr))
}

func parseFile(c *gin.Context) error {
	contentType := c.Request.Header.Get("Content-Type")
	_, params, _ := mime.ParseMediaType(contentType)
	boundary, ok := params["boundary"]
	if !ok {
		log.Error(zap.Error(errors.New("no multipart boundary param in Content-Type")))
		return errors.New("no multipart boundary param in Content-Type")
	}

	// 准备重写数据
	bodyBuf := &bytes.Buffer{}
	wr := multipart.NewWriter(bodyBuf)
	mr := multipart.NewReader(c.Request.Body, boundary)
	for {
		p, err := mr.NextPart() // p的类型为Part
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error(zap.Error(err))
			break
		}

		fileByte, err := ioutil.ReadAll(p)
		if err != nil {
			log.Error(zap.Error(err))
			break
		}

		pName := p.FormName()
		fileName := p.FileName()
		if len(fileName) < 1 {
			if pName == "param" {
				formData, err := decryptString(string(fileByte))
				if err != nil {
					log.Error(zap.Error(err))
					break
				}

				for k, v := range formData {
					val := valuesStr(v)
					err = wr.WriteField(k, val)
					if err != nil {
						log.Error(zap.Error(err))
						break
					}
				}
			} else {
				_ = wr.WriteField(pName, string(fileByte))
			}
		} else {
			tmp, err := wr.CreateFormFile(pName, fileName)
			if err != nil {
				log.Error(zap.Error(err))
				continue
			}
			_, _ = tmp.Write(fileByte)
		}
	}

	// 写结尾标志
	_ = wr.Close()
	c.Request.Header.Set("Content-Type", wr.FormDataContentType())
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBuf.Bytes()))

	return nil
}

func parseForm(c *gin.Context) error {
	// 读取数据 body处理
	payload, err := c.GetRawData()
	if err != nil {
		log.Error(zap.Error(err))
		return err
	}

	// /解密body数据 请求的json是"encryptString= value含有gcm的12字节nonce,实际长度大于32
	if payload != nil && len(payload) > 20 {

		values, err := url.ParseQuery(string(payload))
		if err != nil {
			log.Error(zap.Error(err))
			return err
		}

		payloadText := values.Get("param")
		if len(payloadText) > 0 {
			mapData, err := decryptString(payloadText)
			if err != nil {
				log.Error(zap.Error(err))
				return err
			}
			for k, v := range mapData {
				values.Add(k, valuesStr(v))
			}

			formData := values.Encode()
			payload = []byte(formData)
		}
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	return nil
}

func parseJson(c *gin.Context) error {
	// 读取数据 body处理
	payload, err := c.GetRawData()
	if err != nil {
		log.Error(zap.Error(err))
		return err
	}
	// 解密body数据 请求的json是{"param":{value}}
	if payload != nil && len(payload) > 20 {
		var jsonData protocol.EncryptJson

		if err = json.Unmarshal(payload, &jsonData); err != nil {
			log.Error(zap.Error(err))
			return err
		}
		payloadText := jsonData.Param
		if len(payloadText) > 0 {
			payload, err = utils.AesDecrypt([]byte(payloadText))
			if err != nil {
				log.Error(zap.Error(err))
				return err
			}
		}
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(payload))

	return nil
}

// 处理get url的解密
func parseQuery(c *gin.Context) error {
	param := c.Query("param")
	if len(param) < 1 {
		log.Debug(zap.Any("uri", c.Request.URL.Path))
		return nil
	}

	queryData, err := decryptString(param)
	if err != nil {
		log.Error(zap.Error(err))
		return err
	}

	var args []string
	var logs []string
	for k, v := range queryData {
		val := valuesStr(v)
		args = append(args, fmt.Sprintf("%s=%s", k, url.QueryEscape(val)))
		logs = append(logs, fmt.Sprintf("%s=%s", k, val))
	}
	c.Request.URL.RawQuery = strings.Join(args, "&")
	return nil
}

func decryptString(param string) (map[string]interface{}, error) {
	formData := make(map[string]interface{}, 0)
	if len(param) < 1 {
		log.Error(zap.Error(errors.New("len param error")))
		return formData, nil
	}
	plaintext, err := utils.AesDecrypt([]byte(param))
	if err != nil {
		log.Error(zap.Error(err))
		return formData, err
	}

	if len(plaintext) < 3 {
		log.Error(zap.Error(errors.New("len(plaintext) < 3")))
		return formData, nil
	}

	err = json.Unmarshal(plaintext, &formData)
	if err != nil {
		log.Error(zap.Error(err))
		return formData, err
	}
	return formData, nil
}
func valuesStr(v interface{}) string {
	val := ""
	switch v.(type) {
	case float64:
		// val, _ := decimal.NewFromString(fmt.Sprintf("%.10f", v))
		fl, _ := v.(float64)
		val = strconv.FormatFloat(fl, 'f', -1, 64)
	default:
		val = fmt.Sprintf("%v", v)
	}
	return val
}
func isJsonResponse(c *gin.Context) bool {
	contentType := c.Writer.Header().Get("Content-Type")
	return strings.Contains(contentType, constant.Types[constant.TypeJSON])
}
