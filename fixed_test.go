package fixed_test

import (
	"bytes"
	"encoding/json"
	. "github.com/cryptowrold/fixed"
	"math"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	f0 := NewFromString("123.456")
	f1 := NewFromString("123.456")

	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}

	if f0.UInt() != 123 {
		t.Error("should be equal", f0.UInt(), 123)
	}

	if f0.String() != "123.456" {
		t.Error("should be equal", f0.String(), "123.456")
	}

	f0 = NewFromFloat(1)
	f1 = NewFromFloat(.5).Add(NewFromFloat(.5))
	f2 := NewFromFloat(.3).Add(NewFromFloat(.3)).Add(NewFromFloat(.4))

	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}
	if !f0.Equal(f2) {
		t.Error("should be equal", f0, f2)
	}

	f0 = NewFromFloat(.999)
	if f0.String() != "0.999" {
		t.Error("should be equal", f0, "0.999")
	}
}

func TestNewI(t *testing.T) {
	f := NewFromUintWithExponent(123, 1)
	if f.String() != "12.3" {
		t.Error("should be equal", f, "12.3")
	}
	f = NewFromUint(123)
	if f.String() != "123" {
		t.Error("should be equal", f, "123")
	}
	f = NewFromUintWithExponent(123456789012, 9)
	if f.String() != "123.45678901" {
		t.Error("should be equal", f, "123.45678901")
	}
	f = NewFromUintWithExponent(123456789012, 9)
	if f.StringN(7) != "123.4567890" {
		t.Error("should be equal", f.StringN(7), "123.4567890")
	}

}

func TestSign(t *testing.T) {
	f0 := NewFromString("0")
	if f0.Sign() != 0 {
		t.Error("should be equal", f0.Sign(), 0)
	}
	f0 = NewFromString("NaN")
	if f0.Sign() != 0 {
		t.Error("should be equal", f0.Sign(), 0)
	}
	// f0 = NewFromString("-100")
	// if f0.Sign() != -1 {
	// 	t.Error("should be equal", f0.Sign(), -1)
	// }
	f0 = NewFromString("100")
	if f0.Sign() != 1 {
		t.Error("should be equal", f0.Sign(), 1)
	}

}

func TestMaxValue(t *testing.T) {
	f0 := NewFromString("1234567890")
	if f0.String() != "1234567890" {
		t.Error("should be equal", f0, "1234567890")
	}
	// f0 = NewFromString("123456789012")
	// if f0.String() != "NaN" {
	// 	t.Error("should be equal", f0, "NaN")
	// }
	// f0 = NewFromString("-12345678901")
	// if f0.String() != "-12345678901" {
	// 	t.Error("should be equal", f0, "-12345678901")
	// }
	// f0 = NewFromString("-123456789012")
	// if f0.String() != "NaN" {
	// 	t.Error("should be equal", f0, "NaN")
	// }
	f0 = NewFromString("9999999999")
	if f0.String() != "9999999999" {
		t.Error("should be equal", f0, "9999999999")
	}
	f0 = NewFromString("9.9999999")
	if f0.String() != "9.9999999" {
		t.Error("should be equal", f0, "9.9999999")
	}
	f0 = NewFromString("9999999999.9999999")
	if f0.String() != "9999999999.9999999" {
		t.Error("should be equal", f0, "9999999999.9999999")
	}
	f0 = NewFromString("9999999999.12345678901234567890")
	if f0.String() != "9999999999.12345678" {
		t.Error("should be equal", f0, "9999999999.12345678")
	}

}

func TestFloat(t *testing.T) {
	f0 := NewFromString("123.456")
	f1 := NewFromFloat(123.456)

	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}

	f1 = NewFromFloat(0.0001)

	if f1.String() != "0.0001" {
		t.Error("should be equal", f1.String(), "0.0001")
	}

	f1 = NewFromString(".1")
	f2 := NewFromString(NewFromFloat(f1.Float()).String())
	if !f1.Equal(f2) {
		t.Error("should be equal", f1, f2)
	}

}

func TestInfinite(t *testing.T) {
	f0 := NewFromString("0.10")
	f1 := NewFromFloat(0.10)

	if !f0.Equal(f1) {
		t.Error("should be equal", f0, f1)
	}

	f2 := NewFromFloat(0.0)
	for i := 0; i < 3; i++ {
		f2 = f2.Add(NewFromFloat(.10))
	}
	if f2.String() != "0.3" {
		t.Error("should be equal", f2.String(), "0.3")
	}

	f2 = NewFromFloat(0.0)
	for i := 0; i < 10; i++ {
		f2 = f2.Add(NewFromFloat(.10))
	}
	if f2.String() != "1" {
		t.Error("should be equal", f2.String(), "1")
	}

}

