package main

import (
	"testing"
)

func TestExecuteAdc(t *testing.T) {
	c := Cpu{}

	// Add 0x10 to 0x10 with carry
	c.A = 0x10
	c.Flags.C = true
	ExecuteAdc(&c, 0x10)
	c.Print()
	if c.A != 0x21 {
		t.Fail()
	}

	// Test Zero
	c.A = 0x0
	c.Flags.C = false
	ExecuteAdc(&c, 0x0)
	c.Print()
	if !c.Flags.Z {
		t.Fail()
	}

	// Test Overflow
	c.A = 0xff
	c.Flags.C = true
	ExecuteAdc(&c, 0x0)
	c.Print()
	if !c.Flags.V || !c.Flags.C || !c.Flags.Z {
		t.Fail()
	}
}

func TestExecuteAnd(t *testing.T) {
	c := Cpu{}

	// Test simple and
	c.A = 0xff
	ExecuteAnd(&c, 0xff)
	c.Print()
	if c.A != 0xff || !c.Flags.N {
		t.Fail()
	}

	// Test and resulting in zero
	c.A = 0xff
	ExecuteAnd(&c, 0x00)
	c.Print()
	if c.A != 0x00 || !c.Flags.Z {
		t.Fail()
	}
}

func TestExecuteOra(t *testing.T) {
	c := Cpu{}

	// Test simple or
	c.A = 0x00
	ExecuteOra(&c, 0xff)
	c.Print()
	if c.A != 0xff || !c.Flags.N {
		t.Fail()
	}

	// Test or resulting in zero
	c.A = 0x00
	ExecuteOra(&c, 0x00)
	c.Print()
	if c.A != 0x00 || !c.Flags.Z {
		t.Fail()
	}
}

func TestExecuteCl(t *testing.T) {
	c := Cpu{}
	c.Flags.C, c.Flags.D, c.Flags.I, c.Flags.V = true, true, true, true
	ExecuteClc(&c, 0)
	ExecuteCld(&c, 0)
	ExecuteCli(&c, 0)
	ExecuteClv(&c, 0)
	c.Print()
	if c.Flags.C || c.Flags.D || c.Flags.I || c.Flags.V {
		t.Fail()
	}
}
