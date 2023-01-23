package main

import "fmt"

const (
	Int     = 4
	uInt    = 4
	uInt32  = 4
	Int32   = 4
	Float32 = 4
	Ptr     = 4
	Range   = 4
	Int64   = 8
	uint64  = 8
	Float64 = 8
	String  = 8
	Int8    = 1
	uInt8   = 1
	Byte    = 1
	Bool    = 1
	Int16   = 2
)

type Struct struct {
	data []Var
	sz   int
}

type Var struct {
	name     string
	loc      string
	typ      string
	sz       int
	isConst  bool
	isStatic bool
}

type Func struct {
	name     string
	params   []Var
	returnT  string
	returnS  int
	reciever Struct
	body     Scope
}

type Scope struct {
	vars    []Var
	bytes   uint
	nasm    string
	data    *string
	bss     *string
	externs *string
	_main   *string
}

type Env struct {
	enums   []string
	funcs   []Func
	structs []Struct
	scope   Scope
	data    string
	bss     string
	text    string
	externs string
	_main   string
	strMap  map[string]string
	fltMap  map[float64]string
}

type Compiler struct {
	env  Env
	head Node
	ptr  *Node
	tks  []Token
}

func (cr *Compiler) compile(head Node, ptr *Node, tks []Token) (string, []error, []error) {
	var warns, errs []error
	var nasm string

	cr.head = head
	cr.ptr = ptr
	cr.tks = tks

	cr.env.strMap = make(map[string]string)
	cr.env.fltMap = make(map[float64]string)
	cr.getLiterals()

	__x__, __y__ := cr.rcompile()
	warns, errs = append(warns, __x__...), append(errs, __y__...)

	nasm += cr.env.externs
	nasm += "section .data ; {\n"
	nasm += cr.env.data
	nasm += "; }\n\nsection .bss ; {\n"
	nasm += cr.env.bss
	nasm += "; }\n\nsection .text ; {\n\tglobal main\nmain: ; {\n"
	nasm += cr.env._main
	nasm += "; }\n\n"
	nasm += cr.env.text

	return nasm, warns, errs
}

func (cr *Compiler) getLiterals() {
	var s uint
	var f uint

	for _, t := range cr.tks {
		if t.tokenType == TkString {
			cr.env.strMap[t.sdata] = "_str" + string(rune(s))
			cr.env.data += cr.env.strMap[t.sdata] + " db `" + t.sdata + "`, 0\n"
			s++
		} else if t.tokenType == TkFloat {
			cr.env.fltMap[t.fdata] = "_flt" + string(rune(f))
			cr.env.data += cr.env.fltMap[t.fdata] + " dq "
			cr.env.data += fmt.Sprintf("%f", t.fdata)
			cr.env.data += "\n"
			f++
		}
	}
}

/**
Code:
var := 1 + 2 * 3

AST Structure:

VarDefinition
	|____name1
	|____... (other names, if any)
	|____Body/Value
			|____ +
				|____ *
				|   |____ 3
				|   |____ 2
				|____1

*/
func (cr *Compiler) compileDefine() ([]error, []error) {
	var warns, errs []error

	return warns, errs
}

func (cr *Compiler) rcompile() ([]error, []error) {
	var warns, errs []error

	switch cr.ptr.nt {
	case NdDefine:
		__x__, __y__ := cr.compileDefine()
		warns, errs = append(warns, __x__...), append(errs, __y__...)
	}

	return warns, errs
}
