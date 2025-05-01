package main

import (
	"fmt"
	"os"
	"os/exec"
)

func ejecutarHerramienta(nombre, argumento string) error {
	cmd := exec.Command(nombre, argumento) // Llama al ejecutable
	output, err := cmd.CombinedOutput()    // Obtiene la salida
	if err != nil {
		return fmt.Errorf("error al ejecutar %s: %v", nombre, err)
	}
	fmt.Println("Salida de", nombre, ":", string(output))
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Se necesitan al menos 2 argumentos: herramienta y argumento.")
		os.Exit(1)
	}

	herramienta := os.Args[1]
	argumento := os.Args[2]

	err := ejecutarHerramienta(herramienta, argumento)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
