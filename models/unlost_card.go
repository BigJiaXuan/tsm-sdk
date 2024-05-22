package models

type UnLostCardResponse struct {
	UnLostCard struct {
		Retcode string `json:"retcode"`
		Errmsg  string `json:"errmsg"`
		Account string `json:"account"`
	} `json:"unlost_card"`
}
