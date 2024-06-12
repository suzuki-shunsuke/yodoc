package config

type Config struct {
	Src   string
	Dest  string
	Delim *Delim
	Tasks []*Task
}

type Delim struct {
	Left  string
	Right string
}

func (d *Delim) GetLeft() string {
	if d == nil {
		return ""
	}
	return d.Left
}

func (d *Delim) GetRight() string {
	if d == nil {
		return ""
	}
	return d.Right
}

type Task struct {
	Name         string
	Shell        []string
	Run          string
	Script       string
	Dir          string
	Env          map[string]string
	BeforeScript string `yaml:"before_script"`
	AfterScript  string `yaml:"after_script"`
	Checks       []*Check
}

type Check struct {
	Expr string
}
