package models

type QueryAccInfoTotalResponse struct {
	Retcode string `json:"retcode"`
	Errmsg  string `json:"errmsg"`

	QueryAccinfoTotal `json:"query_accinfo_total"`
}
type QueryAccinfoTotal struct {
	Acctype    string `json:"acctype"`
	TotalCount string `json:"total_count"`
	TotalAmt   string `json:"total_amt"`
}
