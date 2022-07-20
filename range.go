package gocms

import (
	"fmt"
	"strconv"
	"strings"
)

// Struct names derived from docs here:
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Range
type Range struct {
	Start int
	End   int
	Size  int
}

func (r Range) ContentRangeHeader(name string) string {
	return fmt.Sprintf("%s %d-%d/%d", name, r.Start, r.End, r.Size)
}

// Parses Range: unit=x-y into the range value's Start and End members
// Ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Range
func (r *Range) ParseHeader(header, unit string) (err error) {
	if header[:len(unit)] != unit {
		return fmt.Errorf("invalid range unit; expected %s", unit)
	}

	sets := strings.Split(header[len(unit)+1:], ",")
	if len(sets) > 1 {
		return fmt.Errorf("multiple ranges are not supported")
	}

	set := strings.TrimSpace(sets[0])
	if set[:1] == "-" {
		return fmt.Errorf("negative ranges are not supported: %s", set)
	}

	start, end, found := strings.Cut(set, "-")
	if !found {
		return fmt.Errorf("malformed range: %s", set)
	}

	if r.Start, err = strconv.Atoi(start); err != nil {
		return fmt.Errorf("malformed start value (%s): %v", start, err)
	}

	if r.End, err = strconv.Atoi(end); err != nil {
		return fmt.Errorf("malformed end value (%s): %v", end, err)
	}

	return
}
