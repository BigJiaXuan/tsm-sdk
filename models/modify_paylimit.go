package models

type ModifyPayLimit struct {
	PaylimiteModify struct {
		Retcode string `json:"retcode"`
		Errmsg  string `json:"errmsg"`
		Account string `json:"account"`
		Acctype string `json:"acctype"`
	} `json:"paylimite_modify"`
}
