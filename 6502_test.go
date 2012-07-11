package main

import (
    "testing"
    "fmt"
)

/**********************************************************************
 *
 * Most tests derived from: https://github.com/Efseykho/6502-emulator
 *
 **********************************************************************/

var cpu Cpu

func runTest(program []Word) {
    programCounter = 0

    Ram.Init()

    for index,value := range program {
        if err := Ram.Write(index, value); err != nil {
            panic(err.Error())
        }
    }

    for programCounter < len(program) {
        cpu.Step()
    }
}

func step(times int) {
    for i := 0; i < times; i++ {
        cpu.Step()
    }
}

func dumpState() {
    fmt.Printf("X: 0x%x Y: 0x%x A: 0x%x PC: %d\n", cpu.X, cpu.Y, cpu.A, programCounter)
}

func setupZeroPageMemory() ([]Word) {
    return []Word{
        0xA9, 0x30, // LDA #$30
        0x85, 0xFA, // STA $FA
    }
}

func setupZeroPageIndexedXMemory() ([]Word) {
    return []Word{
        0xA9, 0x30, // LDA #$30
        0x85, 0xFA, // STA $FA
        0xA2, 0x03, // LDX $#03
    }
}

func setupZeroPageIndexedYMemory() ([]Word) {
    return []Word{
        0xA9, 0x30, // LDA #$30
        0x85, 0xFA, // STA $FA
        0xA0, 0x03, // LDY $#03
    }
}

func setupIndexedIndirectMemory() ([]Word) {
    return []Word{
        0xA9, 0xFA,        // LDA #$FA
		0x85, 0xDA,        // STA $DA
		0xA9, 0xEA,        // LDA #$EA
		0x85, 0xDB,        // STA $DB
		0xA2, 0x27,	       // LDX $#27
		0xA9, 0xCC,        // LDA #$CC
		0x8D, 0xFA, 0xEA,  // STA $#EAFA
    }
}

func setupIndirectIndexedMemory() ([]Word) {
    return []Word{
        0xA9, 0xFB,        // LDA #$FB
		0x85, 0xDC,        // STA $DC
		0xA9, 0xEA,        // LDA #$EA
		0x85, 0xDD,        // STA $DD
		0xA0, 0x27,	       // LDY $#27
		0xA9, 0xCD,        // LDA #$CD
		0x8D, 0x22, 0xEB,  // STA $#EB22
    }
}

func setupAbsoluteMemory() ([]Word) {
    return []Word{
        0xA9, 0xFC,        // LDA #$FC
		0x8D, 0x23, 0xEB,  // STA $#EB23
    }
}

func setupAbsoluteIndexedYMemory() ([]Word) {
    return []Word{
        0xA9, 0xFD,        // LDA #$FD
		0x8D, 0x24, 0xEB,  // STA $24,$EB
		0xA0, 0x27,       //  LDY $#27
    }
}

func setupAbsoluteIndexedXMemory() ([]Word) {
    return []Word{
        0xA9, 0xFE,        // LDA #$FD
		0x8D, 0x25, 0xEB,  // STA $25,$EB
		0xA2, 0x27,       //  LDX $#27
    }
}

func TestAdc(test *testing.T) {
    cpu.Reset()

    program := []Word{
        0xA9, 0x01, // LDA #$01
        0x69, 0x40, // ADC $#40
    }

    runTest(program)

    if cpu.A != 0x41 {
        test.Errorf("ADC: A is 0x%x not 0x41", cpu.A)
    }

    runTest(append(setupZeroPageMemory(),
                0xA9 , 0x01, // LDA #$01
                0x18,        // CLC
                0x65, 0xFA,  // ADC $#FA
                ))

    if cpu.A != 0x31 {
        test.Errorf("ADC: A is 0x%x not 0x31\n", cpu.A)
    }

    runTest(append(setupZeroPageIndexedXMemory(),
                0xA9 , 0x01, // LDA #$01
                0x18,        // CLC
                0x75, 0xF7,  // ADC $#4F,X
                ))

    if cpu.A != 0x31 {
        test.Errorf("ADC: A is 0x%x not 0x31\n", cpu.A)
    }

    runTest(append(setupIndexedIndirectMemory(),
                0xA9 , 0x01, // LDA #$01
                0x18,        // CLC
                0x61, 0xB3,  // ADC ($#B3,X)
                ))

    if cpu.A != 0xcd {
        test.Errorf("ADC: A is 0x%x not 0xcd\n", cpu.A)
    }

    runTest(append(setupIndirectIndexedMemory(),
                0xa9 , 0x01, // LDA #$01
                0x18,        // CLC
                0x71, 0xdc,  // ADC ($DC),Y
                ))

    if cpu.A != 0xce {
        test.Errorf("ADC: A is 0x%x not 0xce\n", cpu.A)
    }

    runTest(append(setupAbsoluteMemory(),
                0xa9 , 0x01, // LDA #$01
                0x18,        // CLC
                0x6d, 0x23, 0xeb,  // ADC $EB, $23
                ))

    if cpu.A != 0xfd {
        test.Errorf("ADC: A is 0x%x not 0xfd\n", cpu.A)
    }

    runTest(append(setupAbsoluteIndexedXMemory(),
                0xa9 , 0x01, // LDA #$01
                0x18,        // CLC
                0x7d, 0xfe, 0xea, // ADC $EAFE,X
                ))

    if cpu.A != 0xff {
        test.Errorf("ADC: A is 0x%x not 0xff\n", cpu.A)
    }

    runTest(append(setupAbsoluteIndexedYMemory(),
                0xa9 , 0x01, // LDA #$01
                0x18,        // CLC
                0x79, 0xfd, 0xea, // ADC $EAFD,Y
                ))

    if cpu.A != 0xfe {
        test.Errorf("ADC: A is 0x%x not 0xfe\n", cpu.A)
    }
}

