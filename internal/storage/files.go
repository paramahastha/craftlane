package storage

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

func RandomSecret() string {
	b := make([]byte, 32)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b)
}
func PIDAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	p, e := os.FindProcess(pid)
	if e != nil {
		return false
	}
	return p.Signal(os.Signal(nil)) == nil
}
