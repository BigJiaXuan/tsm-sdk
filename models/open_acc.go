package models

type OpenAccResponse struct {
	OpenAcc struct {
		Retcode string `json:"retcode"`
		Errmsg  string `json:"errmsg"`
		Sno     string `json:"sno"`
		Account string `json:"account"`
	} `json:"open_acc"`
}
