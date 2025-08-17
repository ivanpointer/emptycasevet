package allowheadercomment

func g(x int) {
	switch x {
	case 1: // header-only comment should be allowed when flag is enabled
	default:
	}
}
