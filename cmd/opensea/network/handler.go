package network

import (
	// "encoding/hex"
	// "encoding/json"
	"errors"
	"fmt"

	// "strings"
	// "time"

	"gitlab.com/meta-node/meta-node/pkg/smart_contract"

	// . "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/meta-node/meta-node/cmd/client/command"

	// "gitlab.com/meta-node/client/utils"
 "gitlab.com/meta-node/meta-node/pkg/logger"
	"gitlab.com/meta-node/meta-node/pkg/network"
	pb "gitlab.com/meta-node/meta-node/pkg/proto"
	"gitlab.com/meta-node/meta-node/pkg/receipt"
	"gitlab.com/meta-node/meta-node/pkg/state"
	"gitlab.com/meta-node/meta-node/pkg/transaction"
)

var (
	ErrorCommandNotFound = errors.New("command not found")
)
type Receipt1 struct {
	Hash  common.Hash
	Value interface{}
}
type Handler struct {
	accountStateChan chan state.IAccountState
	chData chan interface{}
}

func NewHandler(
	accountStateChan chan state.IAccountState,
	chData chan interface{},
) *Handler {
	return &Handler{
		accountStateChan: accountStateChan,
		chData :chData ,
	}
}
func (h *Handler) GetChData() chan interface{} {
	return h.chData
}


func (h *Handler) HandleRequest(request network.IRequest) (err error) {
	cmd := request.GetMessage().GetCommand()
	logger.Debug("handling command: " + cmd)
	switch cmd {
	case command.InitConnection:
		return h.handleInitConnection(request)
	case command.AccountState:
		return h.handleAccountState(request)
	case command.Receipt:
		return h.handleReceipt(request)
	case command.TransactionError:
		return h.handleTransactionError(request)
	case command.EventLogs:
		return h.handleEventLogs(request)
	}
	return ErrorCommandNotFound
}

/*
handleInitConnection will receive request from connection
then init that connection with data in request then
add it to connection manager
*/
func (h *Handler) handleInitConnection(request network.IRequest) (err error) {
	conn := request.GetConnection()
	initData := &pb.InitConnection{}
	err = request.GetMessage().Unmarshal(initData)
	if err != nil {
		return err
	}
	address := common.BytesToAddress(initData.Address)
	logger.Debug(fmt.Sprintf(
		"init connection from %v type %v", address, initData.Type,
	))
	conn.Init(address, initData.Type, initData.PublicConnectionAddress)
	return nil
}

/*
handleAccountState will receive account state from connection
then push it to account state chan
*/
func (h *Handler) handleAccountState(request network.IRequest) (err error) {
	accountState := &state.AccountState{}
	err = accountState.Unmarshal(request.GetMessage().GetBody())
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("Receive Account state: \n%v", accountState))
	h.accountStateChan <- accountState
	return nil
}


/*
handleTransactionError will receive transaction error from parent node connection
then print it out
*/
func (h *Handler) handleTransactionError(request network.IRequest) (err error) {
	transactionErr := &transaction.TransactionError{}
	err = transactionErr.Unmarshal(request.GetMessage().GetBody())
	if err != nil {
		return err
	}
	logger.Info("Receive transaction error:", transactionErr)

	return nil
}

func (h *Handler) handleEventLogs(request network.IRequest) error {
	eventLogs := smart_contract.EventLogs{}
	err := eventLogs.Unmarshal(request.GetMessage().GetBody())
	if err != nil {
		logger.Error("Handle Event Logs Error", err)
		return err
	}
	eventLogList := eventLogs.GetEventLogList()
	for _, eventLog := range eventLogList {
		logger.Info("EventLogs: ", eventLog.String())
	}
	return nil
}
/*
handleReceipt will receive receipt from connection
then print it out
*/
func (h *Handler) handleReceipt(request network.IRequest) (error) {
	
	receipt1 := &receipt.Receipt{}
	err := receipt1.Unmarshal(request.GetMessage().GetBody())
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Receive receipt: %v", receipt1))
	logger.Info(fmt.Sprintf("Receive To address: %v", request.GetMessage().GetToAddress()))
	fmt.Println("receipt1 là:",receipt1.GetStatus())
	if receipt1.GetStatus() == 1 || receipt1.GetStatus() == 0 {
		// callback:=h.handleReceipt1(receipt1)
		fmt.Println("hlshldfhsklafklsf")
		h.chData <- Receipt1{
			receipt1.GetTransactionHash(),
			common.Bytes2Hex(receipt1.GetReturn()),
			// callback,
		}
	} else {
		log.Warn("Call Error !!! - ", common.Bytes2Hex(receipt1.GetReturn()))
	}

	return nil
}
// func GetTransactionsByHash(hash string) map[string]interface{} {
// 	result:=make(map[string]interface{})
// 	db, err := sqlx.Connect("sqlite3", "./database/doc_2022-12-26_09-16-03.meta_findsdk.db")
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("error when connect sqlite %", err))
// 		panic(fmt.Sprintf("error when connect sqlite %v", err))
// 	}


