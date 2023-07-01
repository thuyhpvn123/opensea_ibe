//
//  MerkleTrie.hpp
//  MetaNode_iOS_Miner
//
//  Created by Nemesis on 24/06/2022.
//

#ifndef merkle_trie_hpp
#define merkle_trie_hpp

#include <stdio.h>
#include <iostream>
#include <fstream>
#include <stdint.h>
#include <deque>
#include <sys/stat.h>
#include <sys/types.h>
#include "bigint.h"

using namespace std;
using namespace intx;

typedef void (*HashFunctionType)(const unsigned char* input,unsigned int inputByteLen,unsigned char* output);

enum nodeType
{
    nodeType_Internal = 0,
    nodeType_Leaf = 1,
    nodeType_NULL = 127,
};

struct membuf : std::streambuf
{
    membuf(char* begin, char* end) {
        this->setg(begin, begin, end);
    }
};

class TrieNode
{
public:
    unsigned char nodeType;         // Internal of Leaf
    short shortcutLength;           // Number of bits in shortcut
    TrieNode* parent;               // Pointer to Parent Node
    TrieNode** children;            // Up to 16 Children from 0x0 ... 0xF (0 -> 15 in DEC) wrap to 4 bits
    uint8_t* shortcut;              // Shortcut path in bytes
    uint8_t* hash;                  // Hash 32 bytes
    uint256_t* pValue;              // Pointer to Value.
    
    TrieNode();
    TrieNode(unsigned char type, TrieNode* parent);
    void makeChildren();
    void makeZeroHashIfNil(int len);
    
    uint16_t childrenMask();

    string generateDebugString(int layler, int hashLength);
    
    ~TrieNode();
};

uint8_t getBit(const uint8_t* data, const long& bitIndex); // bitIndex is 0,1,...,6,7
void setBit(uint8_t* data, const long& bitIndex, uint8_t value); // bitIndex is 0,1,...,6,7
short setBit16Rev(short data, const long& bitIndex, const short bitVal); // Reverse bitIndex 15,14,...,2,1,0
uint16_t getBit16(const uint16_t& src, const long& bitIndex);
/*
uint8_t setBit8(uint8_t data, const long& bitIndex, const uint8_t bitVal);
short setBit16(short data, const long& bitIndex, const short bitVal);
unsigned short setBit16u(unsigned short data, const long& bitIndex, const unsigned short bitVal);
*/

short get4BitIndex(const uint8_t* hashKey, const short& bitIndex);
void copyBits(const uint8_t* src, uint8_t* des, const long& bitIndex, const long& length, const long& desSt);
long match4Bits(const uint8_t* src,const uint8_t* des, const long& bitIndex, const long& max); // des start with 0

class MerkleTrie
{
public:
    TrieNode* root;
    int32_t hashLength;             // General hash length in bytes, Default is 32.
    int32_t dataLength;             // General hash length in bytes, Default is 32.
    bool duplicatedValue;           // Duplicate value of each node if it brings into the trie
    
    HashFunctionType hashFunct;     // Pointer to Hash Function.
    
    uint16_t intVersion;
    
    static TrieNode* rootNodeFromIStream(istream& ifs, const int& hashlength, uint16_t* pLoadedVer, uint32_t* pTotalNode);
    static TrieNode* rootNodeFromFile(const char* path, const int& hashlength, uint16_t* pLoadedVer, uint32_t* pTotalNode);
    
    MerkleTrie();
    MerkleTrie(int hashlength, bool dupval, HashFunctionType hfunct);
    MerkleTrie(const char* path, int hashlength, HashFunctionType hfunct);
    MerkleTrie(istream& ifs, int hashlength, HashFunctionType hfunct);
    
    void reHash(TrieNode* node, bool toRoot);
    
    TrieNode* makeLeaf(const uint8_t* hashKey, const uint8_t* hashValue, const uint256_t* pValue, const short& pathInd, const short& pathLen, TrieNode* parent);
    void overwriteLeaf(TrieNode* node, const uint8_t* hashKey, const uint8_t* hashValue, const uint256_t* pValue, const short& pathInd, const short& pathLen);  // pathLen < 0 means Unchanged
    
    void set(const uint8_t* hashKey, const uint8_t* hashValue, const uint256_t* pValue);
    void set(const uint256_t* key, const uint256_t* value);
    
    TrieNode* get(const uint8_t* hashKey);
    TrieNode* get(const uint256_t& key);
    
    bool remove(const uint8_t* hashKey);
    bool remove(const uint256_t& key);
    
    void clear();
    
    vector<TrieNode*> allPairs();
    
    
    //Ordinary write file
    int writeToOStream(ostream& ofs);
    int writeToFile(const char* path);
    
    //FileWrite
    int nameIndex(const char* path, int len);
    string subPartPath(const string& path, int auxCode, int splitLayer, int nodePerFile, int part);
    int writeToMultiFiles(const char* path, int splitLayer, int nodePerFile); // DO NOT CALL THIS DIRECTLY
    
    int readFromIStream(istream& ifs);
    int readFromFile(const char* path);
    
    vector<TrieNode*> getNodesAtlayer(int layer, bool inclNULL);
    uint32_t getTotalNode();
    uint32_t getTotalLeaf();
    uint32_t estMaxRAMSize();   // Return Maximum Estimated Memory (on RAM) in Bytes
    
    unsigned long Debug_overwriteHashTotal;
    
    string generateDebugString(string split);
    
    ~MerkleTrie();
    
private:
    uint32_t totalNode;
    uint32_t totalLeaf;
    
    int writeRootToOStream(ostream& ofs, TrieNode* rootNode);
    TrieNode* readRootFromIStream(istream& ifs, uint16_t loadedVersion);
    
    int writePartOStream(ostream& ofs, const TrieNode** nodes, int totalNode);
    int writePartFiles(const char* path, TrieNode** nodes, int totalNode, int splitLayer, int nodePerFile, int part);
    
    void freeData();
    void removeNode(TrieNode* node);
};

#endif /* MerkleTrie_hpp */
