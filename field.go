package boneless

type Field struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Label   string `json:"label"`
	Sort    bool   `json:"sort"`
	Column  int    `json:"column"`
	Min     string `json:"min"`
	Max     string `json:"max"`
	Step    string `json:"step"`
	Format  string `json:"format"`
	Options string `json:"options"`
	ClassId string `json:"class_id"`
	Field   string `json:"field"`
}
