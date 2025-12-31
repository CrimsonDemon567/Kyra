package bytecode

import (
	"encoding/binary"
	"fmt"
	"kyra/internal/parser"
	"math"
)

// ---------------------------
// Function table
// ---------------------------

type Function struct {
	Name  string
	Args  []parser.FuncArg
	Chunk *Chunk
}

var functionTable []*Function

func resetFunctions() {
	functionTable = []*Function{}
}

func addFunction(fn *Function) int {
	functionTable = append(functionTable, fn)
	return len(functionTable) - 1
}

// ---------------------------
// Emit functions
// ---------------------------

func emitFunctionDef(parent *Chunk, fn *parser.FuncDef) {
	ch := NewChunk()

	for _, arg := range fn.Args {
		slot := ch.addConst(arg.Name)
		ch.Names[arg.Name] = slot
	}

	for _, st := range fn.Body {
		emitStmt(ch, st)
	}

	ch.emit(OP_RET)

	fnID := addFunction(&Function{
		Name:  fn.Name,
		Args:  fn.Args,
		Chunk: ch,
	})

	slot := parent.addConst(fnID)
	parent.emit(OP_CONST)
	parent.emitInt(slot)
}

func emitFunctionExpr(parent *Chunk, fn *parser.FuncExprDef) {
	ch := NewChunk()

	for _, arg := range fn.Args {
		slot := ch.addConst(arg.Name)
		ch.Names[arg.Name] = slot
	}

	for _, st := range fn.Body {
		emitStmt(ch, st)
	}

	ch.emit(OP_RET)

	fnID := addFunction(&Function{
		Name:  fn.Name,
		Args:  fn.Args,
		Chunk: ch,
	})

	slot := parent.addConst(fnID)
	parent.emit(OP_CONST)
	parent.emitInt(slot)
}

func emitFunctionOneLiner(parent *Chunk, fn *parser.FuncOneLiner) {
	ch := NewChunk()

	for _, arg := range fn.Args {
		slot := ch.addConst(arg.Name)
		ch.Names[arg.Name] = slot
	}

	emitExpr(ch, fn.Expr)
	ch.emit(OP_RET)

	fnID := addFunction(&Function{
		Name:  fn.Name,
		Args:  fn.Args,
		Chunk: ch,
	})

	slot := parent.addConst(fnID)
	parent.emit(OP_CONST)
	parent.emitInt(slot)
}

// ---------------------------
// Register function emitters
// ---------------------------

func init() {
	RegisterStmtEmitter(func(c *Chunk, stmt parser.Stmt) bool {
		switch s := stmt.(type) {

		case *parser.FuncDef:
			emitFunctionDef(c, s)
			return true

		case *parser.FuncExprDef:
			emitFunctionExpr(c, s)
			return true

		case *parser.FuncOneLiner:
			emitFunctionOneLiner(c, s)
			return true
		}

		return false
	})
}

// ---------------------------
// Encode full module with functions
// ---------------------------

func encodeModuleWithFunctions(chunk *Chunk) []byte {
	out := []byte{}

	out = append(out, 'K', 'B', 'C', 2)

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(functionTable)))
	out = append(out, buf...)

	for _, fn := range functionTable {
		out = append(out, encodeChunk(fn.Chunk)...)
	}

	out = append(out, encodeChunk(chunk)...)

	return out
}

func encodeChunk(c *Chunk) []byte {
	out := []byte{}
	buf := make([]byte, 4)

	binary.LittleEndian.PutUint32(buf, uint32(len(c.Constants)))
	out = append(out, buf...)

	for _, v := range c.Constants {
		switch val := v.(type) {
		case string:
			out = append(out, 1)
			str := []byte(val)
			binary.LittleEndian.PutUint32(buf, uint32(len(str)))
			out = append(out, buf...)
			out = append(out, str...)

		case float64:
			out = append(out, 2)
			fb := make([]byte, 8)
			binary.LittleEndian.PutUint64(fb, math.Float64bits(val))
			out = append(out, fb...)

		case int:
			out = append(out, 3)
			binary.LittleEndian.PutUint32(buf, uint32(val))
			out = append(out, buf...)

		default:
			panic(fmt.Sprintf("Unknown constant type: %T", v))
		}
	}

	binary.LittleEndian.PutUint32(buf, uint32(len(c.Code)))
	out = append(out, buf...)
	out = append(out, c.Code...)

	return out
}
