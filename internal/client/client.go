package client

import (
	pb "github.com/iyudintsev/caesr-go-client/proto"
)

func DoSmth() pb.CaesrResponse {
	var resp pb.CaesrResponse
	resp.Transcript = "hello"
	return resp
}
