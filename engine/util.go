package engine

import (
	"fmt"
	"github.com/mumax/3/cuda"
	"github.com/mumax/3/data"
	"github.com/mumax/3/mag"
	"github.com/mumax/3/util"
	"log"
	"math"
	"os"
)

func init() {
	DeclFunc("expect", Expect, "Used for automated tests: checks if a value is close enough to the expected value")
	DeclFunc("fprintln", Fprintln, "Print to file")
	DeclFunc("sign", sign, "Signum function")
	DeclPure("vector", MakeVector, "Constructs a vector with given components")
	DeclConst("mu0", mag.Mu0, "Permittivity of vaccum (Tm/A)")
	DeclFunc("print", myprint, "Print to standard output")
	DeclFunc("LoadFile", LoadFile, "Load a .dump file")
	DeclFunc("Index2Coord", Index2Coord, "Convert cell index to x,y,z coordinate in meter")
}

// Constructs a vector
func MakeVector(x, y, z float64) data.Vector {
	return data.Vector{x, y, z}
}

// Test if have lies within want +/- maxError,
// and print suited message.
func Expect(msg string, have, want, maxError float64) {
	if math.IsNaN(have) || math.IsNaN(want) || math.Abs(have-want) > maxError {
		log.Fatal(msg, ":", " have: ", have, " want: ", want, "±", maxError)
	} else {
		log.Println(msg, ":", have, "OK")
	}
	// note: we also check "want" for NaN in case "have" and "want" are switched.
}

// Append msg to file. Used to write aggregated output of many simulations in one file.
func Fprintln(filename string, msg ...interface{}) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	util.FatalErr(err)
	defer f.Close()
	_, err = fmt.Fprintln(f, msg...)
	util.FatalErr(err)
}

// Read a magnetization state from .dump file.
func LoadFile(fname string) *data.Slice {
	s, _ := data.MustReadFile(fname)
	return s
}

// Download a quantity to host,
// or just return its data when already on host.
func Download(q Slicer) *data.Slice {
	buf, recycle := q.Slice()
	if recycle {
		defer cuda.Recycle(buf)
	}
	if buf.CPUAccess() {
		return buf
	} else {
		return buf.HostCopy()
	}
}

// print with special formatting for some known types
func myprint(msg ...interface{}) {
	for i, m := range msg {

		// Slicer: print formatted slice
		if s, ok := m.(Slicer); ok {
			s, r := s.Slice()
			msg[i] = s.HostCopy()
			if r {
				cuda.Recycle(s)
			}
			continue
		}

		// Tabledata: print average
		if m, ok := m.(TableData); ok {
			msg[i] = m.TableData()
			continue
		}
	}
	log.Println(msg...)
}

// converts cell index to coordinate, internal coordinates
func Index2Coord(ix, iy, iz int) data.Vector {
	m := Mesh()
	n := m.Size()
	c := m.CellSize()
	x := c[X] * (float64(ix) - 0.5*float64(n[X]-1))
	y := c[Y] * (float64(iy) - 0.5*float64(n[Y]-1))
	z := c[Z] * (float64(iz) - 0.5*float64(n[Z]-1))
	return data.Vector{x, y, z}
}

const (
	X = 0
	Y = 1
	Z = 2
)
