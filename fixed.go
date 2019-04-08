package fixed

// release under the terms of file license.txt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// Fixed is a fixed precision 38.24 number (supports 10.8 digits). It supports NaN.
type Fixed struct {
	fp uint64
}

// the following constants can be changed to configure a different number of decimal places - these are
// the only required changes. only 18 significant digits are supported due to NaN

const (
	nPlaces = 8
	scale = uint64(10 * 10 * 10 * 10 * 10 * 10 * 10 * 10)
	zeros = "00000000"
	max = float64(9999999999.99999999)
	nan = uint64(1<<64 - 1)
)

var (
	NaN   = Fixed{fp: nan}
	ZERO  = Fixed{fp: 0}
	ONE   = Fixed{fp: 1e8}
	TWO   = Fixed{fp: 2e8}
	THREE = Fixed{fp: 3e8}
	FOUR  = Fixed{fp: 4e8}
	FIVE  = Fixed{fp: 5e8}
	SIX   = Fixed{fp: 6e8}
	SEVEN = Fixed{fp: 7e8}
	EIGHT = Fixed{fp: 8e8}
	NINE  = Fixed{fp: 9e8}
	TEN   = Fixed{fp: 10e8}
	MAX   = Fixed{fp: 999999999999999999}
)

var errOverflow = errors.New("integer overflow")
var errNegativeNum = errors.New("negative number")
var errTooLarge = errors.New("significand too large")
var errFormat = errors.New("invalid encoding")

// NewS creates a new Fixed from a string, returning NaN if the string could not be parsed
func NewS(s string) Fixed {
	f, err := NewSErr(s)
	if err != nil {
		panic(fmt.Sprintf("newSErr(%s) err: %s", s, err))
	}
	return f
}

// NewSErr creates a new Fixed from a string, returning NaN, and error if the string could not be parsed
func NewSErr(s string) (Fixed, error) {
	if strings.HasPrefix(s, "-") {
		return NaN, errNegativeNum
	}
	if strings.ContainsAny(s, "eE") {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return NaN, err
		}
		return NewF(f), nil
	}
	if "NaN" == s {
		return NaN, nil
	}
	period := strings.Index(s, ".")
	var i uint64
	var f uint64
	var err error
	if period == -1 {
		i, err = strconv.ParseUint(s, 10, 64)
	} else {
		i, err = strconv.ParseUint(s[:period], 10, 64)
		fs := s[period+1:]
		fs = fs + zeros[:maxInt(0, nPlaces-len(fs))]
		f, err = strconv.ParseUint(fs[0:nPlaces], 10, 64)
	}
	if err != nil {
		return NaN, err
	}
	if float64(i) > max {
		return NaN, errTooLarge
	}
	return Fixed{fp: i*scale + f}, nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// NewF creates a Fixed from an float64, rounding at the 8th decimal place
func NewF(f float64) Fixed {
	if math.IsNaN(f) {
		return Fixed{fp: nan}
	}
	if f >= max || f < 0 {
		panic(errOverflow)
	}

	return Fixed{fp: uint64(f*float64(scale))}
}

// NewUI creates a Fixed for an integer, moving the decimal point n places to the left
// For example, NewUI(123,1) becomes 12.3. If n > 8, the value is truncated
func NewUI(i uint64, n uint) Fixed {
	if n > nPlaces {
		i = i / uint64(math.Pow10(int(n-nPlaces)))
		n = nPlaces
	}

	i = i * uint64(math.Pow10(int(nPlaces-n)))

	return Fixed{fp: i}
}

// NewUIFromOriginal creates a Fixed for an fixed original integer, moving the decimal point n places to the left
// For example, NewUIFromOriginal(123) becomes 0.00000123.
func NewUIFromOriginal(i uint64) Fixed {
	return NewUI(i , nPlaces)
}

func (f Fixed) IsNaN() bool {
	return f.fp == nan
}

func (f Fixed) IsZero() bool {
	return f.Equal(ZERO)
}

// Sign returns:
//
//	-1 if f <  0
//	 0 if f == 0 or NaN
//	+1 if f >  0
//
func (f Fixed) Sign() int {
	if f.IsNaN() {
		return 0
	}
	return f.Cmp(ZERO)
}

// Float converts the Fixed to a float64
func (f Fixed) Float() float64 {
	if f.IsNaN() {
		return math.NaN()
	}
	return float64(f.fp) / float64(scale)
}

// Add adds f0 to f producing a Fixed. If either operand is NaN, NaN is returned
func (f Fixed) Add(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}

	var result uint64
	result = f.fp + f0.fp

	if (result > f.fp) != (f0.fp > 0) {
		panic(errOverflow)
	}
	return Fixed{fp: result}
}

// Sub subtracts f0 from f producing a Fixed. If either operand is NaN, NaN is returned
func (f Fixed) Sub(f0 Fixed) Fixed {
	// check overflow
	if f.LessThan(f0) {
		panic(errOverflow)
	}
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	return Fixed{fp: f.fp - f0.fp}
}

// Mul multiplies f by f0 returning a Fixed. If either operand is NaN, NaN is returned
func (f Fixed) Mul(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}

	fp_a := f.fp / scale
	fp_b := f.fp % scale

	fp0_a := f0.fp / scale
	fp0_b := f0.fp % scale

	var result uint64

	if fp0_a != 0 {
		result = fp_a*fp0_a*scale + fp_b*fp0_a
	}
	if fp0_b != 0 {
		result = result + (fp_a * fp0_b) + ((fp_b)*fp0_b)/scale
	}

	// check overflow
	if (fp_a >= 1 && fp_b >= 1) && (fp0_a > 1 && fp0_b > 1) {
		if result/f.fp != fp0_a {
			panic(errOverflow)
		}
	}
	return Fixed{fp: result}
}

