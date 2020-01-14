package nolol

import (
	"fmt"

	"github.com/dbaumgarten/yodk/pkg/nolol/nast"
	"github.com/dbaumgarten/yodk/pkg/parser"
	"github.com/dbaumgarten/yodk/pkg/parser/ast"
)

// Printer can generate the nolol-code corresponding to a nolol ast
type Printer struct {
	yololPrinter parser.Printer
	Indentation  string
}

// NewPrinter creates a new Printer
func NewPrinter() *Printer {
	return &Printer{
		Indentation: "\t",
	}
}

// Print returns the nolol-code for the given ast
func (p *Printer) Print(prog *nast.Program) (string, error) {
	indentLevel := 0

	indentation := func(amount int) string {
		ind := ""
		for i := 0; i < amount; i++ {
			ind += p.Indentation
		}
		return ind
	}

	p.yololPrinter.UnknownHandlerFunc = func(node ast.Node, visitType int) (string, error) {
		switch n := node.(type) {
		case *nast.GoToLabelStatement:
			if visitType == ast.SingleVisit {
				return "goto " + n.Label, nil
			}
			return "", nil

		case *nast.MultilineIf:
			switch visitType {
			case ast.PreVisit:
				return indentation(indentLevel) + "if ", nil
			case ast.InterVisit1:
				indentLevel++
				return " then\n", nil
			case ast.InterVisit2:
				return indentation(indentLevel-1) + "else\n", nil
			case ast.PostVisit:
				indentLevel--
				return indentation(indentLevel) + "end\n", nil
			default:
				return "", nil
			}
		case *nast.WhileLoop:
			switch visitType {
			case ast.PreVisit:
				return indentation(indentLevel) + "while ", nil
			case ast.InterVisit1:
				indentLevel++
				return " do\n", nil
			case ast.PostVisit:
				indentLevel--
				return indentation(indentLevel) + "end\n", nil
			default:
				return "", nil
			}
		case *nast.StatementLine:
			switch visitType {
			case ast.PreVisit:
				out := indentation(indentLevel)
				if n.Label != "" {
					out += n.Label + "> "
				}
				return out, nil
			default:
				return "", nil
			}
		case *nast.ConstDeclaration:
			switch visitType {
			case ast.PreVisit:
				return "const " + n.Name + " = ", nil
			case ast.PostVisit:
				return "", nil
			}
		case *nast.Program:
			return "", nil
		}
		return "", fmt.Errorf("Unknown node-type: %T", node)
	}

	return p.yololPrinter.PrintCommented(prog, prog.Comments)
}
