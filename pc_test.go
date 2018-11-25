package poet_core_api

import (
	"context"
	"github.com/spacemeshos/poet-core-api/pcrpc"
	"testing"
)

// TODO: improve test, create inner tests loop for different n & h params
func TestProverAndVerifier(t *testing.T) {
	ctx := context.Background()
	prover, cleanUp := NewProverClient(DefaultRPCHostPort)
	defer cleanUp()
	defer prover.Clean(ctx, &pcrpc.CleanRequest{})
	verifier, cleanUp := NewVerifierClient(DefaultRPCHostPort)
	defer cleanUp()

	var (
		x = []byte("this is a commitment")
		n = uint32(6)
		h = "sha256"
	)

	d := &pcrpc.DagParams{X: x, N: n, H: h}

	_, err := prover.Compute(ctx, &pcrpc.ComputeRequest{D: d})
	if err != nil {
		t.Fatal(err)
	}

	// verify NIP

	nipRes, err := prover.GetNIP(ctx, &pcrpc.GetNIPRequest{})
	if err != nil {
		t.Fatal(err)
	}

	verifyNIPRes, err := verifier.VerifyNIP(ctx, &pcrpc.VerifyNIPRequest{D: d, P: nipRes.Proof})
	if err != nil {
		t.Fatal(err)
	}
	if !verifyNIPRes.Verified {
		t.Fatal("NIP wasn't verified.")
	}

	// verify random challenge proof

	rndChallengeRes, err := verifier.GetRndChallenge(ctx, &pcrpc.GetRndChallengeRequest{D: d})
	if err != nil {
		t.Fatal(err)
	}

	proofRes, err := prover.GetProof(ctx, &pcrpc.GetProofRequest{C: rndChallengeRes.C})
	if err != nil {
		t.Fatal(err)
	}

	verifyProofRes, err := verifier.VerifyProof(ctx, &pcrpc.VerifyProofRequest{D: d, P: proofRes.Proof, C: rndChallengeRes.C})
	if err != nil {
		t.Fatal(err)
	}
	if !verifyProofRes.Verified {
		t.Fatal("random challenge proof wasn't verified.")
	}
}
