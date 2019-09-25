package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type YololGenerator struct {
	programm string
}

var operatorPriority = map[string]int{
	"or":  0,
	"and": 0,
	"==":  1,
	">=":  1,
	"<=":  1,
	">":   1,
	"<":   1,
	"+":   2,
	"-":   2,
	"*":   3,
	"/":   3,
	"^":   3,
	"%":   3,
	"not": 4,
}

func (y *YololGenerator) Visit(node Node, visitType int) error {
	switch n := node.(type) {
	case *Line:
		if visitType == PostVisit {
			y.programm += "\n"
		}
		if visitType > 0 {
			y.programm += " "
		}
		break
	case *Assignment:
		if visitType == PreVisit {
			y.programm += n.Variable + n.Operator
		}
		break
	case *IfStatement:
		y.generateIf(visitType)
		break
	case *GoToStatement:
		y.programm += "goto " + strconv.Itoa(n.Line)
		break
	case *Dereference:
		y.genDeref(n)
		break
	case *StringConstant:
		y.programm += "\"" + n.Value + "\""
		break
	case *NumberConstant:
		if strings.HasPrefix(n.Value, "-") {
			y.programm += " "
		}
		y.programm += fmt.Sprintf(n.Value)
		break
	case *BinaryOperation:
		y.generateBinaryOperation(n, visitType)
		break
	case *UnaryOperation:
		_, childBinary := n.Exp.(*BinaryOperation)
		if visitType == PreVisit {
			op := n.Operator
			if op == "not" {
				op = " " + op + " "
			}
			if op == "-" {
				op = " " + op
			}
			y.programm += op
			if childBinary {
				y.programm += "("
			}
		}
		if visitType == PostVisit {
			if childBinary {
				y.programm += ")"
			}
		}
		break
	case *FuncCall:
		if visitType == PreVisit {
			y.programm += n.Function + "("
		} else {
			y.programm += ")"
		}
		break
	case *Programm:
		//do noting
		break
	default:
		y.programm += fmt.Sprintf("Unknown ast-node: %T%v", node, node)
	}
	return nil
}

func (y *YololGenerator) generateBinaryOperation(o *BinaryOperation, visitType int) {
	lPrio := priorityForExpression(o.Exp1)
	rPrio := priorityForExpression(o.Exp2)
	_, rBinary := o.Exp2.(*BinaryOperation)
	myPrio := priorityForExpression(o)
	switch visitType {
	case PreVisit:
		if lPrio < myPrio {
			y.programm += "("
		}
		break
	case InterVisit1:
		if lPrio < myPrio {
			y.programm += ")"
		}
		op := o.Operator
		if op == "and" || op == "or" {
			op = " " + op + " "
		}
		y.programm += op
		if rBinary && rPrio <= myPrio {
			y.programm += "("
		}
		break
	case PostVisit:
		if rBinary && rPrio <= myPrio {
			y.programm += ")"
		}
		break
	}
}

func priorityForExpression(e Expression) int {
	switch ex := e.(type) {
	case *BinaryOperation:
		return operatorPriority[ex.Operator]
	default:
		return 10
	}
}

func (y *YololGenerator) generateIf(visitType int) {
	switch visitType {
	case PreVisit:
		y.programm += "if "
	case InterVisit1:
		y.programm += " then "
	case InterVisit2:
		y.programm += " else "
	case PostVisit:
		y.programm += " end"
	default:
		y.programm += " "
	}
}

func (y *YololGenerator) genDeref(d *Dereference) {
	txt := ""
	if d.PrePost == "Pre" {
		txt += " " + d.Operator
	}
	txt += d.Variable
	if d.PrePost == "Post" {
		txt += d.Operator + " "
	}
	y.programm += txt
}

func (y *YololGenerator) Generate(prog *Programm) string {
	y.programm = ""
	prog.Accept(y)
	// during the generation duplicate spaces might appear. Remove them
	return strings.Replace(y.programm, "  ", " ", -1)
}
