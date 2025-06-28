package config

import (
	"fmt"

	"github.com/Liphium/magic/mconfig"
)

// This is the function called once you run the project
func Run(ctx *mconfig.Context) {
	fmt.Println("Generating config..")
}

func Start() {
	fmt.Println("Hello magic!")
}
