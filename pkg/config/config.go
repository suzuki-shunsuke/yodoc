package config

type Config struct {
	Src   string `json:"src"`
	Dest  string `json:"dest"`
	Delim *Delim `json:"delim,omitempty"`
}
