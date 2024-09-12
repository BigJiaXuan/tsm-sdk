package models

type ModifyAccInfo struct {
	ModifyAcc struct {
		Account string `json:"account"`
		Retcode string `json:"retcode"`
		Errmsg  string `json:"errmsg"`
	} `json:"modify_acc"`
}
