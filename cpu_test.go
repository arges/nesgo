package main

import (
	"testing"
)

func TestExecuteAdc(t *testing.T) {
	c := Cpu{}

	// Add 0x10 to 0x10
	c.A = 0x10
	ExecuteAdc(&c, 0x10)
	c.Print()
	if c.A != 0x20 {
		t.Fail()
	}

	// Test Zero
	c.A = 0x0
	ExecuteAdc(&c, 0x0)
	c.Print()
	if !c.Flags.Z {
		t.Fail()
	}
}
