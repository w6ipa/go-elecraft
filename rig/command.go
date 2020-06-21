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

type ReqResponseParser interface {
	ReqResponder
	Parse(buffer []byte) interface{}
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
			return out, nil
		case <-timeout.C:
			if offset > 0 {
				return out, nil
			}
			return buffer, io.EOF
		}

	}

	return out, fmt.Errorf("data channel closed")
}

type om struct {
	parser *regexp.Regexp
}

func NewOM() ReqResponseParser {
	return &om{
		//regexp.MustCompile("OM (A|-)(P|-)(X|F|-)(S|-)(D|-)(F|-)(f|T|-)(L|B|-)(V|X|-)(R|I|-)(-|0)(-|1|2);")
		parser: regexp.MustCompile(`OM ([APXSDFfLVRTBI012-]{12});`),
	}
}

func (c *om) Request(param interface{}) ([]byte, error) {
	return []byte("OM;"), nil
}

func (c *om) Response(buffer []byte) (advance bool, resp []byte, remaining []byte) {

	indexes := c.parser.FindSubmatchIndex(buffer)
	if indexes == nil {
		// response prefix not found
		advance = true
		return
	}
	advance = false
	resp = buffer[indexes[2]:indexes[3]]
	remaining = buffer[indexes[1]:]
	return
}

func (c *om) Parse(buffer []byte) interface{} {
	re := regexp.MustCompile("([A-])([P-])([XF-])([S-])([D-])([F-])([fT-])([LB-])([VX-])([RI-])([0-][12-])")
	ix := re.FindStringSubmatch(string(buffer))
	if ix == nil {
		return nil
	}
	res := make(map[string]bool)

	if ix[1] == "A" {
		res["ATU"] = true
	}
	if ix[2] == "P" {
		res["PA"] = true
	}
	if ix[3] == "X" {
		res["RXI/O"] = true
	}
	if ix[3] == "F" {
		res["FL"] = true
	}
	if ix[4] == "S" {
		res["SUB"] = true
	}
	if ix[5] == "D" {
		res["DVR"] = true
	}
	if ix[6] == "F" {
		res["BPF"] = true
	}
	if ix[7] == "f" {
		res["bpf"] = true
	}
	if ix[7] == "T" {
		res["ATU"] = true
	}
	if ix[8] == "L" {
		res["LNA"] = true
	}
	if ix[8] == "B" {
		res["BAT"] = true
	}
	if ix[9] == "V" {
		res["SYN"] = true
	}
	if ix[9] == "X" {
		res["XVTR"] = true
	}
	if ix[10] == "R" {
		res["K3S"] = true
	}
	if ix[10] == "I" {
		res["IO"] = true
	}
	if ix[11] == "01" {
		res["KX2"] = true
	}
	if ix[11] == "02" {
		res["KX3"] = true
	}
	return res
}

func Rig(om map[string]bool) (rig string) {
	var ok bool
	if _, ok = om["K3S"]; ok {
		rig = "K3S"
	}
	if _, ok = om["KX2"]; ok {
		rig = "KX2"
	}
	if _, ok = om["KX3"]; ok {
		rig = "KX3"
	}
	return rig
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

type tb struct {
	parser *regexp.Regexp
}

func NewTB() ReqResponder {
	return &tb{
		parser: regexp.MustCompile("TB([0-9])([0-3][0-9]|40)"),
	}
}

func (c *tb) Request(param interface{}) ([]byte, error) {
	return []byte("TB;"), nil
}

func (c *tb) Response(buffer []byte) (advance bool, resp []byte, remaining []byte) {

	indexes := c.parser.FindSubmatchIndex(buffer)
	if indexes == nil {
		// response prefix not found
		advance = true
		return
	}
	// this is the length of remaining to be sent. ignore for now.
	_, err := strconv.Atoi(string(buffer[indexes[2]:indexes[3]]))
	if err != nil {
		// need to do something with error
		advance = false
		return
	}
	// length of received.
	dataLen, err := strconv.Atoi(string(buffer[indexes[4]:indexes[5]]))
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
