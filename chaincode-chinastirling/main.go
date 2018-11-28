package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strings"
        "strconv"
        "bytes"
)

type ChinastirlingIOTCertChaincode struct {
}

type ChinastirlingIOTCert struct {
	ObjectType         string `json:"docType"` //docType 用于区分状态数据库中的各种类型
	AccountId          string `json:"accountId"`
	BizId              string `json:"bizId"`
	RawDataPacket      string `json:"rawDataPacket"`
	MacAddress         string `json:"macAddress"`
	DeviceSerialNo     int    `json:"deviceSerialNo"` //这个我用int而非string，是因为我在有个智能合约是传入serialNo的上下限，返回对应记录，所以该字段要可比较
	Interval           string `json:"interval"`
	SettingTemperature string `json:"settingTemperature"`
	SettingMode        string `json:"settingMode"`
	SettingStatus      string `json:"settingStatus"`
	SwitchCount        string `json:"switchCount"`
	TemperatureSensor1 string `json:"temperatureSensor1"`
	TemperatureSensor2 string `json:"temperatureSensor2"`
	TemperatureSensor3 string `json:"temperatureSensor3"`
	TemperatureSensor4 string `json:"temperatureSensor4"`
	TemperatureSensor5 string `json:"temperatureSensor5"`
	Voltage            string `json:"voltage"`
	Ampere             string `json:"ampere"`
	Watt               string `json:"watt"`
	Datetime           string `json:"datetime"`
	Latitude           string `json:"latitude"`
	Longitude          string `json:"longitude"`
	Comment            string `json:"comment"`
	TxHash             string `json:"txHash"`
	Timestamp          string `json:"timestamp"`
}

//返回结果定义
type ReturnValue struct {
	TxHash    string `json:"txHash"`
	Timestamp string `json:"timestamp"`
}

func main() {
	//主函数中启动该智能合约
	err := shim.Start(new(ChinastirlingIOTCertChaincode))
	if err != nil {
		fmt.Printf("Error starting ChinastirlingIOTCert chaincode: %s", err)
	}
}

// ========================================================================
// 初始化方法，在本例子中，我们无需初始化
// ========================================================================
func (t *ChinastirlingIOTCertChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Print("chaincode initial...")
	return shim.Success(nil)
}

