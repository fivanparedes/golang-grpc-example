package main

import (
	base "base/pkg"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conexion, _ := grpc.Dial(
		// Dirección del servidor
		"localhost:7000",

		// Indica que se debe conectar usando TCP sin SSL
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		// Bloquea el hilo hasta que la conexión se establezca
		grpc.WithBlock(),
	)

	// Crea un nuevo cliente gRPC sobre la conexión
	cliente := base.NewBaseClient(conexion)

	ctx := context.Background()
	cliente.Guardar(ctx, &base.ParametroGuardar{Clave: "1", Valor: "AAA"})
	cliente.Guardar(ctx, &base.ParametroGuardar{Clave: "2", Valor: "BBB"})
	cliente.Guardar(ctx, &base.ParametroGuardar{Clave: "3", Valor: "CCC"})
	cliente.Guardar(ctx, &base.ParametroGuardar{Clave: "4", Valor: "DDD"})
	cliente.Guardar(ctx, &base.ParametroGuardar{Clave: "5", Valor: "EEE"})
	cliente.Guardar(ctx, &base.ParametroGuardar{Clave: "6", Valor: "FFF"})
	cliente.Guardar(ctx, &base.ParametroGuardar{Clave: "7", Valor: "GGG"})
	cliente.Guardar(ctx, &base.ParametroGuardar{Clave: "8", Valor: "HHH"})
	cliente.Guardar(ctx, &base.ParametroGuardar{Clave: "9", Valor: "III"})
	cliente.Guardar(ctx, &base.ParametroGuardar{Clave: "10", Valor: "JJJ"})
	resultado, _ := cliente.Obtener(ctx, &base.ParametroObtenerEliminar{Clave: "4"})
	fmt.Println(resultado)
	cliente.Eliminar(ctx, &base.ParametroObtenerEliminar{Clave: "4"})
	cliente.Eliminar(ctx, &base.ParametroObtenerEliminar{Clave: "5"})
}
