package tsm_sdk

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/BigJiaXuan/tsm-sdk/models"
	idvalidator "github.com/guanguans/id-validator"
	"strconv"
)

// GetAccessToken
//
//	@Description:获取访问令牌
//	@receiver c
//	@param ctx
//	@return string token
//	@return error
func (c *Client) GetAccessToken(ctx context.Context) (string, error) {
	accessToken := "0000000000000000000000000000000000000000000000000000000000000000" +
		"0000000000000000000000000000000000000000000000000000000000000000"
	method := "synjones.authorize.access_token"
	request := "{\"authorize_access_token\": {}}"
	resp, err := c.DoRequest(ctx, method, accessToken, request)
	if err != nil {
		return "", err
	}
	type T struct {
		AuthorizeAccessToken struct {
			Retcode     string `json:"retcode"`
			Errmsg      string `json:"errmsg"`
			AccessToken string `json:"access_token"`
			ExpiresIn   string `json:"expires_in"`
		} `json:"authorize_access_token"`
	}
	var t T
	err = json.Unmarshal([]byte(resp), &t)
	// 判断 request解码携带的retcode是否为0 成功
	if t.AuthorizeAccessToken.Retcode != "0" {
		return "", errors.New(t.AuthorizeAccessToken.Errmsg)
	}
	return t.AuthorizeAccessToken.AccessToken, err
}

// QueryCard
//
//		@Description: 获取身份标识对应的所有卡信息 默认都获取
//		@receiver c
//		@param ctx
//		@param token token
//		@param idType 学工号：sno 证件号：id 手机号：phone 账号：acc 卡号：cardid 二维码：barcode
//		@param id 学工号、证件号、手机号、帐号、卡号、二维码
//		@param getMerge 是否获取合并账户标志位 0:不获取 1:获取
//		@param getIdInfo 是否获取身份信息 0:不获取 1:获取
//		@param pwd 查询密码 传递则校验密码 0:不获取 1:获取
//		@param getAccInfo 是否获取电子账户列表 0:不获取 1:获取
//		@param getSchCode 是否获取核算单位 0:不获取 1:获取
//		@param getCardId 是否获取物理卡号 0:不获取 1:获取
//	    @return []models.Card
func (c *Client) QueryCard(ctx context.Context, token string,
	idType string, id string, pwd string) ([]models.Card, error) {
	accessToken := token
	method := "synjones.onecard.query.card"
	type QueryCard struct {
		IdType     string `json:"idtype"`
		Id         string `json:"id"`
		GetMerge   string `json:"get_merge"`
		GetIdInfo  string `json:"get_idinfo"`
		Pwd        string `json:"pwd,omitempty"`
		GetAccInfo string `json:"get_accinfo"`
		GetSchCode string `json:"getschcode"`
		GetCardId  string `json:"getcardid"`
	}
	type R struct {
		QueryCard `json:"query_card"`
	}

	request := R{QueryCard: QueryCard{
		IdType:     idType,
		Id:         id,
		GetMerge:   "1",
		GetIdInfo:  "1",
		Pwd:        pwd,
		GetAccInfo: "1",
		GetSchCode: "1",
		GetCardId:  "1",
	}}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	// 发送请求
	resp, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return nil, err
	}

	var t models.QueryCardResponse
	err = json.Unmarshal([]byte(resp), &t)
	if t.QueryCard.Retcode != "0" {
		return nil, errors.New(t.QueryCard.Errmsg)
	}
	return t.QueryCard.Card, nil
}