// 	transCtrl := hdlTransaction.NewTransactionController(db)

//     transactionInDbKq,err := transCtrl.GetTransactionByHash(hash)

//     if err != nil {
//         return map[string]interface{}{
//             "id": -1,
//         }
//     }
// 	transactionInDb:=transactionInDbKq.Data.(models.Header).Data.(hdlTransaction.TransactionModel)
// 	bTransactionInDb,err:=json.Marshal(transactionInDb)
// 	err=json.Unmarshal(bTransactionInDb,result)
//     return result
// }

// func (h *Handler)handleReceipt1(receipt *receipt.Receipt)  map[string]interface{}{
// 	callback:=make(map[string]interface{})

//     hash := receipt.GetTransactionHash().Hex()

//     fromAddress := receipt.GetFromAddress()
//     toAddress := strings.ToLower((receipt.GetToAddress().Hex())[2:])
//     amount := receipt.GetAmount()
//     status := receipt.GetStatus()
//     returnValue := hex.EncodeToString(receipt.GetReturn())

//     commandResponse := "send-transaction"

//     fmt.Printf("status-rc: %d\n", status)
// 		response := make(map[string]interface{})
// 		response["hash"] = hash
// 		response["from"] = fromAddress
// 		response["to"] = toAddress
// 		response["value"] = amount
// 		response["status"] = status
	
// 		var newTransaction map[string]interface{}
	
	
// 		var returnData []interface{}	
// 			// b, err := json.Marshal(transactionInDb)		
// 			// if err != nil {
// 			// 	fmt.Println("error:", err)
// 			// }
// 			// err = json.Unmarshal(b, &newTransaction)

// 			// newTransaction = transactionInDb
	
// 				commandResponse = "excute-smart-contract"
	
// 				functionCall := transactionInDb.FunctionCall
// 				sc := dappCtrl.GetSmartContractByAddress(toAddress)
// 				var abiData []interface{} 

// 				err:=json.Unmarshal([]byte(sc["abiData"].(string)),&abiData)
// 				if err != nil {
// 					logger.Error(fmt.Sprintf("error when Unmarshal send-transaction %", err))
// 					panic(fmt.Sprintf("error when Unmarshal send-transaction %v", err))
// 				}
// 				for i := 0; i < len(abiData); i++ {
// 					item := abiData[i].(map[string]interface{} )
// 					if item["type"] == "function" && item["name"] == functionCall {
// 						encb, err := hex.DecodeString(returnValue)
// 						if err != nil {
// 							fmt.Printf("invalid hex %s: %v", returnValue, err)
// 						}
// 						fmt.Println("3333333")
// 						fmt.Println("item:",item)
// 						b,err :=json.Marshal(item)
// 						var itemArr []interface{}
// 						for _,v:= range b{
// 							itemArr=append(itemArr,v)
// 						}

// 						fmt.Println("44444444")

// 						if err != nil {
// 							logger.Error(fmt.Sprintf("error when Marshal in send-transaction %", err))
// 							panic(fmt.Sprintf("error when Marshal in send-transaction %v", err))
// 						}
// 						fmt.Println("string(b):",string(b))

// 						abiParser, err := JSON(strings.NewReader(sc["abiData"].(string)))

// 						fmt.Println("5555555")

// 						if err != nil {
// 							logger.Error(fmt.Sprintf("error when JSON send-transaction %", err))
// 							panic(fmt.Sprintf("error when JSON send-transaction %v", err))
// 						}
// 						fmt.Println("66666666")

// 						returnData, err = abiParser.Unpack(functionCall, encb)
// 						if err != nil {
// 							fmt.Printf("test %d (%v) failed: %v", i, returnValue, err)
// 						}
// 						break
// 					}
// 				}
// 				fmt.Println("Decode out là:",returnData)

