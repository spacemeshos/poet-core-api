package poet_core_api

import (
	"context"
	"github.com/spacemeshos/poet-core-api/pcrpc"
	"google.golang.org/grpc"
	"log"
	"time"
)

func NewProverClient(target string) (pcrpc.PoetCoreProverClient, func()) {
	conn := newClientConn(target)

	cleanUp := func() {
		conn.Close()
	}

	return pcrpc.NewPoetCoreProverClient(conn), cleanUp
}

func NewVerifierClient(target string) (pcrpc.PoetVerifierClient, func()) {
	conn := newClientConn(target)

	cleanUp := func() {
		conn.Close()
	}

	return pcrpc.NewPoetVerifierClient(conn), cleanUp
}

func newClientConn(target string) *grpc.ClientConn {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		log.Fatalf("unable to connect to RPC server: %v", err)
	}

	return conn
}
