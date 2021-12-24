/*
 * @Author: Administrator
 * @IDE: GoLand
 * @Date: 2021-12-24 13:58
 * @LastEditors: Administrator
 * @LastEditTime: 2021-12-24 13:58
 * @FilePath: ql-mp/utils/hashid.go
 */

package utils

import (
	"github.com/speps/go-hashids/v2"
)

var hd *hashids.HashIDData

const salt = "__ql-mini_salt_magic_"

func init() {
	hd = hashids.NewData()
	hd.Salt = salt
	hd.MinLength = 8
}

func HashIDEncode(data int64) (string, error) {
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}
	return h.EncodeInt64([]int64{data})
}

func HashIDDecode(data string) (int64, error) {
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return 0, err
	}

	e, err := h.DecodeInt64WithError(data)
	if err != nil {
		return 0, err
	}

	return e[0], nil
}
