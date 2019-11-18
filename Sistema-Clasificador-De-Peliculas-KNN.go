package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"github.com/bradfitz/slice"
	"os"
	"strconv"
	"strings"
)

type Peliculas struct {
	Nombre           string `json:"nombre"`
	Clasificacion    string `json:"clasification"`
	Genero           string `json:"gender"`
	Annio            string `json:"year"`
	Gusto            string `json:"like"`
	HostRegisterPort string `json:"hostRegisterPort"`
	HostNotifyPort   string `json:"hostNotifyPort"`
}

type Notificacion struct {
	Valor_distancia string `json:"valorDistancia"`
	Peliculas       string `json:"peliculas"`
}

type NotificarCliente struct {
	Pelicula string   `json:"pelicula"`
	Puertos  []string `json:"puertos"`
}

var k string
var modo string
var ports []string
var portChan = make(chan []Notificacion, 2)

func Algoritmo_Frecuencia(notify []Notificacion) string {
	if k == "entrenamiento" {
		return "true"
	}

	result := notify
	slice.Sort(result[:], func(i, j int) bool {
		return result[i].Valor_distancia < result[j].Valor_distancia
	})
	cont_1 := 0
	cont_2 := 0
	kInt, _ := strconv.Atoi(k)
	for i := 0; i < kInt; i++ {
		if result[i].Peliculas == "true" {
			cont_1 = cont_1 + 1
		} else {
			cont_2 = cont_1 + 1
		}
	}
	if cont_1 > cont_2 {
		return "true"
	} else {
		return "false"
	}
}
func AlgoritmoKNN(peliculax, peliculay Peliculas) float64 {
	a1, _ := strconv.Atoi(peliculax.Clasificacion)
	a2, _ := strconv.Atoi(peliculay.Clasificacion)
	a := a1 - a2
	b1, _ := strconv.Atoi(peliculax.Genero)
	b2, _ := strconv.Atoi(peliculay.Genero)
	b := b1 - b2
	c1, _ := strconv.Atoi(peliculax.Annio)
	c2, _ := strconv.Atoi(peliculay.Annio)
	c := c1 - c2
	result := math.Sqrt(math.Pow(float64(a), 2) + math.Pow(float64(b), 2) + math.Pow(float64(c), 2))
	return result
}

func notify(port string, peliculaIn Peliculas) Notificacion {
	remotehost := fmt.Sprintf(":%s", port)
	conn, _ := net.Dial("tcp", remotehost)
	defer conn.Close()
	peliculaJson, _ := json.Marshal(peliculaIn)
	fmt.Fprintln(conn, string(peliculaJson))
	r := bufio.NewReader(conn)
	msg, err := r.ReadString('\n')
	var notifyIn Notificacion
	json.Unmarshal([]byte(msg), &notifyIn)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	return notifyIn
}

func tellEverybody(pelicula, peliculaIn Peliculas) {

	var notifies []Notificacion

	for _, port := range ports {
		if strings.Compare(port, pelicula.HostNotifyPort) != 0 {
			notifies = append(notifies, notify(port, peliculaIn))
		}
	}
	portChan <- notifies
}

func handleRegister(conn net.Conn, pelicula Peliculas) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	msg, _ := r.ReadString('\n')
	var peliculaIn Peliculas
	json.Unmarshal([]byte(msg), &peliculaIn)
	if len(ports) == 0 {
		ports = append(ports, pelicula.HostNotifyPort)
	}
	tellEverybody(pelicula, peliculaIn)
	resultadoRegisterNuevo := fmt.Sprintf("%f", AlgoritmoKNN(pelicula, peliculaIn))
	notifyOut := <-portChan
	notifyOut = append(notifyOut, Notificacion{resultadoRegisterNuevo, pelicula.Gusto})
	ganador := Algoritmo_Frecuencia(notifyOut)
	notifyClient := NotificarCliente{ganador, ports}
	notifyClientJson, _ := json.Marshal(notifyClient)
	fmt.Fprintln(conn, string(notifyClientJson))
	ports = append(ports, peliculaIn.HostNotifyPort)
}

