package poet_core_api

import (
	"context"
	"fmt"
	"github.com/spacemeshos/poet-core-api/pcrpc"
	"testing"
)

type testCase struct {
	x []byte
	n uint32
	h string
}

func genTestCases(x []byte, nMin int, nMax int, hFuncs []string) []*testCase {
	var testCases []*testCase
	for i := nMin; i <= nMax; i++ {
		for _, hFunc := range hFuncs {
			testCases = append(testCases, &testCase{x: x, n: uint32(i), h: hFunc})
		}
	}
	return testCases
}

func TestPoetCoreService(t *testing.T) {
	prover, cleanUp := NewProverClient(DefaultRPCHostPort)
	defer cleanUp()
	verifier, cleanUp := NewVerifierClient(DefaultRPCHostPort)
	defer cleanUp()

	testCases := genTestCases([]byte("this is a commitment"), 1, 5, []string{"sha256", "scrypt"})
	for _, testCase := range testCases {
		testCaseStr := fmt.Sprintf("x:%q n:%d h:%s", testCase.x, testCase.n, testCase.h)

		success := t.Run(testCaseStr, func(t1 *testing.T) {
			testProverAndVerifier(t, testCase, prover, verifier)
		})

		if !success {
			break
		}
	}
}

func testProverAndVerifier(t *testing.T, tc *testCase, prover pcrpc.PoetCoreProverClient, verifier pcrpc.PoetVerifierClient) {
	ctx := context.Background()
	d := &pcrpc.DagParams{X: tc.x, N: tc.n, H: tc.h}

	_, err := prover.Compute(ctx, &pcrpc.ComputeRequest{D: d})
	defer prover.Clean(ctx, &pcrpc.CleanRequest{})
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
