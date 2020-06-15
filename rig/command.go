package rig

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"
)

type Requester interface {
	Request(param interface{}) ([]byte, error)
}

type ReqResponder interface {
	Requester
	Response(buffer []byte) (advance bool, resp []byte, remaining []byte)
}

func (c *Connection) SendCommand(cmd interface{}, param interface{}) (out []byte, err error) {
	req, ok := cmd.(Requester)
	if !ok {
		return nil, fmt.Errorf("not implementing requester")
	}

	if c.port == nil {
		return nil, fmt.Errorf("port is not opened")
	}
	cmdBuf, err := req.Request(param)
	if err != nil {
		return nil, err
	}

	c.port.Write(cmdBuf)

	rsp, ok := cmd.(ReqResponder)
	if !ok {
		return
	}

	timeout := time.NewTimer(100 * time.Millisecond)

	buffer := make([]byte, 100)
	offset := 0
Loop:
	for {
		select {
		case newData, ok := <-c.dataChan:
			copy(buffer[offset:], newData)
			offset += len(newData)

			if !ok {
				break Loop
			}

			advance, resp, remaining := rsp.Response(buffer[:offset])
			if advance {
				continue
			}
			if len(resp) > 0 {
				return resp, nil
			}

			if len(remaining) > 0 {
				return nil, fmt.Errorf("unexpected characters: %s", remaining)
			}
		case <-timeout.C:
			if offset > 0 {
				return out, nil
			}
			return buffer, io.EOF
		}

	}

	return nil, fmt.Errorf("Maximum number of runs reached")
}

type tbx struct {
	parser *regexp.Regexp
}

func NewTBX() ReqResponder {
	return &tbx{
		parser: regexp.MustCompile("TBX([0-3][0-9]|40)"),
	}
}

func (c *tbx) Request(param interface{}) ([]byte, error) {
	return []byte("TBX;"), nil
}

func (c *tbx) Response(buffer []byte) (advance bool, resp []byte, remaining []byte) {

	indexes := c.parser.FindSubmatchIndex(buffer)
	if indexes == nil {
		// response prefix not found
		advance = true
		return
	}
	dataLen, err := strconv.Atoi(string(buffer[indexes[2]:indexes[3]]))
	if err != nil {
		// need to do something with error
		advance = false
		return
	}

	dataStart := indexes[1]
	dataEnd := dataStart + dataLen

	// need to wait for remaining char - even the terminating ;
	if dataEnd+1 > len(buffer) {
		advance = true
		return
	}
	resp = buffer[indexes[1]:dataEnd]
	remaining = buffer[dataEnd+1:]
	return
}

type ttx struct{}

func NewTTx() Requester {
	return &ttx{}
}

func (c *ttx) Request(param interface{}) ([]byte, error) {
	p, ok := param.(string)
	if !ok {
		return nil, fmt.Errorf("invalid parameter")
	}
	cmd := fmt.Sprintf("TT%s;", p)
	return []byte(cmd), nil
}
