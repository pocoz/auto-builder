package types

const ActionPush = "push"

type Payload struct {
	Events []*Event
}

type Event struct {
	ID        string      `json:"id"`
	TimeStamp string      `json:"timestamp"`
	Action    string      `json:"action"`
	Target    Target      `json:"target"`
	Request   Request     `json:"request"`
	Actor     interface{} `json:"actor"`
}

type Target struct {
	MediaType  string `json:"mediaType"`
	Size       int64  `json:"size"`
	Digest     string `json:"digest"`
	Length     int64  `json:"length"`
	Repository string `json:"repository"`
	URL        string `json:"url"`
	Tag        string `json:"tag"`
}

type Request struct {
	ID        string `json:"id"`
	Addr      string `json:"addr"`
	Host      string `json:"host"`
	Method    string `json:"method"`
	UserAgent string `json:"useragent"`
}

type Source struct {
	Addr       string `json:"addr"`
	InstanceID string `json:"instanceID"`
}

type Configs struct {
	Auth       Auth      `json:"auth"`
	ConfigList []*Config `json:"config_list"`
}

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Config struct {
	Image        string   `json:"image"`
	Environments []string `json:"environments"`
	Cmd          []string `json:"cmd"`
	Volumes      []string `json:"volumes"`
}
