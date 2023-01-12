package main

import (
	"errors"
)

const (
	NdStart         = "Start"
	NdIdent         = "Ident"
	NdEllipsis      = "Ellipsis"
	NdFunctionDef   = "Function Def"
	NdFunctionParam = "Function Param"
	NdBody          = "Function NdBody"
	NdFunctionCall  = "Function Call"
	NdStructDef     = "Struct Def"
	NdStructCall    = "Struct Call"
	NdReturn        = "Return"
	NdDecl          = "Declare"
	NdDefine        = "Define"
	NdCall          = "Call"
	NdFor           = "For"
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
		if n.nt != NdStart && !endNode {
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

/** Unused
func (p *Parser) barf() error {
	if p.tI < 1 {
		return errors.New("Could not barf to negative index")
	}

	p.tI--
	p.t = &p.tks[p.tI]

	return nil
}
*/

func (p *Parser) parse() (*Node, []error) {
	var errs []error

	p.hd.nt = NdStart
	p.hd.sdata = "_start"
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
	case TkIntKw, TkUIntKw, TkUInt32Kw, TkInt32Kw, TkFloat32Kw:
		bsize = 4
	case TkInt64Kw, TKUint64Kw, TkFloat64Kw, TkStringKw:
		bsize = 8
	case TkInt8Kw, TkUInt8Kw, TkByteKw, TkBoolKw:
		bsize = 1
	case TkInt16Kw, TkUInt16Kw:
		bsize = 2
	case TkLBrack:
		// [10, int]

		// Skip '['
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
		nameN.gt = &name
		nameN.sdata = name.sdata
	}

	// Skip ':'
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	vardcl.bytes = p.typeToBSize() * uint(len(names))
	if vardcl.bytes == 0 {
		if p.t.tokenType == TkEllipsis {
			el := vardcl.child()
			el.nt = NdEllipsis
			el.gt = p.t
			el.sdata = "..."
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
		if err := p.assert([]int8{TkConstKw}); err != nil {
			errs = append(errs, err)
			return errs
		}

		cnst := vardcl.child()
		cnst.gt = p.t
		cnst.nt = NdSect
		cnst.sdata = "const"

		// Skip 'const'
		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return errs
		}
	}

	if p.t.tokenType == TkAssign {
		errs = append(errs, p.rparse(into, []int8{TkWhiteSpace, TkSemicolon})...)
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
		nameN.gt = &name
		nameN.sdata = name.sdata
	}

	// Skip ':='
	if err := p.eat(); err != nil {
		errs = append(errs, err)
		return errs
	}

	errs = append(errs, p.rparse(into, []int8{TkWhiteSpace, TkSemicolon})...)

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

		// Skip type, go to '{'
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

	errs = append(errs, p.rparse(body, []int8{TkRBrace})...)

	return errs
}

func (p *Parser) parseExpr(into *Node, until []int8) []error {
	var errs []error

	var queue []Node
	var stack []Node

	for {
		if p.t.tokenType == TkInt || p.t.tokenType == TkFloat {
			lit := Node{}
			lit.nt = NdLit
			lit.gt = p.t
			queue = append(queue, lit)
		} else if p.t.tokenType == TkIdent {
			tk, err := p.peek()
			if err != nil {
				errs = append(errs, err)
				return errs
			}

			if tk.tokenType == TkLParen {
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
			} else {
				call := Node{}
				call.nt = NdCall
				call.gt = p.t
				call.sdata = "VarCall"
				queue = append(queue, call)
			}
		} else if p.t.tokenType == TkAdd || p.t.tokenType == TkSub || p.t.tokenType == TkMul || p.t.tokenType == TkQuo || p.t.tokenType == TkRem || p.t.tokenType == TkEql || p.t.tokenType == TkNeq || p.t.tokenType == TkLeq || p.t.tokenType == TkLss || p.t.tokenType == TkGeq || p.t.tokenType == TkGtr || p.t.tokenType == TkLand || p.t.tokenType == TkLor {
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
				break
			}
		} else {
			break
		}

		if err := p.eat(); err != nil {
			errs = append(errs, err)
			return errs
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1].gt.tokenType == TkLParen {
			errs = append(errs, errors.New("Mismatched parenthesis: "+p.t.toString()))
			break
		}

		queue = append(queue, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	var i = len(queue) - 1
	errs = append(errs, p.rcompute(into, queue, &i)...)

	return errs
}

func (p *Parser) rcompute(into *Node, stack []Node, i *int) []error {
	var errs []error

	var passes uint = 0

	for ; *i >= 0; *i-- {
		if passes >= 2 {
			break
		}

		n := stack[*i]

		switch n.nt {
		case NdLit, NdCall, NdFunctionCall:
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

			body := elstmt.child()
			body.nt = NdBody
			body.gt = p.t
			body.sdata = "Body"

			if err := p.eat(); err != nil {
				errs = append(errs, err)
				break
			}

			errs = append(errs, p.rparse(body, []int8{TkRBrace})...)
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

			body := elstmt.child()
			body.nt = NdBody
			body.gt = p.t
			body.sdata = "Body"

			errs = append(errs, p.rparse(body, []int8{TkRBrace})...)
		}

		if err := p.eat(); err != nil {
			errs = append(errs, err)
			break
		}
	}

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
		case TkImportKw:
			errs = append(errs, p.parseImport(into, until)...)

		case TkIdent:
		getNames:
			var names []Token

			tk, err := p.peek()
			if err != nil {
				errs = append(errs, err)
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
				inc := into.child()
				inc.nt = NdIncDec
				inc.gt = p.t
				inc.sdata = p.t.sdata

				// Move to ++ || -- || **
				if err := p.eat(); err != nil {
					errs = append(errs, err)
				}
			} else if tk.tokenType == TkWhiteSpace || tk.tokenType == TkSemicolon {
				call := into.child()
				call.nt = NdCall
				call.gt = p.t
				call.sdata = "VarCall"
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
