package models

type QuerySnoNameFromBarcodeResponse struct {
	QuerySnoNameFromBarcode `json:"query_snoname_frombarcode"`
}
type QuerySnoNameFromBarcode struct {
	Retcode string `json:"retcode"`
	Errmsg  string `json:"errmsg"`
	Sno     string `json:"sno"`
	Name    string `json:"name"`
	Paytype string `json:"paytype"`
	Payacc  string `json:"payacc"`
	Accname string `json:"accname"`
}