func TestSbc(test *testing.T) {
    runTest([]Word{
        0xA9, 0x20,
        0x18,
        0xE9, 0x10, // SBC $#10
    })

    switch {
    case cpu.A != 0x10:
        test.Errorf("A was 0x%x, expected 0x0F\n", cpu.A)
    case !cpu.Carry:
        test.Error("Carry bit was not set")
    case cpu.Overflow:
        test.Error("Overflow bit was set")
    }

    /*
    runTest(append(setupZeroPageMemory(),
        0xA9, 0x01,
        0x18,
        0xE5, 0xFA, // SBC $#FA,
    ))

    switch {
    case cpu.A != 0x2F:
        test.Errorf("A was 0x%x, expected 0x2F\n", cpu.A)
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Overflow:
        test.Error("Overflow bit was set")
    }
    */

    runTest(append(setupZeroPageIndexedXMemory(),
        0xA9, 0x01,
        0x18,
        0xF5, 0xF7, // SBC $#4F,X
    ))

    runTest(append(setupIndexedIndirectMemory(),
        0xA9, 0x01,
        0x18,
        0xE1, 0xB3, // SBC ($#B3,X)
    ))

    runTest(append(setupIndirectIndexedMemory(),
        0xA9, 0x01,
        0x18,
        0xF1, 0xDC, // SBC ($DC),Y
    ))

    runTest(append(setupAbsoluteMemory(),
        0xA9, 0x01,
        0x18,
        0xED, 0x23, 0xEB, // SBC $EB,$23
    ))

    runTest(append(setupAbsoluteIndexedYMemory(),
        0xA9, 0x01,
        0x18,
        0xF9, 0xFD, 0xEA, // SBC $EAFD,Y
    ))

    runTest(append(setupAbsoluteIndexedXMemory(),
        0xA9, 0x01,
        0x18,
        0xFD, 0xFE, 0xEA, // SBC $EAFE,X
    ))
}

func TestLd(test *testing.T) {
    cpu.Reset()

    runTest([]Word{
        0xA9, 0x0E, // LDA #$0E
    })

    if cpu.A != 0x0E {
        test.Errorf("LDA: 0x%X is not 0x0E\n", cpu.A)
    }

    if cpu.Zero {
        test.Error("False!")
    }

    if cpu.Negative {
        test.Error("False!")
    }

    cpu.Reset()

    runTest([]Word{
        0xA9, 0x00, // LDA #$00
    })

    if cpu.A != 0x00 {
        test.Errorf("LDA: 0x%X is not 0x00\n", cpu.A)
    }

    if !cpu.Zero {
        test.Error("False!")
    }

    if cpu.Negative {
        test.Error("False!")
    }

    runTest([]Word{
        0xa9, 0xfe, // LDA #$FE
    })

    if cpu.A != 0xfe {
        test.Errorf("LDA: 0x%x is not 0xfe\n", cpu.A)
    }

    if cpu.Zero {
        test.Error("False!")
    }

    if !cpu.Negative {
        test.Error("False!")
    }

    runTest([]Word{
        0xa2, 0xfe, // LDX #$FE
    })

    if cpu.X != 0xfe {
        test.Errorf("LDX: 0x%x is not 0xfe\n", cpu.X)
    }

    if cpu.Zero {
        test.Error("False!")
    }

    if !cpu.Negative {
        test.Error("False!")
    }

    runTest([]Word{
        0xa0, 0x00, // LDY #$00
    })

    if cpu.Y != 0x00 {
        test.Errorf("LDY: 0x%x is not 0x00\n", cpu.Y)
    }

    if !cpu.Zero {
        test.Error("False!")
    }

    if cpu.Negative {
        test.Error("False!")
    }
}

func TestSta(test *testing.T) {
    cpu.Reset()

    runTest([]Word{
        0xa9, 0xfa,        // LDA #$FA
		0x85, 0xda,        // STA $DA
    })

    if *Ram[0x00da] != 0xfa {
        test.Errorf("STA: 0x%x is not 0xfa\n", *Ram[0x00da])
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0xa9, 0x12, // LDA #$12
        0x95, 0xf7, // STA $F7,X
        ))

    if *Ram[0x00fa] != 0x12 {
        test.Errorf("STA: 0x%x is not 0x12\n", *Ram[0x00fa])
    }

    cpu.Reset()

    runTest(append(setupIndexedIndirectMemory(),
        0xa9, 0x12, // LDA #$12
        0x81, 0xb3, // STA ($B3,X)
        ))

    if *Ram[0xEAFA] != 0x12 {
        test.Errorf("STA: 0x%X is not 0x12\n", *Ram[0xeafa])
    }

    runTest(append(setupIndirectIndexedMemory(),
        0xa9, 0x12, // LDA #$12
        0x91, 0xdc, // STA ($DC),Y
        ))

    if *Ram[0xEB22] != 0x12 {
        test.Errorf("STA: 0x%x is not 0x12\n", *Ram[0xeb22])
    }

    runTest(append(setupAbsoluteMemory(),
        0xa9, 0x12, // LDA #$12
        0x8d, 0x23, 0xeb, // STA $23,$EB
        ))

    if *Ram[0xEB23] != 0x12 {
        test.Errorf("STA: 0x%x is not 0x12\n", *Ram[0xeb23])
    }

    runTest(append(setupAbsoluteIndexedYMemory(),
        0xa9, 0x12, // LDA #$12
        0x99, 0xfd, 0xea, // STA $#EAFD,Y
        ))

    if *Ram[0xeb24] != 0x12 {
        test.Errorf("STA: 0x%x is not 0x12\n", *Ram[0xeb24])
    }

    runTest(append(setupAbsoluteIndexedXMemory(),
        0xa9, 0x12, // LDA #$12
        0x9d, 0xfe, 0xea, // STA $#EAFE,X
        ))

    if *Ram[0xeb25] != 0x12 {
        test.Errorf("STA: 0x%x is not 0x12\n", *Ram[0xeb25])
    }
}

