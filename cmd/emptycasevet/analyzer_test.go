package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tdir := analysistest.TestData()
	analysistest.Run(t, tdir, Analyzer, "a")
}

func TestAnalyzerAllowHeaderComment(t *testing.T) {
	// enable the flag to allow inline header comments
	_ = Analyzer.Flags.Set("allow_header_comment", "true")
	defer Analyzer.Flags.Set("allow_header_comment", "false")
	
	tdir := analysistest.TestData()
	analysistest.Run(t, tdir, Analyzer, "b")
}
