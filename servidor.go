package main

import (
	"log"
	"net"
	"strings"
)

type servidor struct {
	salas    map[string]*sala
	comandos chan comando
}

func nuevo_servidor() *servidor {
	return &servidor{
		salas:    make(map[string]*sala),
		comandos: make(chan comando),
	}
}

func (servidor_ *servidor) correr() {

	for comando := range servidor_.comandos {

		switch comando.id {
		case comando_apodo:
			servidor_.apodo(comando.cliente, comando.args)
		case comando_unir:
			servidor_.unir(comando.cliente, comando.args)
		case comando_salas:
			servidor_.listar_salas(comando.cliente)
		case comando_enviar:
			servidor_.enviar(comando.cliente, comando.args, comando.archivo)
		case comando_recibir:
			servidor_.recibir(comando.cliente, comando.args)
		case comando_salir:
			servidor_.salir(comando.cliente)
		}
	}
}

func (servidor_ *servidor) nuevo_cliente(conexion_ net.Conn) *cliente {
	log.Printf("Nuevo cliente se ha conectado: %s", conexion_.RemoteAddr().String())

	return &cliente{
		conexion: conexion_,
		apodo:    "anonimo",
		comandos: servidor_.comandos,
	}
}

func (servidor_ *servidor) apodo(cliente_ *cliente, args_ []string) {

	if len(args_) < 2 {
		log.Println("Se requiere apodo. use: /apodo nombre")
		return
	}

	cliente_.apodo = args_[1]

	log.Printf("El cliente: %s tiene como apodo %s", cliente_.conexion.RemoteAddr().String(), cliente_.apodo)

}

func (servidor_ *servidor) unir(cliente_ *cliente, args_ []string) {

	if len(args_) < 2 {
		log.Println("Se requiere nombre de la sala. use: /unir nombre_sala")
		return
	}

	nombre_sala := args_[1]

	sala_, ok := servidor_.salas[nombre_sala]
	if !ok {
		sala_ = &sala{
			nombre:   nombre_sala,
			miembros: make(map[net.Addr]*cliente),
		}

		servidor_.salas[nombre_sala] = sala_
	}

	sala_.miembros[cliente_.conexion.RemoteAddr()] = cliente_

	servidor_.salir_sala_actual(cliente_)

	cliente_.sala = sala_

	log.Println("Se ha unido a la sala " + sala_.nombre + " el cliente: " + cliente_.apodo)

}

func (servidor_ *servidor) listar_salas(cliente_ *cliente) {

	var salas_ []string

	for nombre_ := range servidor_.salas {
		salas_ = append(salas_, nombre_)
	}
	log.Printf("El cliente: %s esta solicitando el listado de salas disponibles.", cliente_.conexion.RemoteAddr().String())
	cliente_.conexion.Write([]byte(strings.Join(salas_, " , ")))

}

func (servidor_ *servidor) enviar(cliente_ *cliente, args_ []string, archivo_ []byte) {

	if len(args_) < 2 {
		log.Println("Se requiere ruta, use: /enviar ruta")
		return
	}

	cliente_.sala.transmision_fichero(cliente_, args_, archivo_)

	log.Println("Termino de enrutar")

}

func (servidor_ *servidor) recibir(cliente_ *cliente, args_ []string) {

	if len(args_) > 1 {
		log.Println("El comando no requiere parametros /recibir")
		return
	}
	validador := cliente_.sala.validar_recibir(cliente_)

	if !validador {
		cliente_.conexion.Write([]byte("no"))
	} else {
		cliente_.conexion.Write([]byte("si"))
		log.Printf("El cliente %s esta recibiendo archivos", cliente_.conexion.RemoteAddr().String())
		cliente_.recibir = true
	}

}

func (servidor_ *servidor) salir(cliente_ *cliente) {

	servidor_.salir_sala_actual(cliente_)

	cliente_.conexion.Close()

	log.Printf("El cliente se ha desconectado %s", cliente_.conexion.RemoteAddr().String())
}

func (servidor_ *servidor) salir_sala_actual(cliente_ *cliente) {
	if cliente_.sala != nil {
		delete(servidor_.salas[cliente_.sala.nombre].miembros, cliente_.conexion.RemoteAddr())
		log.Println("El cliente: " + cliente_.conexion.LocalAddr().String() + " ha salido de la sala")
	}
}

func enviar_fichero_cliente(cliente_ *cliente, args_ []string, archivo_ []byte) {

	cliente_.conexion.Write([]byte(args_[2]))
	cliente_.conexion.Write([]byte(args_[3]))
	cliente_.conexion.Write(archivo_)

	cliente_.recibir = false

	log.Println("Se ha enviado un archivo a: " + cliente_.apodo)
}