func TestStx(test *testing.T) {
    runTest([]Word{
        0xa2, 0x20, // LDX #$20
        0x86, 0x35, // STX $#35
    })

    if *Ram[0x35] != 0x20 {
        test.Errorf("STX: 0x%x is not 0x20\n", *Ram[0x35])
    }

    runTest([]Word{
        0xa0, 0x02, // LDY #$02
        0x96, 0x35, // STX $#35,Y
    })

    if *Ram[0x35 + 0x02] != 0x20 {
        test.Errorf("STX: 0x%x is not 0x20\n", *Ram[0x35 + 0x02])
    }

    runTest([]Word{
        0x8e, 0x01, 0x30, // STX $01 $30
    })

    if *Ram[0x3001] != 0x20 {
        test.Errorf("STX: 0x%x is not 0x20\n", *Ram[0x3001])
    }
}

func TestSty(test *testing.T) {
    runTest([]Word{
        0xa0, 0x20, // LDY #$20
        0x84, 0x35, // STY $#35
    })

    if *Ram[0x35] != 0x20 {
        test.Errorf("STY: 0x%x is not 0x20\n", *Ram[0x35])
    }

    runTest([]Word{
        0xa2, 0x02, // LDX #$02
        0x94, 0x35, // STY $#35,X
    })

    if *Ram[0x35 + 0x02] != 0x20 {
        test.Errorf("STY: 0x%x is not 0x20\n", *Ram[0x35 + 0x02])
    }

    runTest([]Word{
        0x8c, 0x01, 0x30, // STY $01 $30
    })

    if *Ram[0x3001] != 0x20 {
        test.Errorf("STY: 0x%x is not 0x20\n", *Ram[0x3001])
    }
}

func TestJmp(test *testing.T) {
    program := append(setupAbsoluteMemory(),
        0x18, // CLC
        0x4c, 0x23, 0xeb, // JMP $#eb23
        )

    program = append(program,
        0xa9, 0xfb,        // LDA #$FB
		0x85, 0xdc,        // STA $DC
		0xa9, 0xea,        // LDA #$EA
		0x85, 0xdd,        // STA $DD
		0xa0, 0x27,	       // LDY $#27
		0xa9, 0xcd,        // LDA #$CD
		0x8d, 0x22, 0xeb,  // STA $#EB22
        0x6c, 0xdc, 0x00,
        )

    runTest(program)

    if programCounter != 0xeb23 {
        test.Errorf("JMP: Program counter was 0x%x, should be 0xeb23\n", programCounter)
    }
}

func TestRegisterInstructions(test *testing.T) {
    runTest(append(setupAbsoluteMemory(),
        0xaa,  // TAX
        ))

    if cpu.X != 0xfc {
        test.Errorf("TXX: Value was 0x%x, expected 0xfc\n", cpu.X)
    }

    if cpu.Zero {
        test.Errorf("Zero bit was set")
    }

    if !cpu.Negative {
        test.Errorf("Negative bit not set")
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0x8a, // TXA
        ))

    if cpu.A != 0x03 {
        test.Errorf("TXX: Value was 0x%x, expected 0x03\n", cpu.A)
    }

    if cpu.Zero {
        test.Errorf("Zero bit was set")
    }

    if cpu.Negative {
        test.Errorf("Negative bit not set")
    }

    runTest(append(setupAbsoluteMemory(),
        0xa8, // TAY
        ))

    if cpu.Y != 0xfc {
        test.Errorf("TXX: Value was 0x%x, expected 0xfc\n", cpu.Y)
    }

    if cpu.Zero {
        test.Errorf("Zero bit was set")
    }

    if !cpu.Negative {
        test.Errorf("Negative bit not set")
    }

    runTest(append(setupZeroPageIndexedYMemory(),
        0x98, // TYA
        0xba, // TSX
        ))

    if cpu.A != 0x03 {
        test.Errorf("TXX: Value was 0x%x, expected 0x03\n", cpu.A)
    }

    if cpu.Zero {
        test.Errorf("Zero bit was set")
    }

    if cpu.Negative {
        test.Errorf("Negative bit not set")
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0x9a, // TXS
        ))
}

