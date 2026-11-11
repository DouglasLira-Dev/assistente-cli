# Assistente CLI

Assistente de linha de comando para consulta e organização de conhecimento pessoal.
Centraliza artigos, templates e conversas em um único lugar, com busca local (SQLite FTS5)
e respostas geradas por LLM. Feito para desenvolvedores que preferem o terminal ao navegador.

## Stack

- **Go** — linguagem principal
- **Cobra** — framework de CLI (comandos e flags)
- **SQLite (FTS5)** — armazenamento local com busca full-text
- **Markdown** — formato dos artigos e templates
- **LLM via API** — Groq como provedor principal, OpenRouter como fallback

## Status

Em desenvolvimento inicial, estrutura base concluída.

## Pré-requisitos

- **Go 1.27+** ([download](https://go.dev/dl/))
- **gcc** — necessário para CGO (driver SQLite). No Ubuntu/Debian:
  ```bash
  sudo apt install build-essential
  ```
- **Chaves de API** - Groq (principal) e OpenRouter (fallback), obtidas em console.groq.com e openrouter.ai

## Como compilar

```bash
go build -o assistente ./cmd/assistente
```

O binário `assistente` será gerado na raiz do repositório (ignorado pelo Git).

## Como configurar

1. Copie o template de configuração:

```bash
cp .env.example .env
```

2. Preencha as chaves no `.env`:

```env
GROQ_API_KEY=sua_chave_aqui
OPENROUTER_API_KEY=sua_chave_aqui
```

3. As variáveis de caminho (`DATA_DIR`, `ARTICLES_DIR`, `TEMPLATES_DIR`) já vêm com valores padrão e podem ser ajustadas conforme necessário.

O `.env` **nunca** vai para o Git, está protegido pelo `.gitignore`.

## Estrutura de pastas

```
assistente-cli/
├── cmd/
│   └── assistente/           # Ponto de entrada do binário
│       └── main.go
├── internal/                 # Código privado da aplicação
│   ├── config/               # Carregamento do .env e defaults
│   ├── command/              # Comandos Cobra
│   ├── service/              # Regras de negócio
│   └── repository/           # Acesso a dados (SQLite)
├── data/                     # Dados locais (parcialmente ignorado)
│   ├── articles/             # Artigos versionados
│   ├── templates/            # Templates versionados
│   ├── chat.db               # Conversas (ignorado)
│   └── events.jsonl          # Log de eventos (ignorado)
├── .env.example              # Template de configuração
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Uso

```bash
# Ver ajuda geral
./assistente --help

# Ver versão
./assistente --version
```

Os comandos serão documentados conforme forem implementados.

## Licença

Distribuído sob a licença MIT. Veja [LICENSE](LICENSE) para mais detalhes.