package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"syscall"

	"example.com/kvs/cli/utils"
	"example.com/kvs/resp"
)

var host string
var port int64

func main() {
	const (
		defaultHost = "127.0.0.1"
		usageHost   = "Server hostname"
		defaultPort = 6379
		usagePort   = "Server port"
	)

	flag.StringVar(&host, "h", defaultHost, usageHost)
	flag.Int64Var(&port, "p", defaultPort, usagePort)
	flag.Parse()

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	re := resp.NewResp(conn)
	writer := resp.NewWriter(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s> ", address)
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		args := utils.Parse(line)

		var values []resp.Value
		for _, a := range args {
			values = append(values, resp.Bulk(a))
		}

		writer.Write(resp.Array(values))
		res, err := re.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			if errors.Is(err, syscall.ECONNRESET) {
				break
			}

			panic(err)
		}

		switch res.Typ {
		case "bulk":
			fmt.Printf("\"%s\"\n", res.Bulk)
		case "string":
			fmt.Println(res.Str)
		case "integer":
			fmt.Printf("(integer) %d\n", res.Num)
		case "null":
			fmt.Println("(null)")
		case "error":
			fmt.Println(res.Str)
		}
	}
}
