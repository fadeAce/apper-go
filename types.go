package apper_go

type ConfJ struct {
	Sites map[string]SiteJ `json:"sites"`
}
type ApperConfJ struct {
	Database string `json:"database"`
}

type SingleJ struct {
	Type string
	Rule string
	Key  string
}

type SiteJ struct {
	Single []SingleJ `json:"singles"`
}

type Command struct {
	Configs ConfJ  `json:"config"`
	Cmd     string `json:"cmd"`
}
//传 Start 用到
type Nats_data struct {
	Conf Conf
	Type string
}
//得到 Start 用到
type Nats_data1 struct {
	Key string
	TXID string
}