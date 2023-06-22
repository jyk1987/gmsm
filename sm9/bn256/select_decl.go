//go:build (amd64 && !purego) || (arm64 && !purego)
// +build amd64,!purego arm64,!purego

package bn256

import "golang.org/x/sys/cpu"

var supportAVX2 = cpu.X86.HasAVX2

// If cond is 0, sets res = b, otherwise sets res = a.
//
//go:noescape
func gfP12MovCond(res, a, b *gfP12, cond int)

// If cond is 0, sets res = b, otherwise sets res = a.
//
//go:noescape
func curvePointMovCond(res, a, b *curvePoint, cond int)

// If cond is 0, sets res = b, otherwise sets res = a.
//
//go:noescape
func twistPointMovCond(res, a, b *twistPoint, cond int)