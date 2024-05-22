package models

type UnFreezeCardResponse struct {
	UnFrozenCard struct {
		Retcode string `json:"retcode"`
		Errmsg  string `json:"errmsg"`
		Account string `json:"account"`
	} `json:"unfrozen_card"`
}
