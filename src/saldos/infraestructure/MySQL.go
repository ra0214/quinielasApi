package infraestructure

import (
	"database/sql"
	"fmt"
	"log"
	"quinielas/src/config"
	"quinielas/src/saldos/domain"
	"quinielas/src/shared/money"
)

type MySQL struct {
	conn *config.Conn_MySQL
}

var _ domain.ISaldo = (*MySQL)(nil)

func NewMySQL() domain.ISaldo {
	conn := config.GetDBPool()
	if conn.Err != "" {
		log.Fatalf("Error al configurar el pool de conexiones: %v", conn.Err)
	}
	return &MySQL{conn: conn}
}

func (mysql *MySQL) InitSaldoCliente(idCliente int32) error {
	query := "INSERT INTO saldos (id_cliente, saldo_favor, saldo_deuda) VALUES (?, 0, 0) ON DUPLICATE KEY UPDATE id_cliente=id_cliente"
	_, err := mysql.conn.ExecutePreparedQuery(query, idCliente)
	if err != nil {
		return fmt.Errorf("error al inicializar saldo: %v", err)
	}
	return nil
}

func (mysql *MySQL) GetSaldoByClienteID(idCliente int32) (*domain.Saldo, error) {
	query := "SELECT id_cliente, saldo_favor, saldo_deuda FROM saldos WHERE id_cliente = ?"
	row := mysql.conn.FetchRow(query, idCliente)

	var s domain.Saldo
	err := row.Scan(&s.IDCliente, &s.SaldoFavor, &s.SaldoDeuda)
	if err != nil {
		if err == sql.ErrNoRows {
			// Si no existe el registro de saldo aún, lo creamos en 0
			_ = mysql.InitSaldoCliente(idCliente)
			return domain.NewSaldo(idCliente), nil
		}
		return nil, fmt.Errorf("error al consultar saldo: %v", err)
	}

	return &s, nil
}

func (mysql *MySQL) GetAllSaldos() ([]domain.Saldo, error) {
	// LEFT JOIN: incluimos todos los clientes aunque no tengan fila en saldos
	// (aparecen con saldo 0 y pueden editarse).
	query := `
		SELECT c.id_cliente,
		       COALESCE(s.saldo_favor, 0) AS saldo_favor,
		       COALESCE(s.saldo_deuda, 0) AS saldo_deuda
		FROM clientes c
		LEFT JOIN saldos s ON s.id_cliente = c.id_cliente
		ORDER BY c.id_cliente ASC`
	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener saldos: %v", err)
	}
	defer rows.Close()

	var saldos []domain.Saldo
	for rows.Next() {
		var s domain.Saldo
		if err := rows.Scan(&s.IDCliente, &s.SaldoFavor, &s.SaldoDeuda); err != nil {
			return nil, fmt.Errorf("error al escanear fila de saldo: %v", err)
		}
		saldos = append(saldos, s)
	}

	return saldos, nil
}

func (mysql *MySQL) GetSaldosFiltrados(filtro string) ([]domain.Saldo, error) {
	var query string
	switch filtro {
	case "deben":
		query = `
			SELECT c.id_cliente,
			       COALESCE(s.saldo_favor, 0) AS saldo_favor,
			       COALESCE(s.saldo_deuda, 0) AS saldo_deuda
			FROM clientes c
			LEFT JOIN saldos s ON s.id_cliente = c.id_cliente
			WHERE COALESCE(s.saldo_deuda, 0) > 0
			ORDER BY COALESCE(s.saldo_deuda, 0) DESC`
	case "favor":
		query = `
			SELECT c.id_cliente,
			       COALESCE(s.saldo_favor, 0) AS saldo_favor,
			       COALESCE(s.saldo_deuda, 0) AS saldo_deuda
			FROM clientes c
			LEFT JOIN saldos s ON s.id_cliente = c.id_cliente
			WHERE COALESCE(s.saldo_favor, 0) > 0
			ORDER BY COALESCE(s.saldo_favor, 0) DESC`
	case "ceros":
		query = `
			SELECT c.id_cliente,
			       COALESCE(s.saldo_favor, 0) AS saldo_favor,
			       COALESCE(s.saldo_deuda, 0) AS saldo_deuda
			FROM clientes c
			LEFT JOIN saldos s ON s.id_cliente = c.id_cliente
			WHERE COALESCE(s.saldo_favor, 0) = 0 AND COALESCE(s.saldo_deuda, 0) = 0`
	default:
		return mysql.GetAllSaldos()
	}

	rows, err := mysql.conn.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error al filtrar saldos: %v", err)
	}
	defer rows.Close()

	var saldos []domain.Saldo
	for rows.Next() {
		var s domain.Saldo
		if err := rows.Scan(&s.IDCliente, &s.SaldoFavor, &s.SaldoDeuda); err != nil {
			return nil, fmt.Errorf("error al escanear fila de saldo filtrado: %v", err)
		}
		saldos = append(saldos, s)
	}

	return saldos, nil
}

func (mysql *MySQL) UpdateSaldoCliente(idCliente int32, saldoFavor money.Money, saldoDeuda money.Money) error {
	query := "INSERT INTO saldos (id_cliente, saldo_favor, saldo_deuda) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE saldo_favor = ?, saldo_deuda = ?"
	_, err := mysql.conn.ExecutePreparedQuery(query, idCliente, saldoFavor, saldoDeuda, saldoFavor, saldoDeuda)
	if err != nil {
		return fmt.Errorf("error al actualizar saldo: %v", err)
	}
	return nil
}