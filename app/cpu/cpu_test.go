package cpu

import (
	"testing"
)

type MockBus struct {
	memory [65536]byte
}

func (m *MockBus) Read(addr uint16) uint8       { return m.memory[addr] }
func (m *MockBus) Write(addr uint16, val uint8) { m.memory[addr] = val }

func setup() (*CPU, *MockBus) {
	bus := &MockBus{}
	cpu := CreateCPU(bus)
	// Default state after boot is often different, but for testing we can start at 0
	cpu.PC = 0x0000
	cpu.SP = 0xFFFE
	cpu.interruptsEnabled = false
	return cpu, bus
}

func TestOpcodes(t *testing.T) {
	t.Run("LD R8, n8", func(t *testing.T) {
		cpu, bus := setup()
		bus.Write(0x0000, 0x06) // LD B, 0x42
		bus.Write(0x0001, 0x42)
		
		cpu.Step()
		if cpu.registers.B != 0x42 {
			t.Errorf("Expected B=0x42, got %02X", cpu.registers.B)
		}
		if cpu.PC != 0x0002 {
			t.Errorf("Expected PC=0x0002, got %04X", cpu.PC)
		}
	})

	t.Run("ADD A, R8", func(t *testing.T) {
		cpu, bus := setup()
		cpu.registers.A = 0x0F
		cpu.registers.B = 0x01
		bus.Write(0x0000, 0x80) // ADD A, B

		cpu.Step()
		if cpu.registers.A != 0x10 {
			t.Errorf("Expected A=0x10, got %02X", cpu.registers.A)
		}
		// Half-carry should be set (0x0F + 0x01)
		if (cpu.registers.F & 0x20) == 0 {
			t.Error("Expected Half-Carry flag to be set")
		}
	})

	t.Run("JP a16", func(t *testing.T) {
		cpu, bus := setup()
		bus.Write(0x0000, 0xC3) // JP 0x1234
		bus.Write(0x0001, 0x34)
		bus.Write(0x0002, 0x12)

		cpu.Step()
		if cpu.PC != 0x1234 {
			t.Errorf("Expected PC=0x1234, got %04X", cpu.PC)
		}
	})

	t.Run("PUSH/POP AF", func(t *testing.T) {
		cpu, bus := setup()
		cpu.registers.A = 0xDE
		cpu.registers.F = 0xF0 // All 4 flags set, lower 4 bits zero
		
		bus.Write(0x0000, 0xF5) // PUSH AF
		bus.Write(0x0001, 0xF1) // POP AF
		
		cpu.Step() // PUSH AF
		if cpu.SP != 0xFFFC {
			t.Errorf("Expected SP=0xFFFC, got %04X", cpu.SP)
		}
		
		cpu.registers.A = 0x00
		cpu.registers.F = 0x00
		
		cpu.Step() // POP AF
		if cpu.registers.A != 0xDE || cpu.registers.F != 0xF0 {
			t.Errorf("Expected AF=0xDEF0, got A=%02X F=%02X", cpu.registers.A, cpu.registers.F)
		}
		if cpu.SP != 0xFFFE {
			t.Errorf("Expected SP=0xFFFE, got %04X", cpu.SP)
		}
	})

	t.Run("Block 0: INC/DEC/LD R16", func(t *testing.T) {
		cpu, bus := setup()
		bus.Write(0x0000, 0x01) // LD BC, 0x1234
		bus.Write(0x0001, 0x34)
		bus.Write(0x0002, 0x12)
		bus.Write(0x0003, 0x03) // INC BC
		bus.Write(0x0004, 0x04) // INC B
		bus.Write(0x0005, 0x05) // DEC B
		
		cpu.Step()
		if cpu.registers.getBC() != 0x1234 {
			t.Errorf("Expected BC=0x1234, got %04X", cpu.registers.getBC())
		}
		cpu.Step()
		if cpu.registers.getBC() != 0x1235 {
			t.Errorf("Expected BC=0x1235, got %04X", cpu.registers.getBC())
		}
		cpu.Step()
		if cpu.registers.B != 0x13 {
			t.Errorf("Expected B=0x13, got %02X", cpu.registers.B)
		}
		cpu.Step()
		if cpu.registers.B != 0x12 {
			t.Errorf("Expected B=0x12, got %02X", cpu.registers.B)
		}
	})

	t.Run("Block 1: LD R8, R8", func(t *testing.T) {
		cpu, bus := setup()
		cpu.registers.B = 0x55
		bus.Write(0x0000, 0x48) // LD C, B
		
		cpu.Step()
		if cpu.registers.C != 0x55 {
			t.Errorf("Expected C=0x55, got %02X", cpu.registers.C)
		}
	})

	t.Run("Block 2: Arithmetic/Logic", func(t *testing.T) {
		cpu, bus := setup()
		cpu.registers.A = 0xFF
		cpu.registers.B = 0x01
		bus.Write(0x0000, 0x90) // SUB B
		bus.Write(0x0001, 0xA0) // AND B
		bus.Write(0x0002, 0xB0) // OR B
		bus.Write(0x0003, 0xA8) // XOR B
		bus.Write(0x0004, 0xB8) // CP B
		
		cpu.Step() // SUB B: A = 0xFE, Z=0, N=1, H=0, C=0
		if cpu.registers.A != 0xFE || (cpu.registers.F & 0x40) == 0 {
			t.Errorf("SUB failed: A=%02X F=%02X", cpu.registers.A, cpu.registers.F)
		}
		
		cpu.Step() // AND B: A = 0xFE & 0x01 = 0x00, Z=1, N=0, H=1, C=0
		if cpu.registers.A != 0x00 || (cpu.registers.F & 0x80) == 0 || (cpu.registers.F & 0x20) == 0 {
			t.Errorf("AND failed: A=%02X F=%02X", cpu.registers.A, cpu.registers.F)
		}

		cpu.registers.A = 0x10
		cpu.Step() // OR B: A = 0x10 | 0x01 = 0x11, Z=0, N=0, H=0, C=0
		if cpu.registers.A != 0x11 || (cpu.registers.F & 0xF0) != 0 {
			t.Errorf("OR failed: A=%02X F=%02X", cpu.registers.A, cpu.registers.F)
		}

		cpu.Step() // XOR B: A = 0x11 ^ 0x01 = 0x10, Z=0, N=0, H=0, C=0
		if cpu.registers.A != 0x10 {
			t.Errorf("XOR failed: A=%02X", cpu.registers.A)
		}

		cpu.Step() // CP B: A=0x10, B=0x01. 0x10 - 0x01 = 0x0F. Z=0, N=1, H=1, C=0
		if cpu.registers.A != 0x10 || (cpu.registers.F & 0x40) == 0 {
			t.Errorf("CP failed: A=%02X F=%02X", cpu.registers.A, cpu.registers.F)
		}
	})

	t.Run("Block 3: CALL/RET", func(t *testing.T) {
		cpu, bus := setup()
		bus.Write(0x0000, 0xCD) // CALL 0x1000
		bus.Write(0x0001, 0x00)
		bus.Write(0x0002, 0x10)
		bus.Write(0x1000, 0xC9) // RET
		
		cpu.Step()
		if cpu.PC != 0x1000 || cpu.SP != 0xFFFC {
			t.Errorf("CALL failed: PC=%04X SP=%04X", cpu.PC, cpu.SP)
		}
		
		cpu.Step()
		if cpu.PC != 0x0003 || cpu.SP != 0xFFFE {
			t.Errorf("RET failed: PC=%04X SP=%04X", cpu.PC, cpu.SP)
		}
	})

	t.Run("CB Opcodes", func(t *testing.T) {
		cpu, bus := setup()
		cpu.registers.B = 0x01
		bus.Write(0x0000, 0xCB) 
		bus.Write(0x0001, 0x40) // BIT 0, B
		bus.Write(0x0002, 0xCB)
		bus.Write(0x0003, 0x80) // RES 0, B
		bus.Write(0x0004, 0xCB)
		bus.Write(0x0005, 0xC0) // SET 0, B
		bus.Write(0x0006, 0xCB)
		bus.Write(0x0007, 0x38) // SRL B
		
		cpu.Step() // BIT 0, B (B=1, bit 0 is set, so Z=0)
		if (cpu.registers.F & 0x80) != 0 {
			t.Errorf("BIT 0, B failed: F=%02X", cpu.registers.F)
		}
		
		cpu.Step() // RES 0, B (B becomes 0)
		if cpu.registers.B != 0x00 {
			t.Errorf("RES 0, B failed: B=%02X", cpu.registers.B)
		}

		cpu.Step() // SET 0, B (B becomes 1)
		if cpu.registers.B != 0x01 {
			t.Errorf("SET 0, B failed: B=%02X", cpu.registers.B)
		}

		cpu.Step() // SRL B (B becomes 0, C=1, Z=1)
		if cpu.registers.B != 0x00 || (cpu.registers.F & 0x80) == 0 || (cpu.registers.F & 0x10) == 0 {
			t.Errorf("SRL B failed: B=%02X F=%02X", cpu.registers.B, cpu.registers.F)
		}
	})

	t.Run("Special Opcodes: DAA, CPL, SCF, CCF", func(t *testing.T) {
		cpu, bus := setup()
		// DAA
		cpu.registers.A = 0x09
		bus.Write(0x0000, 0xC6) // ADD A, 1
		bus.Write(0x0001, 0x01)
		bus.Write(0x0002, 0x27) // DAA
		cpu.Step() // A=0x0A
		cpu.Step() // A=0x10
		if cpu.registers.A != 0x10 {
			t.Errorf("DAA failed: A=%02X", cpu.registers.A)
		}

		// CPL
		cpu.registers.A = 0x55 // 0101 0101
		bus.Write(0x0003, 0x2F) // CPL
		cpu.Step() // A=0xAA
		if cpu.registers.A != 0xAA {
			t.Errorf("CPL failed: A=%02X", cpu.registers.A)
		}

		// SCF / CCF
		bus.Write(0x0004, 0x37) // SCF
		bus.Write(0x0005, 0x3F) // CCF
		cpu.Step() // C=1
		if (cpu.registers.F & 0x10) == 0 {
			t.Errorf("SCF failed: F=%02X", cpu.registers.F)
		}
		cpu.Step() // C=0
		if (cpu.registers.F & 0x10) != 0 {
			t.Errorf("CCF failed: F=%02X", cpu.registers.F)
		}
	})

	t.Run("Rotates: RLCA, RRCA, RLA, RRA", func(t *testing.T) {
		cpu, bus := setup()
		cpu.registers.A = 0x80
		bus.Write(0x0000, 0x07) // RLCA
		cpu.Step() // A=0x01, C=1
		if cpu.registers.A != 0x01 || (cpu.registers.F & 0x10) == 0 {
			t.Errorf("RLCA failed: A=%02X F=%02X", cpu.registers.A, cpu.registers.F)
		}

		bus.Write(0x0001, 0x0F) // RRCA
		cpu.Step() // A=0x80, C=1
		if cpu.registers.A != 0x80 || (cpu.registers.F & 0x10) == 0 {
			t.Errorf("RRCA failed: A=%02X F=%02X", cpu.registers.A, cpu.registers.F)
		}

		cpu.registers.A = 0x80
		cpu.registers.F = 0x00 // Carry=0
		bus.Write(0x0002, 0x17) // RLA
		cpu.Step() // A=0x00, C=1
		if cpu.registers.A != 0x00 || (cpu.registers.F & 0x10) == 0 {
			t.Errorf("RLA failed: A=%02X F=%02X", cpu.registers.A, cpu.registers.F)
		}

		cpu.registers.A = 0x01
		cpu.registers.F = 0x10 // Carry=1
		bus.Write(0x0003, 0x1F) // RRA
		cpu.Step() // A=0x80, C=1
		if cpu.registers.A != 0x80 || (cpu.registers.F & 0x10) == 0 {
			t.Errorf("RRA failed: A=%02X F=%02X", cpu.registers.A, cpu.registers.F)
		}
	})
}

