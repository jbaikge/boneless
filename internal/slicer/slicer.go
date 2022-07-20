package slicer

type Slicer struct {
	start     int
	end       int
	processed int
	total     int
}

func NewSlicer(start, end int) *Slicer {
	return &Slicer{
		start: start,
		end:   end,
	}
}

func (s *Slicer) Add(chunk int) {
	s.processed = s.total
	s.total += chunk
}

func (s *Slicer) End() int {
	// No need to chunk if we haven't reached the start point yet
	if s.total < s.start {
		return 0
	}

	chunkSize := s.total - s.processed
	end := s.end + 1
	diff := s.total - end
	// Ending index is beyond what we've seen so far, give the max value
	if diff <= 0 {
		return chunkSize
	}
	// End is within sight, give a value that leads up to the end
	if diff < chunkSize {
		return chunkSize - diff
	}
	// Chunking is over, return zero with Start() for an empty slice
	return 0
}

func (s *Slicer) Start() int {
	// Haven't reached the starting point yet
	if s.total < s.start {
		return 0
	}

	// Start point is within this chunk somewhere
	chunkSize := s.total - s.processed
	diff := s.total - s.start
	if diff < chunkSize {
		return chunkSize - diff
	}

	// Either start from the beginning or End() will return zero for an empty slice
	return 0
}

func (s *Slicer) Total() int {
	return s.total
}
