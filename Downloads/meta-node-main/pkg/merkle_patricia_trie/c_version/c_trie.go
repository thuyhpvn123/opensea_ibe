package c_version

/*
#cgo CFLAGS: -w
#cgo CXXFLAGS: -std=c++17 -w
#cgo LDFLAGS: -L./c_merkle_patricia_trie/build/ -lc_merkle_patricia_trie -lstdc++
#cgo CPPFLAGS: -I./c_merkle_patricia_trie/3rdparty/ -I./c_merkle_patricia_trie/3rdparty/intx/include/ -I./c_merkle_patricia_trie/include -I/usr/include
#include "merkle_trie_linker.hpp"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
)

func GetStorageRoot(storage map[string][]byte) common.Hash {
	gBStorages := []byte{}
	lenStorage := 0
	for k, v := range storage {
		b := common.FromHex(k)
		b = append(b, v[:]...)
		gBStorages = append(gBStorages, b...)
		lenStorage++
	}
	cBStorage := (*C.uchar)(C.CBytes(gBStorages))

	cBHash := C.GetRootHash(
		cBStorage,
		(C.int)(lenStorage),
	)
	defer C.free(unsafe.Pointer(cBStorage))
	defer C.free(unsafe.Pointer(cBHash))
	bHash := C.GoBytes(unsafe.Pointer(cBHash), (C.int)(32))
	return common.BytesToHash(bHash)
}
