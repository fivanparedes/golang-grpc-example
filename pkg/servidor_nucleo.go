package pkg

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"math/big"
	"sync"

	grpc "google.golang.org/grpc"
)

// Cantidad de nodos
const NODOS = 3

// Mapas de nodos (id : puerto)
var MapaNodos = map[string]string{
	"0": "12345",
	"1": "12346",
	"2": "12347",
}

// Define donde almacenar
func hash(clave string) (resultado string) {
	num := big.NewInt(0)
	nodos := big.NewInt(NODOS)
	nodo := big.NewInt(0)
	h := md5.New()
	h.Write([]byte(clave))
	hexstr := hex.EncodeToString(h.Sum(nil))
	num.SetString(hexstr, 16)
	nodo.Mod(num, nodos)
	return nodo.String()
}

// Conexion a servidor que tiene la clave
func conectarServidor(id string) BaseClient {
	// Obtener el puerto del nodo
	puerto := MapaNodos[id]
	// Crear conexion
	conn, err := grpc.Dial("localhost:"+puerto, grpc.WithDefaultServiceConfig("nucleo"))
	if err != nil {
		panic(err)
	}
	// Crear cliente
	cliente := NewBaseClient(conn)
	return cliente
}

// La implementación del servidor
type Servidor struct {
	UnimplementedBaseServer
	// id nodo
	IdNodo string
	// Cerrojo
	lock sync.RWMutex
	// Base de datos
	kv map[string]string
}

// Crea servidor
func NuevoServidor(idNodo string) Servidor {
	return Servidor{
		IdNodo: idNodo,
		lock:   sync.RWMutex{},
		kv:     make(map[string]string),
	}
}

// Guardar en proto
func (s Servidor) Guardar(ctx context.Context, p *ParametroGuardar) (*ResultadoGuardar, error) {
	// Obtener el nodo al que indica la clave
	nodoDestino := hash(p.Clave)
	if s.IdNodo != nodoDestino {
		cliente := conectarServidor(nodoDestino)
		return cliente.Guardar(context.Background(), p)
	} else {
		var existeClave int32
		existeClave = 0
		// Obtener cerrojo
		s.lock.RLock()

		// Si existe retorno 1
		if _, existe := s.kv[p.Clave]; existe {
			existeClave = 1
		}

		// Modifico o agrego clave
		s.kv[p.Clave] = p.Valor

		// Libero cerrojo
		s.lock.RUnlock()
		return &ResultadoGuardar{Valor: existeClave, Error: ""}, nil
	}
}

// Obtener en proto
func (s Servidor) Obtener(ctx context.Context, p *ParametroObtenerEliminar) (*ResultadoObtenerEliminar, error) {
	// Obtener el nodo al que indica la clave
	nodoDestino := hash(p.Clave)
	if s.IdNodo != string(nodoDestino) {
		cliente := conectarServidor(nodoDestino)
		return cliente.Obtener(context.Background(), p)
	} else {
		// Si no hay datos indica Error
		error := ""

		// Obtener cerrojo de lectura
		s.lock.RLock()

		// Obtener clave
		valor, existe := s.kv[p.Clave]

		// Si no existe indicar Error
		if !existe {
			error = "La clave no existe"
		}

		// Liberar cerrojo de lectura
		s.lock.RUnlock()

		// Retornar el resultado
		return &ResultadoObtenerEliminar{Valor: valor, Error: error}, nil
	}
}

// Eliminar en proto
func (s Servidor) Eliminar(ctx context.Context, p *ParametroObtenerEliminar) (*ResultadoObtenerEliminar, error) {
	// Obtener el nodo al que indica la clave
	nodoDestino := hash(p.Clave)
	if s.IdNodo != string(nodoDestino) {
		cliente := conectarServidor(nodoDestino)
		return cliente.Eliminar(context.Background(), p)
	} else {
		// Si no hay datos indica Error
		error := ""

		// Obtener cerrojo de lectura
		s.lock.RLock()

		// Obtener clave
		valor, existe := s.kv[p.Clave]

		// Si no existe indicar Error
		if !existe {
			error = "La clave no existe"
		} else {
			delete(s.kv, p.Clave)
		}

		// Liberar cerrojo de lectura
		s.lock.RUnlock()

		// Retornar el resultado
		return &ResultadoObtenerEliminar{Valor: valor, Error: error}, nil
	}
}