func registerServer(hostRegisterPort string, pelicula Peliculas) {
	host := fmt.Sprintf(":%s", hostRegisterPort)
	ln, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer ln.Close()
	for {
		conn, errAccept := ln.Accept()
		if errAccept != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRegister(conn, pelicula)
	}
}
func registerClient(remotePort2 string, pelicula *Peliculas) {
	remotehost := fmt.Sprintf(":%s", remotePort2)
	peliculaJson, _ := json.Marshal(pelicula)
	conn, err := net.Dial("tcp", remotehost)
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Fprintln(conn, string(peliculaJson))

	r := bufio.NewReader(conn)
	msg, _ := r.ReadString('\n')

	var notifyIn NotificarCliente
	json.Unmarshal([]byte(msg), &notifyIn)
	if modo != "E" {
		pelicula.Gusto = notifyIn.Pelicula
	}
	ports = append(notifyIn.Puertos, pelicula.HostNotifyPort)
	switch pelicula.Gusto {
	case "true":
		fmt.Println("Pelicula para su gusto! Debe verla!")
	case "false":
		fmt.Printf("No es muy recomendable para sus gustos :/")
	}
}
func handleNotify(conn net.Conn, pelicula Peliculas) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	msg, _ := r.ReadString('\n')
	var peliculaIn Peliculas
	json.Unmarshal([]byte(msg), &peliculaIn)
	resultado := fmt.Sprintf("%f", AlgoritmoKNN(pelicula, peliculaIn))
	notify := Notificacion{
		resultado,
		pelicula.Gusto,
	}
	notifyJson, _ := json.Marshal(notify)
	ports = append(ports, peliculaIn.HostNotifyPort)
	fmt.Fprintln(conn, string(notifyJson))

}
func notifyServer(hostNotifyPort string, pelicula Peliculas) {
	host := fmt.Sprintf(":%s", hostNotifyPort)
	ln, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println("Error al escuchar la Notificación:", err.Error())
		os.Exit(1)
	}
	defer ln.Close()
	for {
		conn, errAccept := ln.Accept()
		if errAccept != nil {
			fmt.Println("Error al aceptar la notificación: ", err.Error())
			os.Exit(1)
		}
		go handleNotify(conn, pelicula)
	}
}

func main() {
	var nombre string
	var clasificacion string
	var genero string
	var annio string
	var gusto string
	var hostRegisterPort string
	var hostNotifyPort string

	ginName := bufio.NewReader(os.Stdin)
	fmt.Print("Nombre de la pelicula: ")
	nombre_In, _ := ginName.ReadString('\n')
	nombre = strings.TrimSpace(nombre_In)

	ginClasification := bufio.NewReader(os.Stdin)
	fmt.Print("Clasificacion de la pelicula (Numero del 1 al 10): ")
	clasificacion_In, _ := ginClasification.ReadString('\n')
	clasificacion = strings.TrimSpace(clasificacion_In)

	ginGender := bufio.NewReader(os.Stdin)
	fmt.Print("Genero de la pelicula (Accion, Suspenso, Terror): ")
	genero_In, _ := ginGender.ReadString('\n')
	genero = strings.TrimSpace(genero_In)

	switch strings.ToLower(genero) {
	case "accion":
		genero = "0"
	case "suspenso":
		genero = "1"
	case "terror":
		genero = "2"
	}

	ginYear := bufio.NewReader(os.Stdin)
	fmt.Print("Año de la pelicula: ")
	annio_In, _ := ginYear.ReadString('\n')
	annio = strings.TrimSpace(annio_In)

	ginPort := bufio.NewReader(os.Stdin)
	fmt.Print("Puerto de registro del Host: ")
	hostRegisterPort_In, _ := ginPort.ReadString('\n')
	hostRegisterPort = strings.TrimSpace(hostRegisterPort_In)

	gin_Port := bufio.NewReader(os.Stdin)
	fmt.Print("Puerto de notificacion del Host: ")
	hostNotifyPort_In, _ := gin_Port.ReadString('\n')
	hostNotifyPort = strings.TrimSpace(hostNotifyPort_In)

	ginModo := bufio.NewReader(os.Stdin)
	fmt.Print("¿Entrenar al sistema?(E) ó Ponerlo a prueba(P): ")
	modo, _ = ginModo.ReadString('\n')
	modo = strings.TrimSpace(modo)

	if strings.ToLower(modo) == "e" {
		gingusto := bufio.NewReader(os.Stdin)
		fmt.Print("¿Le gustó o no? (si o no): ")
		gusto, _ = gingusto.ReadString('\n')
		gusto = strings.TrimSpace(gusto)
		if strings.ToLower(gusto) == "si" {
			gusto = "true"
		} else {
			gusto = "false"
		}
		k = "entrenamiento"
	} else {
		ginK := bufio.NewReader(os.Stdin)
		fmt.Print("Inserte el valor de k : ")
		k, _ = ginK.ReadString('\n')
		k = strings.TrimSpace(k)
		gusto = "false"
	}

	pelicula := Peliculas{
		nombre, clasificacion, genero, annio, gusto, hostRegisterPort, hostNotifyPort}

	go registerServer(hostRegisterPort, pelicula)
	go notifyServer(hostNotifyPort, pelicula)

	gin2 := bufio.NewReader(os.Stdin)
	fmt.Print("Introduzca el puerto a conectarse: ")
	remotePort2, _ := gin2.ReadString('\n')
	remotePort2 = strings.TrimSpace(remotePort2)

	if len(remotePort2) > 0 {
		registerClient(remotePort2, &pelicula)
	}
	for {

	}
}
