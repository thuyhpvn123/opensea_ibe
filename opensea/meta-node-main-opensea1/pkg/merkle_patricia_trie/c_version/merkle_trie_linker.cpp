#include "merkle_trie_linker.hpp"
#include "merkle_trie.h"
#include "util.h"

char* GetRootHash(
    unsigned char* storage_p,
    int storage_length
){
    MerkleTrie* s = new MerkleTrie(32,true,c_merkle_trie::keccak_256);
    for (int i = 0; i < storage_length; i++) {
        uint256_t key = c_merkle_trie::from_big_endian((uint8_t *)storage_p+(i*64));
        uint256_t value = c_merkle_trie::from_big_endian((uint8_t *)storage_p+(i*64)+32);
        s->set(&key, &value);
    }
    char* rs = new char[32];
    std::memcpy(rs, s->root->hash, 32);
    delete s;
    return rs;
};