// 				// returnData,_ = decodeAbi(returnValue, []byte(abi))
// 		// if callback != nil {
// 			data := map[string]interface{}{
// 				"success": true,
// 				"data":    newTransaction,
// 			}
// 			if newTransaction == nil {
// 				data["data"] = response
// 			}
// 			if returnData != nil {
// 				data["returnData"] = returnData
// 			}
// 			callback = map[string]interface{}{
// 				"command": commandResponse,
// 				"data":    data,
// 			}
// 			// return callback
// 		// }
// 		callback=(map[string]interface{}{
// 			"command": "TransactionError",
// 			"data": map[string]interface{}{
// 				"success": false,
// 				"data": "error",
// 			},
// 		})
// 		// }
// 	return callback
// }


// func decodeAbi(rawInput string, abiJSON []byte) ([]interface{}, error) {
//     var outputParameters []interface{}
//     abiParser, err := JSON(strings.NewReader(string(abiJSON)))
//     if err != nil {
//         return nil, err
//     }

//     for _, item := range abiParser.Methods {
// 		itemType:= item.Inputs[0].Type.String()
//         switch itemType {
//         case "tuple":
//             outputParameters = append(outputParameters, item.Inputs[0].Type.T)
//         case "tuple[]":
//             outputParameters = append(outputParameters, &[]item.Inputs[0].Type.T)
//         default:
//             // outputParameters = append(outputParameters, abiParser.ArgumentToType(&Argument{
//             //     Type: item.Inputs[0].Type.String(),
//             // }))
// 			t, err := NewType(itemType, "", nil)
// 			if err != nil {
// 				return nil, err
// 			}
// 			outputParameters = append(outputParameters, t)
//         }
//     }
// 	// arg.Type.T 
// 	encb, err := hex.DecodeString(rawInput)
// 			if err != nil {
// 				t.Fatalf("invalid hex %s: %v", rawInput, err)
// 			}
// 	result, err := abiParser.Unpack(rawInput, encb)

//     // result, err := abiParser.Unpack(rawInput, outputParameters[:]...)
//     if err != nil {
//         return nil, err
//     }

//     jsonResult := []interface{}{}
//     for _, item := range result {
//         jsonResult = append(jsonResult, getJSONResult(item))
//     }

//     return jsonResult, nil
// 	// data, err := hexutil.Decode(rawInput)
//     // if err != nil {
//     //     return nil, err
//     // }
//     // results, err := Methods.Unpack(outputParameters, data)
//     // if err != nil {
//     //     return nil, err
//     // }
//     // var result []interface{}
//     // for _, item := range results {
//     //     jsonResult, err := json.Marshal(item)
//     //     if err != nil {
//     //         return nil, err
//     //     }
//     //     var decoded interface{}
//     //     if err := json.Unmarshal(jsonResult, &decoded); err != nil {
//     //         return nil, err
//     //     }
//     //     result = append(result, decoded)
//     // }
//     // return result, nil
// }

// func getJSONResult(item interface{}) interface{} {
//     switch value := item.(type) {
//     case []*network.IRequest:
//         jsonArray := []interface{}{}
//         for _, v := range value {
//             jsonArray = append(jsonArray, getJSONResult(v))
//         }
//         return jsonArray
//     default:
//         jsonData := map[string]interface{}{
//             "type":  fmt.Sprintf("%T", item),
//             "value": hexutil.Encode(item.([]byte)),
// 			// "value":  item.Value,

//         }
//         return jsonData
//     }
// }



// func getJSONResult(item Type) interface{} {
//     if itemList, ok := item.Value.([]Type); ok {
//         jsonArray := make([]interface{}, 0)
//         for _, value := range itemList {
//             jsonArray = append(jsonArray, getJSONResult(value))
//         }
//         return jsonArray
//     } else {
//         jsonObj := make(map[string]interface{})
//         jsonObj["type"] = item.TypeAsString
//         jsonObj["value"] = item.Value
//         return jsonObj
//     }
// }

// func getJSONResult(item interface{}) interface{} {
//     if value, ok := item.([]interface{}); ok {
//         jsonArray := make([]interface{}, 0)
//         for _, v := range value {
//             jsonArray = append(jsonArray, getJSONResult(v))
//         }
//         return jsonArray
//     } else {
//         json := make(map[string]interface{})
//         json["type"] = item.(Type).typeAsString()
//         json["value"] = item.(Type).value
//         return json
//     }
// }

