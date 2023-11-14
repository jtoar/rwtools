package fw

import (
	"fmt"
	"os/exec"
)

func PrintErrMsg() {
	fmt.Println("Expected one of: clean")
}

func Clean() {
	exec.Command("git clean --force -xd --exclude .env").Run()
	exec.Command("yarn install").Run()
	exec.Command("yarn check").Run()
}
