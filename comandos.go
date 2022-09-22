package main

type comando_id int

const (
	comando_apodo = iota
	comando_unir
	comando_salas
	comando_recibir
	comando_enviar
	comando_salir
)

type comando struct {
	id      comando_id
	cliente *cliente
	args    []string
	archivo []byte
}
