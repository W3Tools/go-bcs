package bcs_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/fardream/go-bcs/bcs"
)

func TestSetBigIntToUint256(t *testing.T) {
	tests := []struct {
		Name         string
		StringNumber string
		BCSBytes     []byte
	}{
		{
			Name:         "Set 0 to uint256",
			StringNumber: "0",
			BCSBytes:     []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			Name:         "Set 16 to uint256",
			StringNumber: "16",
			BCSBytes:     []byte{16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			Name:         "Set 256 to uint256",
			StringNumber: "256",
			BCSBytes:     []byte{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			Name:         "Set 300 to uint256",
			StringNumber: "300",
			BCSBytes:     []byte{44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			Name:         "Set 177279138216529409561387389144142902470 to uint256",
			StringNumber: "177279138216529409561387389144142902470",
			BCSBytes:     []byte{198, 244, 62, 55, 220, 101, 84, 121, 182, 167, 101, 52, 87, 184, 94, 133, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			bigInt, ok := new(big.Int).SetString(tt.StringNumber, 10)
			if !ok {
				t.Errorf("big.Int SetString result expected true, but got %v", ok)
			}

			uint256, ok := bcs.NewUint256FromBigInt(bigInt)
			if !ok {
				t.Fatalf("unable to set uint256 from bigint")
			}

			bs, err := uint256.MarshalBCS()
			if err != nil {
				t.Fatalf("failed to marshal bcs, msg: %v", err)
			}
			if len(bs) != 32 {
				t.Errorf("bcs bytes length expect 32, but got %v", len(bs))
			}
			if !bytes.Equal(bs, tt.BCSBytes) {
				t.Errorf("bcs bytes expected %v, but got %v", tt.BCSBytes, bs)
			}

			newUint256 := new(bcs.Uint256)
			_, err = bcs.Unmarshal(tt.BCSBytes, &newUint256)
			if err != nil {
				t.Fatalf("failed to unmarshal bcs, msg: %v", err)
			}

			if newUint256.String() != uint256.String() {
				t.Errorf("new uint256 to string expected %s, but got %s", uint256.String(), newUint256.String())
			}
			if newUint256.String() != tt.StringNumber {
				t.Errorf("new uint256 to string expected %s, but got %s", tt.StringNumber, newUint256.String())
			}
		})
	}

	tests2 := []struct {
		Name              string
		StringNumberArray []string
		BCSBytes          []byte
	}{
		{
			Name:              "Set array 0 to uint256",
			StringNumberArray: []string{"0", "1"},
			BCSBytes:          []byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			Name:              "Set array 16 to uint256",
			StringNumberArray: []string{"16", "17"},
			BCSBytes:          []byte{2, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 17, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			Name:              "Set array 256 to uint256",
			StringNumberArray: []string{"256", "257"},
			BCSBytes:          []byte{2, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			Name:              "Set array 300 to uint256",
			StringNumberArray: []string{"300"},
			BCSBytes:          []byte{1, 44, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			Name:              "Set array 177279138216529409561387389144142902470 to uint256",
			StringNumberArray: []string{"177279138216529409561387389144142902470", "177279138216529409561387389144142902470"},
			BCSBytes:          []byte{2, 198, 244, 62, 55, 220, 101, 84, 121, 182, 167, 101, 52, 87, 184, 94, 133, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 198, 244, 62, 55, 220, 101, 84, 121, 182, 167, 101, 52, 87, 184, 94, 133, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests2 {
		t.Run(tt.Name, func(t *testing.T) {
			u256Array := []*bcs.Uint256{}
			for _, n := range tt.StringNumberArray {
				u, ok := bcs.NewUint256(n)
				if !ok {
					t.Fatalf("new uint256 result expected true, but got %v", ok)
				}
				u256Array = append(u256Array, u)
			}

			bs, err := bcs.Marshal(u256Array)
			if err != nil {
				t.Fatalf("failed to marshal bcs, msg: %v", err)
			}
			if !bytes.Equal(bs, tt.BCSBytes) {
				t.Errorf("bcs bytes expected %v, but got %v", tt.BCSBytes, bs)
			}

			newUint256 := new([]*bcs.Uint256)
			_, err = bcs.Unmarshal(tt.BCSBytes, &newUint256)
			if err != nil {
				t.Fatalf("failed to unmarshal bcs, msg: %v", err)
			}
		})
	}
}
