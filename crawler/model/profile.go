package model

type Video struct {
	Id           string
	Name         string
	UpdateStatus string
	Type         []string
	Time         string
	Region       string
	Description  string
	ImgUrl       string
	BaiduDbank   string
}

type Dbank struct {
	VideoId  string
	SharePwd string
	ShareUrl string
}
