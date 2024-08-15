package bcs

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
)

// Uint256 is like `u256` in move.
type Uint256 struct {
	lo Uint128
	hi Uint128
}

var (
	_ json.Marshaler   = (*Uint256)(nil)
	_ json.Unmarshaler = (*Uint256)(nil)
	_ Marshaler        = (*Uint256)(nil)
	_ Unmarshaler      = (*Uint256)(nil)
)

var maxU256 = (&big.Int{}).Lsh(big.NewInt(1), 256)

func NewUint256(str string) (*Uint256, bool) {
	v := new(big.Int)
	b, ok := v.SetString(str, 10)
	if !ok {
		return nil, false
	}

	return NewUint256FromBigInt(b)
}

// NewUint256FromBigInt creates a new Uint256 from a big.Int.
func NewUint256FromBigInt(v *big.Int) (*Uint256, bool) {
	b := new(Uint256)
	if ok := b.SetBigInt(v); !ok {
		return nil, ok
	}

	return b, true
}

func NewUint256FromUint64(v uint64) (*Uint256, bool) {
	return NewUint256FromBigInt(new(big.Int).SetUint64(v))
}

func IsUint256(bigInt *big.Int) bool {
	return bigInt.Sign() >= 0 && bigInt.Cmp(maxU256) < 0
}

// Big converts Uint256 to *big.Int.
func (i *Uint256) Big() *big.Int {
	loBig := i.lo.Big()
	hiBig := i.hi.Big()

	// Left shift by 128 bits, equivalent to multiplying by 2^128
	hiBig.Lsh(hiBig, 128)

	// Add the high and low parts together
	return hiBig.Add(hiBig, loBig)
}

// SetBigInt sets Uint256 from big.Int.
func (i *Uint256) SetBigInt(v *big.Int) bool {
	if !IsUint256(v) {
		return false
	}

	// Convert big.Int to a fixed-size byte slice (32 bytes)
	bs := v.Bytes()
	if len(bs) > 32 {
		return false // v is larger than Uint256
	}

	// Pad with leading zeros if necessary
	r := make([]byte, 32)
	copy(r[32-len(bs):], bs)

	hi := Uint128{
		hi: binary.BigEndian.Uint64(r[0:8]),
		lo: binary.BigEndian.Uint64(r[8:16]),
	}

	lo := Uint128{
		hi: binary.BigEndian.Uint64(r[16:24]),
		lo: binary.BigEndian.Uint64(r[24:32]),
	}

	i.hi = hi
	i.lo = lo

	return true
}

func (i *Uint256) Cmp(j *Uint256) int {
	switch {
	case i.hi.Cmp(&j.hi) > 0 || (i.hi.Cmp(&j.hi) == 0 && i.lo.Cmp(&j.lo) > 0):
		return 1
	case i.hi.Cmp(&j.hi) == 0 && i.lo.Cmp(&j.lo) == 0:
		return 0
	default:
		return -1
	}
}

func (u *Uint256) String() string {
	return u.Big().String()
}

// JSON Marshaler And Unmarshaler

// MarshalJSON: converts Uint256 to JSON.
func (i *Uint256) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Big().String())
}

// UnmarshalJSON: converts JSON to Uint256.
func (i *Uint256) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	v := new(big.Int)
	if err := v.UnmarshalJSON([]byte(str)); err != nil {
		return err
	}

	if ok := i.SetBigInt(v); !ok {
		return fmt.Errorf("bcs/uint256: cannot set %v into *bcs.Uint256", v.String())
	}

	return nil
}

// BCS Marshaler And Unmarshaler

// MarshalBCS: converts Uint256 to BCS-encoded bytes
func (i *Uint256) MarshalBCS() ([]byte, error) {
	// Initialize a byte slice of length 32
	r := make([]byte, 32)

	// Use binary.LittleEndian.PutUint64 to write Uint128's low and high parts into r sequentially
	binary.LittleEndian.PutUint64(r, i.lo.lo)
	binary.LittleEndian.PutUint64(r[8:], i.lo.hi)
	binary.LittleEndian.PutUint64(r[16:], i.hi.lo)
	binary.LittleEndian.PutUint64(r[24:], i.hi.hi)

	return r, nil
}

// UnmarshalBCS: converts BCS-encoded bytes to Uint256.
func (i *Uint256) UnmarshalBCS(r io.Reader) (int, error) {
	buf := make([]byte, 32)
	n, err := r.Read(buf)
	if err != nil {
		return n, err
	}
	if n != 32 {
		return n, fmt.Errorf("failed to read 32 bytes for Uint256 (read %d bytes)", n)
	}

	// Decode little-endian byte slices into Uint128 parts
	i.lo = Uint128{
		lo: binary.LittleEndian.Uint64(buf[0:8]),
		hi: binary.LittleEndian.Uint64(buf[8:16]),
	}
	i.hi = Uint128{
		lo: binary.LittleEndian.Uint64(buf[16:24]),
		hi: binary.LittleEndian.Uint64(buf[24:32]),
	}

	return n, nil
}
