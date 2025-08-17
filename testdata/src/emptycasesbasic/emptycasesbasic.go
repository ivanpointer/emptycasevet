package emptycasesbasic

func f(x any) {
	switch x.(type) {
	case int:
		// ok: has comment
	case string: // want "empty case body; did you mean `case a, b:`?"
	default:
	}

	switch x := 1; x {
	case 1: // want "empty case body; did you mean `case a, b:`?"
	case 2:
		// allowed with comment
	default:
		// empty default allowed
	}
}
