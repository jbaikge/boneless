package slicer

import (
	"testing"

	"github.com/zeebo/assert"
)

func TestAll(t *testing.T) {
	slicer := NewSlicer(0, 99)
	for i := 0; i < 10; i++ {
		slicer.Add(10)
		assert.Equal(t, 0, slicer.Start())
		assert.Equal(t, 10, slicer.End())
	}
}

func TestOne(t *testing.T) {
	slicer := NewSlicer(54, 54)
	for i := 0; i < 10; i++ {
		slicer.Add(10)
		switch i {
		case 5:
			assert.Equal(t, 4, slicer.Start())
			assert.Equal(t, 5, slicer.End())
		default:
			assert.Equal(t, 0, slicer.Start())
			assert.Equal(t, 0, slicer.End())
		}
	}
}

func TestFront(t *testing.T) {
	slicer := NewSlicer(0, 11)
	for i := 0; i < 10; i++ {
		slicer.Add(10)
		switch i {
		case 0:
			assert.Equal(t, 0, slicer.Start())
			assert.Equal(t, 10, slicer.End())
		case 1:
			assert.Equal(t, 0, slicer.Start())
			assert.Equal(t, 2, slicer.End())
		default:
			assert.Equal(t, 0, slicer.Start())
			assert.Equal(t, 0, slicer.End())
		}
	}
}

func TestMiddle(t *testing.T) {
	slicer := NewSlicer(47, 73)
	for i := 0; i < 10; i++ {
		slicer.Add(10)
		switch i {
		case 4:
			assert.Equal(t, 7, slicer.Start())
			assert.Equal(t, 10, slicer.End())
		case 5, 6:
			assert.Equal(t, 0, slicer.Start())
			assert.Equal(t, 10, slicer.End())
		case 7:
			assert.Equal(t, 0, slicer.Start())
			assert.Equal(t, 4, slicer.End())
		default:
			assert.Equal(t, 0, slicer.Start())
			assert.Equal(t, 0, slicer.End())
		}
	}
}

func TestBack(t *testing.T) {
	slicer := NewSlicer(88, 99)
	for i := 0; i < 10; i++ {
		slicer.Add(10)
		switch i {
		case 8:
			assert.Equal(t, 8, slicer.Start())
			assert.Equal(t, 10, slicer.End())
		case 9:
			assert.Equal(t, 0, slicer.Start())
			assert.Equal(t, 10, slicer.End())
		default:
			assert.Equal(t, 0, slicer.Start())
			assert.Equal(t, 0, slicer.End())
		}
	}
}

func TestTotal(t *testing.T) {
	slicer := NewSlicer(4, 5)
	for i := 0; i < 10; i++ {
		slicer.Add(10)
	}
	assert.Equal(t, 100, slicer.Total())
}
