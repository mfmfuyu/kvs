package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"syscall"

	"example.com/kvs/cmd"
	"example.com/kvs/resp"
)

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING": cmd.Ping,
	"SET":  cmd.Set,
	"GET":  cmd.Get,
}

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(conn)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	res := resp.NewResp(conn)
	writer := resp.NewWriter(conn)

	for {
		value, err := res.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			if errors.Is(err, syscall.ECONNRESET) {
				break
			}

			panic(err)
		}

		if value.Typ != "array" {
			continue
		}

		if len(value.Array) == 0 {
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		handler, ok := Handlers[command]
		if !ok {
			strArgs := []string{}
			for i := range args {
				strArgs = append(strArgs, fmt.Sprintf("'%s'", args[i].Bulk))
			}

			errMsg := fmt.Sprintf("ERR unknown command '%s', with args beginning with: %s", command, strings.Join(strArgs, " "))
			writer.Write(resp.Error(errMsg))

			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
