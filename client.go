package tsm_sdk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/BigJiaXuan/tsm-sdk/models"
	"github.com/golang-module/dongle"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client interface {
	DoRequest(ctx context.Context, method, accessToken, request string) (response string, err error)
	// GetAccessToken 获取访问令牌
	GetAccessToken(ctx context.Context) (string, error)
	// QueryCard 获取身份标识对应的所有卡信息
	QueryCard(ctx context.Context, token string, idType string, id string, pwd string) ([]models.Card, error)
	// QueryAccInfo 查询电子账户信息
	QueryAccInfo(ctx context.Context, token string, account string, accType string) ([]models.AccInfo, error)
	// QueryAccInfoTotal 电子账户汇总
	QueryAccInfoTotal(ctx context.Context, token string, accType string) (models.QueryAccinfoTotal, error)
	// QuerySnoNameFromBarcode 通过二维码查询工号信息
	QuerySnoNameFromBarcode(ctx context.Context, token string, barcode string) (models.QuerySnoNameFromBarcode, error)
	// CheckPassword 查询密码校验
	CheckPassword(ctx context.Context, token string, idType, id, pwd string) (models.CheckPassword, error)
	// ModifyPassword 修改校园卡查询密码
	ModifyPassword(ctx context.Context, token string, account, newPassword string) (bool, error)
	// LostCard 校园卡挂失
	LostCard(ctx context.Context, token string, account string) (bool, error)
	// UnLostCard 校园卡解挂
	UnLostCard(ctx context.Context, token string, account string) (bool, error)
	// FreezeCard 校园卡冻结
	FreezeCard(ctx context.Context, token string, account string) (bool, error)
	// UnFreezeCard 校园卡解冻
	UnFreezeCard(ctx context.Context, token string, account string) (bool, error)
	// ModifyPayLimit 电子账户支付限额设置
	ModifyPayLimit(ctx context.Context, token string, account string, accType string, singleLimit int64, dayCostLimit int64, nonPwdLimit int64) (bool, error)
	// OpenAcc 开户(虚拟卡开户)
	OpenAcc(ctx context.Context, token string, sno, name, idNo, schoolCode, deptCode, cardType, pidCode, inDate, expDate, photoImage, phone string) (account int64, err error)
	// GetBarCode 获取二维码
	GetBarCode(ctx context.Context, token string, account, payType, payAcc string) (models.GetBarCode, error)
	// ModifyAccInfo 校园卡信息变更
	ModifyAccInfo(ctx context.Context, token string, account, accType, expire, autoTranLimit, autoTranAmt, autoTranFlag string) (models.ModifyAccInfo, error)
}

type client struct {
	config *Config
	cipher *dongle.Cipher
}

type Config struct {
	URL      string
	TDes     string
	AppKey   string
	SvrPkcs8 string
	//SvrPrivate string
	//TsmPublic string
}

func NewClient(config *Config) Client {
	return &client{
		config: config,
		cipher: newCipher(config.TDes),
	}
}

