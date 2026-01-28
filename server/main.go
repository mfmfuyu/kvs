package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"
	"syscall"

	"example.com/kvs/resp"
	"example.com/kvs/server/client"
	"example.com/kvs/server/cmd"
	"example.com/kvs/server/request"
)

var Handlers = map[string]func(*request.Request){
	"PING":        cmd.Ping,
	"SET":         cmd.Set,
	"GET":         cmd.Get,
	"EXPIRE":      cmd.Expire,
	"TTL":         cmd.Ttl,
	"SUBSCRIBE":   cmd.Subscribe,
	"UNSUBSCRIBE": cmd.UnSubscribe,
	"PUBLISH":     cmd.Publish,
}
var port int64

func main() {
	const (
		defaultPort = 6379
		usage       = "TCP port number to listen on"
	)

	flag.Int64Var(&port, "port", defaultPort, usage)
	flag.Parse()

	address := fmt.Sprintf(":%d", port)

	l, err := net.Listen("tcp", address)
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

	client := client.New(writer)

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
			client.Write(resp.Error(errMsg))

			continue
		}

		request := &request.Request{
			Args:   args,
			Client: client,
		}

		handler(request)
	}
}
