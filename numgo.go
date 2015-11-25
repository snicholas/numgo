package numgo

import (
	"fmt"
	"strings"
	"sync"
)

type Arrayf struct {
	sync.RWMutex
	shape   []uint64
	strides []uint64
	data    []float64
}

// Create creates an Arrayf object with dimensions given in order from outer-most to inner-most
// All values will default to zero
func Create(shape ...int) (a *Arrayf) {
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v <= 0 {
			return nil
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	a = new(Arrayf)
	a.shape = sh
	a.data = make([]float64, sz)

	a.strides = make([]uint64, len(sh)+1)
	tmp := uint64(1)
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= sh[i-1]
	}
	a.strides[0] = tmp
	return
}

// Internal function to create using the shape of another array
func create(shape ...uint64) (a *Arrayf) {
	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v <= 0 {
			return nil
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	a = new(Arrayf)
	a.shape = sh
	a.data = make([]float64, sz)

	a.strides = make([]uint64, len(sh)+1)
	tmp := uint64(1)
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= sh[i-1]
	}
	a.strides[0] = tmp
	return
}

// Full creates an Arrayf object with dimensions givin in order from outer-most to inner-most
// All elements will be set to 'val' in the retuen
func Full(val float64, shape ...int) (a *Arrayf) {
	a = Create(shape...)
	if a == nil {
		return nil
	}
	a.AddC(val)
	return
}

func full(val float64, shape ...uint64) (a *Arrayf) {
	a = create(shape...)
	if a == nil {
		return nil
	}
	a.AddC(val)
	return
}

// String Satisfies the Stringer interface for fmt package
func (a *Arrayf) String() (s string) {

	a.RLock()
	defer a.RUnlock()

	stride := a.shape[len(a.shape)-1]

	for i, k := uint64(0), 0; i+stride <= uint64(len(a.data)); i, k = i+stride, k+1 {

		t := ""
		for j, v := range a.strides {
			if i%v == 0 && j < len(a.strides)-2 {
				t += "["
			}
		}

		s += strings.Repeat(" ", len(a.shape)-len(t)-1) + t
		s += fmt.Sprint(a.data[i : i+stride])

		t = ""
		for j, v := range a.strides {
			if (i+stride)%v == 0 && j < len(a.strides)-2 {
				t += "]"
			}
		}

		s += t + strings.Repeat(" ", len(a.shape)-len(t)-1)
		if i+stride != uint64(len(a.data)) {
			s += "\n"
			if len(t) > 0 {
				s += "\n"
			}
		}
	}
	return
}

// Arange Creates an array in one of three different ways, depending on input:
// One (stop):         Arrayf from zero to positive value or negative value to zero
// Two (start,stop):   Arrayf from start to stop, with increment of 1 or -1, depending on inputs
// Three (start, stop, step): Arrayf from start to stop, with increment of step
//
// Any inputs beyond three values are ignored
func Arange(vals ...float64) (a *Arrayf) {
	var start, stop, step float64 = 0, 0, 1

	switch len(vals) {
	case 0:
		return nil
	case 1:
		if vals[0] <= 0 {
			start, stop, step = vals[0], 0, -1
		} else {
			stop = vals[0]
		}
	case 2:
		if vals[1] < vals[0] {
			step = -1
		}
		start, stop = vals[0], vals[1]
	default:
		if vals[1] < vals[0] && vals[2] >= 0 || vals[1] > vals[0] && vals[2] <= 0 {
			return nil
		}
		start, stop, step = vals[0], vals[1], vals[2]
	}

	a = Create(int((stop - start) / step))
	for i, v := 0, start; i < len(a.data); i, v = i+1, v+step {
		a.data[i] = v
	}
	return
}

// Reshape Changes the size of the array axes.  Values are not changed or moved.
// This must not change the size of the array.
// Incorrect dimensions will return a nil pointer
func (a *Arrayf) Reshape(shape ...int) *Arrayf {
	if a == nil {
		return nil
	}

	a.Lock()
	defer a.Unlock()

	var sz uint64 = 1
	sh := make([]uint64, len(shape))
	for i, v := range shape {
		if v <= 0 {
			return nil
		}
		sz *= uint64(v)
		sh[i] = uint64(v)
	}

	if sz != uint64(len(a.data)) {
		return nil
	}

	a.strides = make([]uint64, len(sh)+1)
	tmp := uint64(1)
	for i := len(a.strides) - 1; i > 0; i-- {
		a.strides[i] = tmp
		tmp *= sh[i-1]
	}
	a.strides[0] = tmp
	a.shape = sh

	return a
}

// C will return a deep copy of the source array.
func (a *Arrayf) C() (b *Arrayf) {
	if a == nil {
		return nil
	}

	b = create(a.shape...)
	for i, v := range a.data {
		b.data[i] = v
	}
	return
}
