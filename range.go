package gocms

import (
	"encoding/json"
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

// Builds Content-Range header: <unit> <start>-<end>/<size>
// Ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Range
func (r Range) ContentRangeHeader(unit string) string {
	return fmt.Sprintf("%s %d-%d/%d", unit, r.Start, r.End, r.Size)
}

// Parses Range: <unit>=<start>-<end> into the range's Start and End members
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
		return fmt.Errorf("malformed start value (%s): %w", start, err)
	}

	if r.End, err = strconv.Atoi(end); err != nil {
		return fmt.Errorf("malformed end value (%s): %w", end, err)
	}

	if r.End < r.Start {
		return fmt.Errorf("invalid range, end before start (%d < %d)", r.End, r.Start)
	}

	return
}

// Supports two different combinations of parameters:
// _start, _end: Zero-based indexes for the starting and ending items
// _page, _per_page: paging with _page starting at 1; _per_page defaults to 10
// range: JSON encoded 2-element array with zero-based indexes for start and end
// If no useful parameters are available, defaults to the first 10 items
func (r *Range) ParseParams(params map[string]string) (err error) {
	var start, end, page, perPage int

	startStr, hasStart := params["_start"]
	endStr, hasEnd := params["_end"]
	pageStr, hasPage := params["_page"]
	perPageStr, hasPerPage := params["_per_page"]
	rangeStr, hasRange := params["range"]

	if hasStart && !hasEnd {
		return fmt.Errorf("missing _end with _start")
	}

	if !hasStart && hasEnd {
		return fmt.Errorf("missing _start with _end")
	}

	if hasStart && hasEnd {
		if start, err = strconv.Atoi(startStr); err != nil {
			return fmt.Errorf("parsing _start: %w", err)
		}
		if end, err = strconv.Atoi(endStr); err != nil {
			return fmt.Errorf("parsing _end: %w", err)
		}
	}

	// Process range JSON and let it turn into start/end and fall through that
	// processing
	if hasRange {
		bounds := make([]int, 0, 2)
		if err = json.Unmarshal([]byte(rangeStr), &bounds); err != nil {
			return fmt.Errorf("invalid range: %w", err)
		}
		if num := len(bounds); num != 2 {
			return fmt.Errorf("expect exactly 2 range elements; got %d", num)
		}
		hasStart, hasEnd = true, true
		start, end = bounds[0], bounds[1]
	}

	if hasStart && hasEnd {
		if start < 0 {
			return fmt.Errorf("start index is less than zero")
		}
		if end < 0 {
			return fmt.Errorf("end index is less than zero")
		}
		if start > end {
			return fmt.Errorf("start index is greater than _end")
		}

		r.Start, r.End = start, end
		return
	}

	if hasPage {
		if page, err = strconv.Atoi(pageStr); err != nil {
			return fmt.Errorf("parsing _page: %w", err)
		}
		if page < 1 {
			return fmt.Errorf("_page is less than one")
		}
	}

	if hasPerPage {
		if perPage, err = strconv.Atoi(perPageStr); err != nil {
			return fmt.Errorf("parsing _per_page: %w", err)
		}
		if perPage < 1 {
			return fmt.Errorf("_per_page is less than one")
		}
	}

	if hasPage || hasPerPage {
		if !hasPage {
			page = 1
		}
		if !hasPerPage {
			perPage = 10
		}
		r.Start = (page - 1) * perPage
		r.End = page*perPage - 1
		return
	}

	// Default to 10 items on page 1
	r.End = 9
	return
}
