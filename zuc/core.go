package zuc

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

// constant D for ZUC-128
var kd = [16]uint32{
	0x44D7, 0x26BC, 0x626B, 0x135E, 0x5789, 0x35E2, 0x7135, 0x09AF,
	0x4D78, 0x2F13, 0x6BC4, 0x1AF1, 0x5E26, 0x3C4D, 0x789A, 0x47AC,
}

var sbox0 = [256]byte{
	0x3e, 0x72, 0x5b, 0x47, 0xca, 0xe0, 0x00, 0x33, 0x04, 0xd1, 0x54, 0x98, 0x09, 0xb9, 0x6d, 0xcb,
	0x7b, 0x1b, 0xf9, 0x32, 0xaf, 0x9d, 0x6a, 0xa5, 0xb8, 0x2d, 0xfc, 0x1d, 0x08, 0x53, 0x03, 0x90,
	0x4d, 0x4e, 0x84, 0x99, 0xe4, 0xce, 0xd9, 0x91, 0xdd, 0xb6, 0x85, 0x48, 0x8b, 0x29, 0x6e, 0xac,
	0xcd, 0xc1, 0xf8, 0x1e, 0x73, 0x43, 0x69, 0xc6, 0xb5, 0xbd, 0xfd, 0x39, 0x63, 0x20, 0xd4, 0x38,
	0x76, 0x7d, 0xb2, 0xa7, 0xcf, 0xed, 0x57, 0xc5, 0xf3, 0x2c, 0xbb, 0x14, 0x21, 0x06, 0x55, 0x9b,
	0xe3, 0xef, 0x5e, 0x31, 0x4f, 0x7f, 0x5a, 0xa4, 0x0d, 0x82, 0x51, 0x49, 0x5f, 0xba, 0x58, 0x1c,
	0x4a, 0x16, 0xd5, 0x17, 0xa8, 0x92, 0x24, 0x1f, 0x8c, 0xff, 0xd8, 0xae, 0x2e, 0x01, 0xd3, 0xad,
	0x3b, 0x4b, 0xda, 0x46, 0xeb, 0xc9, 0xde, 0x9a, 0x8f, 0x87, 0xd7, 0x3a, 0x80, 0x6f, 0x2f, 0xc8,
	0xb1, 0xb4, 0x37, 0xf7, 0x0a, 0x22, 0x13, 0x28, 0x7c, 0xcc, 0x3c, 0x89, 0xc7, 0xc3, 0x96, 0x56,
	0x07, 0xbf, 0x7e, 0xf0, 0x0b, 0x2b, 0x97, 0x52, 0x35, 0x41, 0x79, 0x61, 0xa6, 0x4c, 0x10, 0xfe,
	0xbc, 0x26, 0x95, 0x88, 0x8a, 0xb0, 0xa3, 0xfb, 0xc0, 0x18, 0x94, 0xf2, 0xe1, 0xe5, 0xe9, 0x5d,
	0xd0, 0xdc, 0x11, 0x66, 0x64, 0x5c, 0xec, 0x59, 0x42, 0x75, 0x12, 0xf5, 0x74, 0x9c, 0xaa, 0x23,
	0x0e, 0x86, 0xab, 0xbe, 0x2a, 0x02, 0xe7, 0x67, 0xe6, 0x44, 0xa2, 0x6c, 0xc2, 0x93, 0x9f, 0xf1,
	0xf6, 0xfa, 0x36, 0xd2, 0x50, 0x68, 0x9e, 0x62, 0x71, 0x15, 0x3d, 0xd6, 0x40, 0xc4, 0xe2, 0x0f,
	0x8e, 0x83, 0x77, 0x6b, 0x25, 0x05, 0x3f, 0x0c, 0x30, 0xea, 0x70, 0xb7, 0xa1, 0xe8, 0xa9, 0x65,
	0x8d, 0x27, 0x1a, 0xdb, 0x81, 0xb3, 0xa0, 0xf4, 0x45, 0x7a, 0x19, 0xdf, 0xee, 0x78, 0x34, 0x60,
}

