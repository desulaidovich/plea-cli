package main

import "plea-cli/cmd"

func main() {
	if err := cmd.Run(); err != nil {
		panic("failed to run cli: " + err.Error())
	}
}