// QueryAccInfo
//
//	@Description: 查询电子账户信息
//	@receiver c
//	@param ctx
//	@param token
//	@param account 一卡通账号
//	@param accType 电子账户类型 ###:为卡账户 其他为电子账户类型 为空时返回所有电子账户
//	@return []models.AccInfo
//	@return error
func (c *Client) QueryAccInfo(ctx context.Context, token string,
	account string, accType string) ([]models.AccInfo, error) {
	accessToken := token
	method := "synjones.onecard.query.accinfo"
	type QueryAccInfo struct {
		Account string `json:"account"`
		AccType string `json:"acctype,omitempty"`
	}
	type R struct {
		QueryAccInfo `json:"query_accinfo"`
	}
	request := R{QueryAccInfo: QueryAccInfo{
		Account: account,
		AccType: accType,
	}}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	// 发送请求
	response, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return nil, err
	}
	var t models.QueryAccInfoResponse
	err = json.Unmarshal([]byte(response), &t)
	if err != nil {
		return nil, err
	}
	if t.QueryAccinfo.Retcode != "0" {
		return nil, errors.New(t.QueryAccinfo.Errmsg)
	}
	return t.QueryAccinfo.Accinfo, nil
}

// QueryAccInfoTotal
//
//	@Description: 电子账户汇总
//	@receiver c
//	@param ctx
//	@param token
//	@param accType 电子账户类型
//	@return models.QueryAccinfoTotal
//	@return error
func (c *Client) QueryAccInfoTotal(ctx context.Context, token string,
	accType string) (models.QueryAccinfoTotal, error) {
	accessToken := token
	method := "synjones.onecard.query.accinfo.total"
	type QueryAccInfoTotal struct {
		AccType string `json:"acctype"`
	}
	type R struct {
		QueryAccInfoTotal `json:"query_accinfo_total"`
	}
	r := R{QueryAccInfoTotal{AccType: accType}}
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return models.QueryAccinfoTotal{}, err
	}
	response, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return models.QueryAccinfoTotal{}, err
	}
	var t models.QueryAccInfoTotalResponse
	err = json.Unmarshal([]byte(response), &t)
	if err != nil {
		return models.QueryAccinfoTotal{}, err
	}
	if t.Retcode != "0" {
		return models.QueryAccinfoTotal{}, errors.New(t.Errmsg)
	}
	return t.QueryAccinfoTotal, nil
}

// QuerySnoNameFromBarcode
//
//	@Description: 通过二维码查询工号信息
//	@receiver c
//	@param ctx
//	@param token
//	@param barcode 二维码
//	@return models.QuerySnoNameFromBarcode
//	@return error
func (c *Client) QuerySnoNameFromBarcode(ctx context.Context, token string,
	barcode string) (models.QuerySnoNameFromBarcode, error) {
	accessToken := token
	method := "synjones.onecard.query.snoname.frombarcode"
	type QuerySnoNameFromBarcode struct {
		Barcode string `json:"barcode"`
	}
	type R struct {
		QuerySnoNameFromBarcode `json:"query_snoname_frombarcode"`
	}
	r := R{QuerySnoNameFromBarcode{Barcode: barcode}}
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return models.QuerySnoNameFromBarcode{}, err
	}
	response, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return models.QuerySnoNameFromBarcode{}, err
	}
	var res models.QuerySnoNameFromBarcodeResponse
	err = json.Unmarshal([]byte(response), &res)
	if err != nil {
		return models.QuerySnoNameFromBarcode{}, err
	}
	if res.QuerySnoNameFromBarcode.Retcode != "0" {
		return models.QuerySnoNameFromBarcode{}, errors.New(res.QuerySnoNameFromBarcode.Errmsg)
	}
	return res.QuerySnoNameFromBarcode, nil
}