var sbox1 = [256]byte{
	0x55, 0xc2, 0x63, 0x71, 0x3b, 0xc8, 0x47, 0x86, 0x9f, 0x3c, 0xda, 0x5b, 0x29, 0xaa, 0xfd, 0x77,
	0x8c, 0xc5, 0x94, 0x0c, 0xa6, 0x1a, 0x13, 0x00, 0xe3, 0xa8, 0x16, 0x72, 0x40, 0xf9, 0xf8, 0x42,
	0x44, 0x26, 0x68, 0x96, 0x81, 0xd9, 0x45, 0x3e, 0x10, 0x76, 0xc6, 0xa7, 0x8b, 0x39, 0x43, 0xe1,
	0x3a, 0xb5, 0x56, 0x2a, 0xc0, 0x6d, 0xb3, 0x05, 0x22, 0x66, 0xbf, 0xdc, 0x0b, 0xfa, 0x62, 0x48,
	0xdd, 0x20, 0x11, 0x06, 0x36, 0xc9, 0xc1, 0xcf, 0xf6, 0x27, 0x52, 0xbb, 0x69, 0xf5, 0xd4, 0x87,
	0x7f, 0x84, 0x4c, 0xd2, 0x9c, 0x57, 0xa4, 0xbc, 0x4f, 0x9a, 0xdf, 0xfe, 0xd6, 0x8d, 0x7a, 0xeb,
	0x2b, 0x53, 0xd8, 0x5c, 0xa1, 0x14, 0x17, 0xfb, 0x23, 0xd5, 0x7d, 0x30, 0x67, 0x73, 0x08, 0x09,
	0xee, 0xb7, 0x70, 0x3f, 0x61, 0xb2, 0x19, 0x8e, 0x4e, 0xe5, 0x4b, 0x93, 0x8f, 0x5d, 0xdb, 0xa9,
	0xad, 0xf1, 0xae, 0x2e, 0xcb, 0x0d, 0xfc, 0xf4, 0x2d, 0x46, 0x6e, 0x1d, 0x97, 0xe8, 0xd1, 0xe9,
	0x4d, 0x37, 0xa5, 0x75, 0x5e, 0x83, 0x9e, 0xab, 0x82, 0x9d, 0xb9, 0x1c, 0xe0, 0xcd, 0x49, 0x89,
	0x01, 0xb6, 0xbd, 0x58, 0x24, 0xa2, 0x5f, 0x38, 0x78, 0x99, 0x15, 0x90, 0x50, 0xb8, 0x95, 0xe4,
	0xd0, 0x91, 0xc7, 0xce, 0xed, 0x0f, 0xb4, 0x6f, 0xa0, 0xcc, 0xf0, 0x02, 0x4a, 0x79, 0xc3, 0xde,
	0xa3, 0xef, 0xea, 0x51, 0xe6, 0x6b, 0x18, 0xec, 0x1b, 0x2c, 0x80, 0xf7, 0x74, 0xe7, 0xff, 0x21,
	0x5a, 0x6a, 0x54, 0x1e, 0x41, 0x31, 0x92, 0x35, 0xc4, 0x33, 0x07, 0x0a, 0xba, 0x7e, 0x0e, 0x34,
	0x88, 0xb1, 0x98, 0x7c, 0xf3, 0x3d, 0x60, 0x6c, 0x7b, 0xca, 0xd3, 0x1f, 0x32, 0x65, 0x04, 0x28,
	0x64, 0xbe, 0x85, 0x9b, 0x2f, 0x59, 0x8a, 0xd7, 0xb0, 0x25, 0xac, 0xaf, 0x12, 0x03, 0xe2, 0xf2,
}

// constant D-0 for ZUC-256
var zuc256_d0 = [16]byte{
	0x22, 0x2F, 0x24, 0x2A, 0x6D, 0x40, 0x40, 0x40,
	0x40, 0x40, 0x40, 0x40, 0x40, 0x52, 0x10, 0x30,
}

