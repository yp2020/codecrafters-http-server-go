package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
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
	ContentEncoding  = "Content-Encoding: "
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
	//fmt.Printf("Request: %s %s \n", request.Method, request.URL.Path)
	if request.Method == "POST" {
		handlerPostRequest(request, conn)
	} else if request.Method == "GET" {
		handlerGetRequest(request, conn)
	} else {
		fmt.Println("Not Support " + request.Method + " Method")
	}
}

func handlerPostRequest(request *http.Request, conn net.Conn) {
	path := request.URL.Path
	if strings.Contains(path, "/files/") {
		fileName := strings.TrimPrefix(path, "/files/")
		//// 测试情况下 dir 设置为当前目录 win 用 \\
		//dir, _ := os.Getwd()
		//filePath := dir + "\\" + fileName
		dir := os.Args[2]
		filePath := dir + fileName
		fmt.Println("filePath: ", filePath)
		// 向文件中写入内容
		value, _ := io.ReadAll(request.Body)
		if err := os.WriteFile(filePath, value, 0666); err != nil {
			fmt.Println("err")
		}
		conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
	}

}

func handlerGetRequest(request *http.Request, conn net.Conn) {
	switch path := request.URL.Path; {
	case path == "/":
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		return
	case path == "/user-agent":
		head := request.Header
		value := head.Get("User-Agent")
		length := len(value)
		contentLength := ContentLength + strconv.Itoa(length)
		res := StatusOK + ContentTypeText + CRLF + contentLength + CRLF + CRLF + value
		conn.Write([]byte(res))
		return

	case strings.Contains(path, "/files"):
		// /files/foo
		url := request.URL.Path
		fileName := strings.TrimPrefix(url, "/files/")
		// 这里的路径是由参数带进来的
		dir := os.Args[2]
		//fmt.Println("fileName: ", fileName)
		//fmt.Println("dir+fileName: ", dir+fileName)
		data, err := os.ReadFile(dir + fileName)
		if err != nil {
			fmt.Println("err: ", err)
			conn.Write([]byte(Status404 + CRLF))
			return
		}
		res := StatusOK + ContentTypeOctet + CRLF + ContentLength + strconv.Itoa(len(data)) + CRLF + CRLF + string(data)
		conn.Write([]byte(res))
	case strings.Contains(path, "/echo"):
		// /echo/{str}
		path := request.URL.Path
		str := strings.Split(path, "/")[2]
		length := len(str)
		contentType := ContentTypeText + CRLF
		contentLength := ContentLength + strconv.Itoa(length) + CRLF
		// 获取 accept-encoding
		encodingStr := request.Header.Get("accept-encoding")
		//fmt.Println("encodingStr: ", encodingStr)

		var res string
		if checkValidEncoding(encodingStr) {
			var buf bytes.Buffer
			writer := gzip.NewWriter(&buf)
			writer.Write([]byte(str))
			writer.Close()
			content := buf.String()
			res = StatusOK + ContentEncoding + "gzip" + CRLF + contentType + ContentLength + fmt.Sprint(len(content)) + CRLF + CRLF + content
		} else {
			res = StatusOK + contentType + contentLength + CRLF + str
		}
		//fmt.Println("res: ", res)
		conn.Write([]byte(res))
	default:
		conn.Write([]byte(Status404 + CRLF))
	}
}

func checkValidEncoding(encodingStr string) bool {
	encodingArr := strings.Split(encodingStr, ", ")
	//fmt.Println("encodingArr: ", encodingArr)
	for _, value := range encodingArr {
		if value == "gzip" {
			return true
		}
	}
	return false
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
	}
}
