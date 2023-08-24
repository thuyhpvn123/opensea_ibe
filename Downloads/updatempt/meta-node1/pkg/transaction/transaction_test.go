package transaction

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestUnmarshal(t *testing.T) {
	bData := common.FromHex("0a202de047c2ac64543d00a1e4e5d3697789267fb8b411faf8514ce05181b0d8d33112200ad8ae03e75608e16ae3050e4acff7bc03a0bc07709b9207cd3ab72e5b9ad61e1a308dbcf904ac4379a75bc8be85b17c4310649e415a12f6b826238313eab23609b82fcbc32267efd9068913676c476d85d6221424846ec0311bd180cd08cb9db50a15f36b8479602a061787188e426432016438a09c01408094ebdc036a20e232daf168bcc8aa4038cef6b70682f24e8a1df8e221c6693d9ef1ee697aeebd7220ca2d8a920ba5f7326fac626766ef35c092e592ec4945bec1e702ba8e2862625e7a60963fc0e99fa4f2c36527737fb2351a99ec291a813763b59b6a4c97bbb5bdd2f1bf22a8efcb6c9009b59d89529b73052d145c6a8c24eeb67011f22271d5a2d72a0cef4ce872ccd63dfef0316d4c260bc0a205c94d81706bbec5b3c20c1c0f16cd")
	transaction := &Transaction{}
	transaction.Unmarshal(bData)
	fmt.Printf("Transaction %v", transaction)
}