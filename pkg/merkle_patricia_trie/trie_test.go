package merkle_patricia_trie

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/storage"
	"google.golang.org/protobuf/proto"
)

func initEmptyTrie() *Trie {
	return New(&FullNode{
		Value: nil,
		flag:  NewFlag(),
	})
}

func initTrie() *Trie {
	root := &ValueNode{
		Value: []byte("123"),
		flag: Flag{
			cacheHash: []byte("456"),
			dirty:     false,
		},
	}
	return New(root)
}

func TestHashRoot(t *testing.T) {
	trie := initTrie()
	root, hash, err := trie.HashRoot()
	assert.Nil(t, err)
	assert.NotNil(t, hash)

	expected := fmt.Sprintf("root %v\n hash %v\n err %v\n", root, hash, err)
	assert.Equal(t, expected, "root └ VALUENODE hash:0x0000000000000000000000000000000000000000000000000000000000343536-dirty:false-value:313233\n\n hash 0x0000000000000000000000000000000000000000000000000000000000343536\n err <nil>\n")
	trie2 := &Trie{}
	root, hash, err = trie2.HashRoot()

	assert.NotNil(t, err)
	assert.Equal(t, hash, common.Hash{})
	assert.Nil(t, root)
}

func TestConcat(t *testing.T) {
	testz := concat(common.Hex2Bytes("f1"), common.Hex2Bytes("f1")...)
	assert.Equal(t, testz, []byte{241, 241})
}

func TestTrieSetGet(t *testing.T) {
	trie := initEmptyTrie()
	key := common.FromHex("0")
	value := common.FromHex("f")

	err := trie.Set(key, value)
	assert.Nil(t, err)

	getValue, err := trie.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, getValue, value)

	root, hash, err := trie.HashRoot()
	fmt.Printf("root %v\n hash %v\n err %v\n test %v\n", root, hash, err, crypto.Keccak256Hash(nil))
}

func TestTrieDelete(t *testing.T) {
	trie := initEmptyTrie()
	befDel := fmt.Sprintf("%v", trie)
	fmt.Println(befDel)

	key := common.FromHex("f1f1f1")
	value := common.FromHex("f2f2f2f2")

	err := trie.Set(key, value)
	assert.Nil(t, err)

	fmt.Println(fmt.Sprintf("%v", trie))
	err2 := trie.Delete(key)
	assert.Nil(t, err2)
	aftDel := fmt.Sprintf("%v", trie)
	fmt.Println(aftDel)
	assert.Equal(t, befDel, aftDel)

	trie.Set(common.FromHex("a72"), value)
	trie.Set(common.FromHex("a7"), value)
	trie.Set(common.FromHex("a71"), value)
	trie.Set(common.FromHex("a71"), value)
	trie.Set(common.FromHex("a1"), value)
	trie.Set(common.FromHex("a1a11bf"), value)

	err = trie.Delete(common.FromHex("a72"))
	assert.Nil(t, err)
	err = trie.Delete(common.FromHex("a1a11bf"))
	assert.Nil(t, err)
	err = trie.Delete(common.FromHex("a71"))
	assert.Nil(t, err)
	err = trie.Delete(common.FromHex("a7"))
	assert.Nil(t, err)
	err = trie.Delete(common.FromHex("a71"))
	assert.Nil(t, err)
	err = trie.Delete(common.FromHex("a1"))
	assert.Nil(t, err)
}

func TestGetTrie(t *testing.T) {
	trie := initEmptyTrie()
	key := common.FromHex("0")
	value := common.FromHex("f")
	trie.Set(key, value)

	data, _ := trie.Get(key)
	assert.Equal(t, value, data)
}

func TestCommitTrie(t *testing.T) {
	store := storage.NewMemoryDb()
	trie := New(&FullNode{
		Value: nil,
		flag:  NewFlag(),
	})

	key := common.FromHex("123456")
	value := common.FromHex("A")
	trie.Set(key, value)
	trie.Commit(store)
	hashedKey := crypto.Keccak256(value)
	v, ok := store.Get(hashedKey)
	protoNode := &pb.MPTNode{}
	err := proto.Unmarshal(v, protoNode)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, common.Bytes2Hex(protoNode.Data), "0a")
	assert.Nil(t, ok)

	err1 := trie.Set([]byte("134567"), []byte("B"))
	if err1 != nil {
		t.Error(err1.Error())
	}
	trie.Commit(store)
	hashedKey = crypto.Keccak256([]byte("B"))
	v1, ok1 := store.Get(hashedKey)
	protoNode1 := &pb.MPTNode{}
	err = proto.Unmarshal(v1, protoNode1)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, string(protoNode1.Data), "B")
	assert.Nil(t, ok1)

}
func TestFetchNodeFromStorage(t *testing.T) {
	store := storage.NewMemoryDb()
	valueNode := &ValueNode{
		Value: []byte("123"),
		flag: Flag{
			cacheHash: []byte("456"),
			dirty:     false,
		},
	}
	shortNode := ShortNode{
		Key:   []byte("789"),
		Value: valueNode,
		flag: Flag{
			cacheHash: []byte("456"),
			dirty:     true,
		},
	}
	fullNode := FullNode{}
	fullNode.Children[0] = &shortNode
	fullNode.Value = []byte("123")
	fullNode.flag.dirty = true

	trie := New(&fullNode)

	trie.Commit(store)
	_, err := FetchNodeFromStorage(fullNode.Hash(), store)
	assert.Nil(t, err)

	value := common.FromHex("A")
	err = trie.Set(common.FromHex("a71"), value)
	assert.Nil(t, err)
}

func TestCopyTrie(t *testing.T) {
	key1 := []byte("111")
	value1 := []byte("v1")
	key2 := []byte("222")
	value2 := []byte("v2")
	value3 := []byte("v3")

	trie := New(NewEmtyFullNode())
	trie.Set(key1, value1)
	_, oldRoot, _ := trie.HashRoot()
	trie2 := trie.Copy()
	trie2.Set(key2, value2)
	t1v1, _ := trie.Get(key1)
	t2v1, _ := trie2.Get(key1)
	_, newRoot, _ := trie.HashRoot()

	t1v2, _ := trie.Get(key2)
	t2v2, _ := trie2.Get(key2)
	assert.Equal(t, t1v1, t2v1)
	assert.NotEqual(t, t1v2, t2v2)
	assert.Equal(t, oldRoot, newRoot)

	trie2.Set(key1, value3)
	t2v3, _ := trie2.Get(key1)
	assert.NotEqual(t, t1v1, t2v3)
}
