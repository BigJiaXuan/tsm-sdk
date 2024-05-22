package models

type FreezeCardResponse struct {
	FrozenCard struct {
		Retcode string `json:"retcode"`
		Errmsg  string `json:"errmsg"`
		Account string `json:"account"`
	} `json:"frozen_card"`
}
