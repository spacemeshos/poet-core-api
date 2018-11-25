package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	api "github.com/spacemeshos/poet-core-api"
	"github.com/spacemeshos/poet-core-api/pcrpc"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"log"
	"os"
	"strconv"
	"strings"
)

type EnumValue struct {
	Enum     []string
	Default  string
	selected string
}

func (e *EnumValue) Set(value string) error {
	for _, enum := range e.Enum {
		if enum == value {
			e.selected = value
			return nil
		}
	}

	return fmt.Errorf("allowed values are %s", strings.Join(e.Enum, ", "))
}

func (e EnumValue) String() string {
	if e.selected == "" {
		return e.Default
	}
	return e.selected
}

var computeCommand = cli.Command{
	Name:      "compute",
	Category:  "prover",
	Usage:     "compute a new dag",
	ArgsUsage: "x n h",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "x",
			Usage: "the commitment input",
		},
		cli.Uint64Flag{
			Name:  "n",
			Usage: "time parameter",
			Value: 1,
		},
		cli.GenericFlag{
			Name:  "hf",
			Usage: "hash function",
			Value: &EnumValue{
				Enum:    []string{"sha256", "scrypt"},
				Default: "sha256",
			},
		},
	},
	Action: compute,
}

func compute(ctx *cli.Context) error {
	var (
		x []byte
		n uint32
		h string
	)

	client, cleanUp := api.NewProverClient(	ctx.GlobalString("rpcserver"))
	defer cleanUp()

	args := ctx.Args()

	switch {
	case ctx.IsSet("x"):
		x = []byte(ctx.String("x"))
	case args.Present():
		x = []byte(args.First())
		args = args.Tail()
	default:
		return fmt.Errorf("x argument missing")
	}

	switch {
	case ctx.IsSet("n"):
		n = uint32(ctx.Uint64("n"))
	case args.Present():
		i, err := strconv.ParseInt(args.First(), 10, 32)
		if err != nil {
			return fmt.Errorf("invalid n argument")
		}
		n = uint32(i)
		args = args.Tail()
	default:
		return fmt.Errorf("n argument missing")
	}

	switch {
	case ctx.IsSet("hf"):
		h = ctx.Generic("hf").(*EnumValue).String()
	case args.Present():
		h = args.First()
		args = args.Tail()
	default:
		return fmt.Errorf("hf argument missing")
	}

	res, err := client.Compute(context.Background(), &pcrpc.ComputeRequest{D: &pcrpc.DagParams{
		X: x,
		N: n,
		H: h,
	}})
	if err != nil {
		return err
	}

	printJSON(struct {
		Phi string `json:"phi"`
	}{
		Phi: fmt.Sprintf("%x", res.Phi),
	},
	)
	return nil
}

var getNIPCommand = cli.Command{
	Name:     "getnip",
	Category: "prover",
	Usage:    "get the dag non-interactive proof",
	Action:   getNIP,
}

func getNIP(ctx *cli.Context) error {
	client, cleanUp := api.NewProverClient(	ctx.GlobalString("rpcserver"))
	defer cleanUp()

	res, err := client.GetNIP(context.Background(), &pcrpc.GetNIPRequest{})
	if err != nil {
		return err
	}

	var l []string
	for _, labels := range res.Proof.L {
		l = append(l, fmt.Sprintf("%x", labels))
	}
	printJSON(struct {
		Phi string   `json:"phi"`
		L   []string `json:"l"`
	}{
		Phi: fmt.Sprintf("%x", res.Proof.Phi),
		L:   l,
	},
	)
	return nil
}

var getProofCommand = cli.Command{
	Name:      "getproof",
	Category:  "prover",
	Usage:     "get a proof to a challenge",
	ArgsUsage: "c",
	Flags: []cli.Flag{
		cli.StringSliceFlag{
			Name:  "c",
			Usage: "the challenge",
		},
	},
	Action: getProof,
}

func getProof(ctx *cli.Context) error {

	var (
		c []string
	)

	client, cleanUp := api.NewProverClient(	ctx.GlobalString("rpcserver"))
	defer cleanUp()

	args := ctx.Args()

	switch {
	case ctx.IsSet("c"):
		c = ctx.StringSlice("c")
	case args.Present():
		c = args
	default:
		return fmt.Errorf("c argument missing")
	}

	res, err := client.GetProof(context.Background(), &pcrpc.GetProofRequest{C: c})
	if err != nil {
		return err
	}

	printJSON(res)

	return nil
}

var cleanCommand = cli.Command{
	Name:     "clean",
	Category: "prover",
	Usage:    "clean the dag",
	Action:   clean,
}

func clean(ctx *cli.Context) error {
	client, cleanUp := api.NewProverClient(	ctx.GlobalString("rpcserver"))
	defer cleanUp()

	res, err := client.Clean(context.Background(), &pcrpc.CleanRequest{})

	if err != nil {
		return err
	}

	printJSON(res)
	return nil
}

var verifyCommand = cli.Command{
	Name:     "verify",
	Category: "verifier",
	Usage:    "verify a proof",
	Action:   verify,
}

func verify(ctx *cli.Context) error {
	client, cleanUp := api.NewVerifierClient(ctx.GlobalString("rpcserver"))
	defer cleanUp()

	res, err := client.VerifyProof(context.Background(), &pcrpc.VerifyProofRequest{})

	if err != nil {
		return err
	}

	printJSON(res)
	return nil
}

func printJSON(resp interface{}) {
	b, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	out.WriteString("\n")
	out.WriteTo(os.Stdout)
}
