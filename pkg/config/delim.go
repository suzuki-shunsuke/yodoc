package config

type Delim struct {
	Left  string `json:"left,omitempty"`
	Right string `json:"right,omitempty"`
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
