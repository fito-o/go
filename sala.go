package main

import (
	"log"
	"net"
	"time"
)

type sala struct {
	nombre   string
	miembros map[net.Addr]*cliente
}

func (sala_ *sala) transmision_fichero(remitente_ *cliente, args_ []string, archivo_ []byte) {

	for nombre_, miembro_ := range sala_.miembros {
		if nombre_ != remitente_.conexion.RemoteAddr() {
			enviar_fichero_cliente(miembro_, args_, archivo_)
			time.Sleep(2 * time.Second)
		}
	}

}

func (sala_ *sala) validar_recibir(remitente_ *cliente) bool {

	if sala_ == nil {
		log.Printf("El cliente: %s se esta conectando a una sala que no existe", remitente_.conexion.RemoteAddr().String())
		return false
	} else {
		for nombre_ := range sala_.miembros {
			if nombre_ == remitente_.conexion.RemoteAddr() {
				return true
			}
		}
	}
	return false
}
