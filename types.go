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

type Conf struct {
	Sites map[string]Site `yaml:"sites"`
}

type ApperConf struct {
	Database string `yaml:"database"`
}

type Single struct {
	Type string
	Rule string
	Key  string
	Url  string
}

type Site struct {
	Single []Single `yaml:"singles"`
}

//传 Start 用到
type NatsData struct {
	Conf Conf   `json:"config"`
	Type string `json:"cmd"`
}

//得到 Start 用到
type NatsData1 struct {
	Key  string
	TXID string
}
