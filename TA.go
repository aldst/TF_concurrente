package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
)

var k int
var searchNode [5]float64 //nodo de busqueda
var nodes []Node
var addrs []string

const numCoord = 4

type Node struct {
	age, imc, hemoglobine, hem_high, dist float64
}

func read_file() {

	var c1, c2, c3, c4, dist float64

	file, error := os.Open("./dataset.txt")

	if error != nil {
		fmt.Println("Hubo un error")
	}

	defer func() {
		file.Close()
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		value := strings.Split(line, ",")
		for j := 0; j < numCoord+1; j++ {

			valueF, _ := strconv.ParseFloat(value[j], 32)
			//fmt.Println(reflect.TypeOf(valueF))

			if j == 0 {
				c1 = valueF
			}
			if j == 1 {
				c2 = valueF
			}
			if j == 2 {
				c3 = valueF
			}
			if j == 3 {
				c4 = valueF
			}
		}
		dist = math.Sqrt(math.Pow(searchNode[0]-c1, 2) + math.Pow(searchNode[1]-c2, 2) + math.Pow(searchNode[2]-c3, 2) + math.Pow(searchNode[3]-c4, 2))

		nodes = append(nodes, Node{c1, c2, c3, c4, dist})

	}

}

func inicialize() {
	fmt.Print("Ingrese el numero de nodos a analizar: ")
	fmt.Scanf("%d", &k)

	for i := 0; i < numCoord; i++ {

		if i == 0 {
			fmt.Print("Ingrese su edad [entre 12 y 50]: ")
		}
		if i == 1 {
			fmt.Print("Ingrese el indice de masa corporal [entre 12.00 y 100.00]: ")
		}
		if i == 2 {
			fmt.Print("Ingrese su nivel de hemoglobina: ")
		}
		if i == 3 {
			fmt.Print("Ingrese su nivel de hemoglobina ajustado por altitud:  ")
		}
		fmt.Scanf(" \n%f", &(searchNode[i]))
	}
	read_file()

}

func sortDistances() {
	sort.Slice(nodes, func(p, q int) bool {
		return nodes[p].dist > nodes[q].dist
	})
}

func findTeam() {
	var low []Node
	var medium []Node

	for i := 0; i < k; i++ {
		if nodes[i].hem_high == 0 {
			low = append(low, Node{nodes[i].age, nodes[i].imc, nodes[i].hemoglobine, nodes[i].hem_high, nodes[i].dist})
		}
		if nodes[i].hem_high == 1 {
			medium = append(medium, Node{nodes[i].age, nodes[i].imc, nodes[i].hemoglobine, nodes[i].hem_high, nodes[i].dist})
		}
	}

	for i := 0; i < len(low); i++ {
		fmt.Println(low[i])
	}
	fmt.Println()
	for i := 0; i < len(medium); i++ {
		fmt.Println(medium[i])
	}
	fmt.Println()

	if len(low) > len(medium) {
		fmt.Println("La evaluación de desempeño del asistente es bajo")
	} else {
		fmt.Println("La evaluación de desempeño del asistente es regular")
	}
}
func main() {

	myip := "192.168.1.6"
	fmt.Printf("Soy %s\n", myip)
	go registerServer(myip)
	go hotServer(myip)
	// falta agregarle el hotServer

	gin := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese direccion remota: ")
	remoteIp, _ := gin.ReadString('\n')
	remoteIp = strings.TrimSpace(remoteIp)

	if remoteIp != "" {
		registerSend(remoteIp, myip)
	}

	go func() {
		inicialize()
		sortDistances()
		hotSend()
	}()

	notifyServer(myip)
}

func hotServer(hostAddr string) {
	host := fmt.Sprintf("%s:8002", hostAddr)
	ln, _ := net.Listen("tcp", host)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handleHot(conn)
	}
}

func handleHot(conn net.Conn) {
	defer conn.Close()

	findTeam()
}

func notifyServer(hostAddr string) {
	host := fmt.Sprintf("%s:8001", hostAddr)
	ln, _ := net.Listen("tcp", host)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handleNotify(conn)
	}
}

func handleNotify(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	remoteIp, _ := r.ReadString('\n')
	remoteIp = strings.TrimSpace(remoteIp)

	for _, addr := range addrs {
		if addr == remoteIp {
			return
		}
	}
	addrs = append(addrs, remoteIp)
	fmt.Println(addrs)

}

func hotSend() {
	idx := rand.Intn(len(addrs))
	fmt.Printf("Enviando a %s\n", addrs[idx])
	remote := fmt.Sprintf("%s:8002", addrs[idx])
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	println("prueba")
	findTeam()
}

// funciones del servidor
func registerServer(hostAddr string) {
	host := fmt.Sprintf("%s:8000", hostAddr)
	ln, _ := net.Listen("tcp", host)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		go handleRegister(conn)
	}
}

func handleRegister(conn net.Conn) {
	defer conn.Close()

	// Recibimos addr del nuevo nodo
	r := bufio.NewReader(conn)
	remoteIp, _ := r.ReadString('\n')
	remoteIp = strings.TrimSpace(remoteIp)

	// respondemos enviando lista de direcciones de nodos actuales
	byteAddrs, _ := json.Marshal(addrs)
	fmt.Fprintf(conn, "%s\n", string(byteAddrs))

	// notificar a nodos actuales de llegada de nuevo nodo
	for _, addr := range addrs {
		notifySend(addr, remoteIp)
	}

	// Agregamos nuevo nodo a la lista de direcciones
	for _, addr := range addrs {
		if addr == remoteIp {
			return
		}
	}
	addrs = append(addrs, remoteIp)
	fmt.Println(addrs)
}

func notifySend(addr, remoteIp string) {
	remote := fmt.Sprintf("%s:8001", addr)
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()
	fmt.Fprintln(conn, remoteIp)
}

// funciones del servidor remoto
func registerSend(remoteAddr, hostAddr string) {
	remote := fmt.Sprintf("%s:8000", remoteAddr)
	conn, _ := net.Dial("tcp", remote)
	defer conn.Close()

	// Enviar direccion
	fmt.Fprintln(conn, hostAddr)

	// Recibir lista de direcciones
	r := bufio.NewReader(conn)
	strAddrs, _ := r.ReadString('\n')
	var respAddrs []string
	json.Unmarshal([]byte(strAddrs), &respAddrs)

	// agregamos direcciones de nodos a propia libreta
	for _, addr := range respAddrs {
		if addr == remoteAddr {
			return
		}
	}
	addrs = append(respAddrs, remoteAddr)
	fmt.Println(addrs)
}
