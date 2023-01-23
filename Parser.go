package main

import (
	"errors"
)

const (
	NdPack          = "Package"
	NdIdent         = "Ident"
	NdEllipsis      = "Ellipsis"
	NdFunctionDef   = "Function Def"
	NdFunctionParam = "Function Param"
	NdFunctionRec   = "Function Reciever"
	NdBody          = "Body || Value"
	NdFunctionCall  = "Function Call"
	NdStructDef     = "Struct Def"
	NdStructCall    = "Struct Call"
	NdReturn        = "Return"
	NdAssign        = "Assign"
	NdDecl          = "Declare"
	NdDefine        = "Define"
	NdCall          = "Call"
	NdFor           = "For"
	NdForInit       = "For Init"
	NdForBool       = "For Bool"
	NdForInc        = "For Inc"
	NdWhile         = "While"
	NdIf            = "If"
	NdElIf          = "Else If"
	NdElse          = "Else"
	NdImport        = "Import"
	NdBreak         = "Break"
	NdContinue      = "Continue"
	NdAttribute     = "Attribute"
	NdSect          = "Section"
	NdIncDec        = "Inc || Dec"
	NdBinOp         = "Bin Op"
	NdLit           = "Lit"
	NdType          = "Type"
	NdArrLit        = "Array Lit"
	NdInln          = "Inline Asm"
	NdIndex         = "Index"
	NdEnum          = "Enum"
	NdEnumCall      = "Enum Call"
	NdRangeDef      = "Range Def"
	NdDefer         = "Defer"
	NdClean         = "Clean"
	NdClink         = "C-Link"
	NdAbbreviation  = "Abreviation"
	NdMatch         = "Match"
	NdTernaryOp     = "Ternary Operation"
	NdName          = "Name"
)

type Node struct {
	nt    string
	kids  []Node
	gt    *Token
	sdata string
	bytes uint
}

func (n *Node) child() *Node {
	n.kids = append(n.kids, Node{})
	return &n.kids[len(n.kids)-1]
}

func (n *Node) toString() string {
	return rtoString(n, 0, true, []uint{})
}

func rtoString(n *Node, indent uint, endNode bool, indents []uint) string {
	var s string

	var in uint = 0

	for ; in < indent; in++ {
		for _, x := range indents {
			if in == x {
				s += "| "
			}
		}

		s += "  "
	}

	if endNode {
		s += "└──"
	} else {
		s += "├──"
	}

	for i := uint(0); i < indent/2; i++ {
		s += "──"
	}

	s += " "

	// don't print nil pointer
	if n.gt == nil {
		if n.nt != NdPack && !endNode {
			indents = append(indents, in)
		}
		s += "({" + n.nt + "}, \"" + n.sdata + "\")\n"
	} else {
		s += "(" + n.gt.toString() + ", {" + n.nt + "}, \"" + n.sdata + "\")\n"
		if !endNode {
			indents = append(indents, in)
		}
	}

	if len(n.kids) > 0 {
		for i := range n.kids {
			s += rtoString(&n.kids[i], indent+1, i == len(n.kids)-1, indents)
		}
	}

	return s
}

type Parser struct {
	hd  *Node
	tks []Token
	t   *Token
	tI  uint
}

func (p *Parser) set(_tks []Token) {
	p.tks = _tks
	p.tI = 0
	p.t = &p.tks[p.tI]
	p.hd = new(Node)
}

func (p *Parser) eat() error {
	if p.tI+1 >= uint(len(p.tks)) {
		return errors.New("Expected token, found __eof__.")
	}

	p.tI++
	p.t = &p.tks[p.tI]

	return nil
}

func (p *Parser) assert(expect []int8) error {
	eatErr := p.eat()
	if eatErr != nil {
		return eatErr
	}

	var msg = "Expected: "

	for _i, ttype := range expect {
		if ttype == p.t.tokenType {
			return nil
		}

		if _i != len(expect)-1 {
			msg += tokenTypeStrings[ttype] + " || "
		} else {
			msg += tokenTypeStrings[ttype]
		}
	}

	msg += ", found: " + p.t.toString()

	return errors.New(msg)
}

func (p *Parser) peek() (*Token, error) {
	if p.tI+1 >= uint(len(p.tks)) {
		return nil, errors.New("Expected token, found __eof__")
	}

	return &p.tks[p.tI+1], nil
}

func (p *Parser) parse() (*Node, []error) {
	var errs []error

	p.hd.nt = NdPack
	p.hd.sdata = "_start"

	for p.t.tokenType != TkPackageKw {
		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return p.hd, errs
		}
	}

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return p.hd, errs
	}

	p.hd.gt = p.t

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return p.hd, errs
	}

	errs = append(errs, p.rparse(p.hd, []int8{TkEof})...)

	return p.hd, errs
}

