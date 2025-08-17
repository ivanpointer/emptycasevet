// Package emptycasevet provides a go/analysis Analyzer that reports
// non-default switch cases with empty bodies.
package emptycasevet

import (
	"flag"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// allowHeaderComment controls whether an inline comment on the case header line
// (same line as the colon) should be considered as a non-empty body.
var allowHeaderComment bool

// Analyzer reports empty non-default case clauses in switch/type switch statements.
var Analyzer = &analysis.Analyzer{
	Name: "emptycase",
	Doc:  "reports non-default switch cases with empty bodies",
	Run: func(pass *analysis.Pass) (any, error) {
		for _, f := range pass.Files {
			ast.Inspect(f, func(n ast.Node) bool {
				switch s := n.(type) {
				case *ast.SwitchStmt:
					checkCases(pass, f, s.Body.List)
				case *ast.TypeSwitchStmt:
					checkCases(pass, f, s.Body.List)
				}
				return true
			})
		}
		return nil, nil
	},
	Flags: flag.FlagSet{},
}

func init() {
	Analyzer.Flags.BoolVar(&allowHeaderComment, "allow_header_comment", false, "consider an inline comment on the case header line as a non-empty body")
}

func checkCases(pass *analysis.Pass, file *ast.File, stmts []ast.Stmt) {
	for i, st := range stmts {
		cc, ok := st.(*ast.CaseClause)
		if !ok {
			continue
		}
		// default case is allowed to be empty
		if cc.List == nil {
			continue
		}

		// determine the boundary of this clause: start at after colon, end at start of next clause or end of block
		var nextPos token.Pos
		if i+1 < len(stmts) {
			nextPos = stmts[i+1].Pos()
		} else {
			// best effort: use the end of the case clause node; if empty body, comments may be outside cc.End
			nextPos = cc.End()
		}

		// non-default case can't be empty unless it contains a comment explaining intent
		if len(cc.Body) == 0 && !hasInlineOrInnerComment(pass, file, cc, nextPos) {
			pass.Reportf(cc.Case, "empty case body; did you mean `case a, b:`?")
		}
	}
}

func hasInlineOrInnerComment(pass *analysis.Pass, file *ast.File, cc *ast.CaseClause, boundary token.Pos) bool {
	for _, cg := range file.Comments {
		// comments inside body
		if cg.Pos() >= cc.Pos() && cg.End() <= cc.End() {
			return true
		}
		// comments at or after the colon but before the next case/default clause
		if cg.Pos() > cc.Colon && cg.Pos() < boundary {
			caseLine := pass.Fset.PositionFor(cc.Colon, false).Line
			commentLine := pass.Fset.PositionFor(cg.Pos(), false).Line
			if commentLine > caseLine {
				// on a subsequent line -> always allowed
				return true
			}
			if allowHeaderComment && commentLine == caseLine {
				// inline header comment allowed only if flag is set
				return true
			}
		}
	}
	return false
}
