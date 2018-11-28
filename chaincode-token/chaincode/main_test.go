package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

func checkGet(t *testing.T, stub *shim.MockStub, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("get"), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Get", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Get", name, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println("Get value", name, "was not", value, "as expected, but was ", string(res.Payload))
		t.FailNow()
	}
}

func checkGetStandard(t *testing.T, stub *shim.MockStub, name string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("getstandard"), []byte(name)})
	if res.Status != shim.OK {
		fmt.Println("Get", name, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Get", name, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println("Get value", name, "was not", value, "as expected, but was ", string(res.Payload))
		t.FailNow()
	}
}

func checkBalanceOf(t *testing.T, stub *shim.MockStub, account string, value string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("balanceof"), []byte(account)})
	if res.Status != shim.OK {
		fmt.Println("balanceof", account, "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Get", account, "failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != value {
		fmt.Println("Get value", account, "was not", value, "as expected, but was ", string(res.Payload))
		t.FailNow()
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func Test_Init(t *testing.T) {
	token := new(TokenChaincode)
	stub := shim.NewMockStub("token", token)

	checkInit(t, stub, [][]byte{[]byte("init")})
	checkBalanceOf(t, stub, "0x00", "100000000")
}

func Test_Update(t *testing.T) {
	token := new(TokenChaincode)
	stub := shim.NewMockStub("token", token)

	checkInvoke(t, stub, [][]byte{[]byte("update"), []byte("mike"), []byte("1"), []byte("+")})
	checkGet(t, stub, "mike", "1")
}

func Test_Transfer(t *testing.T) {
	token := new(TokenChaincode)
	stub := shim.NewMockStub("token", token)

	checkInvoke(t, stub, [][]byte{[]byte("transfer"), []byte("0x00"), []byte("tom"), []byte("1")})
	checkBalanceOf(t, stub, "tom", "1")
}

func Test_Multi_Updates(t *testing.T) {
	token := new(TokenChaincode)
	stub := shim.NewMockStub("token", token)

	checkInvoke(t, stub, [][]byte{[]byte("update"), []byte("mike"), []byte("2"), []byte("+")})
	checkInvoke(t, stub, [][]byte{[]byte("update"), []byte("mike"), []byte("2"), []byte("+")})
	checkGet(t, stub, "mike", "4")
}

func Test_PutStandard(t *testing.T) {
	token := new(TokenChaincode)
	stub := shim.NewMockStub("token", token)

	checkInvoke(t, stub, [][]byte{[]byte("putstandard"), []byte("mike"), []byte("1")})
	checkGetStandard(t, stub, "mike", "1")
}