func (p *Parser) parseImport(into *Node, until []int8) []error {
	var errs []error

	im := into.child()
	im.nt = NdImport
	im.sdata = "Import"
	im.gt = p.t

	if err := p.assert([]int8{TkIdent}); err != nil {
		errs = append(errs, err)
		return errs
	}

moreImp:
	for p.t.tokenType != TkComma && p.t.tokenType != TkWhiteSpace && p.t.tokenType != TkSemicolon {
		for _, u := range until {
			if u == p.t.tokenType|TkEof {
				return errs
			}
		}

		if p.t.tokenType == TkIdent {
			id := im.child()
			id.nt = NdIdent
			id.sdata = p.t.sdata
			id.gt = p.t
		} else {
			errs = append(errs, errors.New("Import statements can only be identifiers, did not expect: "+p.t.toString()))
		}

		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return errs
		}
	}

	if p.t.tokenType == TkComma {
		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return errs
		}
		goto moreImp
	}

	return errs
}

func (p *Parser) typeToBSize() uint {
	var bsize uint
	switch p.t.tokenType {
	case TkIntKw, TkUIntKw, TkUInt32Kw, TkInt32Kw, TkFloat32Kw, TkPtrKw, TkRangeKw:
		bsize = 4
	case TkInt64Kw, TKUint64Kw, TkFloat64Kw, TkStringKw:
		bsize = 8
	case TkInt8Kw, TkUInt8Kw, TkByteKw, TkBoolKw:
		bsize = 1
	case TkInt16Kw, TkUInt16Kw:
		bsize = 2
	case TkLBrack:
		p.eat()

		if p.t.tokenType == TkInt {
			count := p.t.idata

			err := p.assert([]int8{TkComma})
			if err != nil {
				break
			}

			p.eat()

			sz := p.typeToBSize()

			bsize = sz * uint(count)

			errx := p.assert([]int8{TkRBrack})
			if errx != nil {
				break
			}
		} else {
			bsize = 8
		}
	}
	return bsize
}

// x: int = 20
func (p *Parser) parseVarDecl(into *Node, until []int8, names []Token) []error {
	var errs []error

	vardcl := into.child()
	vardcl.nt = NdDecl

	for _, name := range names {
		nameN := vardcl.child()
		nameN.nt = NdName
		nameN.gt = &name
		nameN.sdata = name.sdata
	}

	// Skip ':'
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	typ := vardcl.child()
	typ.nt = NdType
	typ.gt = p.t

	vardcl.bytes = p.typeToBSize() * uint(len(names))
	if vardcl.bytes == 0 {
		if p.t.tokenType == TkEllipsis {
			el := vardcl.child()
			el.nt = NdEllipsis
			el.gt = p.t
			el.sdata = "..."
		} else if p.t.tokenType == TkRangeKw {
			rng := vardcl.child()
			rng.nt = NdRangeDef
			rng.gt = p.t
			rng.sdata = "rng"
		} else {
			tp := vardcl.child()
			tp.nt = NdStructDef
			tp.gt = p.t
			tp.sdata = "struct"
		}
	}

	// Skip type
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	if p.t.tokenType == TkComma {
		if err := p.assert([]int8{TkConstKw, TkStaticKw}); err != nil {
			errs = append(errs, err)
			return errs
		}

		cnst := vardcl.child()
		cnst.gt = p.t
		cnst.nt = NdSect
		cnst.sdata = p.t.sdata

		// Skip 'const' || 'static'
		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return errs
		}
	}

	if p.t.tokenType == TkAssign {
		// Skip '='
		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return errs
		}

		body := vardcl.child()
		body.nt = NdBody

		errs = append(errs, p.rparse(body, []int8{TkWhiteSpace, TkSemicolon})...)
	}

	return errs
}

// y := 30
func (p *Parser) parseVarDef(into *Node, until []int8, names []Token) []error {
	var errs []error

	vardcl := into.child()
	vardcl.nt = NdDefine

	for _, name := range names {
		nameN := vardcl.child()
		nameN.nt = NdName
		nameN.gt = &name
		nameN.sdata = name.sdata
	}

	// Skip ':='
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	body := vardcl.child()
	body.nt = NdBody

	errs = append(errs, p.rparse(body, []int8{TkWhiteSpace, TkSemicolon})...)

	return errs
}

