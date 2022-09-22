package main

import (
	"log"
	"net"
)

func main() {

	servidor := nuevo_servidor()
	go servidor.correr()

	escuchador, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Incapaz de conectar al servidor: %s", err.Error())
	}
	defer escuchador.Close()
	log.Printf("Servidor funcionando en :8888")

	for {
		conexion, err := escuchador.Accept()

		if err != nil {
			log.Printf("Incapaz de aceptar la conexion: %s", err.Error())
			continue
		}
		cliente_ := servidor.nuevo_cliente(conexion)
		go cliente_.leer_entrada()
	}

}
