package domain

import "time"

// IPurga elimina los concursos (ediciones) que ya cerraron y todo su
// contenido, para que la base de datos no crezca sin limite.
type IPurga interface {
	// PurgarConcursosAnteriores elimina las ediciones cuya fecha de cierre
	// (o de inicio, si no tienen cierre) es anterior a "limite" y, en
	// cascada, sus quinielas, aportes, premios, repartos y movimientos.
	// Devuelve cuantas ediciones elimino. No toca la tabla de saldos.
	PurgarConcursosAnteriores(limite time.Time) (int, error)
}
