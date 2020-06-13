package rig

import (
	"fmt"
	"io"
	"testing"
)

var lines = []string{
	"AT",
	"BX12 ABCDE",
	"FGH",
	"IJK",
	";",
}

type mockPort struct {
	index int
}

func (p *mockPort) Read(b []byte) (n int, err error) {
	if p.index > len(lines)-1 {
		return
	}
	s := copy(b, lines[p.index])
	p.index++
	return s, nil
}

func (p *mockPort) Write(buff []byte) (n int, err error) {
	return
}

func (p *mockPort) Close() error {
	return nil
}

func newMockPort() io.ReadWriteCloser {
	return &mockPort{
		index: 0,
	}
}

func TestRead(t *testing.T) {

	conn := &Connection{
		size:     100, // shouldn't need more than 100 byte buffer
		dataChan: make(chan []byte),
		port:     newMockPort(),
	}

	go conn.read()
	tbx := NewTBX()

	out, err := conn.SendCommand(tbx, nil)
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()

	fmt.Printf("%s\n", out)

}
