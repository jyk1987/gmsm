//go:build !amd64 || generic
// +build !amd64 generic

package zuc

func xorKeyStream(c *zucState32, dst, src []byte) {
	genericXorKeyStream(c, dst, src)
}
