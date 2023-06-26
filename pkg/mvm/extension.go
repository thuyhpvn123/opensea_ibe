package mvm

/*
#cgo CFLAGS: -w
#cgo CXXFLAGS: -std=c++17 -w
#cgo LDFLAGS: -L./linker/build/lib/static -lmvm_linker -L./c_mvm/build/lib/static -lmvm -lstdc++
#cgo CPPFLAGS: -I./linker/build/include
#include "mvm_linker.hpp"
#include <stdlib.h>
*/
import "C"
import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"unsafe"

	"gitlab.com/meta-node/meta-node/pkg/logger"
	"gitlab.com/meta-node/meta-node/pkg/smart_contract/argument_encode"
)

//export ExtensionCallGetApi
func ExtensionCallGetApi(
	bytes *C.uchar,
	size C.int,
) (
	data_p *C.uchar,
	data_size C.int,
) {
	bCallData := C.GoBytes(unsafe.Pointer(bytes), size)
	logger.Debug("Calling get api data ", hex.EncodeToString(bCallData))
	url := argument_encode.DecodeStringInput(bCallData[4:], 0)
	response, err := http.Get(url)
	if err != nil {
		logger.Warn("Error when call get api to ", url, err)
		return
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Warn("Error when call get api to ", url, err)
		return
	}
	encodedRespone := argument_encode.EncodeSingleString(string(responseData))
	logger.Debug("Extension call get api result ", encodedRespone)
	data_size = C.int(len(encodedRespone))
	data_p = (*C.uchar)(C.CBytes(encodedRespone))
	return
}

//export ExtensionExtractJsonField
func ExtensionExtractJsonField(
	bytes *C.uchar,
	size C.int,
) (
	data_p *C.uchar,
	data_size C.int,
) {
	bCallData := C.GoBytes(unsafe.Pointer(bytes), size)
	logger.Debug("Extension extract json field ", hex.EncodeToString(bCallData))
	jsonMap := make(map[string]interface{})
	jsonStr := argument_encode.DecodeStringInput(bCallData[4:], 0)
	field := argument_encode.DecodeStringInput(bCallData[4:], 1)
	var fieldData interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonMap)
	var data string
	// process json map
	if err == nil {
		fieldData = jsonMap[field]
	} else {
		// process json array
		jsonArr := []interface{}{}
		err = json.Unmarshal([]byte(jsonStr), &jsonArr)
		if err != nil {
			logger.Warn("Error when extract json field ", jsonStr, field, err)
			return
		}
		intField, err := strconv.Atoi(field)
		if err != nil {
			logger.Warn("Error when extract json field ", jsonStr, field, err)
			return
		}
		fieldData = jsonArr[intField]
	}

	if reflect.ValueOf(fieldData).Kind() == reflect.Map || reflect.ValueOf(fieldData).Kind() == reflect.Array {
		bData, _ := json.Marshal(fieldData)
		data = string(bData)
	} else {
		data = fmt.Sprintf("%v", fieldData)
		// reformat boolean
		if data == "false" {
			data = "0"
		}
		if data == "true" {
			data = "1"
		}
	}

	encodedData := argument_encode.EncodeSingleString(data)
	data_size = C.int(len(encodedData))
	data_p = (*C.uchar)(C.CBytes(encodedData))
	return
}
