package gocms

type Field struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Label   string `json:"label"`
	Sort    bool   `json:"sort"`
	Min     string `json:"min"`
	Max     string `json:"max"`
	Step    string `json:"step"`
	Format  string `json:"format"`
	Options string `json:"options"`
}
