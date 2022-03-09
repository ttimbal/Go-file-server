package utils

import "fmt"

const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorPurple  = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorHelp    = ColorWhite
	ColorSuccess = ColorGreen
	ColorError   = ColorRed
	ColorWarning = ColorYellow
)

func PrintReset() {
	fmt.Print(ColorReset)
}

func PrintError(args ...string) {
	fmt.Print(ColorError)
	for _, arg := range args {
		fmt.Println(arg)
	}
	fmt.Println(ColorReset)
}

func PrintWarning(args ...string) {
	fmt.Print(ColorWarning)
	fmt.Println(args)
	fmt.Println(ColorReset)
}
func PrintSuccess(args ...string) {
	fmt.Print(ColorSuccess)
	fmt.Println(args)
	fmt.Println(ColorReset)
}
func PrintHelp(args ...string) {
	fmt.Print(ColorHelp)
	for _, arg := range args {
		fmt.Println(arg)
	}
	fmt.Println(ColorReset)
}
