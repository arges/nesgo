package main

import (
	"fmt"
	"math/bits"
)

type Flags struct {
	N bool // Sign (negative) flag
	V bool // Overflow flag
	B bool // BRK is being executed
	D bool // Decimal mode status
	I bool // Interrupt Enable/Disable
	Z bool // Zero flag
	C bool // Carry flag
}

type Cpu struct {
	A      uint8         // A Register
	X      uint8         // X Register
	Y      uint8         // Y Register
	SP     uint8         // Stack Pointer
	PC     uint16        // Program Counter
	Flags  Flags         // CPU Flags
	Memory [0xffff]uint8 // Memory space
}

// https://wiki.nesdev.com/w/index.php/CPU_addressing_modes
type AddressingMode int

const (
	MODE_IMP AddressingMode = iota // implicit (no address operands)
	MODE_IMM                       // immediate mode
	MODE_ZEP                       // zero page
	MODE_ZPX                       // zero page indexed by x
	MODE_ZPY                       // zero page indexed by y
	MODE_IZX                       // indexed indirect by x
	MODE_IZY                       // indexed indirect by y
	MODE_ABS                       // absolute
	MODE_ABX                       // absolute indexed by x
	MODE_ABY                       // absolute indexed by y
	MODE_IND                       // indirect
	MODE_REL                       // relative
	MODE_ACC                       // accumulator
)

// getValueByMode returns the value to be used based on addressing mode plus a
// formatting string
func (c *Cpu) getValueByMode(mode AddressingMode) (uint8, string) {

	switch mode {
	case MODE_IMM:
		value := c.Memory[c.PC+1]
		return value, "#$%02x"
	case MODE_ZEP:
		index := c.Memory[c.PC+1]
		value := c.Memory[index]
		return value, "$%02x"
	case MODE_ZPX:
		index := (c.Memory[c.PC+1] + c.X)
		value := c.Memory[index]
		return value, "%02x,X"
	case MODE_ZPY:
		index := (c.Memory[c.PC+1] + c.Y)
		value := c.Memory[index]
		return value, "%02x,Y"
	case MODE_IZX:
		return 0, "(%02x,X)" // FIXME
	case MODE_IZY:
		return 0, "(%02x,Y)" // FIXME
	case MODE_ABS:
		return 0, "a" // FIXME
	case MODE_ABX:
		return 0, "a,x" // FIXME
	case MODE_ABY:
		return 0, "a,y" // FIXME
	case MODE_IND:
		return 0, "(a)" // FIXME
	case MODE_REL:
		value := c.Memory[c.PC+1]
		return value, "%02x"
	case MODE_ACC:
		return 0, "A" // FIXME
	}

	// unknown addressing mode
	return 0, "???"
}

// Execute provides a generic framework for executing a single instruction
func (c *Cpu) Execute(inst Instruction) {

	// get value based on addressing mode
	value, format := c.getValueByMode(inst.Mode)

	// debug
	fmt.Printf(inst.Name+" "+format+" ", value)

	// execute
	inst.Execute(c, value)

	// increment PC
	c.PC = c.PC + inst.Length
}

// Print displays cpu registers and relevant info
func (c *Cpu) Print() {

	// FIXME: this could be written better
	flags := ""
	if c.Flags.N {
		flags = flags + "N"
	}
	if c.Flags.Z {
		flags = flags + "Z"
	}
	if c.Flags.C {
		flags = flags + "C"
	}
	if c.Flags.I {
		flags = flags + "I"
	}
	if c.Flags.D {
		flags = flags + "D"
	}
	if c.Flags.V {
		flags = flags + "V"
	}

	fmt.Printf("PC:%04x A:%02x X:%02x Y:%02x SP:%02x FLAGS:%s\n",
		c.PC, c.A, c.X, c.Y, c.SP, flags)
}

// ExecuteAdc - Add memory to accumulator with carry - A + M + C -> A, C
func ExecuteAdc(c *Cpu, value uint8) {
	// clobber flags
	c.Flags.N, c.Flags.Z, c.Flags.C, c.Flags.V = false, false, false, false

	initial := c.A

	// add with carry
	c.A = c.A + value
	if c.Flags.C {
		c.A = c.A + 1
	}

	// check for carry
	c.Flags.C = false
	if c.A < initial {
		c.Flags.C = true
	}

	// check for overflow


	// update Negative flag
	if (c.A & 0x80) == 0x80 {
		c.Flags.N = true
	}

	// update Zero flag
	if c.A == 0 {
		c.Flags.Z = true
	}
}

// ExecuteAnd - "AND" memory with accumulator - A & M -> A
func ExecuteAnd(c *Cpu, value uint8) {
	// clobber flags
	c.Flags.N, c.Flags.Z = false, false

	// and value
	c.A = c.A & value

	// update Negative flag
	if (c.A & 0x80) == 0x80 {
		c.Flags.N = true
	}

	// update Zero flag
	if c.A == 0 {
		c.Flags.Z = true
	}
}

func ExecuteOra(c *Cpu, value uint8) {
	// clobber flags
	c.Flags.N, c.Flags.Z = false, false

	// or value
	c.A = c.A | value

	// update Negative flag
	if (c.A & 0x80) == 0x80 {
		c.Flags.N = true
	}

	// update Zero flag
	if c.A == 0 {
		c.Flags.Z = true
	}
}

