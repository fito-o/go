package main

import (
	"fmt"
	"os"
)

func escribir_archivo(nombre_ string, archivo_ []byte) {

	data := archivo_

	if os.WriteFile(nombre_, data, 0666) == nil {
		fmt.Println("Archivo escrito correctamente.")
	}

}

func rellenar_string(rotorna_string_ string, longitud_ int) string {

	for {

		longitud_string := len(rotorna_string_)

		if longitud_string < longitud_ {
			rotorna_string_ = rotorna_string_ + ":"
			continue
		}

		break
	}

	return rotorna_string_
}

func convertir_byte_string(cadena_bytes []byte, longitud int) (texto string) {

	for i := 0; i < longitud; i++ {

		if cadena_bytes[i] != 0 && cadena_bytes[i] != 58 {
			texto += string(cadena_bytes[i])
		}

	}

	return texto

}
