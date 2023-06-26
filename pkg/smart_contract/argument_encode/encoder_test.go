package argument_encode

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestEncodeSingleString(t *testing.T) {
	input := `{"id":11,"title":"perfume Oil 12 23232 23 23 23 232 23 "}`
	encoded := EncodeSingleString(input)
	fmt.Println(hex.EncodeToString(encoded))
}
