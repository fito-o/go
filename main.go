package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var lector *bufio.Reader = bufio.NewReader(os.Stdin)

func main() {

	conexion := crear_conexion("192.168.0.5", "8888")

comandos:

	for {

		entrada, err := leer_entrada()

		if err != nil {
			log.Fatal(err)
		}

		switch entrada[0] {

		case "/salir":
			salir(entrada[0], conexion)
			defer conexion.Close()
			break comandos

		case "/apodo":
			apodo(entrada, conexion)

		case "/unir":
			unir(entrada, conexion)

		case "/salas":
			salas(entrada[0], conexion)

		case "/enviar":
			enviar_archivo(entrada, conexion)

		case "/recibir":
			recibir(entrada[0], conexion)

		default:
			fmt.Println("Comando no valido.")

		}

	}
}

func crear_conexion(ip_ string, puerto_ string) net.Conn {

	conexion, err := net.Dial("tcp", ip_+":"+puerto_)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	return conexion
}

func leer_entrada() (args_ []string, err_ error) {

	var texto string

	fmt.Print(">")
	texto, err_ = lector.ReadString('\n')

	if err_ != nil {
		return args_, err_
	}

	texto = strings.Replace(texto, "\r\n", "", 1)
	args_ = strings.Split(texto, " ")

	return args_, nil

}
