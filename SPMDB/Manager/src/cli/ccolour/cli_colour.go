package ccolour

// https://stackoverflow.com/questions/2048509/how-to-echo-with-different-colors-in-the-windows-command-line

const ESC = "\u001B"

type Colour string

const (
	Black       = Colour("30")
	Grey        = Colour("90")
	Red         = Colour("31")
	Pink        = Colour("91")
	Green       = Colour("32")
	LightGreen  = Colour("92")
	Yellow      = Colour("33")
	LightYellow = Colour("93")
	Blue        = Colour("34")
	Violet      = Colour("94")
	Magenta     = Colour("35")
	LightPurple = Colour("95")
	Cyan        = Colour("36")
	LightCyan   = Colour("96")
	White       = Colour("37")
	BrightWhite = Colour("97")

	BgBlack       = Colour("40")
	BgGrey        = Colour("100")
	BgRed         = Colour("41")
	BgPink        = Colour("101")
	BgGreen       = Colour("42")
	BgLightGreen  = Colour("102")
	BgYellow      = Colour("43")
	BgLightYellow = Colour("103")
	BgBlue        = Colour("44")
	BgViolet      = Colour("104")
	BgMagenta     = Colour("45")
	BgLightPurple = Colour("105")
	BgCyan        = Colour("46")
	BgLightCyan   = Colour("106")
	BgWhite       = Colour("47")
	BgBrightWhite = Colour("107")

	EndOfLine = Colour(ESC + "[0m")
)

func ApplyColourToText(s string, textColour Colour, bgColour Colour) string {
	return ESC + "[" + string(textColour) + ";" + string(bgColour) + "m" + s + string(EndOfLine)
}
