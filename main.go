package main

import (
	"fmt"
	configloader "x46/confy/internal/configLoader"
)

func main() {
	fmt.Println("Startring confy...")

	config, err := configloader.InitConfig()
	if err != nil {
		fmt.Println("Error initializing config:", err)
		return
	}

	fmt.Println("Config initialized:", config)
}
