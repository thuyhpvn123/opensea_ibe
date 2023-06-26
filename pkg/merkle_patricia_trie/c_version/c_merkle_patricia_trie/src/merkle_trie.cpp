//
//  merkle_trie_hpp.cpp
//  MetaNode_iOS_Miner
//
//  Created by Nemesis on 24/06/2022.
//

#include "merkle_trie.h"

#pragma mark - TrieNode

#define MAX_TRIE_CHILDREN 16
#define DEFAULT_HASH_LENGTH 32

TrieNode::TrieNode()
{
    nodeType = nodeType_Internal;
    shortcutLength = 0;
    parent = NULL;
    children = NULL;
    shortcut = NULL;
    hash = NULL;
    pValue = NULL;
}

TrieNode::TrieNode(unsigned char type, TrieNode* pParent)
{
    nodeType = type;
    shortcutLength = 0;
    parent = pParent;
    children = NULL;
    shortcut = NULL;
    hash = NULL;
    pValue = NULL;
}

void TrieNode::makeChildren()
{
    if (children) return;
    children = new TrieNode*[MAX_TRIE_CHILDREN];
    for (int i=0; i<MAX_TRIE_CHILDREN; i++) children[i] = NULL;
}

void TrieNode::makeZeroHashIfNil(int len)
{
    if (hash) return;
    if (len < 0) len = 1;
    hash = new uint8_t[len + 1];
    for (int i=0; i<len; i++) hash[i] = 0;
}

uint16_t TrieNode::childrenMask()
{
    if (!children) return 0;
    uint16_t res = 0;
    for (int i=0; i<MAX_TRIE_CHILDREN; i++)
    {
        if (!children[i]) continue;
        res = res | (1 << i);
    }
    return res;
}

string TrieNode::generateDebugString(int layler, int hashLength)
{
    stringstream buffer;
    buffer << "[";
    buffer << "L" << layler << "-";
    buffer << "T" << ((int)nodeType) << "-";
    buffer << "S" << shortcutLength;
    if (shortcutLength > 0) {
        buffer << "x";
        for (int i=0; i<shortcutLength; i+=4)
            buffer << std::hex << std::uppercase << get4BitIndex(shortcut, i);
    }
    buffer << "-";
    if (parent && parent->hash)
    {
        buffer << "Px";
        for (int i=hashLength-2; i<hashLength; i++)
            buffer << std::setw(2) << std::setfill('0') << std::hex << std::uppercase << ((int)parent->hash[i]);
        buffer << "-";
    }
    else
    {
        buffer << "Px-";
    }
    if (children) {
        buffer << "Cx";
        for (int i=0; i<MAX_TRIE_CHILDREN; i++)
            if (children[i]) buffer << std::hex << std::uppercase << i;
        buffer << "-";
    } else {
        buffer << "Cx-";
    }
    if (hash) {
        buffer << "Hx";
        for (int i=0; i<hashLength; i++)
            buffer << std::setw(2) << std::setfill('0') << std::hex << std::uppercase << ((int)hash[i]);
        buffer << "-";
    } else {
        buffer << "Hx-";
    }
    if (pValue) {
        buffer << to_string(*pValue);
    } else {
        buffer << "Vx";
    }
    buffer << "]";
    return buffer.str();
}

TrieNode::~TrieNode()
{
    if (children) delete[] children;
    if (shortcut) delete[] shortcut;
    if (hash) delete[] hash;
    //pValue is managed by MerkleTrie ! No delete!
}

#pragma mark - MerkleTrie

uint8_t getBit(const uint8_t* data, const long& bitIndex) // bitIndex is 0,1,...,6,7
{
    return ((data[bitIndex / 8]) >> (7 - (bitIndex % 8))) & 1;
}

void setBit(uint8_t* data, const long& bitIndex, uint8_t value)
{
    long byteind = bitIndex / 8;
    short bitinbyte = bitIndex % 8;
    if (!value)
        data[byteind] = data[byteind] & (~(1 << (7 - bitinbyte)));
    else
        data[byteind] = data[byteind] | (1 << (7 - bitinbyte));
}

short setBit16Rev(short data, const long& bitIndex, const short bitVal) // Reverse bitIndex 15,14,...,2,1,0
{
    return data | (bitVal << bitIndex);
}

uint16_t getBit16(const uint16_t& src, const long& bitIndex)
{
    return (src >> bitIndex) & 1;
}

/*
uint8_t setBit8(uint8_t data, const long& bitIndex, const uint8_t bitVal)
{
    return data | (bitVal << (8 - bitIndex));
}

short setBit16(short data, const long& bitIndex, const short bitVal)
{
    return data | (bitVal << (16 - bitIndex));
}

unsigned short setBit16u(unsigned short data, const long& bitIndex, const unsigned short bitVal)
{
    return data | (bitVal << (16 - bitIndex));
}
*/

