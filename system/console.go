package system

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var clearFunctions map[string]func() = map[string]func(){
	"linux": func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	},
	"windows": func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	},
}

// ConsoleClear clears the console.
func ConsoleClear() error {
	clear, ok := clearFunctions[runtime.GOOS]

	if !ok {
		return fmt.Errorf("OS not supported")
	}

	clear()

	return nil
}
