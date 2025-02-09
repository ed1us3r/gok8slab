package utils

import "fmt"

// ANSI Color Codes
const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Reset  = "\033[0m"
)

// Colorized output helpers
func Info(msg string) {
	fmt.Println(Blue + "🟦 [INFO]    " + msg + Reset)
}

func Success(msg string) {
	fmt.Println(Green + "🟩 [SUCCESS] " + msg + Reset)
}

func Warning(msg string) {
	fmt.Println(Yellow + "⚠️ [WARNING]  " + msg + Reset)
}

func Error(msg string) {
	fmt.Println(Red + "⛔[ERROR]    " + msg + Reset)
}
