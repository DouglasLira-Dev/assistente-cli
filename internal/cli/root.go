// Package cli concentra a orquestração do assistente.
// Toda a lógica de comandos, flags e configuração de execução vive aqui,
// mantendo cmd/assistente/main.go como um entry point descartável.
package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/DouglasLira-Dev/assistente-cli/internal/config"
)

// jsonOutput controla a saída em JSON para todos os subcomandos.
// É registrada como flag persistente no comando raiz, então todo
// subcomando herda --json automaticamente.
var jsonOutput bool

// cfg é a configuração global do assistente.
// Carregada uma única vez pelo PersistentPreRunE do comando raiz.
var cfg *config.Config

// Execute é o ponto de entrada do CLI.
// Cria o comando raiz, executa e retorna o erro (sem imprimir nada —
// a impressão em stderr e o exit code ficam a cargo do main.go).
func Execute() error {
	rootCmd := newRootCmd()
	return rootCmd.Execute()
}

// newRootCmd constrói a árvore de comandos do assistente.
func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "assistente",
		Short: "Assistente de linha de comando para conhecimento pessoal",
		Long: `Assistente CLI centraliza artigos, templates e conversas em um único lugar.

Fluxo principal: artigo → memória → API.
Consulte artigos locais, faça perguntas ao LLM e reindexe seu acervo.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Help e completion não exigem configuração.
			if cmd.Name() == "help" || cmd.Name() == "completion" {
				return nil
			}
			if help, _ := cmd.Flags().GetBool("help"); help {
				return nil
			}

			loaded, err := config.Load()
			if err != nil {
				return err
			}
			cfg = loaded
			return nil
		},
	}

	// Flag persistente: herdada por todos os subcomandos.
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "saída em formato JSON")

	// Subcomandos (stubs neste módulo).
	rootCmd.AddCommand(newAskCmd())
	rootCmd.AddCommand(newArticleCmd())
	rootCmd.AddCommand(newReindexCmd())

	return rootCmd
}

// newAskCmd cria o subcomando "ask" (stub).
func newAskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ask <pergunta>",
		Short: "Faz uma pergunta ao assistente",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("comando 'ask' ainda não implementado")
		},
	}
}

// newArticleCmd cria o subcomando "article" com seus filhos (stub).
func newArticleCmd() *cobra.Command {
	articleCmd := &cobra.Command{
		Use:   "article",
		Short: "Gerencia artigos locais",
	}

	articleCmd.AddCommand(newArticleAddCmd())
	articleCmd.AddCommand(newArticleListCmd())

	return articleCmd
}

// newArticleAddCmd cria o subcomando "article add" (stub).
func newArticleAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Adiciona um novo artigo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("comando 'article add' ainda não implementado")
		},
	}
}

// newArticleListCmd cria o subcomando "article list" (stub).
func newArticleListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lista artigos existentes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("comando 'article list' ainda não implementado")
		},
	}
}

// newReindexCmd cria o subcomando "reindex" (stub).
func newReindexCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reindex",
		Short: "Reindexa o acervo local",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("comando 'reindex' ainda não implementado")
		},
	}
}