func TestAddSub(t *testing.T) {
	f0 := NewFromString("1")
	f1 := NewFromString("0.3333333")

	f2 := f0.Sub(f1)
	f2 = f2.Sub(f1)
	f2 = f2.Sub(f1)

	if f2.String() != "0.0000001" {
		t.Error("should be equal", f2.String(), "0.0000001")
	}
	f2 = f2.Sub(NewFromString("0.0000001"))
	if f2.String() != "0" {
		t.Error("should be equal", f2.String(), "0")
	}

	f0 = NewFromString("0")
	for i := 0; i < 10; i++ {
		f0 = f0.Add(NewFromString("0.1"))
	}
	if f0.String() != "1" {
		t.Error("should be equal", f0.String(), "1")
	}

}

func TestMulDiv(t *testing.T) {
	f0 := NewFromString("123.456")
	f1 := NewFromString("1000")

	f2 := f0.Mul(f1)
	if f2.String() != "123456" {
		t.Error("should be equal", f2.String(), "123456")
	}
	f0 = NewFromString("123456")
	f1 = NewFromString("0.0001")

	f2 = f0.Mul(f1)
	if f2.String() != "12.3456" {
		t.Error("should be equal", f2.String(), "12.3456")
	}

	// f0 = NewFromString("123.456")
	// f1 = NewFromString("-1000")
	//
	// f2 = f0.Mul(f1)
	// if f2.String() != "-123456" {
	// 	t.Error("should be equal", f2.String(), "-123456")
	// }
	//
	// f0 = NewFromString("-123.456")
	// f1 = NewFromString("-1000")
	//
	// f2 = f0.Mul(f1)
	// if f2.String() != "123456" {
	// 	t.Error("should be equal", f2.String(), "123456")
	// }
	//
	// f0 = NewFromString("123.456")
	// f1 = NewFromString("-1000")
	//
	// f2 = f0.Mul(f1)
	// if f2.String() != "-123456" {
	// 	t.Error("should be equal", f2.String(), "-123456")
	// }
	//
	// f0 = NewFromString("-123.456")
	// f1 = NewFromString("-1000")
	//
	// f2 = f0.Mul(f1)
	// if f2.String() != "123456" {
	// 	t.Error("should be equal", f2.String(), "123456")
	// }

	f0 = NewFromString("10000.1")
	f1 = NewFromString("10000")

	f2 = f0.Mul(f1)
	if f2.String() != "100001000" {
		t.Error("should be equal", f2.String(), "100001000")
	}

	f2 = f2.Div(f1)
	if !f2.Equal(f0) {
		t.Error("should be equal", f0, f2)
	}

	f0 = NewFromString("2")
	f1 = NewFromString("3")

	f2 = f0.Div(f1)
	if f2.String() != "0.66666666" {
		t.Error("should be equal", f2.String(), "0.66666666")
	}

	f0 = NewFromString("1000")
	f1 = NewFromString("10")

	f2 = f0.Div(f1)
	if f2.String() != "100" {
		t.Error("should be equal", f2.String(), "100")
	}

	f0 = NewFromString("1000")
	f1 = NewFromString("0.1")

	f2 = f0.Div(f1)
	if f2.String() != "10000" {
		t.Error("should be equal", f2.String(), "10000")
	}

	f0 = NewFromString("1")
	f1 = NewFromString("0.1")

	f2 = f0.Mul(f1)
	if f2.String() != "0.1" {
		t.Error("should be equal", f2.String(), "0.1")
	}

}
//
// func TestNegatives(t *testing.T) {
// 	f0 := NewFromString("99")
// 	f1 := NewFromString("100")
//
// 	f2 := f0.Sub(f1)
// 	if f2.String() != "-1" {
// 		t.Error("should be equal", f2.String(), "-1")
// 	}
// 	f0 = NewFromString("-1")
// 	f1 = NewFromString("-1")
//
// 	f2 = f0.Sub(f1)
// 	if f2.String() != "0" {
// 		t.Error("should be equal", f2.String(), "0")
// 	}
// 	f0 = NewFromString(".001")
// 	f1 = NewFromString(".002")
//
// 	f2 = f0.Sub(f1)
// 	if f2.String() != "-0.001" {
// 		t.Error("should be equal", f2.String(), "-0.001")
// 	}
// }

func TestOverflow(t *testing.T) {
	f0 := NewFromFloat(1.1234567)
	if f0.String() != "1.1234567" {
		t.Error("should be equal", f0.String(), "1.1234567")
	}
	f0 = NewFromFloat(1.123456789123)
	if f0.String() != "1.12345678" {
		t.Error("should be equal", f0.String(), "1.1234568")
	}
	f0 = NewFromFloat(1.0 / 3.0)
	if f0.String() != "0.33333333" {
		t.Error("should be equal", f0.String(), "0.33333333")
	}
	f0 = NewFromFloat(2.0 / 3.0)
	if f0.String() != "0.66666666" {
		t.Error("should be equal", f0.String(), "0.66666666")
	}

	assert.True(t, assert.Panics(t, func() {
		_ = NewFromString("999999999999999999")
	}))

	// add overflow
	// f0 = NewFromString("99999999999.99")
	// 184467440737 09551615
	// 00000000 0000 00000000
	f0 = NewFromUint(1<<64-1).Sub(NINE)
	f0 = f0.Add(NINE)
	t.Log(f0)
	assert.True(t, assert.Panics(t, func() {
		f0 = f0.Add(TEN)
	}))

	// sub overflow
	f0 = NewFromString("2")
	f1 := NewFromString("3")
	assert.True(t, assert.Panics(t, func() {
		_ = f0.Sub(f1)
	}))

	// mul overflow
	f0 = NewFromString("18.44674408")
	f1 = NewFromString("9999999999.99999999")
	assert.True(t, assert.Panics(t, func() {
		_ = f0.Mul(f1).String()
	}))
}