func TestBranchInstructions(test *testing.T) {
    programCounter = 0

    program := []Word{
        0x38, //SEC
		0x90, 0x45,  //BCC #$45
		0x18, //CLC
		0x90, 0x45,  //BCC #$45
		//BCS******************
		0x38, //SEC
		0xB0, 0x45,  //BCS #$45
		0x18, //CLC = 9th instr
		0xB0, 0x45,  //BCS #$45
		//BEQ*****************
		0xA9, 0x00, //LDA #$00
		0xF0, 0x45,  //BEQ $#45
		//BMI*****************
		0xA9, 0xFA, //LDA #$FA
		0x30, 0x45,  //BMI $#45
		//BNE*****************
		0xA9, 0x0F, //LDA #$0F
		0xD0, 0x45,  //BNE $#45
		//BPL*****************
		0xA9, 0x0F, //LDA #$0F
		0x10, 0x45,  //BPL $#45
		//BVC*****************
		0xA9, 0xAB, //LDA #$AB
		0x50, 0x45,  //BVC $#45
		//BSV*****************
		0xA9, 0xAB, //LDA #$AB
		0x70, 0x45,  //BVS $#45
        }

    for index,value := range program {
        *Ram[index] = value
    }

    step(2)

    if programCounter != 0x3 {
        test.Errorf("Program counter was 0x%x, expected 0x3\n", programCounter)
    }

    step(2)

    if programCounter != 0x4b {
        test.Errorf("Program counter was 0x%x, expected 0x4b\n", programCounter)
    }

    programCounter = 0x06

    step(2)

    if programCounter != 0x4e {
        test.Errorf("Program counter was 0x%x, expected 0x4e\n", programCounter)
    }

    programCounter = 0x09

    step(2)

    if programCounter != 0x0c {
        test.Errorf("Program counter was 0x%x, expected 0x0c\n", programCounter)
    }

    step(2)

    if programCounter != 0x55 {
        test.Errorf("Program counter was 0x%x, expected 0x55\n", programCounter)
    }

    programCounter = 0x10

    step(2)

    if programCounter != 0x59 {
        test.Errorf("Program counter was 0x%x, expected 0x59\n", programCounter)
    }

    programCounter = 0x14

    step(2)

    if programCounter != 0x5d {
        test.Errorf("Program counter was 0x%x, expected 0x5d\n", programCounter)
    }

    programCounter = 0x18

    step(2)

    if programCounter != 0x61 {
        test.Errorf("Program counter was 0x%x, expected 0x61\n", programCounter)
    }

    programCounter = 0x1c

    step(2)

    if programCounter != 0x65 {
        test.Errorf("Program counter was 0x%x, expected 0x65\n", programCounter)
    }

    programCounter = 0x20

    step(2)

    if programCounter != 0x24 {
        test.Errorf("Program counter was 0x%x, expected 0x24\n", programCounter)
    }
}

func TestPhx(test *testing.T) {
    programCounter = 0

    program := []Word{
        0xa9, 0xcc, // LDA #$CC
        0x48,       // PHA
        0xa9, 0x1f, // LDA #$1F
        0x68,       // PLA
        0x08,       // PLP
        0x28,
        }

    for index,value := range program {
        *Ram[index] = value
    }

    step(2)

    if cpu.StackPointer != 0xfe {
        test.Errorf("StackPointer was 0x%x, expected 0xfe\n", cpu.StackPointer)
    }

    if *Ram[cpu.StackPointer] != 0xcc {
        test.Errorf("Memory was 0x%x, expected 0xcc\n", *Ram[cpu.StackPointer])
    }

    step(2)

    if cpu.StackPointer != 0xff {
        test.Errorf("StackPointer was 0x%x, expected 0xff\n", cpu.StackPointer)
    }

    if cpu.A != 0xcc {
        test.Errorf("A was 0x%x, expected 0xcc\n", cpu.A)
    }

    if cpu.Zero {
        test.Error("Zero bit was set\n")
    }

    if !cpu.Negative {
        test.Error("Negative bit was not set\n")
    }

    cpu.SetProcessorStatus(0xde)
    step(1)

    if cpu.StackPointer != 0xfe {
        test.Errorf("StackPointer was 0x%x, expected 0xfe\n", cpu.StackPointer)
    }

    if *Ram[cpu.StackPointer] != 0xde {
        test.Errorf("Memory was 0x%x, expected 0xde\n", *Ram[cpu.StackPointer])
    }

    cpu.SetProcessorStatus(0x00)
    step(1)

    if cpu.ProcessorStatus() != 0xde {
        test.Errorf("ProcessorStatus was 0x%x, expected 0xde\n", cpu.ProcessorStatus())
    }

    if cpu.StackPointer != 0xff {
        test.Errorf("StackPointer was 0x%x, expected 0xff\n", cpu.StackPointer)
    }
}

