package config

type Config struct {
	Src   string
	Dest  string
	Delim *Delim
	Tasks []*Task
}
