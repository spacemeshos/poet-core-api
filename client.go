package poet_core_api

import (
	"github.com/spacemeshos/poet-core-api/pcrpc"
	"google.golang.org/grpc"
	"log"
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
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		log.Fatalf("unable to connect to RPC server: %v", err)
	}

	return conn
}