func TestCmp(test *testing.T) {
    runTest([]Word{
        0xa9, 0x10, // LDA $#10
        0xc9, 0x05, // CMP $#05
    })

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case !cpu.Carry:
        test.Error("Carry bit was not set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupZeroPageMemory(),
        0xa9, 0x10, // LDA $#10
        0x18, // CLC
        0xc5, 0xfa, // CMP $#10
    ))

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case !cpu.Negative:
        test.Error("Negative bit was not set")
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0xa9, 0x10, // LDA $#10
        0x18, // CLC
        0xd5, 0xf7, // CMP $#F7,X
    ))

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case !cpu.Negative:
        test.Error("Negative bit was not set")
    }

    runTest(append(setupIndexedIndirectMemory(),
        0xa9, 0x0c, // LDA #$0C
        0x18,
        0xc1, 0xb3, // CMP
    ))

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case !cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupIndirectIndexedMemory(),
        0xa9, 0x0d,
        0x18,
        0xd1, 0xdc, // CMP ($DC),Y
    ))

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case !cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupAbsoluteMemory(),
        0xa9, 0x0c,
        0x18,
        0xcd, 0x23, 0xeb, // CMP $EB,$23
    ))

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case !cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupAbsoluteIndexedYMemory(),
        0xa9, 0x0d,
        0x18,
        0xd9, 0xfd, 0xea, // CMP $EAFD,Y
    ))

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case !cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupAbsoluteIndexedXMemory(),
        0xa9, 0x0e,
        0x18,
        0xdd, 0xfe, 0xea, // CMP $EAFE,X
    ))

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case !cpu.Negative:
        test.Error("Negative bit was set")
    }
}

func TestCpx(test *testing.T) {
    runTest([]Word{
        0xa2, 0x07, // LDX $#07
        0xe0, 0x01, // CPX $#01
    })

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case !cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was not set")
    }

    runTest(append(setupZeroPageMemory(),
        0xa2, 0x07, // LDX $#07
        0xe4, 0xfa, // CPX $#FA
    ))

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was not set")
    case !cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupAbsoluteMemory(),
        0xa2, 0xff,       // LDX $#FF
        0xec, 0x23, 0xeb, // CPX $EB,$23
    ))

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case !cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was not set")
    }
}

func TestCpy(test *testing.T) {
    runTest([]Word{
        0xa0, 0x07, // LDY $#07
        0xc0, 0x30, // CPY $#30
    })

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case !cpu.Negative:
        test.Error("Negative bit was not set")
    }

    runTest(append(setupZeroPageMemory(),
        0xa0, 0x07, // LDY $#07
        0xc4, 0xfa, // CPY $#FA
    ))

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was not set")
    case !cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupAbsoluteMemory(),
        0xa0, 0xff,       // LDY $#FF
        0xcc, 0x23, 0xeb, // CPY $EB,$23
    ))

    switch {
    case cpu.Zero:
        test.Error("Zero bit was set")
    case !cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was not set")
    }
}

func TestAnd(test *testing.T) {
    runTest([]Word{
        0xa9, 0x2a, // LDA #$2A
        0x29, 0x0a, // AND $#0A
    })

    if cpu.A != 0x0a {
        test.Errorf("A was 0x%x, expected 0x0a\n", cpu.A)
    }

    runTest(append(setupZeroPageMemory(),
        0xa9, 0x10, // LDA #$10
        0x18, // CLC
        0x25, 0xfa, // AND #$10
    ))

    if cpu.A != 0x10 {
        test.Errorf("A was 0x%x, expected 0x10\n", cpu.A)
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0xa9, 0x10, // LDA #$10
        0x18,
        0x35, 0xf7, // AND #$F7,X
    ))

    if cpu.A != 0x10 {
        test.Errorf("A was 0x%x, expected 0x10\n", cpu.A)
    }

    runTest(append(setupIndexedIndirectMemory(),
        0xa9, 0x0c, // LDA #$0C
        0x18,
        0x21, 0xb3, // AND ($B3,X)
    ))

    if cpu.A != 0x0c {
        test.Errorf("A was 0x%x, expected 0x0c\n", cpu.A)
    }

    runTest(append(setupIndirectIndexedMemory(),
        0xa9, 0x0d, // LDA #$0D
        0x18,
        0x31, 0xdc, // AND ($DC),Y
    ))

    if cpu.A != 0x0d {
        test.Errorf("A was 0x%x, expected 0x0d\n", cpu.A)
    }

    runTest(append(setupAbsoluteMemory(),
        0xa9, 0x0c, // LDA #$0C
        0x18,
        0x2d, 0x23, 0xeb, // AND $EB,$23
    ))

    if cpu.A != 0x0c {
        test.Errorf("A was 0x%x, expected 0x0c\n", cpu.A)
    }

    runTest(append(setupAbsoluteIndexedYMemory(),
        0xa9, 0x0d, // LDA #$01
        0x18,
        0x39, 0xfd, 0xea, // AND $EAFD,Y
    ))

    if cpu.A != 0x0d {
        test.Errorf("A was 0x%x, expected 0x0d\n", cpu.A)
    }

    runTest(append(setupAbsoluteIndexedXMemory(),
        0xa9, 0x0e, // LDA #$01
        0x18,
        0x3d, 0xfe, 0xea, // AND $EAFD,Y
    ))

    if cpu.A != 0x0e {
        test.Errorf("A was 0x%x, expected 0x0e\n", cpu.A)
    }
}

