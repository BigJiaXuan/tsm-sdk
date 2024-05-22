package tsm_sdk

import (
	"crypto"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/golang-module/dongle"
	"sort"
	"strings"
)

func newCipher(tDes string) *dongle.Cipher {
	cipher := dongle.NewCipher()
	cipher.SetKey(dongle.Decode.FromString(tDes).ByBase64().ToString())
	cipher.SetMode(dongle.CBC)
	cipher.SetIV(make([]byte, des.BlockSize))
	cipher.SetPadding(dongle.PKCS5)
	return cipher
}

// EncryptRequestParam 加密 request参数
// tDes  3des.txt 中数据
// request 待加密的数据
// Returns 加密后的数据
func (c *Client) EncryptRequestParam(request string) string {
	return dongle.Encrypt.FromString(request).By3Des(c.cipher).ToBase64String()
}

//	获取sign参数
//
// Returns 加密后的数据
func (c *Client) EncryptSignParam(access_token, app_key, method, format, signMethod, timestamp, v, request string) (string, error) {
	params := map[string]string{
		"access_token": access_token,
		"app_key":      app_key,
		"method":       method,
		"timestamp":    timestamp,
		"v":            v,
		"request":      request,
		"format":       format,
		"sign_method":  signMethod,
	}
	// 构造一个切片
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	// 对key进行排序
	sort.Strings(keys)
	var builder strings.Builder
	for _, key := range keys {
		// 拼接成access_tokenxxxapp_keyxxxmethodxxxrequestxxxtimestampxxxvxxx
		builder.WriteString(key)
		builder.WriteString(params[key])
	}
	// 获取最终地签名字符串
	sign := builder.String()

	// 使用rsa-sha1 进行私钥签名
	data := c.signatureSign(sign, c.SvrPkcs8)
	//fmt.Println("签名后的sign", data)
	return data, nil
}

func (c *Client) signatureSign(text, key string) string {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return ""
	}
	private, err := x509.ParsePKCS8PrivateKey(block.Bytes) //之前看java demo中使用的是pkcs8
	if err != nil {
		return ""
	}
	h := crypto.Hash.New(crypto.SHA1) //进行SHA1的散列
	h.Write([]byte(text))
	hashed := h.Sum(nil)
	// 进行rsa加密签名
	signedData, err := rsa.SignPKCS1v15(rand.Reader, private.(*rsa.PrivateKey), crypto.SHA1, hashed)
	data := base64.StdEncoding.EncodeToString(signedData)
	return data
}

// DecryptRequestParam 解析request参数
func (c *Client) DecryptRequestParam(request string) string {
	return dongle.Decrypt.FromBase64String(request).By3Des(c.cipher).ToString()
}
