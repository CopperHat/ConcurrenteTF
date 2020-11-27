package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

const (
	cnum = iota // Valores en secuencia
	num
)

type tmsg struct {
	Code int
	Addr string
	Op   int
}

// IP propia
const localAddr = "localhost:8004"

// IP de compañeros
var addrs = []string{"localhost:8000",
	"localhost:8001",
	"localhost:8002",
	"localhost:8003"}

var chInfo chan map[string]int

// Nodo Nodo
type Nodo struct {
	Peso               int
	Nivel              int
	Izquierda, Derecha *Nodo
}

// Arbol Arbol
type Arbol struct {
	primerNodo *Nodo
	ordenado   []int
	Pasos      int
}

//Recorrer recorre el arbol
func (t *Arbol) Recorrer(nodo *Nodo) {
	t.Pasos++
	if nodo.Izquierda != nil {
		t.Recorrer(nodo.Izquierda)
	}
	t.ordenado = append(t.ordenado, nodo.Peso)
	if nodo.Derecha != nil {
		t.Recorrer(nodo.Derecha)
	}
}

// Agregar agrega nodo a Arbol
func (t *Arbol) Agregar(peso int) {
	if t.primerNodo == nil {
		t.primerNodo = new(Nodo)
		t.primerNodo.Peso = peso
		t.primerNodo.Nivel = 1
		return
	}
	actual := t.primerNodo
	for {
		if peso < actual.Peso {
			if actual.Izquierda == nil {
				actual.Izquierda = new(Nodo)
				actual.Izquierda.Peso = peso
				actual.Izquierda.Nivel = actual.Nivel + 1
				break
			}
			actual = actual.Izquierda
		} else {
			if actual.Derecha == nil {
				actual.Derecha = new(Nodo)
				actual.Derecha.Peso = peso
				actual.Derecha.Nivel = actual.Nivel + 1
				break
			}
			actual = actual.Derecha
		}
	}
}

// Ordenar ordena nodos en Arbol
func (t *Arbol) Ordenar() []int {
	t.Pasos = 0
	t.ordenado = []int{}
	t.Recorrer(t.primerNodo)
	return t.ordenado
}

func checkFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	t := new(Arbol)
	filename := "Arbol.json"

	chInfo = make(chan map[string]int)
	go func() { chInfo <- map[string]int{} }()
	go server()
	time.Sleep(time.Millisecond * 100)
	var op int
	for {
		checkFile(filename)
		file, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println("No se puede abrir archivo")
		}
		data := []Arbol{}
		// Here the magic happens!
		json.Unmarshal(file, &data)

		fmt.Print("Numero a agregar: ")
		fmt.Scanf("%d\n", &op)
		if op == 0 {
			jsonArbol := json.NewEncoder(os.Stdout)
			jsonArbol.SetIndent("", "  ")
			_ = jsonArbol.Encode(t.primerNodo)
		} else if op == 999 {
			os.Exit(0)
		} else {
			t.Agregar(op)
			msg := tmsg{cnum, localAddr, op}
			for _, addr := range addrs {
				send(addr, msg)
			}
		}
	}
}
func server() {
	if ln, err := net.Listen("tcp", localAddr); err != nil {
		log.Panicln("No se puede iniciar en", localAddr)
	} else {
		defer ln.Close()
		fmt.Println("Escuchando a ", localAddr)
		for {
			if conn, err := ln.Accept(); err != nil {
				log.Println("No se puede aceptar ", conn.RemoteAddr())
			} else {
				go handle(conn)
			}
		}
	}
}
func handle(conn net.Conn) {
	defer conn.Close()
	dec := json.NewDecoder(conn)
	var msg tmsg
	if err := dec.Decode(&msg); err != nil {
		log.Println("No es correcto ", conn.RemoteAddr())
	} else {
		fmt.Println(msg)
		switch msg.Code {
		case cnum:
			concensus(conn, msg)
		}
	}
}
func concensus(conn net.Conn, msg tmsg) {
	info := <-chInfo
	info[msg.Addr] = msg.Op
	if len(info) == len(addrs) {
		info = map[string]int{}
	}
	go func() { chInfo <- info }()
}
func send(remoteAddr string, msg tmsg) {
	if conn, err := net.Dial("tcp", remoteAddr); err != nil {
		log.Println("No hay respuesta de ", remoteAddr)
	} else {
		defer conn.Close()
		fmt.Println("Enviando a ", remoteAddr)
		enc := json.NewEncoder(conn)
		enc.Encode(msg)
	}
}
