package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"kyra/internal/lexer"
	"kyra/internal/parser"
	"kyra/internal/bytecode"
	"kyra/kvm"
	"kyra/internal/kar"
)

// ------------------------------------------------------------
// Entry point
// ------------------------------------------------------------

func RunCLI(args []string) {
	if len(args) < 3 {
		usage()
		return
	}

	flag := args[1]

	switch flag {

	// --------------------------------------------------------
	// COMPILER FLAGS (kyrac)
	// --------------------------------------------------------
	case "-kbc":
		// kyrac -kbc file.kyra  → compile only
		if isCompiler(args[0]) {
			buildKBC(args[2])
			return
		}

		// kyra -kbc file.kyra → compile + run
		runKBC(args[2])
		return

	case "-kar":
		// kyrac -kar folder → build archive
		if isCompiler(args[0]) {
			buildKAR(args[2])
			return
		}

		// kyra -kar file.kar → run archive
		runKAR(args[2])
		return

	default:
		fmt.Println("Unknown flag:", flag)
		usage()
	}
}

// ------------------------------------------------------------
// Detect compiler binary
// ------------------------------------------------------------

func isCompiler(exe string) bool {
	base := filepath.Base(exe)
	return base == "kyrac"
}

// ------------------------------------------------------------
// Usage
// ------------------------------------------------------------

func usage() {
	fmt.Println("Kyra Toolchain")
	fmt.Println("")
	fmt.Println("Compiler (kyrac):")
	fmt.Println("  kyrac -kbc <file.kyra>   Compile a single Kyra file into .kbc")
	fmt.Println("  kyrac -kar <folder>      Build a .kar archive from a folder")
	fmt.Println("")
	fmt.Println("Runtime (kyra):")
	fmt.Println("  kyra -kbc <file.kyra>    Compile and execute a single Kyra file")
	fmt.Println("  kyra -kar <file.kar>     Execute a Kyra .kar archive")
}

// ------------------------------------------------------------
// Compiler: KBC
// ------------------------------------------------------------

func buildKBC(path string) {
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	tokens := lexer.New(string(src)).Lex()
	p := parser.New(tokens)
	ast := p.Parse()

	bytecode.ResetFunctions()
	code := bytecode.Emit(ast)

	out := strings.TrimSuffix(path, filepath.Ext(path)) + ".kbc"
	err = os.WriteFile(out, code, 0644)
	if err != nil {
		fmt.Println("Error writing output:", err)
		return
	}

	fmt.Println("Built:", out)
}

// ------------------------------------------------------------
// Compiler: KAR
// ------------------------------------------------------------

func buildKAR(folder string) {
	arc, err := kar.BuildFromFolder(folder)
	if err != nil {
		fmt.Println("Error building KAR:", err)
		return
	}

	out := folder + ".kar"
	err = arc.Save(out)
	if err != nil {
		fmt.Println("Error saving KAR:", err)
		return
	}

	fmt.Println("Built archive:", out)
}

// ------------------------------------------------------------
// Runtime: KBC
// ------------------------------------------------------------

func runKBC(path string) {
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	tokens := lexer.New(string(src)).Lex()
	p := parser.New(tokens)
	ast := p.Parse()

	bytecode.ResetFunctions()
	code := bytecode.Emit(ast)

	kvm.LoadFunctionsFromModule(code)

	vm := kvm.New(code)
	result := vm.Run()

	if result != nil {
		fmt.Println(result)
	}
}

// ------------------------------------------------------------
// Runtime: KAR
// ------------------------------------------------------------

func runKAR(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading KAR:", err)
		return
	}

	arc, err := kar.Decode(data)
	if err != nil {
		fmt.Println("Error decoding KAR:", err)
		return
	}

	main := arc.Get("main.kbc")
	if main == nil {
		fmt.Println("KAR archive has no main.kbc")
		return
	}

	kvm.LoadFunctionsFromModule(main)
	vm := kvm.New(main)
	result := vm.Run()

	if result != nil {
		fmt.Println(result)
	}
}
