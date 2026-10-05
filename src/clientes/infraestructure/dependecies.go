// dependencies.go
package infraestructure

import (
	"quinielas/src/clientes/domain"
)

type Dependencies struct {
	ClienteRepo domain.ICliente
}

func InitDependencies() *Dependencies {
	return &Dependencies{
		ClienteRepo: NewMySQL(),
	}
}