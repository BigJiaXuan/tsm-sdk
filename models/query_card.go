package models

type Card struct {
	Account         string        `json:"account"`
	Name            string        `json:"name"`
	UnsettleAmount  string        `json:"unsettle_amount"`
	DbBalance       string        `json:"db_balance"`
	AccStatus       string        `json:"acc_status"`
	Lostflag        string        `json:"lostflag"`
	Freezeflag      string        `json:"freezeflag"`
	Barflag         string        `json:"barflag"`
	Idflag          string        `json:"idflag"`
	Expdate         string        `json:"expdate"`
	Cardtype        string        `json:"cardtype"`
	Cardname        string        `json:"cardname"`
	Bankacc         string        `json:"bankacc"`
	Sno             string        `json:"sno"`
	Phone           string        `json:"phone"`
	Certtype        string        `json:"certtype"`
	Cert            string        `json:"cert"`
	Createdate      string        `json:"createdate"`
	AutotransLimite string        `json:"autotrans_limite"`
	AutotransAmt    string        `json:"autotrans_amt"`
	AutotransFlag   string        `json:"autotrans_flag"`
	Mscard          string        `json:"mscard"`
	ScardNum        string        `json:"scard_num"`
	ElecAccamt      string        `json:"elec_accamt"`
	MergeAccamt     string        `json:"merge_accamt"`
	MergeFlag       string        `json:"merge_flag"`
	Sex             string        `json:"sex"`
	Areacode        string        `json:"areacode"`
	Areaname        string        `json:"areaname"`
	Pidclass        string        `json:"pidclass"`
	Pidcode         string        `json:"pidcode"`
	Pidname         string        `json:"pidname"`
	Deptcode        string        `json:"deptcode"`
	Deptname        string        `json:"deptname"`
	Cardid          string        `json:"cardid"`
	Schcode         string        `json:"schcode"`
	Cardnum         string        `json:"cardnum"`
	Accinfo         []interface{} `json:"accinfo"`
}

type QueryCardResponse struct {
	QueryCard struct {
		Retcode string `json:"retcode"`
		Errmsg  string `json:"errmsg"`
		Card    []Card `json:"card"`
	} `json:"query_card"`
}
