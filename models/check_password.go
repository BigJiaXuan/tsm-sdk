package models

type CheckPasswordResponse struct {
	CheckPassword `json:"check_pwd"`
}
type CheckPassword struct {
	Retcode    string `json:"retcode"`
	Errmsg     string `json:"errmsg"`
	Account    string `json:"account"`
	Sno        string `json:"sno"`
	Identityid string `json:"identityid"`
}