// func(x: int, y: int): int { return x + y; }
func (p *Parser) parseFuncLit(into *Node, until []int8) []error {
	var errs []error

	fdef := into.child()
	fdef.gt = p.t
	fdef.nt = NdFunctionDef
	fdef.sdata = p.t.sdata

	if err := p.assert([]int8{TkLParen}); err != nil {
		errs = append(errs, err)
		return errs
	}

	param := fdef.child()
	param.gt = p.t
	param.nt = NdFunctionParam
	param.sdata = "Params"

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	for p.t.tokenType != TkRParen {
		errs = append(errs, p.rparse(param, []int8{TkComma, TkRParen})...)
		if p.t.tokenType == TkComma {
			// Skip ','
			if err := p.eat(); err != nil {
				errs = append(errs, err)
				return errs
			}
		}
	}

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	if p.t.tokenType == TkColon {
		// Skip ':'
		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return errs
		}

		fdef.bytes = p.typeToBSize()

		// Skip type
		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return errs
		}

		if p.t.tokenType == TkShr {
			// Skip '>>'
			if err := p.eat(); err != nil {
				errs = append(errs, err)
				return errs
			}

			if p.t.tokenType != TkIdent {
				errs = append(errs, errors.New("Reciever name was not an identifier: "+p.t.toString()))
				return errs
			}

			rec := fdef.child()
			rec.nt = NdFunctionRec
			rec.gt = p.t
			rec.sdata = p.t.sdata

			if err := p.assert([]int8{TkLBrace}); err != nil {
				errs = append(errs, err)
				return errs
			}
		} else if p.t.tokenType != TkLBrace {
			errs = append(errs, errors.New("Function had no body; instead found:"+p.t.toString()))
			return errs
		}
	} else if p.t.tokenType == TkShr {
		// Skip '>>'
		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return errs
		}

		if p.t.tokenType != TkIdent {
			errs = append(errs, errors.New("Reciever name was not an identifier: "+p.t.toString()))
			return errs
		}

		rec := fdef.child()
		rec.nt = NdFunctionRec
		rec.gt = p.t
		rec.sdata = p.t.sdata

		if err := p.assert([]int8{TkLBrace}); err != nil {
			errs = append(errs, err)
			return errs
		}
	} else if p.t.tokenType != TkLBrace {
		errs = append(errs, errors.New("Function did not have body: "+p.t.toString()))
		return errs
	}

	body := fdef.child()
	body.nt = NdBody
	body.gt = p.t
	body.sdata = "NdBody"

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	errs = append(errs, p.rparse(body, []int8{TkRBrace})...)

	return errs
}

