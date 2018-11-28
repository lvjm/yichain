package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strings"
	"time"
)

type DigitalAssetChaincode struct {
}

type TDoc struct {

	TxHash string `json:"txHash"`
	TransactionAt string `json:"transactionAt"`
	AssetAmount string `json:"assetAmount"`
	AssetUnit string `json:"assetUnit"`
	TransactionFromAddress string `json:"transactionFromAddress"`
	TransactionFromAddressAlias string `json:"transactionFromAddressAlias"`
	TransactionToAddress string `json:"transactionToAddress"`
	TransactionToAddressAlias string `json:"transactionToAddressAlias"`
	Comment string `json:"comment"`
	ProvidedGas string `json:"providedGas"`
	ProvidedFee string `json:"providedFee"`
	ActualGas string `json:"actualGas"`
	ActualFee string `json:"actualFee"`
	FeeUnit string `json:"feeUnit"`
	BlockNum string `json:"blockNum"`
	BlockHash string `json:"blockHash"`
	CreateAt string `json:"createAt"`
}

type ReturnValue struct {
	//Tdoc TDoc
	TxHash          string `json:"txhash"`
	BlockNumber     int    `json:"blocknumber"`
	Timestamp       string `json:"timestamp"`
}

func main() {
	err := shim.Start(new(DigitalAssetChaincode))
	if err != nil {
		fmt.Printf("Error starting DigitalAssetChaincode chaincode: %s", err)
	}
}

func (t *DigitalAssetChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Print("chaincode initial...")
	return shim.Success(nil)
}

func (t *DigitalAssetChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
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

func (t *DigitalAssetChaincode) uploadDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var tdoc TDoc

	//
	if len(args) != 16 {
		return shim.Error("Incorrect number of arguments. Expecting 16")
	}

	if len(args[0]) <= 0 {
		return shim.Error("txHash can't be empty!")
	}

	tdoc.TxHash = strings.TrimSpace(args[0])
	tdoc.TransactionAt = strings.TrimSpace(args[1])
	tdoc.AssetAmount = strings.TrimSpace(args[2])
	tdoc.AssetUnit = strings.TrimSpace(args[3])
	tdoc.TransactionFromAddress= strings.TrimSpace(args[4])
	tdoc.TransactionFromAddressAlias = strings.TrimSpace(args[5])
	tdoc.TransactionToAddress = strings.TrimSpace(args[6])
	tdoc.TransactionToAddressAlias = strings.TrimSpace(args[7])
	tdoc.Comment = strings.TrimSpace(args[8])
	tdoc.ProvidedGas = strings.TrimSpace(args[9])
	tdoc.ProvidedFee = strings.TrimSpace(args[10])
	tdoc.ActualGas = strings.TrimSpace(args[11])
	tdoc.ActualFee= strings.TrimSpace(args[12])
	tdoc.FeeUnit = strings.TrimSpace(args[13])
	tdoc.BlockNum = strings.TrimSpace(args[14])
	tdoc.BlockHash = strings.TrimSpace(args[15])
	tdoc.CreateAt = time.Now().Format("2006-01-02 15:04:05")
	val, err := json.Marshal(tdoc)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(tdoc.TxHash, []byte(val))
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

func (t *DigitalAssetChaincode) queryDoc(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting TxHash")
	}

	id := args[0]

	val, err := stub.GetState(id)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(val)
}
