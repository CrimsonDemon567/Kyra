package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"kyra/internal/lexer"
	"kyra/internal/parser"
	"kyra/internal/bytecode"
	"kyra/internal/kvm"
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

	exe := filepath.Base(args[0])
	flag := args[1]

	switch exe {

	// --------------------------------------------------------
	// COMPILER: kyrac
	// --------------------------------------------------------
	case "kyrac":
		switch flag {

		case "-kbc":
			// kyrac -kbc file.kyra → nur kompilieren
			buildKBC(args[2])
			return

		case "-kar":
			// kyrac -kar folder → kompilieren + archiv bauen
			buildKAR(args[2])
			return

		default:
			fmt.Println("Unknown flag:", flag)
			usageCompiler()
			return
		}

	// --------------------------------------------------------
	// RUNTIME: kyra
	// --------------------------------------------------------
	case "kyra":
		switch flag {

		case "-kbc":
			// kyra -kbc file.kbc → nur ausführen
			runKBC(args[2])
			return

		case "-kar":
			// kyra -kar file.kar → archiv ausführen
			runKAR(args[2])
			return

		default:
			fmt.Println("Unknown flag:", flag)
			usageRuntime()
			return
		}

	default:
		fmt.Println("Unknown executable:", exe)
		usage()
	}
}

// ------------------------------------------------------------
// Usage
// ------------------------------------------------------------

func usage() {
	fmt.Println("Kyra Toolchain")
	fmt.Println("")
	usageCompiler()
	usageRuntime()
}

func usageCompiler() {
	fmt.Println("Compiler (kyrac):")
	fmt.Println("  kyrac -kbc <file.kyra>   Compile a single Kyra file into .kbc")
	fmt.Println("  kyrac -kar <folder>      Compile folder and build .kar archive")
}

func usageRuntime() {
	fmt.Println("Runtime (kyra):")
	fmt.Println("  kyra -kbc <file.kbc>     Execute a compiled Kyra bytecode file")
	fmt.Println("  kyra -kar <file.kar>     Execute a Kyra archive")
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
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading .kbc:", err)
		return
	}

	kvm.LoadFunctionsFromModule(data)

	vm := kvm.New(data)
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
		fmt.Println("Error reading .kar:", err)
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
