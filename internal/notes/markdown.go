// Package notes lida com arquivos Markdown com frontmatter YAML.
// Responsável por parsear (leitura) e escrever (gravação) documentos
// no formato usado pelo assistente.
package notes

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Frontmatter representa o bloco YAML no topo de um arquivo Markdown.
type Frontmatter struct {
	Title   string   `yaml:"title"`
	Tags    []string `yaml:"tags,omitempty"`
	Type    string   `yaml:"type,omitempty"`
	Created string   `yaml:"created,omitempty"`
	Updated string   `yaml:"updated,omitempty"`
}

// Document representa um arquivo Markdown completo: frontmatter + corpo.
type Document struct {
	Frontmatter Frontmatter
	Body        string
	Path        string
}

const (
	frontmatterDelim = "---"
	defaultType      = "note"
)

// Parse lê um arquivo Markdown com frontmatter YAML e retorna o Document.
//
// O arquivo precisa começar com "---" na primeira linha. O frontmatter
// termina no próximo "---" isolado. O restante é o corpo.
func Parse(path string) (*Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("notes: abrir %s: %w", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024) // suporta linhas longas

	// Primeira linha precisa ser o delimitador de abertura.
	if !scanner.Scan() {
		return nil, fmt.Errorf("notes: arquivo vazio: %s", path)
	}
	if strings.TrimSpace(scanner.Text()) != frontmatterDelim {
		return nil, fmt.Errorf("frontmatter YAML ausente: o arquivo deve começar com '---'")
	}

	// Coleta as linhas do frontmatter até o delimitador de fechamento.
	var yamlBuf bytes.Buffer
	foundClose := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == frontmatterDelim {
			foundClose = true
			break
		}
		yamlBuf.WriteString(line)
		yamlBuf.WriteByte('\n')
	}
	if !foundClose {
		return nil, fmt.Errorf("notes: frontmatter não fechado em %s (faltou '---')", path)
	}

	// O resto do scanner é o corpo.
	var bodyBuf bytes.Buffer
	for scanner.Scan() {
		bodyBuf.WriteString(scanner.Text())
		bodyBuf.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("notes: ler %s: %w", path, err)
	}

	// Unmarshal do YAML.
	var fm Frontmatter
	if err := yaml.Unmarshal(yamlBuf.Bytes(), &fm); err != nil {
		return nil, fmt.Errorf("notes: YAML inválido em %s: %w", path, err)
	}

	if strings.TrimSpace(fm.Title) == "" {
		return nil, fmt.Errorf("notes: 'title' obrigatório em %s", path)
	}
	if fm.Type == "" {
		fm.Type = defaultType
	}

	return &Document{
		Frontmatter: fm,
		Body:        strings.TrimRight(bodyBuf.String(), "\n"),
		Path:        path,
	}, nil
}

// Write serializa o frontmatter em YAML e grava o arquivo no formato
// "---\n<yaml>---\n<body>". Usado pelo comando article add (M5).
func Write(path string, fm Frontmatter, body string) error {
	if fm.Type == "" {
		fm.Type = defaultType
	}

	yamlBytes, err := yaml.Marshal(fm)
	if err != nil {
		return fmt.Errorf("notes: serializar frontmatter: %w", err)
	}

	var buf bytes.Buffer
	buf.WriteString(frontmatterDelim)
	buf.WriteByte('\n')
	buf.Write(yamlBytes)
	buf.WriteString(frontmatterDelim)
	buf.WriteByte('\n')
	buf.WriteString(body)
	if !strings.HasSuffix(body, "\n") {
		buf.WriteByte('\n')
	}

	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("notes: escrever %s: %w", path, err)
	}
	return nil
}

// SlugFromPath extrai o nome do arquivo sem extensão .md (case-insensitive).
// Exemplo: data/articles/sqlite-fts5.md → sqlite-fts5
func SlugFromPath(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if strings.EqualFold(ext, ".md") {
		return strings.TrimSuffix(base, ext)
	}
	return base
}