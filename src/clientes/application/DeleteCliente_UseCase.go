package application

import (
	"quinielas/src/clientes/domain"
	saldosDomain "quinielas/src/saldos/domain"
)

type DeleteCliente struct {
	repo      domain.ICliente
	saldoRepo saldosDomain.ISaldo
}

func NewDeleteCliente(repo domain.ICliente, saldoRepo saldosDomain.ISaldo) *DeleteCliente {
	return &DeleteCliente{repo: repo, saldoRepo: saldoRepo}
}

func (dc *DeleteCliente) Execute(id int32) error {
	// Regla de negocio: un cliente con saldo pendiente (deuda o a favor) no puede
	// borrarse; solo se borra cuando no debe ni tiene saldo a favor.
	saldo, err := dc.saldoRepo.GetSaldoByClienteID(id)
	if err != nil {
		return err
	}
	if !saldo.SaldoFavor.IsZero() || !saldo.SaldoDeuda.IsZero() {
		return &saldosDomain.SaldoPendienteError{Favor: saldo.SaldoFavor, Deuda: saldo.SaldoDeuda}
	}

	if err := dc.repo.DeleteCliente(id); err != nil {
		return err
	}
	return nil
}