// CheckPassword
//
//	@Description: 查询密码校验
//	@receiver c
//	@param ctx
//	@param token
//	@param idType 身份标识类型 学工号:sno 证件号:id 手机号:phone 帐号:acc 卡号:cardid
//	@param id 身份标识  学工号、证件号、手机号、帐号或卡号
//	@param pwd 校园卡查询密码
//	@return models.CheckPassword
//	@return error
func (c *Client) CheckPassword(ctx context.Context, token string,
	idType, id, pwd string) (models.CheckPassword, error) {
	accessToken := token
	method := "synjones.onecard.check.pwd"
	type CheckPassword struct {
		IdType string `json:"idtype"`
		Id     string `json:"id"`
		Pwd    string `json:"pwd"`
	}
	type R struct {
		CheckPassword `json:"check_pwd"`
	}
	request := R{CheckPassword: CheckPassword{idType, id, pwd}}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return models.CheckPassword{}, err
	}
	response, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return models.CheckPassword{}, err
	}
	var resp models.CheckPasswordResponse
	err = json.Unmarshal([]byte(response), &resp)
	if err != nil {
		return models.CheckPassword{}, err
	}
	if resp.CheckPassword.Retcode != "0" {
		return models.CheckPassword{}, errors.New(resp.CheckPassword.Errmsg)
	}
	return resp.CheckPassword, nil
}

// ModifyPassword
//
//	@Description: 修改校园卡查询密码
//	@receiver c
//	@param ctx
//	@param token
//	@param account 一卡通账户
//	@param newPassword 校园卡查询密码 8位数字+字母
//	@return bool
//	@return error
func (c *Client) ModifyPassword(ctx context.Context, token string,
	account, newPassword string) (bool, error) {
	accessToken := token
	method := "synjones.onecard.modify.pwd"
	type ModifyPassword struct {
		Account     string `json:"account"`
		NewPassword string `json:"newpwd"`
	}
	type R struct {
		ModifyPassword `json:"modify_pwd"`
	}
	r := R{ModifyPassword{account, newPassword}}
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return false, err
	}
	response, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return false, err
	}
	var res models.ModifyPasswordResponse
	err = json.Unmarshal([]byte(response), &res)
	if err != nil {
		return false, err
	}
	if res.ModifyPwd.Retcode != "0" {
		return false, errors.New(res.ModifyPwd.Errmsg)
	}
	return true, nil
}

// LostCard
//
//	@Description: 校园卡挂失
//	@receiver c
//	@param ctx
//	@param token
//	@param account 一卡通账号
//	@return bool
//	@return error
func (c *Client) LostCard(ctx context.Context, token string,
	account string) (bool, error) {
	accessToken := token
	method := "synjones.onecard.lost.card"
	type LostCard struct {
		Account string `json:"account"`
	}
	type R struct {
		LostCard `json:"lost_card"`
	}
	r := R{LostCard{Account: account}}
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return false, err
	}
	response, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return false, err
	}
	var res models.LostCardResponse
	err = json.Unmarshal([]byte(response), &res)
	if err != nil {
		return false, err
	}
	if res.LostCard.Retcode != "0" {
		return false, errors.New(res.LostCard.Errmsg)
	}
	return true, nil
}

// UnLostCard
//
//	@Description: 校园卡解挂
//	@receiver c
//	@param ctx
//	@param token
//	@param account 一卡通账号
//	@return bool
//	@return error
func (c *Client) UnLostCard(ctx context.Context, token string,
	account string) (bool, error) {
	accessToken := token
	method := "synjones.onecard.unlost.card"
	type UnLostCard struct {
		Account string `json:"account"`
	}
	type R struct {
		UnLostCard `json:"unlost_card"`
	}
	r := R{UnLostCard{Account: account}}
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return false, err
	}
	response, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return false, err
	}
	var res models.UnLostCardResponse
	err = json.Unmarshal([]byte(response), &res)
	if err != nil {
		return false, err
	}
	if res.UnLostCard.Retcode != "0" {
		return false, errors.New(res.UnLostCard.Errmsg)
	}
	return true, nil
}

