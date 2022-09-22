package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func salir(texto_ string, conexion_ net.Conn) {

	msj := []byte(texto_)
	conexion_.Write(msj)
	fmt.Println("Se ha desconectado del servidor.")
}

func apodo(texto_ []string, conexion_ net.Conn) {

	if len(texto_) == 2 {

		msj := []byte(texto_[0] + " " + texto_[1])
		conexion_.Write(msj)
		fmt.Println("Su nuevo apodo es: " + texto_[1])

	} else {
		fmt.Println("El comando no se ejecuto correctamente: /apodo nombre")
	}

}

func unir(texto_ []string, conexion_ net.Conn) {

	if len(texto_) == 2 {

		msj := []byte(texto_[0] + " " + texto_[1])
		conexion_.Write(msj)
		fmt.Println("Te has unido a las sala: " + texto_[1])

	} else {
		fmt.Println("El comando no se ejecuto correctamente: /unir canal")
	}

}

func salas(texto_ string, conexion_ net.Conn) {

	msj := []byte(texto_)
	lista_salas := make([]byte, 64)

	conexion_.Write(msj)

	conexion_.Read(lista_salas)

	salas := convertir_byte_string(lista_salas, 64)
	fmt.Println("Las salas disponibles son: " + salas)

}

func enviar_archivo(texto_ []string, conexion_ net.Conn) {

	msj := []byte(strings.Join(texto_[0:], " "))
	ruta := texto_[1]

	validador := true
	archivo, err := os.Open(ruta)
	if err != nil {
		fmt.Println(err)
		validador = false
		return
	}
	defer archivo.Close()

	info_archivo, err := archivo.Stat()
	if err != nil {
		fmt.Println(err)
		validador = false
		return
	}

	if validador {
		conexion_.Write(msj)

		tamano_archivo := rellenar_string(strconv.FormatInt(info_archivo.Size(), 10), 10)
		nombre_archivo := rellenar_string(info_archivo.Name(), 64)

		conexion_.Write([]byte(nombre_archivo))
		conexion_.Write([]byte(tamano_archivo))

		enviar_bufer := make([]byte, info_archivo.Size())

		for {
			_, err = archivo.Read(enviar_bufer)
			if err == io.EOF {
				break
			}
			conexion_.Write(enviar_bufer)
		}
		fmt.Println("El fichero se ha enviado, correctamente!")

	}

}

func recibir(texto_ string, conexion_ net.Conn) {

	msj := []byte(texto_)
	conexion_.Write(msj)

	validador_byte := make([]byte, 2)
	conexion_.Read(validador_byte)

	validador := convertir_byte_string(validador_byte, 2)

	if validador == "no" {

		conexion_.Read(make([]byte, 1024))

		fmt.Println("No perteneces a ninguna sala, por favor unase a una sala, para poder enviar archivos.")

	} else {

		fmt.Println("Se esta recibiendo archivo...")

		buffer_nombre_archivo := make([]byte, 64)
		buffer_tamano_archivo := make([]byte, 10)

		conexion_.Read([]byte(buffer_nombre_archivo))
		conexion_.Read([]byte(buffer_tamano_archivo))

		nombre_archivo := convertir_byte_string(buffer_nombre_archivo, 64)
		tamano_archivo := convertir_byte_string(buffer_tamano_archivo, 10)

		nuevo_fichero, err := os.Create(nombre_archivo)

		if err != nil {
			panic(err)
		}

		defer nuevo_fichero.Close()

		archivo_tamano, err_no := strconv.Atoi(tamano_archivo)

		if err_no != nil {
			panic(err_no)
		}

		buffer_archivo := make([]byte, int(archivo_tamano))

		conexion_.Read(buffer_archivo)

		escribir_archivo(nombre_archivo, buffer_archivo)

		fmt.Println("Se recibio archivo correctamente.")
	}

}
