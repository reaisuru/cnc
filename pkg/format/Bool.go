package format

func ColorizeBool(v bool) string {
	if v {
		return "\x1b[32mtrue\x1b[0m"
	}

	return "\x1b[31mfalse\x1b[0m"
}