func TestNaN(t *testing.T) {
	f0 := NewFromFloat(math.NaN())
	if !f0.IsNaN() {
		t.Error("f0 should be NaN")
	}
	if f0.String() != "NaN" {
		t.Error("should be equal", f0.String(), "NaN")
	}
	f0 = NewFromString("NaN")
	if !f0.IsNaN() {
		t.Error("f0 should be NaN")
	}

	f0 = NewFromString("0.0004096")
	if f0.String() != "0.0004096" {
		t.Error("should be equal", f0.String(), "0.0004096")
	}

}

func TestIntFrac(t *testing.T) {
	f0 := NewFromFloat(1234.5678)
	if f0.UInt() != 1234 {
		t.Error("should be equal", f0.UInt(), 1234)
	}
	if f0.Frac() != .5678 {
		t.Error("should be equal", f0.Frac(), .5678)
	}
}

func TestIntOrigin(t *testing.T) {
	f0 := NewFromFloat(1234.5678)
	if f0.Original() != 123456780000 {
		t.Error("should be equal", f0.Original(), 123456780000)
	}

	f1 := NewFromOriginal(f0.Original())
	if !f1.Equal(f0) {
		t.Error("should be equal", f1.String(), f0)
	}

	f0 = NewFromFloat(.5678)
	if f0.Original() != 56780000 {
		t.Error("should be equal", f0.Original(), 56780000)
	}

	f1 = NewFromOriginal(f0.Original())
	if !f1.Equal(f0) {
		t.Error("should be equal", f1.String(), f0)
	}

}

func TestString(t *testing.T) {
	f0 := NewFromFloat(1234.5678)
	if f0.String() != "1234.5678" {
		t.Error("should be equal", f0.String(), "1234.5678")
	}
	f0 = NewFromFloat(1234.0)
	if f0.String() != "1234" {
		t.Error("should be equal", f0.String(), "1234")
	}
}

func TestStringN(t *testing.T) {
	f0 := NewFromString("1.1")
	s := f0.StringN(2)

	if s != "1.10" {
		t.Error("should be equal", s, "1.10")
	}
	f0 = NewFromString("1")
	s = f0.StringN(2)

	if s != "1.00" {
		t.Error("should be equal", s, "1.00")
	}

	f0 = NewFromString("1.123")
	s = f0.StringN(2)

	if s != "1.12" {
		t.Error("should be equal", s, "1.12")
	}
	f0 = NewFromString("1.123")
	s = f0.StringN(2)

	if s != "1.12" {
		t.Error("should be equal", s, "1.12")
	}
	f0 = NewFromString("1.123")
	s = f0.StringN(0)

	if s != "1" {
		t.Error("should be equal", s, "1")
	}
}

func TestRound(t *testing.T) {
	f0 := NewFromString("1.12345")
	f1 := f0.Round(2)

	if f1.String() != "1.12" {
		t.Error("should be equal", f1, "1.12")
	}

	f1 = f0.Round(5)

	if f1.String() != "1.12345" {
		t.Error("should be equal", f1, "1.12345")
	}
	f1 = f0.Round(4)

	if f1.String() != "1.1235" {
		t.Error("should be equal", f1, "1.1235")
	}

	// f0 = NewFromString("-1.12345")
	// f1 = f0.Round(3)
	//
	// if f1.String() != "-1.123" {
	// 	t.Error("should be equal", f1, "-1.123")
	// }
	// f1 = f0.Round(4)

	// if f1.String() != "-1.1235" {
	// 	t.Error("should be equal", f1, "-1.1235")
	// }
}

func TestEncodeDecode(t *testing.T) {
	b := &bytes.Buffer{}

	f := NewFromString("12345.12345")

	_ = f.WriteTo(b)

	f0, err := ReadFrom(b)
	if err != nil {
		t.Error(err)
	}

	if !f.Equal(f0) {
		t.Error("don't match", f, f0)
	}

	data, err := f.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	f1 := NewFromFloat(0)
	_ = f1.UnmarshalBinary(data)

	if !f.Equal(f1) {
		t.Error("don't match", f, f0)
	}
}

type JStruct struct {
	F Fixed `json:"f"`
}

func TestJSON(t *testing.T) {
	j := JStruct{}

	f := NewFromString("12345678.12345678")
	j.F = f

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	err := enc.Encode(&j)
	if err != nil {
		t.Error(err)
	}

	j.F = ZERO

	dec := json.NewDecoder(&buf)

	err = dec.Decode(&j)
	if err != nil {
		t.Error(err)
	}

	if !j.F.Equal(f) {
		t.Error("don't match", j.F, f)
	}
}
