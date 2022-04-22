/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2022-04-22 16:16
 * @LastEditors: Administrator
 * @LastEditTime: 2022-04-22 16:16
 * @FilePath: ql-gateway/internal/constant/http.go
 */

package constant

type RequestType string

const (
	TypeJSON              RequestType = "json"
	TypeXML               RequestType = "xml"
	TypeUrlencoded        RequestType = "urlencoded"
	TypeForm              RequestType = "form"
	TypeFormData          RequestType = "form-data"
	TypeMultipartFormData RequestType = "multipart-form-data"
)

var Types = map[RequestType]string{
	TypeJSON:              "application/json",
	TypeXML:               "application/xml",
	TypeUrlencoded:        "application/x-www-form-urlencoded",
	TypeForm:              "application/x-www-form-urlencoded",
	TypeFormData:          "application/x-www-form-urlencoded",
	TypeMultipartFormData: "multipart/form-data",
}