short get4BitIndex(const uint8_t* hashKey, const short& bitIndex)
{
    short bitInByteInd = bitIndex % 8;
    if (bitInByteInd <= 4)
    {
        return ((hashKey[bitIndex / 8] << bitInByteInd) & 0xF0) >> 4;
    }
    else
    {
        //return ((hashKey[byteindex] << bitInByteInd) & 0xFF) | ((hashKey[byteindex+1] >> (8-bitInByteInd)) & 0xFF);
        short res = 0;
        res = setBit16Rev(res, 3, getBit(hashKey, bitIndex));
        res = setBit16Rev(res, 2, getBit(hashKey, bitIndex+1));
        res = setBit16Rev(res, 1, getBit(hashKey, bitIndex+2));
        res = setBit16Rev(res, 0, getBit(hashKey, bitIndex+3));
        return res;
    }
}

void copyBits(const uint8_t* src, uint8_t* des, const long& bitIndex, const long& length, const long& desSt)
{
    if (bitIndex % 4 == 0 && desSt % 4 == 0)
    {
        long len = bitIndex + length;
        long desInd = desSt, bibsrc, bibdes;
        for (long i=bitIndex; i<len; i+=4)
        {
            bibsrc = i % 8; bibdes = desInd % 8;
            if (bibdes == 0) des[desInd / 8] = 0;
            if (bibdes == 0 && bibsrc == 0)
            {
                des[desInd / 8] = (des[desInd / 8] | (src[i / 8] & 0xF0));
            }
            else if (bibdes == 0 && bibsrc == 4)
            {
                des[desInd / 8] = (des[desInd / 8] | ((src[i / 8] << 4) & 0xF0));
            }
            else if (bibdes == 4 && bibsrc == 0)
            {
                des[desInd / 8] = (des[desInd / 8] | ((src[i / 8] >> 4) & 0x0F));
            }
            else if (bibdes == 4 && bibsrc == 4)
            {
                des[desInd / 8] = (des[desInd / 8] | (src[i / 8] & 0x0F));
            }
            desInd+=4;
        }
    }
    else
    {
        long len = bitIndex + length;
        long desi = 0;
        for (long i=bitIndex; i<len; i++)
        {
            setBit(des, desi, getBit(src, i));
            desi++;
        }
    }
}

long match4Bits(const uint8_t* src,const uint8_t* des, const long& bitIndex, const long& max) // des start with 0
{
    if (bitIndex % 4 == 0 && max % 4 == 0)
    {
        long desInd = 0, bibsrc, bibdes;
        uint8_t desVal = 0; uint8_t srcVal = 0;
        for (long i=bitIndex; i<max; i+=4)
        {
            bibsrc = i % 8; bibdes = desInd % 8;
            if (bibdes == 0) desVal = des[desInd / 8] >> 4;
            else if (bibdes == 4) desVal = des[desInd / 8] & 0b00001111;
            if (bibsrc == 0) srcVal = src[i / 8] >> 4;
            else if (bibsrc == 4) srcVal = src[i / 8] & 0b00001111;
            if (desVal != srcVal) return i - bitIndex;
            desInd+=4;
        }
        return max - bitIndex;
    }
    else
    {
        long desi = 0;
        uint8_t desVal = 0; uint8_t srcVal = 0;
        for (long i=bitIndex; i<max; i++)
        {
            srcVal = getBit(src, i);
            desVal = getBit(des,desi);
            if (srcVal != desVal) return (i - bitIndex) / 4 * 4;;
            desi++;
        }
        return (max - bitIndex) / 4 * 4;;
    }
}

MerkleTrie::MerkleTrie()
{
    root = NULL;
    hashLength = 32;
    dataLength = 32;
    totalNode = 0;
    totalLeaf = 0;
    duplicatedValue = false;
    hashFunct = NULL;
    intVersion = 1;
    Debug_overwriteHashTotal = 0;
    root = new TrieNode(); totalNode++;
    root->makeZeroHashIfNil(hashLength);
    root->makeChildren();
}

MerkleTrie::MerkleTrie(int hashlength, bool dupval, HashFunctionType hfunct)
{
    root = NULL;
    hashLength = hashlength;
    dataLength = 32;
    totalNode = 0;
    totalLeaf = 0;
    duplicatedValue = dupval;
    hashFunct = hfunct;
    intVersion = 1;
    Debug_overwriteHashTotal = 0;
    root = new TrieNode(); totalNode++;
    root->makeZeroHashIfNil(hashLength);
    root->makeChildren();
}

MerkleTrie::MerkleTrie(const char* path, int hashlength, HashFunctionType hfunct)
{
    root = NULL;
    hashLength = hashlength;
    dataLength = 32;
    totalNode = 0;
    totalLeaf = 0;
    duplicatedValue = true;
    hashFunct = hfunct;
    intVersion = 1;
    Debug_overwriteHashTotal = 0;
    readFromFile(path);
}

MerkleTrie::MerkleTrie(istream& ifs, int hashlength, HashFunctionType hfunct)
{
    root = NULL;
    hashLength = hashlength;
    dataLength = 32;
    totalNode = 0;
    totalLeaf = 0;
    duplicatedValue = true;
    hashFunct = hfunct;
    intVersion = 1;
    Debug_overwriteHashTotal = 0;
    readFromIStream(ifs);
}

