package kvm

import (
	"encoding/binary"
	"fmt"
	"math"
)

// ---------------------------
// Function representation
// ---------------------------

type Function struct {
	Chunk    []byte
	Consts   []interface{}
}

// Global function table
var functionTable []*Function

// ---------------------------
// Loader
// ---------------------------

func loadFunction(fnID int) *Function {
	if fnID < 0 || fnID >= len(functionTable) {
		panic(fmt.Sprintf("Invalid function ID: %d", fnID))
	}
	return functionTable[fnID]
}

// ---------------------------
// Decode function chunks
// ---------------------------

func LoadFunctionsFromModule(code []byte) {
	functionTable = []*Function{}

	if string(code[:3]) != "KBC" {
		panic("Invalid bytecode header")
	}

	version := code[3]
	if version != 2 {
		panic(fmt.Sprintf("Unsupported KBC version: %d", version))
	}

	offset := 4

	fnCount := int(binary.LittleEndian.Uint32(code[offset:]))
	offset += 4

	for i := 0; i < fnCount; i++ {
		fn, consumed := decodeChunk(code[offset:])
		offset += consumed
		functionTable = append(functionTable, fn)
	}
}

// ---------------------------
// Decode a single chunk
// ---------------------------

func decodeChunk(code []byte) (*Function, int) {
	offset := 0
	buf := make([]byte, 4)

	// Constants
	cCount := int(binary.LittleEndian.Uint32(code[offset:]))
	offset += 4

	consts := make([]interface{}, cCount)

	for i := 0; i < cCount; i++ {
		kind := code[offset]
		offset++

		switch kind {
		case 1: // string
			l := int(binary.LittleEndian.Uint32(code[offset:]))
			offset += 4
			str := string(code[offset : offset+l])
			offset += l
			consts[i] = str

		case 2: // float64
			bits := binary.LittleEndian.Uint64(code[offset:])
			offset += 8
			consts[i] = math.Float64frombits(bits)

		case 3: // int
			v := int(binary.LittleEndian.Uint32(code[offset:]))
			offset += 4
			consts[i] = v

		default:
			panic(fmt.Sprintf("Unknown constant type in function chunk: %d", kind))
		}
	}

	// Code
	codeLen := int(binary.LittleEndian.Uint32(code[offset:]))
	offset += 4

	fnCode := make([]byte, codeLen)
	copy(fnCode, code[offset:offset+codeLen])
	offset += codeLen

	return &Function{
		Chunk:  fnCode,
		Consts: consts,
	}, offset
}
