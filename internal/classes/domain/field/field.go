package field

import (
	"errors"
	"regexp"
)

var validFieldName = regexp.MustCompile(`^[a-z0-9_]+$`)

type Field struct {
	fieldType  string
	name       string
	label      string
	sort       bool
	column     int
	min        string
	max        string
	step       string
	format     string
	options    string
	classId    string
	classField string
}

func NewField(
	fieldType string,
	name string,
	label string,
	sort bool,
	column int,
	min string,
	max string,
	step string,
	format string,
	options string,
	classId string,
	classField string,
) (*Field, error) {
	// TODO check if fieldType is valid
	if name == "" {
		return nil, errors.New("empty field name")
	}
	if !validFieldName.MatchString(name) {
		return nil, errors.New("invalid field name, lower alphanumeric only")
	}
	if label == "" {
		return nil, errors.New("empty field label")
	}
	if classId != "" && classField == "" {
		return nil, errors.New("empty class field when class ID set")
	}
	return &Field{
		fieldType:  fieldType,
		name:       name,
		label:      label,
		sort:       sort,
		column:     column,
		min:        min,
		max:        max,
		step:       step,
		format:     format,
		options:    options,
		classId:    classId,
		classField: classField,
	}, nil
}