void MerkleTrie::reHash(TrieNode* node, bool toRoot)
{
    if (!hashFunct) return;
    int stepLength = hashLength*2; // 32 bytes will be splited to 64 chunk of 4 bits
    for (int step=0; step<stepLength; step++)
    {
        if (node->nodeType == nodeType_Leaf && !node->hash)
        {
            node->hash = new uint8_t[hashLength + 1];
            hashFunct((uint8_t*)((void*)node->pValue),dataLength,node->hash);
            Debug_overwriteHashTotal += dataLength;
        }
        else if (node->nodeType == nodeType_Internal)
        {
            int node_count = 0, dataind = 0, byteind = 0, bitind = 0;
            for (int i=0; i<MAX_TRIE_CHILDREN; i++)
                if (node->children[i] && node->children[i]->hash) node_count++;
            uint8_t* data = new uint8_t[hashLength*node_count+5];
            data[0] = data[1] = 0; dataind = MAX_TRIE_CHILDREN / 8;
            for (int i=0; i<MAX_TRIE_CHILDREN; i++)
            {
                if (!node->children[i] || !node->children[i]->hash) continue;
                byteind = i / 8; bitind = i % 8;
                data[byteind] = data[byteind] | (1 << bitind);
                memcpy(data+dataind, node->children[i]->hash, hashLength);
                dataind+=hashLength;
            }
            if (!node->hash) node->hash = new uint8_t[hashLength + 1];
            hashFunct(data, hashLength*node_count + 2, node->hash);
            Debug_overwriteHashTotal += hashLength*node_count;
            delete[] data;
        }
        if (!toRoot) break;
        if (!node->parent) break;
        node = node->parent;
    }
}

TrieNode* MerkleTrie::makeLeaf(const uint8_t* hashKey, const uint8_t* hashValue, const uint256_t* pValue, const short& pathInd, const short& pathLen, TrieNode* parent)
{
    TrieNode* leaf = new TrieNode(nodeType_Leaf, parent);
    totalNode++; totalLeaf++;
    //Copy Path
    if (pathLen > 4)
    {
        leaf->shortcutLength = pathLen;
        leaf->shortcut = new uint8_t[pathLen/8 + 2];
        copyBits(hashKey, leaf->shortcut, pathInd, pathLen, 0);
    }
    // Assign Val
    if (duplicatedValue)
    {
        leaf->pValue = new uint256_t();
        *(leaf->pValue) = *pValue;
    }
    else
    {
        void* pointer = (void*)pValue;
        leaf->pValue = ((uint256_t*)pointer);
    }
    //Duplicate hash val
    leaf->hash = new uint8_t[hashLength + 1];
    memcpy(leaf->hash, hashValue, hashLength);
    return leaf;
}

void MerkleTrie::overwriteLeaf(TrieNode* node, const uint8_t* hashKey, const uint8_t* hashValue, const uint256_t* pValue, const short& pathInd, const short& pathLen) // pathLen < 0 means Unchanged
{
    TrieNode* leaf = node;
    //Copy Path
    if (pathLen > 4)
    {
        if (leaf->shortcut) delete[] leaf->shortcut;
        leaf->shortcutLength = pathLen;
        leaf->shortcut = new uint8_t[pathLen/8 + 2];
        copyBits(hashKey, leaf->shortcut, pathInd, pathLen, 0);
    }
    else if (pathLen >= 0)
    {
        if (leaf->shortcut) delete[] leaf->shortcut;
        leaf->shortcutLength = pathLen; leaf->shortcut = NULL;
    }
    //pathLen < 0 is unchanged
    // Assign Val
    if (duplicatedValue)
    {
        if (leaf->pValue) delete leaf->pValue;
        leaf->pValue = new uint256_t();
        *(leaf->pValue) = *pValue;
    }
    else
    {
        void* pointer = (void*)pValue;
        leaf->pValue = ((uint256_t*)pointer);
    }
    //Duplicate hash val
    if (!leaf->hash) leaf->hash = new uint8_t[hashLength + 1];
    memcpy(leaf->hash, hashValue, hashLength);
}

