package models

type LostCardResponse struct {
	LostCard struct {
		Retcode string `json:"retcode"`
		Errmsg  string `json:"errmsg"`
		Account string `json:"account"`
	} `json:"lost_card"`
}
