package stringx

func NonEmpty(str ...string) string {
	for _, s := range str {
		if len(s) > 0 {
			return s
		}
	}

	return ""
}
