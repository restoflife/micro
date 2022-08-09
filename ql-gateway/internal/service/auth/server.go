/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-16 17:40
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-16 17:40
 * @FilePath: ql-gateway/internal/service/auth/server.go
 */

package auth

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/restoflife/micro/gateway/internal/component/log"
	"github.com/restoflife/micro/gateway/internal/encoding"
	"github.com/restoflife/micro/gateway/internal/errutil"
	"github.com/restoflife/micro/gateway/internal/protocol"
	"github.com/restoflife/micro/gateway/utils"
	"go.uber.org/zap"
	"mime/multipart"
	"path"
	"time"
)

func MakeCaptchaHandler(c *gin.Context) {
	id, b64s, err := utils.DriverDigitFunc()
	if err != nil {
		log.Error(zap.Error(err))
		return
	}
	resp := &protocol.CaptchaResp{
		Id:  id,
		Url: b64s,
	}
	encoding.Ok(c, resp)

}

func MakeRegisterHandler(c *gin.Context) {
	req := &protocol.RegisterReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		encoding.Error(c, errutil.ErrIllegalParameter)
		return
	}
	req.Ip = c.ClientIP()
	uid, _ := utils.GetUUID()
	req.UID = uid
	if err := makeRegisterService(req); err != nil {
		encoding.Error(c, err)
		return
	}
	encoding.Ok(c, ``)

}
func MakeUserListHandler(c *gin.Context) {
	encoding.Error(c, errutil.ErrInternalServer)
	return
}

func MakeUploadHandler(c *gin.Context) {
	filetType := c.PostForm("type")
	file, err := c.FormFile("upload_file")
	if err != nil {
		encoding.Error(c, errutil.ErrIllegalParameter)
		return
	}
	f, err := file.Open()
	if err != nil {
		encoding.Error(c, err)
		return
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)
	switch filetType {
	case "1":
		// 图片
		filetType = "image"
	case "2":
		// 视频
		filetType = "video"
	case "3":
		// 音频
		filetType = "audio"
	default:
		// 文件
		filetType = "file"
	}
	var b = make([]byte, file.Size)
	if _, err = f.Read(b); err != nil {
		encoding.Error(c, err)
		return
	}
	fileExt := path.Ext(file.Filename)
	var key = fmt.Sprintf("%s/%s%d", filetType, utils.MD5String(file.Filename), time.Now().Unix()) + hex.EncodeToString(utils.MD5(b))
	key += fileExt
	result, err := utils.UploadImageByQiNiu(bytes.NewReader(b), key, file.Size)
	if err != nil {
		encoding.Error(c, err)
		return
	}
	encoding.Ok(c, result)
}
