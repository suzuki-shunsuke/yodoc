package config

type Config struct {
	Src         string
	Dest        string
	Delim       *Delim
	Tasks       []*Task
	Env         map[string]string
	AppendedEnv map[string]string `yaml:"appended_env"`
}