type zucState32 struct {
	lfsr [16]uint32 // linear feedback shift register
	r1   uint32
	r2   uint32
}

func (s *zucState32) bitReconstruction() []uint32 {
	result := make([]uint32, 4)
	result[0] = ((s.lfsr[15] & 0x7FFF8000) << 1) | (s.lfsr[14] & 0xFFFF)
	result[1] = ((s.lfsr[11] & 0xFFFF) << 16) | (s.lfsr[9] >> 15)
	result[2] = ((s.lfsr[7] & 0xFFFF) << 16) | (s.lfsr[5] >> 15)
	result[3] = ((s.lfsr[2] & 0xFFFF) << 16) | (s.lfsr[0] >> 15)
	return result
}

func l1(x uint32) uint32 {
	return x ^ bits.RotateLeft32(x, 2) ^ bits.RotateLeft32(x, 10) ^ bits.RotateLeft32(x, 18) ^ bits.RotateLeft32(x, 24)
}

func l2(x uint32) uint32 {
	return x ^ bits.RotateLeft32(x, 8) ^ bits.RotateLeft32(x, 14) ^ bits.RotateLeft32(x, 22) ^ bits.RotateLeft32(x, 30)
}

func (s *zucState32) f32(x0, x1, x2 uint32) uint32 {
	w := s.r1 ^ x0 + s.r2
	w1 := s.r1 + x1
	w2 := s.r2 ^ x2
	u := l1((w1 << 16) | (w2 >> 16))
	v := l2((w2 << 16) | (w1 >> 16))
	s.r1 = binary.BigEndian.Uint32([]byte{sbox0[u>>24], sbox1[(u>>16)&0xFF], sbox0[(u>>8)&0xFF], sbox1[u&0xFF]})
	s.r2 = binary.BigEndian.Uint32([]byte{sbox0[v>>24], sbox1[(v>>16)&0xFF], sbox0[(v>>8)&0xFF], sbox1[v&0xFF]})
	return w
}

func rotateLeft31(x uint32, k int) uint32 {
	return (x<<k | x>>(31-k)) & 0x7FFFFFFF
}

func add31(x, y uint32) uint32 {
	resut := x + y
	return (resut & 0x7FFFFFFF) + (resut >> 31)
}

func (s *zucState32) enterInitMode(w uint32) {
	v := s.lfsr[0]
	v = add31(v, rotateLeft31(s.lfsr[0], 8))
	v = add31(v, rotateLeft31(s.lfsr[4], 20))
	v = add31(v, rotateLeft31(s.lfsr[10], 21))
	v = add31(v, rotateLeft31(s.lfsr[13], 17))
	v = add31(v, rotateLeft31(s.lfsr[15], 15))
	v = add31(v, w)
	if v == 0 {
		v = 0x7FFFFFFF
	}
	for i := 0; i < 15; i++ {
		s.lfsr[i] = s.lfsr[i+1]
	}
	s.lfsr[15] = v
}

func (s *zucState32) enterWorkMode() {
	s.enterInitMode(0)
}

func makeFieldValue3(a, b, c uint32) uint32 {
	return (a << 23) | (b << 8) | c
}

func makeFieldValue4(a, b, c, d uint32) uint32 {
	return (a << 23) | (b << 16) | (c << 8) | d
}

func (s *zucState32) loadKeyIV16(key, iv []byte) {
	for i := 0; i < 16; i++ {
		s.lfsr[i] = makeFieldValue3(uint32(key[i]), kd[i], uint32(iv[i]))
	}
}