void MerkleTrie::set(const uint8_t* hashKey, const uint8_t* hashValue, const uint256_t* pValue)
{
    short pathInd = 0;
    TrieNode* curNode = root;
    int stepLength = hashLength*2; // 32 bytes will be splited to 64 chunk of 4 bits
    for (int step=0; step<stepLength; step++)
    {
        short ind4Bit = get4BitIndex(hashKey, pathInd);
        assert(0 <= ind4Bit && ind4Bit < 16);
        if (!curNode->children)
        {   // Would be Root
            curNode->makeChildren();
            TrieNode* leaf = this->makeLeaf(hashKey, hashValue, pValue, pathInd, 8*hashLength - pathInd, curNode);
            curNode->children[ind4Bit] = leaf;
            reHash(leaf, true);
            break;
        }
        else if (!curNode->children[ind4Bit])
        {   // Not found Path -> New path
            TrieNode* leaf = this->makeLeaf(hashKey, hashValue, pValue, pathInd, 8*hashLength - pathInd, curNode);
            curNode->children[ind4Bit] = leaf;
            reHash(leaf, true);
            break;
        }
        else if (curNode->children[ind4Bit]->nodeType == nodeType_Internal && !curNode->children[ind4Bit]->shortcut)
        {   // No shortcut -> move 4 bits
            curNode = curNode->children[ind4Bit];
            pathInd+=4;
        }
        else // Shortcut Node
        {
            long max = min((long)(8*hashLength),(long)(pathInd + curNode->children[ind4Bit]->shortcutLength));
            long match = match4Bits(hashKey, curNode->children[ind4Bit]->shortcut, pathInd, max);
            if (match >= curNode->children[ind4Bit]->shortcutLength) // Completely matched
            {
                if (curNode->children[ind4Bit]->nodeType == nodeType_Leaf)
                {   // Override Value
                    this->overwriteLeaf(curNode->children[ind4Bit], hashKey, hashValue, pValue, pathInd, -1);
                    reHash(curNode->children[ind4Bit], true);
                    break;
                }
                else
                {   // Move match bit
                    curNode = curNode->children[ind4Bit];
                    pathInd+=match;
                }
            }
            else
            {   // Split Branch -> insert a internal Node
                TrieNode* child = new TrieNode(nodeType_Internal,curNode); totalNode++;
                child->makeChildren();
                TrieNode* grandchild = curNode->children[ind4Bit];
                //copy Data
                child->shortcutLength = match;
                child->shortcut = new uint8_t[match / 8 + 2];
                copyBits(hashKey, child->shortcut, pathInd, match, 0);
                child->makeChildren();
                uint8_t* oldshortcut = grandchild->shortcut;
                grandchild->shortcutLength = grandchild->shortcutLength-match;
                grandchild->shortcut = new uint8_t[grandchild->shortcutLength / 8 + 2];
                copyBits(oldshortcut, grandchild->shortcut, match, grandchild->shortcutLength, 0);
                if (oldshortcut) delete[] oldshortcut;
                //connect node
                short gchildInd = get4BitIndex(grandchild->shortcut, 0);
                assert(0 <= gchildInd && gchildInd < 16);
                curNode->children[ind4Bit] = child;
                child->children[gchildInd] = grandchild;
                grandchild->parent = child;
                curNode = child;
                pathInd+=match;
                //copy hash from grandchild to child
                child->hash = new uint8_t[hashLength + 1];
                memcpy(child->hash, grandchild->hash, hashLength);
            }
        }
    }
}

void MerkleTrie::set(const uint256_t* key, const uint256_t* value)
{
    uint8_t keybyte[32 + 1]; uint8_t valbyte[32 + 1];
    uint8_t keyHash[32 + 1]; uint8_t valHash[32 + 1];
    intx::be::unsafe::store((uint8_t*)keybyte, *key);
    intx::be::unsafe::store((uint8_t*)valbyte, *value);
    if (hashFunct)
    {
        hashFunct(keybyte, hashLength, (uint8_t*)keyHash);
        hashFunct(valbyte, hashLength, (uint8_t*)valHash);
        set(keyHash, valHash, value);
    }
    else
    {
        set((uint8_t*)keybyte, (uint8_t*)valbyte, value);
    }
}

TrieNode* MerkleTrie::get(const uint8_t* hashKey)
{
    TrieNode* res = NULL;
    short pathInd = 0;
    TrieNode* curNode = root;
    int stepLength = hashLength*2; // 32 bytes will be splited to 64 chunk of 4 bits
    for (int step=0; step<stepLength; step++)
    {
        short ind4Bit = get4BitIndex(hashKey, pathInd);
        assert(0 <= ind4Bit && ind4Bit < 16);
        if (!curNode->children)
        {   // Would be Root (Empty) -> return
            break;
        }
        else if (!curNode->children[ind4Bit])
        {   // Not found Path -> return
            break;
        }
        else if (curNode->children[ind4Bit]->nodeType == nodeType_Internal && !curNode->children[ind4Bit]->shortcut)
        {   // No shortcut -> move 4 bits
            curNode = curNode->children[ind4Bit];
            pathInd+=4;
        }
        else // Shortcut Node
        {
            long max = min((long)(8*hashLength),(long)(pathInd + curNode->children[ind4Bit]->shortcutLength));
            long match = match4Bits(hashKey, curNode->children[ind4Bit]->shortcut, pathInd, max);
            if (match >= curNode->children[ind4Bit]->shortcutLength) // Completely matched
            {
                if (curNode->children[ind4Bit]->nodeType == nodeType_Leaf)
                {   // finish path -> return node
                    return curNode->children[ind4Bit];
                }
                else
                {   // Move match bit
                    curNode = curNode->children[ind4Bit];
                    pathInd+=match;
                }
            }
            else
            {
                break;
            }
        }
    }
    return res;
}

TrieNode* MerkleTrie::get(const uint256_t& key)
{
    uint8_t keybyte[32 + 1]; uint8_t keyHash[32 + 1];
    intx::be::unsafe::store((uint8_t*)keybyte, key);
    if (hashFunct)
    {
        hashFunct(keybyte, hashLength, (uint8_t*)keyHash);
        return get(keyHash);
    }
    else
    {
        return get(keybyte);
    }
}

