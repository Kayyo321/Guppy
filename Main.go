package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func usage() {
	fmt.Println("Usage: gup <interaction> <flags> <files>")
	fmt.Println()
	fmt.Println("<interaction>'s")
	fmt.Println("    ├──> build  : compile the following files into an executable.")
	fmt.Println("    ├──> run    : compile the files, and run it after built.")
	fmt.Println("    ├──> status : show that status of the Guppy compiler (i.e. everything is installed).")
	fmt.Println("    └──> edit   : edit a virtual file line by line, then save and or run it.")
	fmt.Println()
	fmt.Println("<flags>")
	fmt.Println("    ├──> '-o'     : build the files, and rename the output file: `-o mygup.exe`")
	fmt.Println("    ├──> '-asm'   : build the files to nasm; instead of '.exe' || '.o'.")
	fmt.Println("    └──> '-warns' : turns warnings either on || off: `-warns off`")
	fmt.Println()
	fmt.Println("<files>")
	fmt.Println("    ├──> Anything that ends with '.gpy' || '.guppy' .")
	fmt.Println("    ├──> Support absolute && relative file path.")
	fmt.Println("    └──> (by default, output name is \"a.exe\", not the first file's name.")
	os.Exit(0)
}

func exitAndStatus() {
	code := 0

	// Check if NASM (Netwide Assembler) is installed (or reachable) on the system.
	nasmV, err := exec.Command("nasm", "-v").Output()
	if err != nil {
		fmt.Println("NASM (Netwide Assembler) is not installed (or is unreachable) for this system.")
		code = 1
	} else {
		fmt.Println("NASM (Netwide Assembler) is installed && reachable on this system:\n\n", nasmV)
	}

	// Check if GCC (GNU C Compiler) is installed (or reachable) on the system.
	gccV, _err := exec.Command("gcc", "-v").Output()
	if _err != nil {
		fmt.Println("GCC (GNU C Compiler) is not installed (or is unreachable) for this system.")
		code = 1
	} else {
		fmt.Println("GCC (GNU C Compiler) is installed && reachable on this system:\n\n", gccV)
	}

	// Check if golang is installed on the system.
	goV, __err := exec.Command("go", "version").Output()
	if __err != nil {
		fmt.Println("Golang is not installed (or is unreachable) for this system")
		fmt.Println("(Not fatal; Guppy is built to an executable).")
	} else {
		fmt.Println("Golang is installed && reachable on this system:\n\n", goV)
	}

	fmt.Println("Guppy needs NASM to compile it's assembly code (based on the source file) to machine code.")
	fmt.Println("Guppy needs GCC to link it's machine code to the CSL (C Standard Library).")
	fmt.Println()
	fmt.Println("If these tools are not installed on your system, run the installer again, or download them seperately.")
	os.Exit(code)
}

// False = good
func status() (bool, bool, bool) {
	// Check if NASM (Netwide Assembler) is installed (or reachable) on the system.
	_, err := exec.Command("nasm", "-v").Output()

	// Check if GCC (GNU C Compiler) is installed (or reachable) on the system.
	_, _err := exec.Command("gcc", "-v").Output()

	// Check if golang is installed on the system.
	_, __err := exec.Command("go", "version").Output()

	return err == nil, _err == nil, __err == nil
}

func build(files []string, flags map[string]string) string {
	// var obj []string
	for _, filename := range files {
		bytes, err := os.ReadFile(filename)
		if err != nil {
			fmt.Println("Error reading: \""+filename+"\":", err)
			os.Exit(1)
		}

		data := string(bytes)

		tokens, errs := lex(data)
		if len(errs) != 0 {
			for _, err := range errs {
				fmt.Println("<!>", err)
			}

			fmt.Println("Guppy exited from the lexer.")
			os.Exit(1)
		}

		var p Parser
		p.set(tokens)

		ast, errs := p.parse()
		if len(errs) != 0 {
			for _, err := range errs {
				fmt.Println("<!>", err)
			}

			fmt.Println("Guppy exited from the parser.")
			os.Exit(1)
		}

		fmt.Println("AST {\n" + ast.toString() + "\n}")

		var c Compiler
		nasm, warns, errs := c.compile(*ast, ast, tokens)
		for _, warn := range warns {
			fmt.Println("<.>", warn)
		}
		if len(errs) != 0 {
			for _, err := range errs {
				fmt.Println("<!>", err)
			}

			fmt.Println("Guppy exited from the compiler.")
			os.Exit(1)
		}

		fmt.Println(nasm)

		// obj = append(obj, outputPath)
	}

	// After compilation, link all objects.
	// Return filepath to exe.

	return ""
}

func run(files []string, flags map[string]string) {
	o := build(files, flags)
	fmt.Println(o)
	// Run the exe.
}

func parseArgs() (map[string]string, []string) {
	flags := make(map[string]string)
	var files []string

	for i, arg := range os.Args {
		if arg[0] == '-' {
			flags[arg] = os.Args[i+1]
		} else {
			fileExt := filepath.Ext(arg)
			if fileExt == ".gpy" || fileExt == ".guppy" {
				files = append(files, arg)
			}
		}
	}

	return flags, files
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	// TODO: impliment flags.
	flags, files := parseArgs()
	fmt.Println(flags)

	switch os.Args[1] {
	case "status":
		exitAndStatus()

	case "build":
		// Check if necessary tools are installed.
		// Golang is not essential.
		x, y, _ := status()
		code := 0
		if x == true {
			fmt.Println("NASM (Netwide Assembler) is not installed (or is unreachable) for this system.")
			code = 1
		}
		if y == true {
			fmt.Println("GCC (GNU C Compiler) is not installed (or is unreachable) for this system.")
			code = 1
		}
		if code != 0 {
			os.Exit(code)
		}

		build(files, flags)

	case "run":
		run(files, flags)

	case "edit":
		// TODO: make edit.
		break

	default:
		fmt.Println("Unknown interaction:", os.Args[1])
		usage()
		os.Exit(1)
	}
}
