package internal

import "testing"

func TestParseAddr(t *testing.T) {
	addr, err := ParseAddr(":10")
	t.Log(err)
	t.Log(addr)
}
