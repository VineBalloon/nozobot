package helpers

func Surround(s string, d string) string {
	out := d + s
	c := []byte(d)

	for i := len(c)/2 - 1; i >= 0; i-- {
		opp := len(c) - 1 - i
		c[i], c[opp] = c[opp], c[i]
	}
	return out + string(c)
}

func Italics(s string) string {
	return Surround(s, "*")
}

func Bold(s string) string {
	return Surround(s, "**")
}

func Code(s string) string {
	return Surround(s, "`")
}

func Noembed(s string) string {
	return "<" + s + ">"
}

func At(s string) string {
	return "<@" + s + ">"
}

func Chan(s string) string {
	return "<#" + s + ">"
}