func (s *zucState32) loadKeyIV32(key, iv, d []byte) {
	iv17 := iv[17] >> 2
	iv18 := ((iv[17] & 0x3) << 4) | (iv[18] >> 4)
	iv19 := ((iv[18] & 0xf) << 2) | (iv[19] >> 6)
	iv20 := iv[19] & 0x3f
	iv21 := iv[20] >> 2
	iv22 := ((iv[20] & 0x3) << 4) | (iv[21] >> 4)
	iv23 := ((iv[21] & 0xf) << 2) | (iv[22] >> 6)
	iv24 := iv[22] & 0x3f
	s.lfsr[0] = makeFieldValue4(uint32(key[0]), uint32(d[0]), uint32(key[21]), uint32(key[16]))
	s.lfsr[1] = makeFieldValue4(uint32(key[1]), uint32(d[1]), uint32(key[22]), uint32(key[17]))
	s.lfsr[2] = makeFieldValue4(uint32(key[2]), uint32(d[2]), uint32(key[23]), uint32(key[18]))
	s.lfsr[3] = makeFieldValue4(uint32(key[3]), uint32(d[3]), uint32(key[24]), uint32(key[19]))
	s.lfsr[4] = makeFieldValue4(uint32(key[4]), uint32(d[4]), uint32(key[25]), uint32(key[20]))
	s.lfsr[5] = makeFieldValue4(uint32(iv[0]), uint32(d[5]|iv17), uint32(key[5]), uint32(key[26]))
	s.lfsr[6] = makeFieldValue4(uint32(iv[1]), uint32(d[6]|iv18), uint32(key[6]), uint32(key[27]))
	s.lfsr[7] = makeFieldValue4(uint32(iv[10]), uint32(d[7]|iv19), uint32(key[7]), uint32(iv[2]))
	s.lfsr[8] = makeFieldValue4(uint32(key[8]), uint32(d[8]|iv20), uint32(iv[3]), uint32(iv[11]))
	s.lfsr[9] = makeFieldValue4(uint32(key[9]), uint32(d[9]|iv21), uint32(iv[12]), uint32(iv[4]))
	s.lfsr[10] = makeFieldValue4(uint32(iv[5]), uint32(d[10]|iv22), uint32(key[10]), uint32(key[28]))
	s.lfsr[11] = makeFieldValue4(uint32(key[11]), uint32(d[11]|iv23), uint32(iv[6]), uint32(iv[13]))
	s.lfsr[12] = makeFieldValue4(uint32(key[12]), uint32(d[12]|iv24), uint32(iv[7]), uint32(iv[14]))
	s.lfsr[13] = makeFieldValue4(uint32(key[13]), uint32(d[13]), uint32(iv[15]), uint32(iv[8]))
	s.lfsr[14] = makeFieldValue4(uint32(key[14]), uint32(d[14]|(key[31]>>4)), uint32(iv[16]), uint32(iv[9]))
	s.lfsr[15] = makeFieldValue4(uint32(key[15]), uint32(d[15]|(key[31]&0x0f)), uint32(key[30]), uint32(key[29]))
}

func newZUCState(key, iv []byte) (*zucState32, error) {
	k := len(key)
	ivLen := len(iv)
	state := &zucState32{}
	switch k {
	default:
		return nil, fmt.Errorf("zuc: invalid key size %d, we support 16/32 now", k)
	case 16: // ZUC-128
		if ivLen != 16 {
			return nil, fmt.Errorf("zuc: invalid iv size %d, expect 16 in bytes", ivLen)
		}
		state.loadKeyIV16(key, iv)
	case 32: // ZUC-256
		if ivLen != 23 {
			return nil, fmt.Errorf("zuc: invalid iv size %d, expect 23 in bytes", ivLen)
		}
		state.loadKeyIV32(key, iv, zuc256_d0[:])
	}

	// initialization
	for i := 0; i < 32; i++ {
		x := state.bitReconstruction()
		w := state.f32(x[0], x[1], x[2])
		state.enterInitMode(w >> 1)

	}

	// work state
	x := state.bitReconstruction()
	state.f32(x[0], x[1], x[2])
	state.enterWorkMode()
	return state, nil
}

func (s *zucState32) genKeyword() uint32 {
	x := s.bitReconstruction()
	z := x[3] ^ s.f32(x[0], x[1], x[2])
	s.enterWorkMode()
	return z
}

func (s *zucState32) genKeywords(words []uint32) {
	if len(words) == 0 {
		return
	}
	for i := 0; i < len(words); i++ {
		words[i] = s.genKeyword()
	}
}
