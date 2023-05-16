package hello

import "net/url"

type (
	private struct {
		x, y, z int
		a       int
		b       string
		c       uint
		X, Y, Z float64
		A       string
		B       complex64
		C       uintptr
	}
	Empty      struct{}
	NotStruct  int
	NoExported struct {
		x, y, z int
		a       int
		b       string
		c       uint
	}
	KitchenSink struct {
		T0  bool
		T1  uint
		T2  uint16
		T3  uint32
		T4  uint64
		T5  int
		T6  int16
		T7  int32
		T8  int64
		T9  float32
		T10 float64
		T11 complex64
		T12 complex128
		T13 string
		T14 map[string]any
		T15 []bool
		T16 func(x int, y string, z bool) (string, error)
		T17 map[[4]complex64]any
		T18 [][]func() any
		T19 InnerStruct
		T20 *InnerStruct
		T21 *url.URL
	}

	InnerStruct struct{ X, Y, Z int }
)
