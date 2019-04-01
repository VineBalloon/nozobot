package helpers

// Italics Surrounds with italics escapes
func Italics(s string) string {
	return "*" + s + "*"
}

// Bold Surrounds with bold escapes
func Bold(s string) string {
	return "**" + s + "**"
}

// Code Surrounds with code escapes
func Code(s string) string {
	return "`" + s + "`"
}

// Spoiler Surrounds with spoiler escapes
func Spoiler(s string) string {
	return "||" + s + "||"
}

// Noembed Surrounds with no-embed escapes
func Noembed(s string) string {
	return "<" + s + ">"
}

// At Surrounds ID with @ escapes
func At(s string) string {
	return "<@" + s + ">"
}

// Chan Surrounds ID with chan escapes
func Chan(s string) string {
	return "<#" + s + ">"
}

// Stringinslice Checks if string is in the slice
func Stringinslice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
