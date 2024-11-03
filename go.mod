module github.com/marcin-brzozowski/openapi-examples

go 1.23.2

require internal/cmd v1.0.0

replace internal/cmd => ./internal/cmd

require github.com/spf13/cobra v1.8.1 // indirect

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)