// Div divides f by f0 returning a Fixed. If either operand is NaN, NaN is returned
func (f Fixed) Div(f0 Fixed) Fixed {
	if f.IsNaN() || f0.IsNaN() {
		return NaN
	}
	return NewF(f.Float() / f0.Float())
}

// Round returns a rounded (half-up, away from zero) to n decimal places
func (f Fixed) Round(n int) Fixed {
	if f.IsNaN() {
		return NaN
	}

	round := .5

	f0 := f.Frac()
	f0 = f0*math.Pow10(n) + round
	f0 = float64(int(f0)) / math.Pow10(n)

	return NewF(float64(f.UInt()) + f0)
}

// Equal returns true if the f == f0. If either operand is NaN, false is returned. Use IsNaN() to test for NaN
func (f Fixed) Equal(f0 Fixed) bool {
	if f.IsNaN() || f0.IsNaN() {
		return false
	}
	return f.Cmp(f0) == 0
}

// GreaterThan tests Cmp() for 1
func (f Fixed) GreaterThan(f0 Fixed) bool {
	return f.Cmp(f0) == 1
}

// GreaterThaOrEqual tests Cmp() for 1 or 0
func (f Fixed) GreaterThanOrEqual(f0 Fixed) bool {
	cmp := f.Cmp(f0)
	return cmp == 1 || cmp == 0
}

// LessThan tests Cmp() for -1
func (f Fixed) LessThan(f0 Fixed) bool {
	return f.Cmp(f0) == -1
}

// LessThan tests Cmp() for -1 or 0
func (f Fixed) LessThanOrEqual(f0 Fixed) bool {
	cmp := f.Cmp(f0)
	return cmp == -1 || cmp == 0
}

// Cmp compares two Fixed. If f == f0, return 0. If f > f0, return 1. If f < f0, return -1. If both are NaN, return 0. If f is NaN, return 1. If f0 is NaN, return -1
func (f Fixed) Cmp(f0 Fixed) int {
	if f.IsNaN() && f0.IsNaN() {
		return 0
	}
	if f.IsNaN() {
		return 1
	}
	if f0.IsNaN() {
		return -1
	}

	if f.fp == f0.fp {
		return 0
	}
	if f.fp < f0.fp {
		return -1
	}
	return 1
}

// String converts a Fixed to a string, dropping trailing zeros
func (f Fixed) String() string {
	s, point := f.toStr()
	if point == -1 {
		return s
	}
	index := len(s) - 1
	for ; index != point; index-- {
		if s[index] != '0' {
			return s[:index+1]
		}
	}
	return s[:point]
}

// StringN converts a Fixed to a String with a specified number of decimal places, truncating as required
func (f Fixed) StringN(decimals int) string {

	s, point := f.toStr()

	if point == -1 {
		return s
	}
	if decimals == 0 {
		return s[:point]
	} else {
		return s[:point+decimals+1]
	}
}

func (f Fixed) toStr() (string, int) {
	fp := f.fp
	if fp == 0 {
		return "0." + zeros, 1
	}
	if fp == nan {
		return "NaN", -1
	}

	b := make([]byte, 24)
	b = itoa(b, fp)

	return string(b), len(b) - nPlaces - 1
}

func itoa(buf []byte, val uint64) []byte {
	i := len(buf) - 1
	idec := i - nPlaces
	for val >= 10 || i >= idec {
		buf[i] = byte(val%10 + '0')
		i--
		if i == idec {
			buf[i] = '.'
			i--
		}
		val /= 10
	}
	buf[i] = byte(val + '0')
	return buf[i:]
}

// UInt return the integer portion of the Fixed, or 0 if NaN
func (f Fixed) UInt() uint64 {
	if f.IsNaN() {
		return 0
	}
	return f.fp / scale
}

// Frac return the fractional portion of the Fixed, or NaN if NaN
func (f Fixed) Frac() float64 {
	if f.IsNaN() {
		return math.NaN()
	}
	return float64(f.fp%scale) / float64(scale)
}

// Original return the original digital of the Fixed,
func (f Fixed) Original() uint64 {
	return f.fp
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface
func (f *Fixed) UnmarshalBinary(data []byte) error {
	fp, n := binary.Uvarint(data)
	if n < 0 {
		return errFormat
	}
	f.fp = fp
	return nil
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (f Fixed) MarshalBinary() (data []byte, err error) {
	var buffer [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buffer[:], f.fp)
	return buffer[:n], nil
}

// WriteTo write the Fixed to an io.Writer, returning the number of bytes written
func (f Fixed) WriteTo(w io.ByteWriter) error {
	x := f.fp
	i := 0
	for x >= 0x80 {
		err := w.WriteByte(byte(x) | 0x80)
		if err != nil {
			return err
		}
		x >>= 7
		i++
	}
	return w.WriteByte(byte(x))
}

// ReadFrom reads a Fixed from an io.Reader
func ReadFrom(r io.ByteReader) (Fixed, error) {
	fp, err := binary.ReadUvarint(r)
	if err != nil {
		return NaN, err
	}
	return Fixed{fp: fp}, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (f *Fixed) UnmarshalJSON(bytes []byte) error {
	s := string(bytes)
	if s == "null" {
		return nil
	}

	fixed, err := NewSErr(s)
	*f = fixed
	if err != nil {
		return fmt.Errorf("error decoding string '%s': %s", s, err)
	}
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (f Fixed) MarshalJSON() ([]byte, error) {
	buffer := make([]byte, 24)
	return itoa(buffer, f.fp), nil
}
