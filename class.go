package gocms

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
)

type Class struct {
	Id          string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	TableLabels string    `json:"table_labels"`
	TableFields string    `json:"table_fields"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Fields      []Field   `json:"fields"`
}

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

func (r *Range) ParseHeader(header, unit string) (err error) {
	if header[:len(unit)] != unit {
		return fmt.Errorf("invalid range unit; expected %s", unit)
	}

	sets := strings.Split(header[len(unit):], ",")
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

type ClassFilter struct {
	Range Range
}

type ClassRepository interface {
	CreateClass(context.Context, *Class) error
	DeleteClass(context.Context, string) error
	GetClassById(context.Context, string) (Class, error)
	GetClassList(context.Context, ClassFilter) ([]Class, Range, error)
	UpdateClass(context.Context, *Class) error
}

type ClassService struct {
	repo ClassRepository
}

func NewClassService(repo ClassRepository) ClassService {
	return ClassService{
		repo: repo,
	}
}

func (s ClassService) ById(ctx context.Context, id string) (Class, error) {
	if _, err := xid.FromString(id); err != nil {
		return Class{}, err
	}
	return s.repo.GetClassById(ctx, id)
}

func (s ClassService) Create(ctx context.Context, class *Class) (err error) {
	// TODO validate internal fields

	if class.Id != "" {
		return fmt.Errorf("class already has an ID")
	}

	// TODO check for existing class with same slug

	now := time.Now()
	class.Id = xid.NewWithTime(now).String()
	class.Created = now
	class.Updated = now

	return s.repo.CreateClass(ctx, class)
}

func (s ClassService) List(ctx context.Context, filter ClassFilter) ([]Class, Range, error) {
	return s.repo.GetClassList(ctx, filter)
}

func (s ClassService) Update(ctx context.Context, class *Class) (err error) {
	if class.Id == "" {
		return fmt.Errorf("class has no ID")
	}

	// TODO check for slug overwrite

	class.Updated = time.Now()

	return s.repo.UpdateClass(ctx, class)
}
