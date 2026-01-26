package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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

	address := net.JoinHostPort(host, strconv.FormatInt(port, 10))

	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

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
		if args == nil {
			continue
		}

		writer.Write(toRespArray(args))
		isSubscribe := strings.ToUpper(args[0]) == "SUBSCRIBE"
		if isSubscribe {
			for {
				res, err := re.Read()
				if err != nil {
					if err == io.EOF {
						continue
					}

					if errors.Is(err, syscall.ECONNRESET) {
						continue
					}

					panic(err)
				}

				fmt.Print(format(res, 0))
				if res.Typ == "error" {
					break
				}
			}
		} else {
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

			fmt.Print(format(res, 0))
		}
	}
}

func toRespArray(args []string) resp.Value {
	var values []resp.Value
	for _, a := range args {
		values = append(values, resp.Bulk(a))
	}

	return resp.Array(values)
}

func format(res resp.Value, nest int) string {
	switch res.Typ {
	case "bulk":
		return fmt.Sprintf("\"%s\"\n", res.Bulk)
	case "string":
		return res.Str + "\n"
	case "integer":
		return fmt.Sprintf("(integer) %d\n", res.Num)
	case "null":
		return "(null)\n"
	case "error":
		return res.Str + "\n"
	case "array":
		var b strings.Builder
		l := utils.Digits(len(res.Array))
		for i, v := range res.Array {
			spaces := l
			if nest > 0 && i > 0 {
				spaces += nest
			}

			s := fmt.Sprintf("%*d) %s", spaces, i+1, format(v, spaces+1+1))
			b.WriteString(s)
		}

		return b.String()
	}

	return ""
}
