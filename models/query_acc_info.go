package models

type AccInfo struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	Balance         string `json:"balance"`
	AutotransLimite string `json:"autotrans_limite"`
	AutotransAmt    string `json:"autotrans_amt"`
	AutotransFlag   string `json:"autotrans_flag"`
	Singlelimit     string `json:"singlelimit"`
	Daycostlimit    string `json:"daycostlimit"`
	Daycostamt      string `json:"daycostamt"`
	Nonpwdlimit     string `json:"nonpwdlimit"`
}

type QueryAccInfoResponse struct {
	QueryAccinfo struct {
		Retcode string    `json:"retcode"`
		Errmsg  string    `json:"errmsg"`
		Account string    `json:"account"`
		Accinfo []AccInfo `json:"accinfo"`
	} `json:"query_accinfo"`
}
