package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/iyudintsev/caesr-go-client/internal/client"
	"github.com/iyudintsev/caesr-go-client/internal/config"
	pb "github.com/iyudintsev/caesr-go-client/proto"
	sherpa "github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	inputWavFile := flag.String("InputWavFile", "", "")
	flag.Parse()

	cfg, err := config.GetConfig()
	if err != nil {
		panic("can't get config")
	}

	conn, err := grpc.NewClient(cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("can't get grpc connection")
	}

	cc := client.NewCaesrClient(conn)
	inputStream := make(chan *pb.CaesrRequest)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("error while connection closing: %s\n", err.Error())
		}
		close(inputStream)
		cancel()
	}()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return cc.Recognize(ctx, inputStream)
	})

	wave := sherpa.ReadWave(*inputWavFile)
	audioChunkNums := len(wave.Samples) / cfg.WindowSize
	for audioChunkIndex := range audioChunkNums {
		left := audioChunkIndex * cfg.WindowSize
		right := (audioChunkIndex + 1) * cfg.WindowSize
		if audioChunkIndex == audioChunkNums-1 {
			right = len(wave.Samples)
		}
		req := &pb.CaesrRequest{SampleRate: 16000}
		for index := left; index < right; index += 1 {
			req.AudioChunk = append(req.AudioChunk, wave.Samples[index])
		}
		inputStream <- req
	}

}
