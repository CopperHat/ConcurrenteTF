package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

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
	accuracy   int
	Pasos      int
}

var t Arbol

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

//PrintArbol Devuelve el arbol
func PrintArbol() {
	fmt.Println()
	t := new(Arbol)
	t.accuracy = 1
	fmt.Println()
	fmt.Printf("Valores ordenados: %v\n", t.Ordenar())
	fmt.Printf("Cantidad: %v\n", t.Pasos)
	fmt.Print("Arbol:\n")
	jsonArbol := json.NewEncoder(os.Stdout)
	jsonArbol.SetIndent("", "  ")
	_ = jsonArbol.Encode(t.primerNodo)
	fmt.Print("Finalizacion: %", t.accuracy*100)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Trabajo Final de Programacion Concurrente")
}

func agregarRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	peso, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Error leyendo numero: %v", err)
		http.Error(w, "No se lee", http.StatusBadRequest)
		return
	}
	t.Agregar(peso)

}

func arbolRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Arbol Ordenado")
	PrintArbol()
}

func main() {

	t.accuracy = 1

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute).Methods("GET")
	router.HandleFunc("/arbol/{id}", agregarRoute).Methods("POST")
	router.HandleFunc("/arbol", arbolRoute).Methods("GET")

	log.Fatal(http.ListenAndServe(":9000", router))

}