func TestInterrupts(t *testing.T) {
	t.Run("VBlank Interrupt", func(t *testing.T) {
		cpu, bus := setup()
		cpu.interruptsEnabled = true
		bus.Write(0xFFFF, 0x01) // Enable VBlank interrupt
		bus.Write(0xFF0F, 0x01) // Request VBlank interrupt
		
		cpu.PC = 0x1234
		cpu.Step()
		
		if cpu.PC != 0x0040 {
			t.Errorf("Expected PC=0x0040 (VBlank handler), got %04X", cpu.PC)
		}
		if cpu.interruptsEnabled {
			t.Error("IME should be disabled after servicing interrupt")
		}
		// Return address should be on stack
		low := bus.Read(cpu.SP)
		high := bus.Read(cpu.SP + 1)
		if (uint16(high) << 8 | uint16(low)) != 0x1234 {
			t.Errorf("Expected 0x1234 on stack, got %04X", uint16(high) << 8 | uint16(low))
		}
	})
}

func TestHalt(t *testing.T) {
	cpu, bus := setup()
	bus.Write(0x0000, 0x76) // HALT
	
	cpu.Step()
	if !cpu.Halted {
		t.Error("CPU should be halted")
	}
	
	// Should stay halted
	cycles := cpu.Step()
	if !cpu.Halted || cycles != 4 {
		t.Error("CPU should remain halted and return 4 cycles")
	}
	
	// Wake up with interrupt
	bus.Write(0xFFFF, 0x01)
	bus.Write(0xFF0F, 0x01)
	
	cpu.Step()
	if cpu.Halted {
		t.Error("CPU should have woken from halt")
	}
}
