package colorise

const (
	Green  = "\033[32m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
)

type Color int

const (
	ColorGreen Color = iota
	ColorRed
	ColorYellow
	ColorReset
)

func ColorString(s string, color Color) string {
	switch color {
	case ColorGreen:
		return Green + s + Reset
	case ColorRed:
		return Red + s + Reset
	case ColorYellow:
		return Yellow + s + Reset
	case ColorReset:
		return Reset + s + Reset
	default:
		return s
	}
}