func TestOra(test *testing.T) {
    runTest([]Word{
        0xa9, 0x10,
        0x09, 0x01, // ORA #$01
    })

    if cpu.A != 0x11 {
        test.Errorf("A was 0x%x, expected 0x11\n", cpu.A)
    }

    runTest(append(setupZeroPageMemory(),
        0xa9, 0x10,
        0x18,
        0x05, 0xfa, // ORA #$10
    ))

    if cpu.A != 0x30 {
        test.Errorf("A was 0x%x, expected 0x30\n", cpu.A)
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0xa9, 0x10,
        0x18,
        0x15, 0xf7, // ORA #$F7,X
    ))

    if cpu.A != 0x30 {
        test.Errorf("A was 0x%x, expected 0x30\n", cpu.A)
    }

    runTest(append(setupIndexedIndirectMemory(),
        0xa9, 0x0c,
        0x18,
        0x01, 0xb3, // ORA ($B3,X)
    ))

    if cpu.A != 0xcc {
        test.Errorf("A was 0x%x, expected 0xcc\n", cpu.A)
    }

    runTest(append(setupIndirectIndexedMemory(),
        0xa9, 0x0d,
        0x18,
        0x11, 0xdc, // ORA ($DC),Y
    ))

    if cpu.A != 0xcd {
        test.Errorf("A was 0x%x, expected 0xcd\n", cpu.A)
    }

    runTest(append(setupAbsoluteMemory(),
        0xa9, 0x0c,
        0x18,
        0x0d, 0x23, 0xeb, // ORA $EB,$23
    ))

    if cpu.A != 0xfc {
        test.Errorf("A was 0x%x, expected 0xfc\n", cpu.A)
    }

    runTest(append(setupAbsoluteIndexedYMemory(),
        0xa9, 0x0d,
        0x18,
        0x19, 0xfd, 0xea, // ORA $EAFD,Y
    ))

    if cpu.A != 0xfd {
        test.Errorf("A was 0x%x, expected 0xfd\n", cpu.A)
    }

    runTest(append(setupAbsoluteIndexedXMemory(),
        0xa9, 0x0e,
        0x18,
        0x1d, 0xfe, 0xea, // ORA $EAFE,X
    ))

    if cpu.A != 0xfe {
        test.Errorf("A was 0x%x, expected 0xfe\n", cpu.A)
    }
}

func TestEor(test *testing.T) {
    runTest([]Word{
        0xa9, 0x10,
        0x49, 0x55, // EOR #$55
    })

    if cpu.A != 0x45 {
        test.Errorf("A was 0x%x, expected 0x45\n", cpu.A)
    }

    runTest(append(setupZeroPageMemory(),
        0xa9, 0x10,
        0x18,
        0x45, 0xfa, // EOR #$10
    ))

    if cpu.A != 0x20 {
        test.Errorf("A was 0x%x, expected 0x20\n", cpu.A)
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0xa9, 0x10,
        0x18,
        0x55, 0xf7, // EOR #$F7,X
    ))

    if cpu.A != 0x20 {
        test.Errorf("A was 0x%x, expected 0x20\n", cpu.A)
    }

    runTest(append(setupIndexedIndirectMemory(),
        0xa9, 0x0c,
        0x18,
        0x41, 0xb3, // EOR ($B3,X)
    ))

    if cpu.A != 0xc0 {
        test.Errorf("A was 0x%x, expected 0xc0\n", cpu.A)
    }

    runTest(append(setupIndirectIndexedMemory(),
        0xa9, 0x0d,
        0x18,
        0x51, 0xdc, // EOR ($DC),Y
    ))

    if cpu.A != 0xc0 {
        test.Errorf("A was 0x%x, expected 0xc0\n", cpu.A)
    }

    runTest(append(setupAbsoluteMemory(),
        0xa9, 0x0c,
        0x18,
        0x4d, 0x23, 0xeb, // EOR $EB,$23
    ))

    if cpu.A != 0xf0 {
        test.Errorf("A was 0x%x, expected 0xf0\n", cpu.A)
    }

    runTest(append(setupAbsoluteIndexedYMemory(),
        0xa9, 0x0d,
        0x18,
        0x59, 0xfd, 0xea, // EOR $EAFD,Y
    ))

    if cpu.A != 0xf0 {
        test.Errorf("A was 0x%x, expected 0xf0\n", cpu.A)
    }

    runTest(append(setupAbsoluteIndexedXMemory(),
        0xa9, 0x0e,
        0x18,
        0x5d, 0xfe, 0xea, // EOR $EAFE,X
    ))

    if cpu.A != 0xf0 {
        test.Errorf("A was 0x%x, expected 0xf0\n", cpu.A)
    }
}

func TestDec(test *testing.T) {
    runTest(append(setupZeroPageMemory(),
        0xc6, 0xfa,
    ))

    switch {
    case *Ram[0x00fa] != 0x2f:
        test.Errorf("Memory at 0x00FA was 0x%x, expected 0x2F\n", *Ram[0x00fa])
    case cpu.Negative:
        test.Error("Negative bit was set")
    case cpu.Zero:
        test.Error("Zero bit was set")
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0xd6, 0xf7,
    ))

    if *Ram[0x00fa] != 0x2f {
        test.Errorf("Memory at 0x00FA was 0x%x, expected 0x2F\n", *Ram[0x00fa])
    }

    runTest(append(setupAbsoluteMemory(),
        0xce, 0x23, 0xeb,
    ))

    switch {
    case *Ram[0xeb23] != 0xfb:
        test.Errorf("Memory at 0xEB23 was 0x%x, expected 0xFB\n", *Ram[0xEB23])
    case !cpu.Negative:
        test.Error("Negative bit was not set")
    }

    runTest(append(setupAbsoluteIndexedXMemory(),
        0xde, 0xfe, 0xea,
    ))

    if *Ram[0xeb25] != 0xfd {
        test.Errorf("Memory at 0xEB25 was 0x%x, expected 0xFD\n", *Ram[0xEB25])
    }
}

