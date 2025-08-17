package main

import (
	emptycasevet "github.com/ivanpointer/emptycasevet"
	"golang.org/x/tools/go/analysis/singlechecker"
)

// Analyzer is kept here for backward compatibility and tests in this folder.
var Analyzer = emptycasevet.Analyzer

func main() { singlechecker.Main(emptycasevet.Analyzer) }
