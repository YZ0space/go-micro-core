package registry

type Service struct {
	Name     string            `json:"name"`
	Version  string            `json:"version"`
	Metadata map[string]string `json:"metadata"`
	Nodes    []*Node           `json:"nodes"`
}

type Node struct {
	Id       string            `json:"id"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Metadata map[string]string `json:"metadata"`
}
