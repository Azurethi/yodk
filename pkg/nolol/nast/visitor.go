package nast

import (
	"github.com/dbaumgarten/yodk/pkg/parser/ast"
)

// Accept is used to implement Acceptor
func (g *GoToLabelStatement) Accept(v ast.Visitor) error {
	return v.Visit(g, ast.SingleVisit)
}

// Accept is used to implement Acceptor
func (p *Program) Accept(v ast.Visitor) error {
	err := v.Visit(p, ast.PreVisit)
	if err != nil {
		return err
	}
	for i := 0; i < len(p.Lines); i++ {
		err = v.Visit(p, i)
		if err != nil {
			return err
		}
		err = p.Lines[i].Accept(v)
		if repl, is := err.(ast.NodeReplacement); is {
			p.Lines = patchLines(p.Lines, i, repl)
			i += len(repl.Replacement) - 1
			err = nil
		}
		if err != nil {
			return err
		}
	}
	return v.Visit(p, ast.PostVisit)
}

// Accept is used to implement Acceptor
func (l *StatementLine) Accept(v ast.Visitor) error {
	err := v.Visit(l, ast.PreVisit)
	if err != nil {
		return err
	}

	for i := 0; i < len(l.Statements); i++ {
		err = v.Visit(l, i)
		if err != nil {
			return err
		}
		err = l.Statements[i].Accept(v)
		if repl, is := err.(ast.NodeReplacement); is {
			l.Statements = ast.PatchStatements(l.Statements, i, repl)
			i += len(repl.Replacement) - 1
			err = nil
		}
		if err != nil {
			return err
		}
	}

	return v.Visit(l, ast.PostVisit)
}

// Accept is used to implement Acceptor
func (l *ConstDeclaration) Accept(v ast.Visitor) error {
	err := v.Visit(l, ast.PreVisit)
	if err != nil {
		return err
	}
	err = l.Value.Accept(v)
	if repl, is := err.(ast.NodeReplacement); is {
		l.Value = repl.Replacement[0].(ast.Expression)
		err = nil
	}
	if err != nil {
		return err
	}
	return v.Visit(l, ast.PostVisit)
}

// Accept is used to implement Acceptor
func (s *Block) Accept(v ast.Visitor) error {
	err := v.Visit(s, ast.PreVisit)
	if err != nil {
		return err
	}
	for i := 0; i < len(s.Lines); i++ {
		err = v.Visit(s, i)
		if err != nil {
			return err
		}
		err = s.Lines[i].Accept(v)
		if repl, is := err.(ast.NodeReplacement); is {
			s.Lines = patchExecutableLines(s.Lines, i, repl)
			i += len(repl.Replacement) - 1
			err = nil
		}
		if err != nil {
			return err
		}
	}
	err = v.Visit(s, ast.PostVisit)
	if err != nil {
		return err
	}
	return nil
}

// Accept is used to implement Acceptor
func (s *MultilineIf) Accept(v ast.Visitor) error {
	err := v.Visit(s, ast.PreVisit)
	if err != nil {
		return err
	}

	for i := 0; i < len(s.Conditions); i++ {
		err = v.Visit(s, i)
		if err != nil {
			return err
		}
		err = s.Conditions[i].Accept(v)
		if repl, is := err.(ast.NodeReplacement); is {
			s.Conditions[i] = repl.Replacement[0].(ast.Expression)
			err = nil
		}
		if err != nil {
			return err
		}
		err = v.Visit(s, ast.InterVisit1)
		if err != nil {
			return err
		}
		err = s.Blocks[i].Accept(v)
		if repl, is := err.(ast.NodeReplacement); is {
			s.Blocks[i] = repl.Replacement[0].(*Block)
			err = nil
		}
		if err != nil {
			return err
		}
	}
	if s.ElseBlock != nil {
		err = v.Visit(s, ast.InterVisit2)
		if err != nil {
			return err
		}
		err = s.ElseBlock.Accept(v)
		if repl, is := err.(ast.NodeReplacement); is {
			s.ElseBlock = repl.Replacement[0].(*Block)
			err = nil
		}
		if err != nil {
			return err
		}
	}
	return v.Visit(s, ast.PostVisit)
}

// Accept is used to implement Acceptor
func (s *WhileLoop) Accept(v ast.Visitor) error {
	err := v.Visit(s, ast.PreVisit)
	if err != nil {
		return err
	}
	err = s.Condition.Accept(v)
	if repl, is := err.(ast.NodeReplacement); is {
		s.Condition = repl.Replacement[0].(ast.Expression)
		err = nil
	}
	if err != nil {
		return err
	}
	err = v.Visit(s, ast.InterVisit1)
	if err != nil {
		return err
	}
	err = s.Block.Accept(v)
	if repl, is := err.(ast.NodeReplacement); is {
		s.Block = repl.Replacement[0].(*Block)
		err = nil
	}
	if err != nil {
		return err
	}
	return v.Visit(s, ast.PostVisit)
}

// Accept is used to implement Acceptor
func (s *WaitStatement) Accept(v ast.Visitor) error {
	err := v.Visit(s, ast.PreVisit)
	if err != nil {
		return err
	}
	err = s.Condition.Accept(v)
	if repl, is := err.(ast.NodeReplacement); is {
		s.Condition = repl.Replacement[0].(ast.Expression)
		err = nil
	}
	if err != nil {
		return err
	}
	return v.Visit(s, ast.PostVisit)
}

// Accept is used to implement Acceptor
func (s *IncludeDirective) Accept(v ast.Visitor) error {
	return v.Visit(s, ast.SingleVisit)
}

func patchLines(old []Line, position int, repl ast.NodeReplacement) []Line {
	newv := make([]Line, 0, len(old)+len(repl.Replacement)-1)
	newv = append(newv, old[:position]...)
	for _, elem := range repl.Replacement {
		if line, is := elem.(Line); is {
			newv = append(newv, line)
		} else {
			panic("Could not patch slice. Wrong type.")
		}
	}
	newv = append(newv, old[position+1:]...)
	return newv
}

func patchExecutableLines(old []ExecutableLine, position int, repl ast.NodeReplacement) []ExecutableLine {
	newv := make([]ExecutableLine, 0, len(old)+len(repl.Replacement)-1)
	newv = append(newv, old[:position]...)
	for _, elem := range repl.Replacement {
		if line, is := elem.(ExecutableLine); is {
			newv = append(newv, line)
		} else {
			panic("Could not patch slice. Wrong type.")
		}
	}
	newv = append(newv, old[position+1:]...)
	return newv
}