// FreezeCard
//
//	@Description: 校园卡冻结
//	@receiver c
//	@param ctx
//	@param token
//	@param account 一卡通账号
//	@return bool
//	@return error
func (c *Client) FreezeCard(ctx context.Context, token string,
	account string) (bool, error) {
	accessToken := token
	method := "synjones.onecard.frozen.card"
	type FrozenCard struct {
		Account string `json:"account"`
	}
	type R struct {
		FrozenCard `json:"frozen_card"`
	}
	r := R{FrozenCard{Account: account}}
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return false, err
	}
	response, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return false, err
	}
	var res models.FreezeCardResponse
	err = json.Unmarshal([]byte(response), &res)
	if err != nil {
		return false, err
	}
	if res.FrozenCard.Retcode != "0" {
		return false, errors.New(res.FrozenCard.Errmsg)
	}
	return true, nil
}

// UnFreezeCard
//
//	@Description: 校园卡解冻
//	@receiver c
//	@param ctx
//	@param token
//	@param account 一卡通账号
//	@return bool
//	@return error
func (c *Client) UnFreezeCard(ctx context.Context, token string,
	account string) (bool, error) {
	accessToken := token
	method := "synjones.onecard.unfrozen.card"
	type UnFrozenCard struct {
		Account string `json:"account"`
	}
	type R struct {
		UnFrozenCard `json:"unfrozen_card"`
	}
	r := R{UnFrozenCard{Account: account}}
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return false, err
	}
	response, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return false, err
	}
	var res models.UnFreezeCardResponse
	err = json.Unmarshal([]byte(response), &res)
	if err != nil {
		return false, err
	}
	if res.UnFrozenCard.Retcode != "0" {
		return false, errors.New(res.UnFrozenCard.Errmsg)
	}
	return true, nil
}

// ModifyPayLimit
//
//	@Description: 电子账户支付限额设置
//	@receiver c
//	@param ctx
//	@param token
//	@param account 一卡通账户
//	@param accType 账户类型
//	@param singleLimit 单次支付限额(分)
//	@param dayCostLimit 当天支付限额(分)
//	@param nonPwdLimit 免密金额(分)
//	@return bool
//	@return error
func (c *Client) ModifyPayLimit(ctx context.Context, token string,
	account string, accType string, singleLimit int64, dayCostLimit int64, nonPwdLimit int64) (bool, error) {
	accessToken := token
	method := " synjones.onecard.paylimite.modify"
	type PayLimit struct {
		Account      string `json:"account"`
		AccType      string `json:"acctype"`
		SingleLimit  int64  `json:"singlelimit"`
		DayCostLimit int64  `json:"daycostlimit"`
		NonPwdLimit  int64  `json:"nonpwdlimit"`
	}
	type R struct {
		PayLimit `json:"paylimit_modify"`
	}
	r := R{PayLimit{
		Account:      account,
		AccType:      accType,
		SingleLimit:  singleLimit,
		DayCostLimit: dayCostLimit,
		NonPwdLimit:  nonPwdLimit,
	}}
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return false, err
	}
	response, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return false, err
	}
	var resp models.ModifyPayLimit
	err = json.Unmarshal([]byte(response), &resp)
	if err != nil {
		return false, err
	}
	if resp.PaylimiteModify.Retcode != "0" {
		return false, errors.New(resp.PaylimiteModify.Errmsg)
	}
	return true, nil
}