void MerkleTrie::removeNode(TrieNode* node)
{
    TrieNode* parent = node->parent; TrieNode* remnode = NULL;
    if (node == root && !parent) // Remove Parent -> Update hash to 0x0000...0000
    {
        for (int i=0; i<MAX_TRIE_CHILDREN; i++) // Check Zero Children
        {
            if (node->children[i]) return;
        }
        if (node->hash) delete[] node->hash;
        node->hash = NULL;
        node->makeZeroHashIfNil(hashLength);
    }
    if (!parent) return;
    int remcount = 0, ind = -1, remind = ind; // rem = Remaining Node, Index
    for (int i=0; i<MAX_TRIE_CHILDREN; i++)
    {
        if (!(parent->children[i])) continue;
        if (parent->children[i] == node)
        {
            ind = i;
            continue;
        }
        remnode = parent->children[i]; remind = i;
        remcount++;
    }
    if (0 <= ind && ind < MAX_TRIE_CHILDREN)
    {
        parent->children[ind] = NULL;
    }
    node->parent = NULL;
    if (remcount == 1 && parent->parent && parent->parent->children) // Merge Up Node
    {
        TrieNode* gparent = parent->parent;
        int prind = -1;
        for (int i=0; i<MAX_TRIE_CHILDREN; i++)
        {
            if (gparent->children[i] == parent) {
                prind = i; break;
            }
        }
        if (0 <= prind && prind < MAX_TRIE_CHILDREN)
        {
            //build shortcut
            int prlen = (parent->shortcutLength > 4 ? parent->shortcutLength : 4);
            int ndlen = (remnode->shortcutLength > 4 ? remnode->shortcutLength : 4);
            int newslen = prlen + ndlen;
            uint8_t* newShortcut = new uint8_t[newslen / 8 + 2]; newShortcut[0] = 0;
            //copy shortcut
            if (prlen < 4 || !parent->shortcut) {
                uint8_t byte = prind;
                newShortcut[0] = newShortcut[0] | (byte << 4);
            } else {
                copyBits(parent->shortcut,newShortcut,0,prlen,0);
            }
            if (ndlen < 4 || !remnode->shortcut) {
                uint8_t byte = remind;
                newShortcut[0] = newShortcut[0] | byte;
            } else {
                copyBits(remnode->shortcut,newShortcut,0,ndlen,prlen);
            }
            if (remnode->shortcut) delete[] remnode->shortcut;
            remnode->shortcut = newShortcut; remnode->shortcutLength = newslen;
            //connect node
            gparent->children[prind] = remnode; remnode->parent = gparent;
            reHash(gparent, true);
            // Finally Delete Parent
            if (parent->nodeType == nodeType_Leaf) totalLeaf--;
            delete parent;
            totalNode--;
        }
    }
    else if (remcount == 0)
    {
        removeNode(parent);
    }
    else
    {
        reHash(parent, true);
    }
    if (node->nodeType == nodeType_Leaf) totalLeaf--;
    delete node;
    totalNode--;
}


bool MerkleTrie::remove(const uint8_t* hashKey)
{
    TrieNode* node = get(hashKey);
    if (!node) return false;
    if (node->nodeType != nodeType_Leaf) return false;
    removeNode(node);
    return true;
}

bool MerkleTrie::remove(const uint256_t& key)
{
    uint8_t keybyte[32 + 1]; uint8_t keyHash[32 + 1];
    intx::be::unsafe::store((uint8_t*)keybyte, key);
    if (hashFunct)
    {
        hashFunct(keybyte, hashLength, (uint8_t*)keyHash);
        return remove(keyHash);
    }
    else
    {
        return remove(keybyte);
    }
}

void MerkleTrie::clear()
{
    freeData();
    root = new TrieNode(); totalNode++;
    root->makeZeroHashIfNil(hashLength);
    root->makeChildren();
}

vector<TrieNode*> MerkleTrie::allPairs()
{
    vector<TrieNode*> res;
    deque< TrieNode*> queue;
    if (root) queue.push_back(root);
    while(queue.size() > 0)
    {
        TrieNode* node = queue.front();
        queue.pop_front();
        if (node->children)
        {
            for (int i=0; i<MAX_TRIE_CHILDREN; i++)
            {
                if (!node->children[i]) continue;
                queue.push_back( node->children[i] );
            }
        }
        if (node->nodeType == nodeType_Leaf)
            res.push_back(node);
    }
    return res;
}

#pragma mark - Ordinary Write

