package rig

import "testing"

func TestOMParse(t *testing.T) {
	str := "APX--F--VR--"
	om := NewOM()

	out := om.Parse([]byte(str))
	t.Fatalf("%+v", out)
}