func (p *Parser) parseExpr(into *Node, until []int8) []error {
	var errs []error

	var queue []Node
	var stack []Node

	for {
		abort := false

		for _, u := range until {
			if u == p.t.tokenType || p.t.tokenType == TkEof || p.t.tokenType == TkWhiteSpace {
				abort = true
			}
		}

		if abort {
			break
		}

		if p.t.tokenType == TkInt || p.t.tokenType == TkChar || p.t.tokenType == TkFloat || p.t.tokenType == TkString || p.t.tokenType == TkNilKw {
			lit := Node{}
			lit.nt = NdLit
			lit.gt = p.t
			queue = append(queue, lit)
		} else if p.t.tokenType == TkEnumKw {
			enum := Node{}
			enum.nt = NdEnumCall
			enum.gt = p.t

			if err := p.assert([]int8{TkPeriod}); err != nil {
				errs = append(errs, err)
				return errs
			}

			if err := p.assert([]int8{TkIdent}); err != nil {
				errs = append(errs, err)
				return errs
			}

			call := enum.child()
			call.nt = NdIdent
			call.gt = p.t

			queue = append(queue, enum)
		} else if p.t.tokenType == TkIdent {
			tk, err := p.peek()
			if err != nil {
				errs = append(errs, err)
				return errs
			}

			for _, u := range until {
				if u == tk.tokenType || tk.tokenType == TkEof {
					call := Node{}
					call.nt = NdCall
					call.gt = p.t
					call.sdata = "VarCall"
					if err := p.eat(); err != nil {
						errs = append(errs, err)
					}
					abort = true
					queue = append(queue, call)
					break
				}
			}

			if abort {
				break
			}

			if tk.tokenType == TkOr { // |
				attr := Node{}
				attr.nt = NdAttribute
				attr.gt = p.t
				attr.sdata = p.t.sdata

				// Move to '|'
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				// Skip '|'
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				errs = append(errs, p.rparse(&attr, []int8{TkWhiteSpace, TkSemicolon})...)
				queue = append(queue, attr)
			} else if tk.tokenType == TkPeriod { // .
				strcall := Node{}
				strcall.nt = NdStructCall
				strcall.gt = p.t
				strcall.sdata = p.t.sdata

				// Move to '.'
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				// Skip '.'
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				errs = append(errs, p.rparse(&strcall, []int8{TkWhiteSpace, TkSemicolon})...)
				queue = append(queue, strcall)
			} else if tk.tokenType == TkLParen { // (
				fcall := Node{}
				fcall.nt = NdFunctionCall
				fcall.gt = p.t
				fcall.sdata = p.t.sdata

				// Move to '('
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				param := fcall.child()
				param.nt = NdFunctionParam
				param.gt = p.t
				param.sdata = "Params"

				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				for p.t.tokenType != TkRParen {
					errs = append(errs, p.rparse(param, []int8{TkComma, TkRParen})...)
					if p.t.tokenType == TkComma {
						// Skip ','
						if err := p.eat(); err != nil {
							errs = append(errs, err)
							break
						}
					}
				}

				queue = append(queue, fcall)
			} else if tk.tokenType == TkInc || tk.tokenType == TkDec || tk.tokenType == TkSqrInc {
				inc := Node{}
				inc.nt = NdIncDec
				inc.gt = tk
				inc.sdata = tk.sdata

				vcall := inc.child()
				vcall.nt = NdCall
				vcall.gt = p.t

				// Move to ++ || -- || **
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				queue = append(queue, inc)
			} else if tk.tokenType == TkLBrack { // [
				index := Node{}
				index.nt = NdIndex
				index.sdata = "Index"

				varname := index.child()
				varname.gt = p.t
				varname.nt = NdCall

				// Move to the index literal
				if err := p.eat(); err != nil {
					errs = append(errs, err)
					break
				}
				if err := p.eat(); err != nil {
					errs = append(errs, err)
					break
				}

				num := index.child()
				num.nt = NdLit

				errs = append(errs, p.rparse(num, []int8{TkRBrack})...)

				pk, err := p.peek()
				if err != nil {
					errs = append(errs, err)
					break
				}

				body := index.child()
				body.nt = NdBody

				if pk.tokenType == TkAssign || pk.tokenType == TkOrAssign || pk.tokenType == TkAddAssign || pk.tokenType == TkAndAssign || pk.tokenType == TkMulAssign || pk.tokenType == TkAndNotAssign || pk.tokenType == TkQuoAssign || pk.tokenType == TkRemAssign || pk.tokenType == TkShlAssign || pk.tokenType == TkShrAssign || pk.tokenType == TkSubAssign || pk.tokenType == TkXorAssign {
					errs = append(errs, p.parseAssign(body, []int8{TkWhiteSpace, TkSemicolon})...)
				}

				queue = append(queue, index)
			} else if tk.tokenType == TkAssign || tk.tokenType == TkOrAssign || tk.tokenType == TkAddAssign || tk.tokenType == TkAndAssign || tk.tokenType == TkMulAssign || tk.tokenType == TkAndNotAssign || tk.tokenType == TkQuoAssign || tk.tokenType == TkRemAssign || tk.tokenType == TkShlAssign || tk.tokenType == TkShrAssign || tk.tokenType == TkSubAssign || tk.tokenType == TkXorAssign {
				asn := Node{}
				errs = append(errs, p.parseAssign(&asn, until)...)
				queue = append(queue, asn)
			} else {
				call := Node{}
				call.nt = NdCall
				call.gt = p.t
				call.sdata = "VarCall"
				queue = append(queue, call)
			}
		} else if p.t.tokenType == TkAdd || p.t.tokenType == TkSub || p.t.tokenType == TkMul || p.t.tokenType == TkQuo || p.t.tokenType == TkRem || p.t.tokenType == TkEql || p.t.tokenType == TkNeq || p.t.tokenType == TkLeq || p.t.tokenType == TkLss || p.t.tokenType == TkGeq || p.t.tokenType == TkGtr || p.t.tokenType == TkLand || p.t.tokenType == TkLor || p.t.tokenType == TkNot {
			o1 := p.t

			for len(stack) > 0 {
				o2 := stack[len(stack)-1]

				if (o1.sdata != "^" && o1.precedence <= o2.gt.precedence) || (o1.sdata == "^" && o1.precedence < o2.gt.precedence) {
					stack = stack[:len(stack)-1]
					queue = append(queue, o2)
					continue
				}

				break
			}

			bo := Node{}
			bo.gt = o1
			bo.nt = NdBinOp
			bo.sdata = "BinOp"
			stack = append(stack, bo)
		} else if p.t.tokenType == TkLParen {
			lpa := Node{}
			lpa.nt = NdLit
			lpa.gt = p.t
			lpa.sdata = "Lpa"
			stack = append(stack, lpa)
		} else if p.t.tokenType == TkRParen {
			match := false

			for len(stack) > 0 && stack[len(stack)-1].gt.tokenType != TkLParen {
				queue = append(queue, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
				match = true
			}

			if !match && len(stack) < 1 {
				errs = append(errs, errors.New("Right parenthesis: "+p.t.toString()))
				break
			}

			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
				if err := p.eat(); err != nil {
					errs = append(errs, err)
					return errs
				}
				break
			}
		} else {
			break
		}

		for _, u := range until {
			if u == p.t.tokenType || p.t.tokenType == TkEof || p.t.tokenType == TkWhiteSpace {
				abort = true
			}
		}

		if abort {
			break
		}

		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return errs
		}
	}

	for len(stack) > 0 {
		queue = append(queue, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	var i = len(queue)
	errs = append(errs, p.rcompute(into, queue, &i)...)

	return errs
}

func (p *Parser) rcompute(into *Node, stack []Node, i *int) []error {
	var errs []error

	var passes uint = 0

	for *i--; *i >= 0; *i-- {
		if passes >= 2 {
			break
		}

		n := stack[*i]

		switch n.nt {
		case NdLit, NdCall, NdFunctionCall, NdStructCall, NdAttribute, NdIndex:
			x := into.child()
			*x = n

		case NdBinOp:
			y := into.child()
			*y = n
			errs = append(errs, p.rcompute(y, stack, i)...)

		default:
			errs = append(errs, errors.New("Didn't expect node: "+n.toString()))
		}

		passes++
	}

	return errs
}

func (p *Parser) parseStructLit(into *Node, until []int8) []error {
	var errs []error

	if err := p.assert([]int8{TkLBrace}); err != nil {
		errs = append(errs, err)
		return errs
	}

	errs = append(errs, p.rparse(into, []int8{TkRBrace})...)

	return errs
}

func (p *Parser) parseIfStmt(into *Node, until []int8) []error {
	var errs []error

	ifstmt := into.child()
	ifstmt.nt = NdIf
	ifstmt.gt = p.t
	ifstmt.sdata = "If Statement"

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	errs = append(errs, p.parseExpr(ifstmt, []int8{TkLBrace})...)

	body := ifstmt.child()
	body.nt = NdBody
	body.gt = p.t
	body.sdata = "Body"

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	errs = append(errs, p.rparse(body, []int8{TkRBrace})...)

	// Skip '}'
	if err := p.eat(); err != nil {
		errs = append(errs, err)
	}

	for p.t.tokenType == TkElseKw {
		tk, err := p.peek()
		if err != nil {
			errs = append(errs, err)
			break
		}

		if tk.tokenType == TkIfKw {
			elstmt := into.child()
			elstmt.nt = NdElIf
			elstmt.gt = p.t
			elstmt.sdata = "Else If Statement"

			errs = append(errs, p.parseExpr(elstmt, []int8{TkLBrace})...)

			ebody := elstmt.child()
			ebody.nt = NdBody
			ebody.gt = p.t
			ebody.sdata = "Body"

			if err := p.eat(); err != nil {
				errs = append(errs, err)
				break
			}

			errs = append(errs, p.rparse(ebody, []int8{TkRBrace})...)
		} else {
			elstmt := into.child()
			elstmt.nt = NdElse
			elstmt.gt = p.t
			elstmt.sdata = "Else Statement"

			if err := p.assert([]int8{TkLBrace}); err != nil {
				errs = append(errs, err)
				break
			}

			if err := p.eat(); err != nil {
				errs = append(errs, err)
				break
			}

			ebody := elstmt.child()
			ebody.nt = NdBody
			ebody.gt = p.t
			ebody.sdata = "Body"

			errs = append(errs, p.rparse(ebody, []int8{TkRBrace})...)
		}

		if err := p.eat(); err != nil {
			errs = append(errs, err)
			break
		}
	}

	// Skip '}'
	if p.t.tokenType == TkRBrace {
		if err := p.eat(); err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (p *Parser) parseForStmt(into *Node, until []int8) []error {
	var errs []error

	forstmt := into.child()
	forstmt.nt = NdFor
	forstmt.gt = p.t
	forstmt.sdata = "For Statement"

	// Skip 'for'
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	// For init
	forinit := forstmt.child()
	forinit.nt = NdForInit

	errs = append(errs, p.rparse(forinit, []int8{TkSemicolon})...)

	// Skip ';'
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	// For bool
	forbool := forstmt.child()
	forbool.nt = NdForBool

	errs = append(errs, p.parseExpr(forbool, []int8{TkSemicolon})...)

	// Skip ';'
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	// For inc
	forinc := forstmt.child()
	forinc.nt = NdForInc

	errs = append(errs, p.rparse(forinc, []int8{TkLBrace})...)

	body := forstmt.child()
	body.nt = NdBody
	body.gt = p.t
	body.sdata = "Body"

	// Skip '{'
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	errs = append(errs, p.rparse(body, []int8{TkRBrace})...)

	// Skip '}'
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	return errs
}

func (p *Parser) parseWhileStmt(into *Node, until []int8) []error {
	var errs []error

	while := into.child()
	while.nt = NdWhile
	while.gt = p.t
	while.sdata = "While Statement"

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	errs = append(errs, p.parseExpr(while, []int8{TkLBrace})...)

	body := while.child()
	body.nt = NdBody
	body.gt = p.t
	body.sdata = "Body"

	// Skip '{'
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	errs = append(errs, p.rparse(body, []int8{TkRBrace})...)

	// Skip '}'
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	return errs
}

func (p *Parser) parseAssign(into *Node, until []int8) []error {
	var errs []error

	asn := into.child()
	asn.nt = NdAssign
	asn.gt = p.t
	asn.sdata = "Assign"

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	bo := asn.child()
	bo.nt = NdBinOp
	bo.gt = p.t

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	body := asn.child()
	body.nt = NdBody

	errs = append(errs, p.rparse(body, []int8{TkWhiteSpace, TkSemicolon})...)

	return errs
}

func (p *Parser) parseEnumLit(into *Node, until []int8) []error {
	var errs []error

	enum := into.child()
	enum.gt = p.t
	enum.nt = NdEnum

	if err := p.assert([]int8{TkLBrace}); err != nil {
		errs = append(errs, err)
		return errs
	}

	for p.t.tokenType != TkRBrace {
		switch p.t.tokenType {
		case TkIdent:
			lit := enum.child()
			lit.gt = p.t
			lit.sdata = p.t.sdata
			lit.nt = NdIdent

		case TkComma, TkWhiteSpace:
			break

		default:
			errs = append(errs, errors.New("Did not expect to find: "+p.t.toString()+" in enum literal."))
		}

		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return errs
		}
	}

	return errs
}

func (p *Parser) parseDeferStmt(into *Node, until []int8) []error {
	var errs []error

	dfr := into.child()
	dfr.nt = NdDefer
	dfr.gt = p.t

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	body := dfr.child()
	body.nt = NdBody

	_until := until
	_until = append(_until, []int8{TkWhiteSpace, TkSemicolon}...)
	errs = append(errs, p.rparse(body, _until)...)

	return errs
}

func (p *Parser) parseCleanStmt(into *Node, until []int8) []error {
	var errs []error

	cln := into.child()
	cln.nt = NdClean
	cln.gt = p.t

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	body := cln.child()
	body.nt = NdBody

	_until := until
	_until = append(_until, []int8{TkWhiteSpace, TkSemicolon}...)
	errs = append(errs, p.rparse(body, _until)...)

	return errs
}

func (p *Parser) parseMatch(into *Node, until []int8) []error {
	var errs []error

	mtch := into.child()
	mtch.nt = NdMatch
	mtch.gt = p.t
	mtch.sdata = "match"

	cond := mtch.child()
	cond.nt = NdBody

	errs = append(errs, p.rparse(cond, []int8{TkLBrace})...)

	// Skip '{'
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	body := mtch.child()
	body.nt = NdBody

	for p.t.tokenType != TkRBrace {
		_cond := body.child()
		_cond.nt = NdLit

		for p.t.tokenType == TkWhiteSpace {
			if err := p.eat(); err != nil {
				errs = append(errs, err)
				break
			}
		}

		if p.t.tokenType == TkEllipsis {
			x := _cond.child()
			x.nt = NdEllipsis
			x.gt = p.t
			x.sdata = "..."

			if err := p.assert([]int8{TkLBrace}); err != nil {
				errs = append(errs, err)
				return errs
			}
		} else {
			errs = append(errs, p.rparse(_cond, []int8{TkLBrace, TkRBrace})...)

			if p.t.tokenType != TkLBrace {
				break
			}
		}

		_body := _cond.child()
		_body.nt = NdBody

		errs = append(errs, p.rparse(_body, []int8{TkRBrace})...)

		if err := p.eat(); err != nil {
			errs = append(errs, err)
			break
		}
	}

	return errs
}

func (p *Parser) parseTernary(into *Node, until []int8) []error {
	var errs []error

	tern := into.child()
	tern.nt = NdTernaryOp
	tern.gt = p.t

	if err := p.assert([]int8{TkLParen}); err != nil {
		errs = append(errs, err)
		return errs
	}

	errs = append(errs, p.parseExpr(tern, []int8{TkRParen})...)

	if err := p.assert([]int8{TkColon}); err != nil {
		errs = append(errs, err)
		return errs
	}

	errs = append(errs, p.rparse(tern, []int8{TkNot})...)

	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	errs = append(errs, p.rparse(tern, []int8{TkWhiteSpace, TkSemicolon})...)

	return errs
}

func (p *Parser) parseClink(into *Node, until []int8) []error {
	var errs []error

	clk := into.child()
	clk.nt = NdClink
	clk.gt = p.t

	if err := p.assert([]int8{TkArrow}); err != nil {
		errs = append(errs, err)
		return errs
	}

	if err := p.assert([]int8{TkIdent}); err != nil {
		errs = append(errs, err)
		return errs
	}

	idn := clk.child()
	idn.nt = NdIdent
	idn.gt = p.t
	idn.sdata = p.t.sdata

	return errs
}

func (p *Parser) rparse(into *Node, until []int8) []error {
	var errs []error

	for {
		for _, u := range until {
			if u == p.t.tokenType || p.t.tokenType == TkEof {
				return errs
			}
		}

		switch p.t.tokenType {
		case TkEof:
			return errs

		case TkImportKw:
			errs = append(errs, p.parseImport(into, until)...)

		case TkIdent:
		getNames:
			var names []Token

			tk, err := p.peek()
			if err != nil {
				errs = append(errs, err)
			}

			for _, u := range until {
				if u == tk.tokenType || tk.tokenType == TkEof {
					call := into.child()
					call.nt = NdCall
					call.gt = p.t
					call.sdata = "VarCall"
					if err := p.eat(); err != nil {
						errs = append(errs, err)
					}
					return errs
				}
			}

			if tk.tokenType == TkComma { // ,
				names = append(names, *p.t)

				// Move to ','
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				// Skip ','
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				goto getNames

			} else if tk.tokenType == TkColon { // :
				names = append(names, *p.t)

				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				errs = append(errs, p.parseVarDecl(into, until, names)...)
			} else if tk.tokenType == TkDefine { // :=
				names = append(names, *p.t)

				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				errs = append(errs, p.parseVarDef(into, until, names)...)
			} else if tk.tokenType == TkOr { // |
				attr := into.child()
				attr.nt = NdAttribute
				attr.gt = p.t
				attr.sdata = p.t.sdata

				// Move to '|'
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				// Skip '|'
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				errs = append(errs, p.rparse(attr, []int8{TkWhiteSpace, TkSemicolon})...)
			} else if tk.tokenType == TkPeriod { // .
				strcall := into.child()
				strcall.nt = NdStructCall
				strcall.gt = p.t
				strcall.sdata = p.t.sdata

				// Move to '.'
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				// Skip '.'
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				errs = append(errs, p.rparse(strcall, []int8{TkWhiteSpace, TkSemicolon})...)
			} else if tk.tokenType == TkLParen { // (
				fcall := into.child()
				fcall.nt = NdFunctionCall
				fcall.gt = p.t
				fcall.sdata = p.t.sdata

				// Move to '('
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				param := fcall.child()
				param.nt = NdFunctionParam
				param.gt = p.t
				param.sdata = "Params"

				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}

				for p.t.tokenType != TkRParen {
					errs = append(errs, p.rparse(param, []int8{TkComma, TkRParen})...)
					if p.t.tokenType == TkComma {
						// Skip ','
						if err := p.eat(); err != nil {
							errs = append(errs, err)
							break
						}
					}
				}
			} else if tk.tokenType == TkInc || tk.tokenType == TkDec || tk.tokenType == TkSqrInc {
				inc := Node{}
				inc.nt = NdIncDec
				inc.gt = tk
				inc.sdata = tk.sdata

				vcall := inc.child()
				vcall.nt = NdCall
				vcall.gt = p.t

				// Move to ++ || -- || **
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}
			} else if tk.tokenType == TkLBrack { // [
				index := into.child()
				index.nt = NdIndex
				index.sdata = "Index"

				varname := index.child()
				varname.gt = p.t
				varname.nt = NdCall

				// Move to the index literal
				if err := p.eat(); err != nil {
					errs = append(errs, err)
					break
				}
				if err := p.eat(); err != nil {
					errs = append(errs, err)
					break
				}

				num := index.child()
				num.nt = NdLit

				errs = append(errs, p.rparse(num, []int8{TkRBrack})...)

				pk, err := p.peek()
				if err != nil {
					errs = append(errs, err)
					break
				}

				body := index.child()
				body.nt = NdBody

				for _, u := range until {
					if u == tk.tokenType || tk.tokenType == TkEof {
						call := into.child()
						call.nt = NdCall
						call.gt = p.t
						call.sdata = "VarCall"
						if err := p.eat(); err != nil {
							errs = append(errs, err)
						}
						return errs
					}
				}

				if pk.tokenType == TkAssign || pk.tokenType == TkOrAssign || pk.tokenType == TkAddAssign || pk.tokenType == TkAndAssign || pk.tokenType == TkMulAssign || pk.tokenType == TkAndNotAssign || pk.tokenType == TkQuoAssign || pk.tokenType == TkRemAssign || pk.tokenType == TkShlAssign || pk.tokenType == TkShrAssign || pk.tokenType == TkSubAssign || pk.tokenType == TkXorAssign {
					errs = append(errs, p.parseAssign(body, []int8{TkWhiteSpace, TkSemicolon})...)
				} else if pk.tokenType != TkWhiteSpace && pk.tokenType != TkSemicolon {
					_until := until
					_until = append(_until, []int8{TkWhiteSpace, TkSemicolon}...)
					errs = append(errs, p.parseExpr(body, _until)...)
				}
			} else if tk.tokenType == TkWhiteSpace || tk.tokenType == TkSemicolon {
				call := into.child()
				call.nt = NdCall
				call.gt = p.t
				call.sdata = "VarCall"
			} else if tk.tokenType == TkAssign || tk.tokenType == TkOrAssign || tk.tokenType == TkAddAssign || tk.tokenType == TkAndAssign || tk.tokenType == TkMulAssign || tk.tokenType == TkAndNotAssign || tk.tokenType == TkQuoAssign || tk.tokenType == TkRemAssign || tk.tokenType == TkShlAssign || tk.tokenType == TkShrAssign || tk.tokenType == TkSubAssign || tk.tokenType == TkXorAssign {
				errs = append(errs, p.parseAssign(into, until)...)
			} else {
				_until := until
				_until = append(_until, []int8{TkWhiteSpace, TkSemicolon}...)
				errs = append(errs, p.parseExpr(into, _until)...)
			}

		case TkFuncKw:
			errs = append(errs, p.parseFuncLit(into, until)...)

		case TkStructKw:
			errs = append(errs, p.parseStructLit(into, until)...)

		case TkIfKw:
			errs = append(errs, p.parseIfStmt(into, until)...)

		case TkForKw:
			errs = append(errs, p.parseForStmt(into, until)...)

		case TkWhileKw:
			errs = append(errs, p.parseWhileStmt(into, until)...)

		case TkDeferKw:
			errs = append(errs, p.parseDeferStmt(into, until)...)

		case TkCleanKw:
			errs = append(errs, p.parseCleanStmt(into, until)...)

		case TkCLinkKw:
			errs = append(errs, p.parseClink(into, until)...)

		case TkMatchKw:
			errs = append(errs, p.parseMatch(into, until)...)

		case TkQuotMrk:
			errs = append(errs, p.parseTernary(into, until)...)

		case TkEnumKw:
			tk, err := p.peek()
			if err != nil {
				errs = append(errs, err)
				return errs
			}

			if tk.tokenType == TkLBrace {
				errs = append(errs, p.parseEnumLit(into, until)...)
			} else if tk.tokenType == TkPeriod {
				enum := into.child()
				enum.nt = NdEnumCall
				enum.gt = p.t

				p.eat() // cannot fail, we've already peeked.

				if err := p.assert([]int8{TkIdent}); err != nil {
					errs = append(errs, err)
					return errs
				}

				call := enum.child()
				call.nt = NdIdent
				call.gt = p.t

			} else {
				errs = append(errs, errors.New("Did not expect to find: "+p.t.toString()+" after enum declaration."))
			}

		case TkLBrack:
			al := into.child()
			al.nt = NdArrLit
			al.gt = p.t

			_until := until
			_until = append(_until, []int8{TkComma, TkRBrack}...)

			for p.t.tokenType != TkRBrack {
				errs = append(errs, p.rparse(al, _until)...)
				if p.t.tokenType == TkComma {
					if err := p.eat(); err != nil {
						errs = append(errs, err)
						break
					}
				} else if p.t.tokenType != TkRBrack {
					errs = append(errs, errors.New("Didn't expect to find token: "+p.t.toString()))
					break
				}
			}

		case TkReturnKw:
			ret := into.child()
			ret.nt = NdReturn
			ret.gt = p.t
			ret.sdata = "Return"

			errs = append(errs, p.rparse(ret, []int8{TkWhiteSpace, TkSemicolon})...)

		case TkBreakKw:
			brk := into.child()
			brk.nt = NdBreak
			brk.gt = p.t
			brk.sdata = "Break"

			if err := p.assert([]int8{TkWhiteSpace, TkSemicolon}); err != nil {
				errs = append(errs, err)
			}

		case TkContinueKw:
			con := into.child()
			con.nt = NdContinue
			con.gt = p.t
			con.sdata = "Continue"

			if err := p.assert([]int8{TkWhiteSpace, TkSemicolon}); err != nil {
				errs = append(errs, err)
			}

		case TkArrow:
			if err := p.assert([]int8{TkString}); err != nil {
				errs = append(errs, err)
				break
			}

			inln := into.child()
			inln.nt = NdInln
			inln.gt = p.t

		case TkEllipsis:
			el := into.child()
			el.nt = NdEllipsis
			el.gt = p.t
			el.sdata = "..."

		case TkInt, TkFloat, TkString, TkChar, TkNilKw, TkLParen, TkNot:
			errs = append(errs, p.parseExpr(into, until)...)

		case TkWhiteSpace, TkSemicolon:
			break

		case TkAbreviation:
			abr := into.child()
			abr.nt = NdAbbreviation
			abr.gt = p.t
			abr.sdata = p.t.sdata

		default:
			errs = append(errs, errors.New("Did not expect to find token: "+p.t.toString()))
		}

		for _, u := range until {
			if u == p.t.tokenType || p.t.tokenType == TkEof {
				return errs
			}
		}

		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return errs
		}
	}
}
