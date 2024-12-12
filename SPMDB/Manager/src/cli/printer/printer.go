package printer

import (
	"fmt"
	"spdb_manager/cli/ccolour"
)

func PrintNormalText(msg string) {
	fmt.Println(msg)
}

func PrintNotice(msg string) {
	fmt.Println(ccolour.ApplyColourToText(msg, ccolour.Grey, ccolour.BgBlack))
}

func PrintError(msg string) {
	fmt.Println(ccolour.ApplyColourToText(msg, ccolour.Red, ccolour.BgBlack))
}

func PrintInfo(msg string) {
	fmt.Println(ccolour.ApplyColourToText(msg, ccolour.Yellow, ccolour.BgBlack))
}

func PrintSuccess(msg string) {
	fmt.Println(ccolour.ApplyColourToText(msg, ccolour.Green, ccolour.BgBlack))
}