int MerkleTrie::writeRootToOStream(ostream& ofs, TrieNode* rootNode)
{
    int index = 0;
    vector< TrieNode*> nodes;
    vector< pair<uint32_t,uint8_t> > prlink;
    nodes.push_back(rootNode); prlink.push_back( make_pair(0, 0));
    while(index < nodes.size())
    {
        TrieNode* node = nodes[index]; index++;
        if (node->children)
        {
            for (int i=0; i<MAX_TRIE_CHILDREN; i++)
            {
                if (!node->children[i]) continue;
                nodes.push_back( node->children[i] );
                prlink.push_back( make_pair(index-1, i) );
            }
        }
    }
    
    uint32_t total_node = (uint32_t)nodes.size();
    uint16_t childmask; TrieNode* node = NULL;
    uint32_t prindex; uint8_t prchildindex;
    
    int len = max(2*hashLength, 2*dataLength);
    uint8_t dump[len];
    for (int i=0; i<len; i++) dump[i] = 0;

    ofs.write((char *)&total_node, sizeof(total_node));
    
    for (uint32_t i=0; i<total_node; i++)
    {
        node = nodes[i];
        ofs.write((char *)&(node->nodeType), sizeof(node->nodeType));
        ofs.write((char *)&(node->shortcutLength), sizeof(node->shortcutLength));
        if (node->shortcutLength > 0)
        {
            ofs.write((char *)(node->shortcut), sizeof(uint8_t) * (node->shortcutLength / 8 + 2));
        }
        prindex = prlink[i].first; prchildindex = prlink[i].second;
        ofs.write((char *)&prindex,sizeof(prindex));
        ofs.write((char *)&prchildindex,sizeof(prchildindex));
        childmask = node->childrenMask();
        ofs.write((char *)&childmask,sizeof(childmask));
        if (node->hash) ofs.write((char *)(node->hash),sizeof(uint8_t) * hashLength);
        else ofs.write((char *)dump,sizeof(uint8_t) * hashLength);
        if (node->nodeType == nodeType_Leaf)
        {
            if (node->pValue) ofs.write((char *)(node->pValue),sizeof(uint256_t));
            else ofs.write((char *)dump,sizeof(uint8_t) * dataLength);
        }
    }
    
    return 0;
}

int MerkleTrie::writeToOStream(ostream& ofs)
{
    ofs.write("MTRIE", 5);
    ofs.write((char *)&intVersion, sizeof(intVersion));
    return writeRootToOStream(ofs, root);
}

int MerkleTrie::writeToFile(const char* path)
{
    string tmp_path = path + (string)"_TMP";
    
    ofstream ofs(tmp_path.c_str(), ios::out | ios::binary);
    
    if (!ofs)
    {
        return 1;
    }
    
    int wres = writeToOStream(ofs);
    
    if (wres != 0) return wres;
    
    ofs.close();
    
    int finalres = 0;
    
    //Delete old file
    try {
        std::remove(path);
    } catch (...) { }
    
    //Rename new file
    try {
        finalres = finalres | std::rename(tmp_path.c_str() ,path);
    } catch (...) { }
    
    return finalres;
}

#pragma mark - MultiFileWrite

int MerkleTrie::nameIndex(const char* path, int len)
{
    for (int i=len-1; i>=0; i--)
    {
        if (path[i] == '/') return i+1;
    }
    return 0;
}

string MerkleTrie::subPartPath(const string& path, int auxCode, int splitLayer, int nodePerFile, int part)
{   // Example Name P012.000.MTrie.dat
    char tmp[64];
    snprintf(tmp, 63, "P%X%X%X.%04d", auxCode, splitLayer, nodePerFile,part);
    return path + ((string)tmp);
}

int MerkleTrie::writePartOStream(ostream& ofs, const TrieNode** nodes, int totalNode)
{
    
    return 0;
}

int MerkleTrie::writePartFiles(const char* path, TrieNode** nodes, int totalNode, int splitLayer, int nodePerFile, int part)
{
    ofstream ofs(path, ios::out | ios::binary);
    if (!ofs)
    {
        return 1;
    }
    vector<TrieNode*> vtnodes;
    if (!nodes)
    {
        vtnodes = getNodesAtlayer(splitLayer, true);
        nodes = vtnodes.data(); totalNode = (int)vtnodes.size();
    }
    // Write Header Data
    
    // Write Part Data
    int res = writePartOStream(ofs, (const TrieNode**)(nodes + (part*nodePerFile)), nodePerFile);
    ofs.close();
    return res;
}

int MerkleTrie::writeToMultiFiles(const char* path, int splitLayer, int nodePerFile) // DO NOT CALL THIS DIRECTLY
{
    // Generate FileName, Path, Folder Name, Tmp File, etc...
    int pathlen = (int)strlen(path);
    int nameind = nameIndex(path, pathlen);
    string fullname = (const char*)(path+nameind);
    string parentpath = path;
    if (nameind<=0) parentpath = "";
    else parentpath = parentpath.substr(0,nameind);
    string partpath = ((string)path) + ((string)".path/");
    
    // Check & Create Folder
    try {
        mkdir(partpath.c_str(), S_IRWXU | S_IRWXG | S_IROTH | S_IXOTH);
    } catch(...) {
        
    }
    
    //Test
    /*
    string ppath1 = subPartPath(partpath, 3, 11, 9, 20);
    string ppath2 = subPartPath(partpath, 10, 13, 14, 15);
    */
    //Test
    
    // Example Name P012.000.MTrie.dat
    int totalchildnode = 1;
    for (int i=0; i<=splitLayer; i++) totalchildnode*=16;
    int totalfile = totalchildnode / nodePerFile;
    if (totalchildnode % nodePerFile != 0) totalfile++;
    vector<TrieNode*> vtnodes = getNodesAtlayer(splitLayer, true);
    for (int i=0; i<totalfile; i++)
    {
        string ppath = subPartPath(partpath, 0, splitLayer, nodePerFile, i);
        writePartFiles(ppath.c_str(), vtnodes.data(), (int)vtnodes.size(), splitLayer, nodePerFile, i);
    }
    
    
    // Write File
    // Delete old Files
    // Rename new Files
    return 0;
}

