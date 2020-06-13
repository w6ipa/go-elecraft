package rig

import (
	"io"
	"log"

	"github.com/tarm/serial"
)

type Connection struct {
	name     string
	baud     int
	size     int
	closed   bool
	dataChan chan []byte
	port     io.ReadWriteCloser
}

func New(port string, speed int) *Connection {
	return &Connection{
		name:     port,
		baud:     speed,
		size:     100, // shouldn't need more than 100 byte buffer
		dataChan: make(chan []byte),
	}
}

func (c *Connection) Open() error {

	s, err := serial.OpenPort(&serial.Config{
		Name: c.name,
		Baud: c.baud,
	})
	if err != nil {
		return err
	}
	c.port = s

	go c.read()

	return nil
}

func (c *Connection) Close() error {
	c.closed = true
	return c.port.Close()
}

func (c *Connection) read() {
	for {
		tmp := make([]byte, c.size)
		if c.closed {
			break
		}
		length, err := c.port.Read(tmp)
		if err != nil {
			log.Print(err)
			break
		}
		if length > 0 {
			c.dataChan <- tmp[0:length]
		}
	}
	close(c.dataChan)
}