func TestInc(test *testing.T) {
    runTest(append(setupZeroPageMemory(),
        0xe6, 0xfa,
    ))

    switch {
    case *Ram[0x00fa] != 0x31:
        test.Errorf("Memory at 0x00FA was 0x%x, expected 0x31\n", *Ram[0x00fa])
    case cpu.Negative:
        test.Error("Negative bit was set")
    case cpu.Zero:
        test.Error("Zero bit was set")
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0xf6, 0xf7,
    ))

    if *Ram[0x00fa] != 0x31 {
        test.Errorf("Memory at 0x00FA was 0x%x, expected 0x31\n", *Ram[0x00fa])
    }

    runTest(append(setupAbsoluteMemory(),
        0xee, 0x23, 0xeb,
    ))

    switch {
    case *Ram[0xeb23] != 0xfd:
        test.Errorf("Memory at 0xEB23 was 0x%x, expected 0xFD\n", *Ram[0xEB23])
    case !cpu.Negative:
        test.Error("Negative bit was not set")
    }

    runTest(append(setupAbsoluteIndexedXMemory(),
        0xfe, 0xfe, 0xea,
    ))

    if *Ram[0xeb25] != 0xff {
        test.Errorf("Memory at 0xEB25 was 0x%x, expected 0xFF\n", *Ram[0xEB25])
    }
}

func TestJsr(test *testing.T) {
    programCounter = 0

    program := []Word{
        0xea, // NOP
        0xea,
        0xea,
        0x20, 0x34, 0x12, // JSR #$1234
        0xea,
        0xea,
        0xea,
    }

    for index,value := range program {
        *Ram[index] = value
    }

    step(3)

    if programCounter != 0x03 {
        test.Errorf("Program counter was 0x%x, expected 0x03\n", programCounter)
    }

    step(1)

    if programCounter != 0x1234 {
        test.Errorf("Program counter was 0x%x, expected 0x1234\n", programCounter)
    }

    // TODO: Add stack logic tests
    // *Ram[0x01ff] == 0x00
    // *Ram[0x01fe] == 0x05

    // Force in some instructions
    *Ram[0x1234] = 0x38 // SEC
    *Ram[0x1235] = 0x60 // RTS

    step(3)

    if programCounter != 0x07 {
        test.Errorf("Program counter was 0x%x, expected 0x07\n", programCounter)
    }
}

func TestLsr(test *testing.T) {
    runTest(append(setupZeroPageMemory(),
        0x18,
        0x46, 0xfa, // LSR #$FA
    ))

    switch {
    case *Ram[0x00fa] != 0x18:
        test.Errorf("Memory at 0x00FA was 0x%x, expected 0x18\n", *Ram[0x00fa])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0x38, // SEC
        0x56, 0xf7, // LSR #$F7,X
    ))

    switch {
    case *Ram[0x00fa] != 0x18:
        test.Errorf("Memory at 0x00FA was 0x%x, expected 0x18\n", *Ram[0x00fa])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupAbsoluteMemory(),
        0x18,
        0x4e, 0x23, 0xeb, // LSR #$EB23
    ))

    switch {
    case *Ram[0xeb23] != 0x7e:
        test.Errorf("Memory at 0xEB23 was 0x%x, expected 0x7E\n", *Ram[0xEB23])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupAbsoluteIndexedXMemory(),
        0x18,
        0x5e, 0xfe, 0xea, // LSR #$EAFE,X
    ))

    switch {
    case *Ram[0xeb25] != 0x7f:
        test.Errorf("Memory at 0xEB25 was 0x%x, expected 0x7F\n", *Ram[0xEB25])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest([]Word{
        0xa9, 0xfe,
        0x18,
        0x4a, // LSR (shift accumulator)
    })

    switch {
    case cpu.A != 0x7f:
        test.Errorf("A was 0x%x, expected 0x7F\n", cpu.A)
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }
}

func TestAsl(test *testing.T) {
    runTest(append(setupZeroPageMemory(),
        0x18,
        0x06, 0xfa, // ASL #$FA
    ))

    switch {
    case *Ram[0x00fa] != 0x60:
        test.Errorf("Memory at 0x00FA was 0x%x, expected 0x60\n", *Ram[0x00fa])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0x38, // SEC
        0x16, 0xf7, // ASL #$F7,X
    ))

    switch {
    case *Ram[0x00fa] != 0x60:
        test.Errorf("Memory at 0x00FA was 0x%x, expected 0x60\n", *Ram[0x00fa])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupAbsoluteMemory(),
        0x18,
        0x0e, 0x23, 0xeb, // ASL #$EB23
    ))

    switch {
    case *Ram[0xeb23] != 0xf8:
        test.Errorf("Memory at 0xEB23 was 0x%x, expected 0xF8\n", *Ram[0xEB23])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case !cpu.Carry:
        test.Error("Carry bit was not set")
    case !cpu.Negative:
        test.Error("Negative bit was not set")
    }

    runTest(append(setupAbsoluteIndexedXMemory(),
        0x18,
        0x1e, 0xfe, 0xea, // ASL #$EAFE,X
    ))

    switch {
    case *Ram[0xeb25] != 0xfc:
        test.Errorf("Memory at 0xEB25 was 0x%x, expected 0xFC\n", *Ram[0xEB25])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case !cpu.Carry:
        test.Error("Carry bit was not set")
    case !cpu.Negative:
        test.Error("Negative bit was not set")
    }

    runTest([]Word{
        0xa9, 0xfe,
        0x18,
        0x0a, // ASL (shift accumulator)
    })

    switch {
    case cpu.A != 0xfc:
        test.Errorf("A was 0x%x, expected 0xfc\n", cpu.A)
    case cpu.Zero:
        test.Error("Zero bit was set")
    case !cpu.Carry:
        test.Error("Carry bit was not set")
    case !cpu.Negative:
        test.Error("Negative bit was not set")
    }
}

