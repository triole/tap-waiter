//go:build freebsd || darwin

package util

// TODO: implement working getFileCreated, not just a dummy
func (util Util) GetFileCreated(_ string) int64 {
	return 0
}
