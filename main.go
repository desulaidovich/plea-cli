package main

import "github.com/desulaidovich/plea-cli/cmd"

func main() {
	if err := cmd.Run(); err != nil {
		panic("failed to run cli: " + err.Error())
	}
}
