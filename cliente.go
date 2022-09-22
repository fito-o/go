package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type cliente struct {
	conexion net.Conn
	apodo    string
	sala     *sala
	comandos chan<- comando
	recibir  bool
}

func (cliente_ *cliente) leer_entrada() {
bucle_leer_entrada:
	for {
		comando_byte := make([]byte, 64)

		cliente_.conexion.Read([]byte(comando_byte))

		texto := convertir_byte_string(comando_byte, 64)
		msj := strings.Trim(texto, "\r\n")
		args := strings.Split(msj, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/recibir":
			cliente_.comandos <- comando{
				id:      comando_recibir,
				cliente: cliente_,
				args:    args,
			}
		case "/apodo":
			cliente_.comandos <- comando{
				id:      comando_apodo,
				cliente: cliente_,
				args:    args,
			}
		case "/unir":
			cliente_.comandos <- comando{
				id:      comando_unir,
				cliente: cliente_,
				args:    args,
			}
		case "/salas":
			cliente_.comandos <- comando{
				id:      comando_salas,
				cliente: cliente_,
				args:    args,
			}
		case "/enviar":
			bufer_nombre_archivo := make([]byte, 64)
			bufer_tamano_archivo := make([]byte, 10)
			cliente_.conexion.Read([]byte(bufer_nombre_archivo))
			cliente_.conexion.Read([]byte(bufer_tamano_archivo))

			tamano_archivo := convertir_byte_string(bufer_tamano_archivo, 10)
			archivo_tamano, err_no := strconv.Atoi(tamano_archivo)

			if err_no != nil {
				panic(err_no)
			}

			buffer_archivo := make([]byte, archivo_tamano)
			cliente_.conexion.Read([]byte(buffer_archivo))

			temp := append(args, convertir_byte_string2(bufer_nombre_archivo, 64))
			buffer := append(temp, convertir_byte_string2(bufer_tamano_archivo, 10))

			cliente_.comandos <- comando{
				id:      comando_enviar,
				cliente: cliente_,
				args:    buffer,
				archivo: buffer_archivo,
			}
		case "/salir":
			cliente_.comandos <- comando{
				id:      comando_salir,
				cliente: cliente_,
				args:    args,
			}
			break bucle_leer_entrada
		default:
			error := fmt.Errorf("comando desconocido: %s", cmd)
			fmt.Println(error)
			break bucle_leer_entrada

		}
	}
}

func convertir_byte_string2(cadena_bytes_ []byte, longitud_ int) (texto_ string) {

	for i := 0; i < longitud_; i++ {
		texto_ += string(cadena_bytes_[i])
	}

	return texto_

}

func convertir_byte_string(cadena_bytes_ []byte, longitud_ int) (texto_ string) {

	for i := 0; i < longitud_; i++ {
		if cadena_bytes_[i] != 0 && cadena_bytes_[i] != 58 {
			texto_ += string(cadena_bytes_[i])
		}
	}

	return texto_

}
