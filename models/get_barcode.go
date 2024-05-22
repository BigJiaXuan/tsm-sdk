package models

type GerBarCodeResponse struct {
	GetBarCode `json:"barcode_get"`
}
type GetBarCode struct {
	Retcode string `json:"retcode"`
	Errmsg  string `json:"errmsg"`
	Account string `json:"account"`
	Acctype string `json:"acctype"`
	Paytype string `json:"paytype"`
	Payacc  string `json:"payacc"`
	Barcode string `json:"barcode"`
	Expires string `json:"expires"`
}
