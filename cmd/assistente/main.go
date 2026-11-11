// Package main é o ponto de entrada do binário assistente.
//
// Este arquivo é intencionalmente mínimo: ele apenas delega a execução
// para o pacote internal/cli, que concentra toda a orquestração testável.
// Qualquer lógica de negócio, parsing de flags ou configuração de comandos
// deve viver em internal/cli — nunca aqui.
package main

import (
	"fmt"
	"os"

	"github.com/DouglasLira-Dev/assistente-cli/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
