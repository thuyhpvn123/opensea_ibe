package controllers

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"gitlab.com/meta-node/meta-node/pkg/bls"
	cm "gitlab.com/meta-node/meta-node/pkg/common"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)
var deckKq []string
var keysArray []string
// "gitlab.com/meta-node/meta-node/cmd/chiabai/core"
func(caller *CallData) VerifySign(call map[string]interface{}) map[string]interface{} {
	fmt.Println("hhhhhhhhh")
	result:=make(map[string]interface{})
	hash, ok := call["hash"].(string)
    if !ok {
        result=(map[string]interface{}{
            "success": true,
            "data":    false,
        })
        return result
    }
	bhash:=common.FromHex(hash)
    sign, ok := call["sign"].(string)
    if !ok {
        result=(map[string]interface{}{
            "success": true,
            "data":    false,
        })
        return result
    }
	bsign:=common.FromHex(sign)

    pubKey, ok := call["pubKey"].(string)
    if !ok {
        result=(map[string]interface{}{
            "success": true,
            "data":    false,
        })
        return result
    }
	bpubKey:=common.FromHex(pubKey)
    pubKeyCm:=cm.PubkeyFromBytes(bpubKey)
    signCm:=cm.SignFromBytes(bsign)    
    success := bls.VerifySign(pubKeyCm, signCm, bhash)
	address:=crypto.Keccak256(bpubKey)[12:]
	// testPubkey1 = common.FromHex("a2702ce6bbfb2e013935781bac50a0e168732bd957861e6fbf185d688c82ade34c9f33fead179decb5953b3382b061df")
	// testSign1   = common.FromHex("a507c03ab7ebb69a4b3adc22a0347bb2466788e6a3baa174a62bd74cdff60dfd6d6ba9ec6237098f1ceef6013bfeff1d0c8be716266710e1493c422293a676e7f168007324a23435d4590896f97f8e3686cf0c280240b9406800c1cec6bafb5d")
	// testHash1   = common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
	// sign1 := Sign(cm.PrivateKeyFromBytes(common.FromHex(testSecret1)), testHash1.Bytes())
	// fmt.Printf("Sign1: %v\n", common.Bytes2Hex(sign1.Bytes()))

	result=(map[string]interface{}{
        "success": true,
        "data":    success,
		"address":hex.EncodeToString(address),
    })
	fmt.Println("verifySign:",success)

    return result

}
func(caller *CallData) GetKeyForPlayer(callMap map[string]interface{})  {
	// verifyKq := caller.verifySign(callMap)

	// if(verifyKq["data"].(bool)==true){
		fmt.Println("GetKeyForPlayer:",callMap)
		cardArr:=caller.GetCards(callMap)
		fmt.Println("Get card done")
		fmt.Println("encrypted-cards",cardArr)
		fmt.Println("encrypted-deck",deckKq)
		fmt.Println("encrypted-keys",keysArray)
		call:=map[string]interface{}{
			"encrypted-cards":cardArr,
			"encrypted-deck": deckKq,
			"encrypted-keys": keysArray,
		}
		keyArr:=caller.FindKeys(call)
		caller.sentToClient("get-key-for-player", keyArr)
	// }else{
	// 	caller.sentToClient("get-key-for-player", "Not Authorised Address")
	// }
	
}
func(caller *CallData) GetCards(callMap map[string]interface{}) interface{} {
	contract := caller.server.contractABI["chiabai"]
	// call:=callMap["value"].(map[string]interface{}) 
	functionName:=callMap["function-name"].(string)
	result:=caller.TryCall(callMap)
	cards := contract.Decode(functionName, result.(string)).(map[string]interface{})[""]
	log.Info("GetCards - Result - ", cards)
	go caller.sentToClient("get-cards", cards)
	return cards
}
	