#pragma mark - Ordinary Read

TrieNode* MerkleTrie::readRootFromIStream(istream& ifs, uint16_t loadedVersion)
{
    uint32_t total_node = 0;
    ifs.read((char *)&total_node, sizeof(total_node));
    
    uint32_t prindex; uint8_t prchildindex; uint16_t childmask;
    
    /*
     unsigned char nodeType;         // Internal of Leaf
     short shortcutLength;           // Number of bits in shortcut
     TrieNode* parent;               // Pointer to Parent Node
     TrieNode** children;            // Up to 16 Children from 0x0 ... 0xF (0 -> 15 in DEC) wrap to 4 bits
     uint8_t* shortcut;              // Shortcut path in bytes
     uint8_t* hash;                  // Hash 32 bytes
     uint256_t* pValue;              // Pointer to Value.
    */
    
    vector<TrieNode*> nodes(total_node);
    for (uint32_t i=0; i<total_node; i++) {
        TrieNode* node = new TrieNode(); nodes[i] = node;
        ifs.read((char *)&(node->nodeType), sizeof(node->nodeType));
        totalNode++; if (node->nodeType == nodeType_Leaf) totalLeaf++;
        ifs.read((char *)&(node->shortcutLength), sizeof(node->shortcutLength));
        if (node->shortcutLength > 0)
        {
            node->shortcut = new uint8_t[node->shortcutLength/8 + 2];
            ifs.read((char *)(node->shortcut), sizeof(uint8_t) * (node->shortcutLength / 8 + 2));
        }
        ifs.read((char *)&prindex, sizeof(prindex));
        ifs.read((char *)&prchildindex, sizeof(prchildindex));
        ifs.read((char *)&childmask, sizeof(childmask));
        if (childmask != 0) node->makeChildren();
        node->hash = new uint8_t[hashLength + 1];
        ifs.read((char *)(node->hash), sizeof(uint8_t) * hashLength);
        if (node->nodeType == nodeType_Leaf)
        {
            node->pValue = new uint256_t();
            ifs.read((char *)(node->pValue), sizeof(uint256_t));
        }
        // Link to Parent & Vice versa
        if (i != 0)
        {
            node->parent = nodes[prindex];
            if (0 <= prchildindex && prchildindex < MAX_TRIE_CHILDREN)
                nodes[prindex]->children[prchildindex] = node;
        }
    }
    return nodes[0];
}

int MerkleTrie::readFromIStream(istream& ifs)
{
    freeData();
    
    int len = max(2*hashLength, 2*dataLength);
    uint8_t dump[len];
    for (int i=0; i<len; i++) dump[len] = 0;
    
    try {
        ifs.read((char *)dump, 5);
        if (strcmp((const char*)dump, "MTRIE") != 0) return 2;
        
        uint16_t loadedVersion = intVersion;
        ifs.read((char *)&loadedVersion, sizeof(loadedVersion));
        
        root = readRootFromIStream(ifs, loadedVersion);
    }
    catch (...) {
        root = new TrieNode(); totalNode++;
        root->makeZeroHashIfNil(hashLength);
        root->makeChildren();
        return 0xE;
    }
    
    return 0;
}

int MerkleTrie::readFromFile(const char* path)
{
    ifstream ifs(path, ios::out | ios::binary);
    
    if(!ifs) return 1;
    
    int res = readFromIStream(ifs);
    
    ifs.close();
    
    return res;
}

#pragma mark - Quick Read

TrieNode* MerkleTrie::rootNodeFromIStream(istream& ifs, const int& hashlength, uint16_t* pLoadedVer, uint32_t* pTotalNode)
{
    TrieNode* rroot = NULL;
    int hashlen = hashlength;
    if (hashlen <= 0) hashlen = DEFAULT_HASH_LENGTH;
    int len = max(2*hashlength, 2*hashlength);
    uint8_t dump[len];
    for (int i=0; i<len; i++) dump[len] = 0;
    
    try {
        ifs.read((char *)dump, 5);
        if (strcmp((const char*)dump, "MTRIE") != 0) return NULL;
        
        uint16_t loadedVersion = 0;
        ifs.read((char *)&loadedVersion, sizeof(loadedVersion));
        
        uint32_t total_node = 0;
        ifs.read((char *)&total_node, sizeof(total_node));
        
        uint32_t prindex; uint8_t prchildindex; uint16_t childmask;
        
        if (total_node > 0) {
            TrieNode* node = new TrieNode();
            ifs.read((char *)&(node->nodeType), sizeof(node->nodeType));
            ifs.read((char *)&(node->shortcutLength), sizeof(node->shortcutLength));
            if (node->shortcutLength > 0)
            {
                node->shortcut = new uint8_t[node->shortcutLength/8 + 2];
                ifs.read((char *)(node->shortcut), sizeof(uint8_t) * (node->shortcutLength / 8 + 2));
            }
            ifs.read((char *)&prindex, sizeof(prindex));
            ifs.read((char *)&prchildindex, sizeof(prchildindex));
            ifs.read((char *)&childmask, sizeof(childmask));
            if (childmask != 0) node->makeChildren();
            node->hash = new uint8_t[hashlen + 1];
            ifs.read((char *)(node->hash), sizeof(uint8_t) * hashlen);
            if (node->nodeType == nodeType_Leaf)
            {
                node->pValue = new uint256_t();
                ifs.read((char *)(node->pValue), sizeof(uint256_t));
            }
            rroot = node;
        }
    }
    catch (...) {
        return NULL;
    }
    
    return rroot;
}

