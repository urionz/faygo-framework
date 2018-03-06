package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	//BASE64字符表,不要有重复
	base64Table        = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	hashFunctionHeader = "zh.ife.iya"
	hashFunctionFooter = "09.O25.O20.78"
)

// CopyBody 解决body内容只能读取一次的问题，临时保存再赋值
var CopyBody = func(r *http.Request) (content []byte, err error) {
	content, err = ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(content))
	return
}

var Md5String = func(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

var SHA1String = func(s string) string {
	t := sha1.New()
	t.Write([]byte(s))
	return hex.EncodeToString(t.Sum(nil))
}

var coder = base64.NewEncoding(base64Table)

var Base64Encode = func(str string) string {
	var src []byte = []byte(hashFunctionHeader + str + hashFunctionFooter)
	return string([]byte(coder.EncodeToString(src)))
}

var Base64Decode = func(str string) (string, error) {
	var src []byte = []byte(str)
	by, err := coder.DecodeString(string(src))
	return strings.Replace(strings.Replace(string(by), hashFunctionHeader, "", -1), hashFunctionFooter, "", -1), err
}
