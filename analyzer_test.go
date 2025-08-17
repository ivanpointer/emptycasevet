package emptycasevet

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer_BasicEmptyCases(t *testing.T) {
	tdir := analysistest.TestData()
	analysistest.Run(t, tdir, Analyzer, "emptycasesbasic")
}

func TestAnalyzer_AllowHeaderComment(t *testing.T) {
	// enable the flag to allow inline header comments
	_ = Analyzer.Flags.Set("allow_header_comment", "true")
	defer Analyzer.Flags.Set("allow_header_comment", "false")

	tdir := analysistest.TestData()
	analysistest.Run(t, tdir, Analyzer, "allowheadercomment")
}
