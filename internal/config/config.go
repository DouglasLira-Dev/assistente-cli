// Package config centraliza o carregamento e validação das configurações
// do assistente, lendo variáveis de ambiente e o arquivo .env (se existir).
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config agrupa todas as configurações do assistente.
// Os campos são exportados para que outros pacotes possam lê-los.
type Config struct {
	GroqKey         string
	OpenRouterKey   string
	DataDir         string
	ArticlesDir     string
	TemplatesDir    string
	GroqModel       string
	OpenRouterModel string
}

// Load carrega a configuração a partir do .env (se existir) e de variáveis
// de ambiente, aplicando defaults para valores ausentes.
//
// Retorna erro apenas se GROQ_API_KEY não estiver configurada.
func Load() (*Config, error) {
	// .env é opcional: se não existir, seguimos com env vars e defaults.
	// Ignoramos apenas erro de "arquivo não encontrado"; outros erros
	// (ex.: arquivo corrompido) são propagados.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("erro ao carregar .env: %w", err)
	}

	cfg := &Config{
		GroqKey:         getEnv("GROQ_API_KEY", ""),
		OpenRouterKey:   getEnv("OPENROUTER_API_KEY", ""),
		DataDir:         getEnv("DATA_DIR", "./data"),
		ArticlesDir:     getEnv("ARTICLES_DIR", "./data/articles"),
		TemplatesDir:    getEnv("TEMPLATES_DIR", "./data/templates"),
		GroqModel:       getEnv("GROQ_MODEL", "llama-3.3-70b-versatile"),
		OpenRouterModel: getEnv("OPENROUTER_MODEL", "meta-llama/llama-3.3-70b-instruct:free"),
	}

	// Validação mínima: GroqKey é obrigatória (provedor principal).
	if cfg.GroqKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY não configurada: copie .env.example para .env e preencha")
	}

	// OpenRouterKey é opcional (fallback).

	return cfg, nil
}

// getEnv lê uma variável de ambiente e retorna fallback se estiver vazia.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}