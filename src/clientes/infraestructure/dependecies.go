// dependencies.go
package infraestructure

import (
	"quinielas/src/clientes/domain"
	saldosDomain "quinielas/src/saldos/domain"
	saldosInfra "quinielas/src/saldos/infraestructure"
)

type Dependencies struct {
	ClienteRepo domain.ICliente
	SaldoRepo   saldosDomain.ISaldo
}

func InitDependencies() *Dependencies {
	return &Dependencies{
		ClienteRepo: NewMySQL(),
		SaldoRepo:   saldosInfra.NewMySQL(),
	}
}