package cmd

import (
	"strconv"
	"time"

	"example.com/kvs/kv"
	"example.com/kvs/resp"
	"example.com/kvs/server/request"
)

func Expire(req *request.Request) {
	if len(req.Args) != 2 {
		req.Client.Write(InvalidArg("expire"))
		return
	}

	key := req.Args[0].Bulk
	rawSeconds := req.Args[1].Bulk
	seconds, err := strconv.ParseInt(rawSeconds, 10, 64)
	if err != nil {
		req.Client.Write(resp.Error(ERR_NUM))
		return
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(seconds) * time.Second)

	ok := kv.SetExpiresAt(key, expiresAt)
	if !ok {
		req.Client.Write(resp.Integer(0))
		return
	}

	req.Client.Write(resp.Integer(1))
}
