package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Types
type task struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
}

type allTasks []task

// Persistence
var tasks = allTasks{
	{
		ID:      1,
		Name:    "Agregar",
		Content: "Agrega Nodo a lista",
	},
	{
		ID:      2,
		Name:    "Actualizar",
		Content: "Crea el Arbol",
	},
	{
		ID:      3,
		Name:    "Arbol",
		Content: "Muestra el Arbol",
	},
}

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
ss
// Agregar agrega nodo a Arbol
func (t *Arbol) Agregar(peso int) {

}

// Ordenar ordena nodos en Arbol
func (t *Arbol) Ordenar() []int {
	t.Pasos = 0
	t.ordenado = []int{}
	t.Recorrer(t.primerNodo)
	return t.ordenado
}

// Devuelve el arbol
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
	reqBody, err := ioutil.ReadAll(r.Body)
	peso = 0
	peso, err = strconv.Atoi(reqBody)
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Task Data")
	}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

}

func ordenarRoute(w http.ResponseWriter, r *http.Request) {
	Ordenar()
	fmt.Fprintf(w, "Arbol Ordenado")
}

func arbolRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ARbol Ordenado")
	PrintArbol()
}

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tareas", getTasks)
	router.HandleFunc("/agregar/{id}", agregarRoute).Methods("POST")
	router.HandleFunc("/ordenar", ordenarRoute)
	router.HandleFunc("/arbol", arbolRoute).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))

}
