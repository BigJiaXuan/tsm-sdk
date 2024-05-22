package models

type ModifyPasswordResponse struct {
	ModifyPwd struct {
		Retcode string `json:"retcode"`
		Errmsg  string `json:"errmsg"`
		Account string `json:"account"`
	} `json:"modify_pwd"`
}