func(caller *CallData) ShuffleCard(callMap map[string]interface{})  {
	contract := caller.server.contractABI["chiabai"]
	// call:=callMap["value"].(map[string]interface{}) 
	functionName:=callMap["function-name"].(string)
	result:=caller.TryCall(callMap)
	deck := contract.Decode(functionName, result.(string)).(map[string]interface{})[""]
	// deckKq=deck.([]string)
	log.Info("ShuffleCard - Result - ", deck)
	go caller.sentToClient("shuffle-card", fmt.Sprint(deck))

}
func(caller *CallData) SetDeck(callMap map[string]interface{})  {
	contract := caller.server.contractABI["chiabai"]
	// call:=callMap["value"].(map[string]interface{}) 
	functionName:=callMap["function-name"].(string)
	result:=caller.TryCall(callMap)
	deck := contract.Decode(functionName, result.(string)).(map[string]interface{})[""]
	log.Info("SetDeck - Result - ", deck)
	go caller.sentToClient("set-Deck", fmt.Sprint(deck))
}
func(caller *CallData) SetPlayers(callMap map[string]interface{})  {
	result:=caller.TryCall(callMap)
	go caller.sentToClient("set-players", result)
}
// find array keys
func (caller *CallData) FindKeys(callMap map[string]interface{}) []string {
	fmt.Println("Find Keys")
	encryptedDeck:=callMap["encrypted-deck"].([]string)
	encryptedCardArr:=callMap["encrypted-cards"].([]string)
	encryptedKey:=callMap["encrypted-keys"].([]string)
	indices := findIndices(encryptedDeck, encryptedCardArr)	
	result := findArray(encryptedKey, indices)
	return result
}

func findIndices(firstArray []string, secondArray []string) []int {
	indices := make([]int, len(secondArray))
	indexMap := make(map[string]bool)

	// Create a map of elements from the first array for efficient lookup
	for _, num := range firstArray {
		indexMap[num] = true
	}

	// Find the indices of elements from the second array in the first array
	for i, num := range secondArray {
		if indexMap[num] {
			indices[i] = findIndex(firstArray, num)
		} else {
			indices[i] = -1 // Element not found in the first array
		}
	}

	return indices
}

func findIndex(arr []string, target string) int {
	for i, num := range arr {
		if num == target {
			return i
		}
	}
	return -1 // Element not found
}
func findArray(firstArray []string, secondArray []int) []string {
	result := make([]string, len(secondArray))

	for i, num := range secondArray {
		if num >= 0 && num < len(firstArray) {
			result[i] = firstArray[num]
		} else {
			result[i] = "nil" // Invalid index, assign -1 as the value
		}
	}

	return result
}
func(caller *CallData) GetSign(callMap map[string]interface{}) cm.Sign {
	privateKey:=callMap["privateKey"].(string)
	addressForSign:=callMap["address"].(string)
// Initialize the random number generator
	rand.Seed(time.Now().UnixNano())

// Generate a random number between 0 and 100
	randomNumber := rand.Intn(101)

	message := common.FromHex(addressForSign + strconv.Itoa(randomNumber))
//vd addressForSign= "0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2" 
	keyPair := bls.NewKeyPair(common.FromHex(privateKey))
//vd privateKey="36e1aa979f98c7154fb2491491ec044ccac099651209ccfbe2561746dbe29ebb"
	hash := crypto.Keccak256(message)
	prikey:= keyPair.GetPrivateKey()

	sign := bls.Sign(prikey, hash)
	pubkey:= keyPair.GetBytesPublicKey()
	add:= keyPair.GetAddress()
	address:=crypto.Keccak256(pubkey)[12:]
	fmt.Println("sign:",sign)
	fmt.Println("pubkey:",hex.EncodeToString(pubkey))
	fmt.Println("address:",add)
	fmt.Println("hash:",hex.EncodeToString(hash))
	fmt.Println("address tu publickey:",hex.EncodeToString(address))
	go caller.sentToClient("get-sign", fmt.Sprint(sign))

	return sign
}
//result :
// sign: a7d09af4337cb7e8630f7e4e0ff05f881bdebff862c73d13f304004c50f5cd3090b464582789e02a79907d3e5a0bfba515743e087ddb4ff177ff05bdc8908770e35c9d115fe3d4b9417136c9eee7bdc75996f1f25f25815d2510df387fbd833a
// pubkey: aff6fc519567a8513179b20dddda1b8ba9556cf536117c28ef24d2cc1895b0a10d95ac31d2f0bc00eb24d971099164ab
// address: 0x45c75cfb8E20a8631C134555FA5d61FCf3E602f2
// hash: 999bf57501565dbd2fdcea36efa2b9aef8340a8901e3459f4a4c926275d36cdb