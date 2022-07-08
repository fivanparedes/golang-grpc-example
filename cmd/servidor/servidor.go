package main

import (
	base "base/pkg"
	"fmt"
	"log"
	"net"
	"os"
	"google.golang.org/grpc"
)

func main() {
	// Uso os.Args para validar argumentos
	argumentos := os.Args[1:]
	if len(argumentos) != 1 {
		fmt.Println("Error: Se debe pasar un solo argumento, debe ser un valor de 0 a 2")
		os.Exit(1)
	}

	puerto, ok := base.MapaNodos[argumentos[0]]
	if !ok {
		fmt.Println("Error: El argumento debe ser un valor de 0 a 2")
		os.Exit(1)
	} 
	
    servicio := base.NuevoServidor(argumentos[0])
    servidorReal := grpc.NewServer()

    base.RegisterBaseServer(servidorReal, servicio)
	listen, err := net.Listen("tcp", "localhost:" + puerto)
	fmt.Println(servicio.IdNodo, " ",  puerto)
   
	if err != nil {
		log.Fatalf("fallo al escuchar: %v", err)
	}
	if err := servidorReal.Serve(listen); err != nil {
		log.Fatalf("fallo al servir: %v", err)
	}
}