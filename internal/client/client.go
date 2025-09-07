package client

import (
	"context"
	"fmt"
	"io"

	pb "github.com/iyudintsev/caesr-go-client/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type CaesrClient struct {
	conn *grpc.ClientConn
}

func NewCaesrClient(conn *grpc.ClientConn) *CaesrClient {
	return &CaesrClient{conn: conn}
}

func (cc *CaesrClient) send(
	ctx context.Context, steam pb.CaesrService_RecognizeClient,
	inputStream chan *pb.CaesrRequest,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case req, ok := <-inputStream:
			if !ok {
				return nil
			}
			if err := steam.Send(req); err != nil {
				return err
			}
		}
	}
}

func (cc *CaesrClient) recv(ctx context.Context, steam pb.CaesrService_RecognizeClient) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			resp, err := steam.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			fmt.Printf("Transcript: %s\n", resp.Transcript)
		}
	}
}

func (cc *CaesrClient) Recognize(ctx context.Context, inputStream chan *pb.CaesrRequest) error {
	client := pb.NewCaesrServiceClient(cc.conn)
	stream, err := client.Recognize(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := stream.CloseSend(); err != nil {
			fmt.Printf("[client] error: %s\n", err.Error())
		}
	}()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return cc.send(ctx, stream, inputStream)
	})
	g.Go(func() error {
		return cc.recv(ctx, stream)
	})

	return g.Wait()
}