func (c *client) DoRequest(ctx context.Context, method, accessToken, request string) (response string, err error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	v := "2.0"
	signMethod := "rsa"
	format := "json"
	// 加密request参数
	requested := c.encryptRequestParam(request)

	sign, err := c.encryptSignParam(
		accessToken,
		c.config.AppKey,
		method,
		format,
		signMethod,
		timestamp,
		v,
		requested)
	if err != nil {
		return "", fmt.Errorf("加密sign参数失败：%v", err)
	}

	formData := url.Values{}
	formData.Set("method", method)
	formData.Set("timestamp", timestamp)
	formData.Set("format", "json")
	formData.Set("app_key", c.config.AppKey)
	formData.Set("access_token", accessToken)
	formData.Set("v", v)
	formData.Set("sign", sign)
	formData.Set("sign_method", signMethod)
	formData.Set("request", requested)
	// 发送请求
	req, err := http.NewRequestWithContext(ctx, "POST", c.config.URL, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("构造http请求失败：%v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送http请求失败：%v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", fmt.Errorf("读取请求的响应失败：%v", err)
	}

	// 先解码 URL编码的字符串
	data, err := url.PathUnescape(string(body))
	if err != nil {
		//c.logger.Error("TSM URL编码解码失败", zap.Error(err))
		return "", fmt.Errorf("URL 解码失败：%v", err)
	}

	// 可能会出现直接返回错误信息的情况
	if !strings.Contains(data, "&") {
		// 可能是直接返回的错误信息
		// 返回的response gbk转utf8
		toUtf8, err := gbkToUtf8([]byte(data))
		if err != nil {
			return "", fmt.Errorf("GBK转UTF8失败：%v", err)
		}
		//c.logger.Error("TSM 返回的不是标准信息，直接返回内容", zap.Error(fmt.Errorf(string(toUtf8))))
		return "", errors.New(string(toUtf8))
	}
	//c.logger.Debug("TSM 返回的数据 解码前", zap.String("data", data))
	// 解析response
	parseResponse, err := c.parseResponse(data)
	if err != nil {
		return "", fmt.Errorf("解密response失败：%v", err)
	}

	// 判断tsm的状态码
	if parseResponse.ErrCode != 0 {
		// 获取错误信息
		errMsg := GetErrorMessage(parseResponse.ErrCode)
		//c.logger.Error("TSM 处理请求失败", zap.Error(errors.New(errMsg)))
		return parseResponse.Request, fmt.Errorf("TSM 返回错误：%s", errMsg)
	}
	// 解密request

	decryptedResp := c.decryptRequestParam(parseResponse.Request)
	//c.logger.Info("TSM 返回的数据 解码后", zap.Any("ErrCode", parseResponse.ErrCode), zap.Any("Request", decryptedResp))

	//c.logger.Info("TSM 请求耗时", zap.String("耗时", fmt.Sprintf("%d毫秒", useTime)))
	return decryptedResp, nil
}

type Params struct {
	ErrCode int64
	Request string
	Sign    string
}

// 将tsm返回的response进行分割
func (c *client) parseResponse(query string) (Params, error) {
	// errcode=0&request=ZbsyvxKuG5LKYaA5hAhZHsX2pHgo%2B4SB7RsmonVdDE5POlmpFnySfhf4RXeMubeOYtT8Hjh3YlUekh13lYX%2F8z5651j2xEOLRB0xobyTgtXCxOoHiFvk5OIFDwdFH8uTX0MYZeroT6nHzf5mey0M1gQIeUBFzg%2Bznk%2Btb1WKDfWuaNo4ipNLMOsXMbADUE2kpYF61%2BpWx6C9450UIe2NK3KNsUz6NE7yVSgzgIxmJ5wAukZYl%2Fe6pSep%2BqX%2B%2Fav4sJkLnVOySJrxSUxlJD%2Frppcrq9PbEVBNjCxf5%2B2v3CCTKyLxAlCkeTI95mg4qT03&sign=dItDagbZxlMaCh%2Byp5vXWAWnu4luIgfAdYP5f8hKZnbzdeDpwVU57FU6DC4edAVCTZAynDHax6XKy42Il5poz63Sn7poPJz2PVXGMVYXj2FaOtx3jn4XCKU6mLUCe7YCryJ9442Pf%2B20jBSrMaZeC6wB1aFFxrKbhYIPqKUWERmA6P4hKN0cRvuuU5YZV3Vyvp6F%2BxRp02lLEQSgca8UR7HC6bcwsz%2FDgkPO5I9YpO47LA8x7inPH2Mycniozc1Mfvj1pHbEz8f2citp7fBBdzlRD1KTIWxPSkaFiZ0gvqJrKuf%2FrCGmD1kbvI5Gf9FPYz3e59KOk3B5wNNKFZdjyQ%3D%3D
	// 先按照 & 分割
	var params Params
	split := strings.Split(query, "&")
	for _, v := range split {
		// 再按照 = 分割
		split := strings.SplitN(v, "=", 2)
		if len(split) == 2 {
			switch split[0] {
			case "errcode":
				code, err := strconv.Atoi(split[1])
				if err != nil {
					return Params{}, err
				}
				params.ErrCode = int64(code)
			case "request":
				params.Request = split[1]
			case "sign":
				params.Sign = split[1]
			}
		}
	}
	return params, nil
}

// gbkToUtf8 将 GBK 编码的字节数组转换为 UTF-8 编码的字节数组
func gbkToUtf8(gbkData []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(gbkData), simplifiedchinese.GBK.NewDecoder())
	return io.ReadAll(reader)
}
