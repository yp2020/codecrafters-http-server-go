package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const StausOK = "HTTP/1.1 200 OK\r\n"
const CRLF = "\r\n"

func Handler(conn net.Conn) {
	defer conn.Close()
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("Error reading request.", err.Error())
		return
	}
	fmt.Printf("Request: %s %s \n", request.Method, request.URL.Path)

	if request.URL.Path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return

	} else if strings.Contains(request.URL.Path, "/echo") {
		// /echo/{str}
		path := request.URL.Path
		str := strings.Split(path, "/")[2]
		fmt.Println("str: ", str)
		length := len(str)
		fmt.Println("len: ", length)
		fmt.Println("string(len): ", string(length))

		contentType := "Content-Type: text/plain" + CRLF
		contentLength := "Content-Length: " + strconv.Itoa(length) + CRLF
		res := StausOK + contentType + contentLength + CRLF + str
		fmt.Println("res: ", res)
		conn.Write([]byte(res))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	Handler(conn)
	//还可以直接用 conn.Read 方法来做解析，先读进来,再用 "\r\n" 切分

}
