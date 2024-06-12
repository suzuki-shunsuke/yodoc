package config

type Config struct {
	Src         string
	Dest        string
	Deliminater *Deliminaters
	Tasks       []*Task
}

type Deliminaters struct {
	Left  string
	Right string
}

type Task struct {
	Name         string
	Shell        string
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
