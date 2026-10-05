package domain

import "time"

// ICliente define el puerto (interfaz) para el repositorio de clientes
type ICliente interface {
	SaveCliente(nombre string, telefono string) (*Cliente, error)
	DeleteCliente(id int32) error
	UpdateCliente(id int32, nombre string, telefono string) error
	GetAll() ([]Cliente, error)
	GetClienteByID(id int32) (*Cliente, error)
	SearchClientesByName(query string) ([]Cliente, error)
}

// Cliente representa la entidad de dominio de la tabla 'clientes'
type Cliente struct {
	ID            int32     `json:"id_cliente"`
	Nombre        string    `json:"nombre"`
	Telefono      string    `json:"telefono"`
	FechaRegistro time.Time `json:"fecha_registro"`
}

// NewCliente crea una nueva instancia de la entidad Cliente
func NewCliente(nombre string, telefono string) *Cliente {
	return &Cliente{
		Nombre:        nombre,
		Telefono:      telefono,
		FechaRegistro: time.Now(),
	}
}

// Métodos setters y auxiliares del dominio
func (c *Cliente) SetNombre(nombre string) {
	c.Nombre = nombre
}

func (c *Cliente) SetTelefono(telefono string) {
	c.Telefono = telefono
}