func TestRol(test *testing.T) {
    runTest(append(setupZeroPageMemory(),
        0x18,
        0x26, 0xFA, // ROL #$FA
    ))

    switch {
    case *Ram[0x00fa] != 0x60:
        test.Errorf("Memory at 0x00FA was 0x%x, expected 0x60\n", *Ram[0x00fa])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0x38,
        0x36, 0xF7, // ROL $#F7,X
    ))

    switch {
    case *Ram[0x00fa] != 0x61:
        test.Errorf("Memory at 0x00FA was 0x%x, expected 0x61\n", *Ram[0x00fa])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupAbsoluteMemory(),
        0x18,
        0x2E, 0x23, 0xEB, // ROL $#EB23
    ))

    switch {
    case *Ram[0xEB23] != 0xF8:
        test.Errorf("Memory at 0xEB23 was 0x%x, expected 0xF8\n", *Ram[0xEB23])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case !cpu.Carry:
        test.Error("Carry bit was not set")
    case !cpu.Negative:
        test.Error("Negative bit was not set")
    }

    runTest(append(setupAbsoluteIndexedXMemory(),
        0x18,
        0x3E, 0xFE, 0xEA, // ROL $#EAFE,X
    ))

    if *Ram[0xEB25] != 0xFC {
        test.Errorf("Memory at 0xEB25 was 0x%X, expected 0xFC\n", *Ram[0xEB25])
    }

    runTest([]Word{
        0xA9, 0xFE,
        0x18,
        0x2A,
    })

    if cpu.A != 0xFC {
        test.Errorf("A was 0x%x, expected 0xFC\n", cpu.A)
    }
}

func TestRor(test *testing.T) {
    runTest(append(setupZeroPageMemory(),
        0x18,
        0x66, 0xFA, // ROR #$FA
    ))

    switch {
    case *Ram[0x00FA] != 0x18:
        test.Errorf("Memory at 0x00FA was 0x%x, expected 0x18\n", *Ram[0x00fa])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupZeroPageIndexedXMemory(),
        0x38,
        0x76, 0xF7, // ROR $#F7,X
    ))

    switch {
    case *Ram[0x00FA] != 0x98:
        test.Errorf("Memory at 0x00FA was 0x%x, expected 0x98\n", *Ram[0x00fa])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case !cpu.Negative:
        test.Error("Negative bit was not set")
    }

    runTest(append(setupAbsoluteMemory(),
        0x18,
        0x6E, 0x23, 0xEB, // ROR $#EB23
    ))

    switch {
    case *Ram[0xEB23] != 0x7E:
        test.Errorf("Memory at 0xEB23 was 0x%x, expected 0x7E\n", *Ram[0xEB23])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest(append(setupAbsoluteIndexedXMemory(),
        0x18,
        0x7E, 0xFE, 0xEA, // ROR $#EAFE,X
    ))

    switch {
    case *Ram[0xEB25] != 0x7F:
        test.Errorf("Memory at 0xEB25 was 0x%x, expected 0x7F\n", *Ram[0xEB25])
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Carry:
        test.Error("Carry bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest([]Word{
        0xA9, 0xFE,
        0x18,
        0x6A,
    })

    if cpu.A != 0x7F {
        test.Errorf("A was 0x%x, expected 0x7F\n", cpu.A)
    }
}

func TestBit(test *testing.T) {
    runTest(append(setupZeroPageMemory(),
        0xA9, 0x10,
        0x18,
        0x24, 0xFA, // BIT #$10
    ))

    switch {
    case cpu.A != 0x10:
        test.Errorf("A was 0x%x, expected 0x10\n", cpu.A)
    case cpu.Zero:
        test.Error("Zero bit was set")
    case cpu.Overflow:
        test.Error("Overflow bit was set")
    case cpu.Negative:
        test.Error("Negative bit was set")
    }

    runTest([]Word{
        0xA9, 0xFF,
        0x85, 0xFA,
        0xA9, 0x00,
        0x24, 0xFA, // BIT $#FA
    })

    switch {
    case cpu.A != 0x00:
        test.Errorf("A was 0x%x, expected 0x00\n", cpu.A)
    case !cpu.Zero:
        test.Error("Zero bit was not set")
    case !cpu.Overflow:
        test.Error("Overflow bit was not set")
    case !cpu.Negative:
        test.Error("Negative bit was not set")
    }

    runTest(append(setupAbsoluteMemory(),
        0xA9, 0x0C,
        0x18,
        0x2C, 0x23, 0xEB, // BIT $EB,$23
    ))

    switch {
    case cpu.A != 0x0C:
        test.Errorf("A was 0x%x, expected 0x0C\n", cpu.A)
    case cpu.Zero:
        test.Error("Zero bit was set")
    case !cpu.Overflow:
        test.Error("Overflow bit was not set")
    case !cpu.Negative:
        test.Error("Negative bit was not set")
    }
}