// ========================================================================
// 方法调用中心代理，多个实际的方法调用转发
// ========================================================================
func (t *ChinastirlingIOTCertChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	if function == "uploadIOTEnvData" {
		return t.uploadIOTEnvData(stub, args)
	} else if function == "queryByBizId" {
		return t.queryByBizId(stub, args)
	}else if function == "queryByAccountId" {
		return t.queryByAccountId(stub, args)
	}else if function == "queryByMacAddress" {
		return t.queryByMacAddress(stub, args)
	}else if function == "queryByRawDataPacket" {
		return t.queryByRawDataPacket(stub, args)
	}else if function == "queryBySerialNoRange" {
		return t.queryByMacAddress(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ========================================================================
// 智能合约方法：从某个华斯特林的IOT设备上传IOT的环境数据包
// ========================================================================
func (t *ChinastirlingIOTCertChaincode) uploadIOTEnvData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//必须是22个参数
	if len(args) != 22 {
		return shim.Error("Incorrect number of arguments.")
	}

	if len(args[0]) <= 0 {
		return shim.Error("Account Id can't be empty!")
	}
	if len(args[1]) <= 0 {
		return shim.Error("Biz Id can't be empty!")
	}

	bizId := strings.TrimSpace(args[1])
	//根据bizId查询对应的记录是否已经存在，如果已经存在，则不存
	iotCertAsBytes, err := stub.GetState(bizId)
	if err != nil {
		return shim.Error("Failed to get cert: " + err.Error())
	} else if iotCertAsBytes != nil {
		fmt.Println("This iot cert already exists: " + bizId)
		return shim.Error("This iot cert already exists: " + bizId)
	}

	//从入参中获取所有的必要参数
	accountId := strings.TrimSpace(args[0])
	rawDataPacket := strings.TrimSpace(args[2])
	macAddress := strings.TrimSpace(args[3])
	deviceSerialNo,err  := strconv.Atoi(args[4]) 
        if err != nil {
		return shim.Error("The fifth argument must be a numeric string")
	}
	interval := strings.TrimSpace(args[5])
	settingTemperature := strings.TrimSpace(args[6])
	settingMode := strings.TrimSpace(args[7])
	settingStatus := strings.TrimSpace(args[8])
	switchCount := strings.TrimSpace(args[9])
	temperatureSensor1 := strings.TrimSpace(args[10])
	temperatureSensor2 := strings.TrimSpace(args[11])
	temperatureSensor3 := strings.TrimSpace(args[12])
	temperatureSensor4 := strings.TrimSpace(args[13])
	temperatureSensor5 := strings.TrimSpace(args[14])
	voltage := strings.TrimSpace(args[15])
	ampere := strings.TrimSpace(args[16])
	watt := strings.TrimSpace(args[17])
	datetime := strings.TrimSpace(args[18])
	latitude := strings.TrimSpace(args[19])
	longitude := strings.TrimSpace(args[20])
	comment := strings.TrimSpace(args[21])
	
	
	// ==== 获得txId和存证时间信息====
	txHash := stub.GetTxID()                 //交易提案时指定的交易ID
	timestamp, err := stub.GetTxTimestamp() //交易被创建时客户端的时间戳，从交易的ChannelHeader中提取
         if err != nil {
                return shim.Error(err.Error())
         }
	

	// ==== 创建对象，并且序列化为json====
	objectType := "iotcert"
	chinastirlingIOTCert := &ChinastirlingIOTCert{objectType, accountId, bizId, rawDataPacket,
		macAddress, deviceSerialNo, interval, settingTemperature, settingMode, settingStatus, switchCount,
		temperatureSensor1, temperatureSensor2, temperatureSensor3, temperatureSensor4, temperatureSensor5,
		voltage, ampere, watt, datetime, latitude, longitude, comment, txHash, timestamp.String()}
	iotcertJSONAsBytes, err := json.Marshal(chinastirlingIOTCert)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== 将iot存证写入couchdb状态库====
	err = stub.PutState(bizId, iotcertJSONAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//构造返回结果
	var returnValue ReturnValue
	returnValue.TxHash = txHash
	returnValue.Timestamp = timestamp.String()
	ret, err := json.Marshal(returnValue)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(ret)

	// ==== 返回结果 ====
	fmt.Println("- end upload iot env data")
	return shim.Success(ret)
}

// ========================================================================
// 基于键：bizId来查询对应的IOT存证信息
// ========================================================================
func (t *ChinastirlingIOTCertChaincode) queryByBizId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
        var jsonResp string
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments.")
	}

	bizId := args[0]
	valAsbytes, err := stub.GetState(bizId)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + bizId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"IOT Cert does not exist: " + bizId + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(valAsbytes)
}

// ============================================================================================================
// 工具方法：基于查询字符串获得查询结果。结果集会报封装在byte array以json形式返回
// ============================================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// ============================================================================================================
// 富查询：基于accountId(chinastirling)查询所有的匹配环境存证
// ============================================================================================================
func (t *ChinastirlingIOTCertChaincode) queryByAccountId(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	accountId := strings.TrimSpace(args[0])

	//queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"iotcert\",\"accountId\":\"%s\"}}", accountId)
	queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"docType\":\"iotcert\"},{\"accountId\":\"%s\"}]}}", accountId)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ========================================================================
// 富查询：基于accountId+macAddress查询对应的一组IOT环境信息存证记录
// ========================================================================
func (t *ChinastirlingIOTCertChaincode) queryByMacAddress(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	accountId := strings.TrimSpace(args[0])
	macAddress := strings.TrimSpace(args[1])

	queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"docType\":\"iotcert\"},{\"accountId\":\"%s\"},{\"macAddress\":\"%s\"}]}}", accountId, macAddress)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ========================================================================
// 富查询：基于accountId+macAddress+rawDataPacket查询对应的一组IOT环境信息存证记录
// ========================================================================
func (t *ChinastirlingIOTCertChaincode) queryByRawDataPacket(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	accountId := strings.TrimSpace(args[0])
	macAddress := strings.TrimSpace(args[1])
	rawDataPacket := strings.TrimSpace(args[2])

	queryString := fmt.Sprintf("{\"selector\":{\"$and\":[{\"docType\":\"iotcert\"},{\"accountId\":\"%s\"},{\"macAddress\":\"%s\"},{\"rawDataPacket\":\"%s\"}]}}", accountId, macAddress, rawDataPacket)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ========================================================================
// 富查询：accountId+macAddress+发送序列号start+发送序列号end,查询对应的一组IOT环境信息存证记录
// ========================================================================
func (t *ChinastirlingIOTCertChaincode) queryBySerialNoRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 4 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	accountId := strings.TrimSpace(args[0])
	macAddress := strings.TrimSpace(args[1])
	startDeviceSerialNo,err := strconv.Atoi(args[2])
        if err != nil {
		return shim.Error("3rd argument must be a numeric string")
	}
	endDeviceSerialNo,err := strconv.Atoi(args[3])
        if err != nil {
		return shim.Error("4th argument must be a numeric string")
	}
	queryString := fmt.Sprintf("{\"selector\":{\"$and\":[ {\"docType\":\"iotcert\"},  {\"accountId\":\"%s\"},  {\"macAddress\":\"%s\"}, { {\"deviceSerialNo\":{\"$gte\":\"%d\"}}  ,{\"deviceSerialNo\":{\"$lte\":\"%d\"}} }  ]}}", accountId, macAddress, startDeviceSerialNo, endDeviceSerialNo)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