func ExecuteAsl(c *Cpu, value uint8) {
	// clobber flags
	c.Flags.N, c.Flags.Z, c.Flags.C = false, false, false

	// determine if carry should be set
	if value&0x80 == 0x80 {
		c.Flags.C = true
	}

	// rotate left by 1 and mask out low bit
	bits.RotateLeft8(value&0xfe, 1)

	// update Negative flag
	if (c.A & 0x80) == 0x80 {
		c.Flags.N = true
	}

	// update Zero flag
	if c.A == 0 {
		c.Flags.Z = true
	}
}

func ExecuteBcc(c *Cpu, value uint8) {
	// branch on C == 0
	if !c.Flags.C {
		// value needs to be signed
		//c.PC = c.PC + value
	}
}

func ExecuteBcs(c *Cpu, value uint8) {
	// branch on C == 1
	if c.Flags.C {
		// value needs to be signed
		//c.PC = c.PC + value

	}
}

type Instruction struct {
	Name    string            // 3-letter name of instruction
	Mode    AddressingMode    // Instruction mode of opcode
	Length  uint16            // Bytes of instruction
	Cycles  int               // Cycles instructions takes to execute (at best)
	Execute func(*Cpu, uint8) // Execution function
}

// http://nesdev.com/6502.txt
var instructionMap = map[uint8]Instruction{

	0x69: Instruction{"ADC", MODE_IMM, 2, 2, ExecuteAdc},
	0x65: Instruction{"ADC", MODE_ZEP, 2, 3, ExecuteAdc},
	0x75: Instruction{"ADC", MODE_ZPX, 2, 4, ExecuteAdc},
	0x60: Instruction{"ADC", MODE_ABS, 3, 4, ExecuteAdc},
	0x70: Instruction{"ADC", MODE_ABX, 3, 4, ExecuteAdc},
	0x79: Instruction{"ADC", MODE_ABY, 3, 4, ExecuteAdc},
	0x61: Instruction{"ADC", MODE_IZX, 2, 6, ExecuteAdc},
	0x71: Instruction{"ADC", MODE_IZY, 2, 5, ExecuteAdc},

	0x29: Instruction{"AND", MODE_IMM, 2, 2, ExecuteAnd},
	0x25: Instruction{"AND", MODE_ZEP, 2, 3, ExecuteAnd},
	0x35: Instruction{"AND", MODE_ZPX, 2, 4, ExecuteAnd},
	0x2D: Instruction{"AND", MODE_ABS, 3, 4, ExecuteAnd},
	0x3D: Instruction{"AND", MODE_ABX, 3, 4, ExecuteAnd},
	0x39: Instruction{"AND", MODE_ABY, 3, 4, ExecuteAnd},
	0x21: Instruction{"AND", MODE_IZX, 2, 6, ExecuteAnd},
	0x31: Instruction{"AND", MODE_IZY, 2, 5, ExecuteAnd},

	0x0a: Instruction{"ASL", MODE_ACC, 1, 2, ExecuteAsl},
	0x06: Instruction{"ASL", MODE_ZEP, 2, 5, ExecuteAsl},
	0x16: Instruction{"ASL", MODE_ZPX, 2, 6, ExecuteAsl},
	0x0e: Instruction{"ASL", MODE_ABS, 3, 6, ExecuteAsl},
	0x1e: Instruction{"ASL", MODE_ABX, 3, 7, ExecuteAsl},

	//
	0x90: Instruction{"BCC", MODE_REL, 2, 2, ExecuteBcc},
	0xb0: Instruction{"BCS", MODE_REL, 2, 2, ExecuteBcs},
	/*0xf0: Instruction{"BEQ", MODE_REL, 2, 2, ExecuteBeq},

	0x24: Instruction{"BIT", MODE_ZEP, 2, 3, ExecuteBit},
	0x2c: Instruction{"BIT", MODE_ABS, 3, 4, ExecuteBit},

	0x30: Instruction{"BMI", MODE_REL, 2, 2, ExecuteBmi},
	0xd0: Instruction{"BNE", MODE_REL, 2, 2, ExecuteBne},
	*/

	//
	0x09: Instruction{"ORA", MODE_IMM, 2, 2, ExecuteOra},
	0x05: Instruction{"ORA", MODE_ZEP, 2, 3, ExecuteOra},
	0x15: Instruction{"ORA", MODE_ZPX, 2, 4, ExecuteOra},
	0x0D: Instruction{"ORA", MODE_ABS, 3, 4, ExecuteOra},
	0x1D: Instruction{"ORA", MODE_ABX, 3, 4, ExecuteOra},
	0x19: Instruction{"ORA", MODE_ABY, 3, 4, ExecuteOra},
	0x01: Instruction{"ORA", MODE_IZX, 2, 6, ExecuteOra},
	0x11: Instruction{"ORA", MODE_IZY, 2, 5, ExecuteOra},
}

func (c *Cpu) Step() {
	inst := instructionMap[c.Memory[c.PC]]
	c.Execute(inst)
	c.Print()
}

func main() {
	c := Cpu{}
	c.Memory[0] = 0x29
	c.Memory[1] = 0x19
	c.A = 0x19

	c.Step()
}
