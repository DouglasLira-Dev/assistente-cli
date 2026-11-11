// Package store gerencia a persistência do assistente em SQLite.
// Este arquivo cuida apenas da fundação: abrir/fechar o banco e aplicar
// o schema de forma idempotente. Nenhuma lógica de negócio vive aqui.
package store

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaSQL string

// Open abre o banco SQLite no caminho informado, aplica os PRAGMAs
// necessários e executa o schema de forma idempotente.
//
// O caminho é decidido pelo chamador (o cli, que lê do config).
// Exemplo: store.Open(filepath.Join(cfg.DataDir, "chat.db"))
func Open(path string) (*sql.DB, error) {
	// PRAGMAs via DSN: mais limpo que executar PRAGMA após abrir.
	// _foreign_keys=on: obrigatório para ON DELETE CASCADE funcionar.
	// _journal_mode=WAL: reduz custo de escrita, cria -wal e -shm.
	// _busy_timeout=5000: evita "database is locked" em concorrência.
	dsn := path + "?_foreign_keys=on&_journal_mode=WAL&_busy_timeout=5000"

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("store: abrir banco: %w", err)
	}

	// Executa o schema (idempotente graças aos IF NOT EXISTS).
	if _, err := db.Exec(schemaSQL); err != nil {
		db.Close()
		return nil, fmt.Errorf("store: aplicar schema: %w", err)
	}

	// Verificação de sanidade: confirma que as tabelas esperadas existem.
	// Se o schema estiver quebrado, falha cedo.
	if err := verifyTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// verifyTables confirma que as tabelas e índices esperados existem
// após a execução do schema.
func verifyTables(db *sql.DB) error {
	expected := []string{
		"conversations",
		"messages",
		"messages_fts",
		"articles",
		"articles_fts",
	}

	for _, table := range expected {
		var count int
		query := `SELECT COUNT(*) FROM sqlite_master WHERE type IN ('table','index') AND name = ?`
		if err := db.QueryRow(query, table).Scan(&count); err != nil {
			return fmt.Errorf("store: verificar tabela %q: %w", table, err)
		}
		if count == 0 {
			return fmt.Errorf("store: tabela %q não encontrada após aplicar schema", table)
		}
	}

	return nil
}