// OpenAcc
//
//	@Description: 开户(虚拟卡开户)
//	@receiver c
//	@param ctx
//	@param token
//	@param sno 学工号
//	@param name 姓名
//	@param idNo 身份证号
//	@param schoolCode  校区代码
//	@param deptCode 部门代码
//	@param cardType 卡类型 800:正式卡，801:临时卡
//	@param pidCode 身份代码
//	@param inDate 入校日期
//	@param expDate 失效日期
//	@param photoImage 照片转换成 base64字符串
//	@param phone 联系电话
//	@return account
//	@return err
func (c *Client) OpenAcc(ctx context.Context, token string,
	sno, name, idNo, schoolCode, deptCode, cardType, pidCode, inDate, expDate, photoImage, phone string) (account int64, err error) {
	accessToken := token
	method := "synjones.onecard.open.acc"
	type OpenAcc struct {
		Sno        string `json:"sno"`
		Name       string `json:"name"`
		Sex        string `json:"sex"`
		IDNo       string `json:"idno"`
		Phone      string `json:"phone"`
		Email      string `json:"email"`
		SchoolCode string `json:"schoolcode"`
		DeptCode   string `json:"depcode"`
		CardType   string `json:"cardtype"`
		Born       string `json:"born"`
		PidCode    string `json:"pidcode"`
		InDate     string `json:"indate"`
		ExpDate    string `json:"expdate"`
		PhotoImage string `json:"photo_image"`
	}
	type R struct {
		OpenAcc `json:"open_acc"`
	}
	var sex string
	var born string
	// 校验身份证号 严格模式校验大陆身份证号
	if idvalidator.IsValid(idNo, true) {
		// 校验通过
		// 获取身份证信息
		info, err := idvalidator.GetInfo(idNo, true)
		if err != nil {
			return 0, err
		}
		// 判断性别
		switch info.Sex {
		case 1:
			sex = "1"
		case 0:
			sex = "2"
		}
		born = info.Birthday.Format("20060102")
	} else {
		return 0, errors.New("身份证不合法")
	}
	r := R{OpenAcc{
		Sno:        sno,
		Name:       name,
		Sex:        sex,
		IDNo:       idNo,
		Phone:      phone,
		Email:      "",
		SchoolCode: schoolCode,
		DeptCode:   deptCode,
		CardType:   cardType,
		Born:       born,
		PidCode:    pidCode,
		InDate:     inDate,
		ExpDate:    expDate,
		PhotoImage: photoImage,
	}}
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return 0, err
	}
	response, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return 0, err
	}
	var resp models.OpenAccResponse
	err = json.Unmarshal([]byte(response), &resp)
	if err != nil {
		return 0, err
	}
	if resp.OpenAcc.Retcode != "0" {
		return 0, errors.New(resp.OpenAcc.Errmsg)
	}
	acc, err := strconv.Atoi(resp.OpenAcc.Account)
	if err != nil {
		return 0, err
	}
	return int64(acc), nil

}

// GetBarCode
//
//	@Description: 获取二维码
//	@receiver c
//	@param ctx
//	@param token
//	@param account 一卡通账户
//	@param payType 支付方式 1、 校园卡账户支付 2、 绑定银行卡支付 3、 自定义银行卡支付
//	@param payAcc 支付账号 paytype 为1时此字段值可为: ###:为卡账户 其他为电子账户类型 paytype 为 2 时此字段 值可为: 空 paytype为3时此字段 值可为: 银行卡号
//	@return models.GetBarCode
//	@return error
func (c *Client) GetBarCode(ctx context.Context, token string,
	account, payType, payAcc string) (models.GetBarCode, error) {
	accessToken := token
	method := "synjones.onecard.barcode.get"
	type BarCode struct {
		Account string `json:"account"`
		PayType string `json:"paytype"`
		PayAcc  string `json:"payacc"`
	}
	type R struct {
		BarCode `json:"barcode_get"`
	}
	r := R{BarCode{
		Account: account,
		PayType: payType,
		PayAcc:  payAcc,
	}}
	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return models.GetBarCode{}, err
	}
	response, err := c.DoRequest(ctx, method, accessToken, string(jsonRequest))
	if err != nil {
		return models.GetBarCode{}, err
	}
	var resp models.GerBarCodeResponse
	err = json.Unmarshal([]byte(response), &resp)
	if err != nil {
		return models.GetBarCode{}, err
	}
	if resp.GetBarCode.Retcode != "0" {
		return models.GetBarCode{}, errors.New(resp.GetBarCode.Errmsg)
	}
	return resp.GetBarCode, nil
}
