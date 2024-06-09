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

const (
	StatusOK         = "HTTP/1.1 200 OK\r\n"
	CRLF             = "\r\n"
	ContentTypeText  = "Content-Type: text/plain"
	ContentTypeOctet = "Content-Type: application/octet-stream"
	ContentLength    = "Content-Length: "
	Status404        = "HTTP/1.1 404 Not Found\r\n"
)

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

	} else if strings.Contains(request.URL.Path, "/files") {
		// /files/foo

		url := request.URL.Path
		fileName := strings.TrimPrefix(url, "/files/")
		dir := os.Args[2]
		fmt.Println("fileName: ", fileName)
		fmt.Println("dir+fileName: ", dir+fileName)
		data, err := os.ReadFile(dir + fileName)
		// 文件名字成功读取到了，但是还是没有搞出来，这是为什么呢？
		// 可能是路径问题
		if err != nil {
			fmt.Println("err: ", err)
			conn.Write([]byte(Status404 + CRLF))
			return
		}
		fmt.Println("success")
		res := StatusOK + ContentTypeOctet + CRLF + ContentLength + strconv.Itoa(len(data)) + CRLF + CRLF + string(data)
		conn.Write([]byte(res))
	} else if request.URL.Path == "/user-agent" {
		header := request.Header
		value := header.Get("User-Agent")
		fmt.Println("value: ", value)
		length := len(value)
		res := StatusOK + ContentTypeText + CRLF + ContentLength + strconv.Itoa(length) + CRLF + CRLF + value
		conn.Write([]byte(res))
	} else if strings.Contains(request.URL.Path, "/echo") {
		// /echo/{str}
		path := request.URL.Path
		str := strings.Split(path, "/")[2]
		fmt.Println("str: ", str)
		length := len(str)
		fmt.Println("len: ", length)
		fmt.Println("string(len): ", string(length))

		contentType := ContentTypeText + CRLF
		contentLength := ContentLength + strconv.Itoa(length) + CRLF
		res := StatusOK + contentType + contentLength + CRLF + str
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

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		// 并发处理 请求
		go Handler(conn)
		//还可以直接用 conn.Read 方法来做解析，先读进来,再用 "\r\n" 切分
	}

}
