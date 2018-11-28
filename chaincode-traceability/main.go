package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strings"
)

type TraceabilityChaincode struct {
}

type TDoc struct {
	// BizId
	// todo: need prefix?
	ID string `json:"id"`
	// PreSerialNo string  `json:"preserialno"`
	// file hash (certInfo)
	Hash string `json:"hash"`
	// accountid of gateway
	AccountID string `json:"accountid"`
	// FileName  string `json:"filename"`
	Comment string `json:"comment"`
}

type ReturnValue struct {
	//Tdoc TDoc
	ContractAddress string `json:"contractaddress"`
	TxHash          string `json:"txhash"`
	BlockNumber     int    `json:"blocknumber"`
	Timestamp       string `json:"timestamp"`
}

func main() {
	err := shim.Start(new(TraceabilityChaincode))
	if err != nil {
		fmt.Printf("Error starting Traceability chaincode: %s", err)
	}
}

func (t *TraceabilityChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Print("chaincode initial...")
	return shim.Success(nil)
}

func (t *TraceabilityChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
	// Handle different functions
	if function == "uploaddoc" {
		return t.uploadDoc(stub, args)
	} else if function == "querydoc" {
		return t.queryDoc(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

func (t *TraceabilityChaincode) uploadDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var tdoc TDoc

	//
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting ID, Hash, AccountID and Comment")
	}

	if len(args[0]) <= 0 {
		return shim.Error("ID can't be empty!")
	}
	if len(args[1]) <= 0 {
		return shim.Error("Hash can't be empty!")
	}

	tdoc.ID = strings.TrimSpace(args[0])
	tdoc.Hash = strings.TrimSpace(args[1])
	tdoc.AccountID = strings.TrimSpace(args[2])
	tdoc.Comment = strings.TrimSpace(args[3])

	val, err := json.Marshal(tdoc)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(tdoc.ID, []byte(val))
	if err != nil {
		return shim.Error(err.Error())
	}

	var returnValue ReturnValue

	returnValue.TxHash = stub.GetTxID()
	timestamp, err := stub.GetTxTimestamp()
	returnValue.Timestamp = timestamp.String()
	if err != nil {
		return shim.Error(err.Error())
	}

	ret, err := json.Marshal(returnValue)
	if err != nil {
		return shim.Error(err.Error())
	}
	//returnValue.ContractAddress = stub.get
	return shim.Success(ret)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (t *TraceabilityChaincode) queryDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting ID")
	}

	id := args[0]

	val, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(val)
}