TrieNode* MerkleTrie::rootNodeFromFile(const char* path, const int& hashlength, uint16_t* pLoadedVer, uint32_t* pTotalNode)
{
    ifstream ifs(path, ios::out | ios::binary);
    
    if(!ifs) return NULL;
    
    TrieNode* res = rootNodeFromIStream(ifs, hashlength, pLoadedVer, pTotalNode);
    
    ifs.close();
    
    return res;
}

#pragma mark - Misc. Info

vector<TrieNode*> MerkleTrie::getNodesAtlayer(int layer, bool inclNULL)
{
    vector<TrieNode*> res;
    TrieNode* nullNodes[MAX_TRIE_CHILDREN+2];
    for (int i=0; i<MAX_TRIE_CHILDREN+1; i++) nullNodes[i] = NULL;
    
    if (layer == 0)
    {
        res.push_back(root);
        return res;
    }
    
    vector< pair<int,TrieNode*> > queue;
    if (root) queue.push_back(make_pair(0,root));
    while(queue.size() > 0)
    {
        TrieNode* node = queue.front().second;
        int layer = queue.front().first;
        queue.erase(queue.begin());
        TrieNode** child = (TrieNode**)nullNodes;
        if (node && node->children) child = (TrieNode**)node->children;
        
        for (int i=0; i<MAX_TRIE_CHILDREN; i++)
        {
            if (layer+1 == layer)
            {
                if (child[i]) res.push_back(node);
                else if (inclNULL) res.push_back(NULL);
            }
            else
            {
                queue.push_back( make_pair(layer+1,child[i]) );
            }
        }
        
    }
 
    return res;
}

uint32_t MerkleTrie::getTotalNode()
{
    return totalNode;
}

uint32_t MerkleTrie::getTotalLeaf()
{
    return totalLeaf;
}

uint32_t MerkleTrie::estMaxRAMSize()
{
    uint32_t intNode = totalNode - totalLeaf;
    uint32_t partsz = sizeof(TrieNode);
    uint32_t size = intNode*partsz;
    // Est Children Array
    partsz = intNode*MAX_TRIE_CHILDREN*sizeof(TrieNode*);
    size+=partsz;
    // Est Shortcut (in byte)
    uint32_t sccount = totalLeaf / (MAX_TRIE_CHILDREN/4); if (sccount < 2) sccount = 2;
    partsz = sccount * sizeof(uint8_t) * hashLength;
    size+=partsz;
    // Est Total Hash Byte
    partsz = totalNode*hashLength;
    size+=partsz;
    //Value
    partsz = totalLeaf * sizeof(uint256_t);
    size+=partsz;
    return size;
}

string MerkleTrie::generateDebugString(string split)
{
    string res = "";
    vector< pair<int,TrieNode*> > queue;
    if (root) queue.push_back( make_pair(0,root) );
    while(queue.size() > 0)
    {
        TrieNode* node = queue.front().second;
        int layer = queue.front().first;
        queue.erase(queue.begin());
        res = res + (node->generateDebugString(layer,hashLength) + split);
        if (!node->children) continue;
        for (int i=0; i<MAX_TRIE_CHILDREN; i++)
        {
            if (!node->children[i]) continue;
            queue.push_back( make_pair(layer+1,node->children[i]) );
            /*
            if (node->children[i]->nodeType != 0 && node->children[i]->nodeType != 1)
            {
                int c = 0;
            }
            */
        }
    }
    return res;
}

#pragma mark - Dealloc

void MerkleTrie::freeData()
{
    deque< TrieNode*> queue;
    if (root) queue.push_back(root);
    while(queue.size() > 0)
    {
        TrieNode* node = queue.front();
        queue.pop_front();
        if (duplicatedValue && node->pValue) delete node->pValue;
        if (node->children)
        {
            for (int i=0; i<MAX_TRIE_CHILDREN; i++)
            {
                if (!node->children[i]) continue;
                queue.push_back( node->children[i] );
            }
        }
        delete node;
    }
    root = NULL;
    totalNode = 0; totalLeaf = 0;
}

MerkleTrie::~MerkleTrie()
{
    freeData();
}
