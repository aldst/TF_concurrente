package main

import (
    "fmt"
    "net"
    "bufio"
    "os"
	"encoding/json"
	"crypto/sha1"
    "encoding/hex"
	"strings"
	"strconv"
)


type Mujer struct {
	id int
	Edad string
	imc string
	hemoglobina string
	hemoglobinaAlt string
}


type Bloque struct {
	hash_actual string
	informacion Mujer
	hash_anterior string
}

var arreglo = make([]Bloque,0)


func cliente() {
    con, _ := net.Dial("tcp", "192.168.1.60:8000")
    defer con.Close()
    r := bufio.NewReader(con)
    gin := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("Dame dato ")
        msg, _ := gin.ReadString('\n')
        fmt.Fprint(con, msg)
		resp, _ := r.ReadString('\n')
		fmt.Print("Resp:", resp)
		fmt.Print(arreglo)
        if len(msg) == 0 || msg[0] == 'x' {
            break
        }
    }
}


var hosts []string = []string{"transacci�n 1",
							"transacción 2",
							"transacción 3",
							"transacción 4",
							"transacción 5"}

func handle(con net.Conn, id int) {
    defer con.Close()
    for {
			poio,_:=json.Marshal(arreglo)
			poio2:=string(poio)
			fmt.Fprintln(con, poio2)
    }
    fmt.Printf("Con%d anemia!\n", id)
}

func servidor() {
    ln, _ := net.Listen("tcp", "192.168.1.60:8000")
    defer ln.Close()
    cont := 0
    for {
		con, _ := ln.Accept()
        go handle(con, cont)
        cont++
    }
}



func generar_bloque() {
		contador := 0
		for {
			//var slice = make([]Bloque,elems)
			gin := bufio.NewReader(os.Stdin)

			fmt.Print("Ingrese Edad ")
			msg2, _ := gin.ReadString('\n')
			fmt.Print("Ingrese imc ")
			msg3, _ := gin.ReadString('\n')
			fmt.Print("Ingrese hemoglobina: ")
			msg4, _ := gin.ReadString('\n')
			fmt.Print("Ingrese hemoglobina altura: ")
			msg5, _ := gin.ReadString('\n')
			s :=Mujer{contador+1, msg2, msg3, msg4,msg5}
			h := sha1.New()
			bytearray := []byte(fmt.Sprintf("%v", s))
			h.Write([]byte(bytearray))
			sha1_hash := hex.EncodeToString(h.Sum(nil))
			if len(arreglo) >= 1{
				hola := Bloque{sha1_hash, s , arreglo[contador-1].hash_actual}
				arreglo = append(arreglo, hola)
			}else{
				hola := Bloque{sha1_hash, s , "---"}
				arreglo = append(arreglo, hola)
			}
			fmt.Println(arreglo[contador])
			contador++
			fmt.Println("Desea ingresar otro pacientes? (Y/N)")
			msg6, _ := gin.ReadString('\n')
			if strings.Compare(msg6, "Y") == 1{
			}else {
				break
			}
		}
}



func main(){
	for {
		fmt.Println("Selecciona una opcion: ")
		fmt.Println("1) Generar un nuevo Bloque")
		fmt.Println("2) Servidor")
		fmt.Println("3) Clientes")

		in:=bufio.NewScanner(os.Stdin)
		in.Scan()
		int,_:=strconv.Atoi(in.Text())
		fmt.Println(int)
		if int == 1 {
			generar_bloque()
		}
		if int == 2 {
			servidor()
		}
		if int == 3 {
			cliente()
		}
	}
}
