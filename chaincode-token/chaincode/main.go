package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type TokenChaincode struct {
}

const compositeIndexName = "acct~amt~op~opp~txid"

func main() {
	err := shim.Start(new(TokenChaincode))
	if err != nil {
		fmt.Printf("Error starting Token chaincode: %s", err)
	}
}

func (t *TokenChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// initialize account
	args := []string{"0x00", "100000000"}
	t.supply(stub, args)
	return shim.Success(nil)
}

func (t *TokenChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
	// Handle different functions
	if function == "transfer" {
		return t.transfer(stub, args)
	} else if function == "balanceof" {
		return t.balanceOf(stub, args)
	} else if function == "update" {
		return t.update(stub, args)
	} else if function == "putstandard" {
		return t.putStandard(stub, args)
	} else if function == "getstandard" {
		return t.getStandard(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// Supply the initial amount of tokens to one account
func (t *TokenChaincode) supply(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Check we have a valid number of args
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments, expecting 2")
	}

	// Extract the args
	account := args[0]
	amountStr := args[1]

	// Amount must be a float
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return shim.Error("Transfer amount was not a number")
	}

	// Make sure amount is larger than 0
	if amount <= 0 {
		return shim.Error("Transfer amount must larger than 0")
	}

	// Supply only positive number of tokens
	op := "+"

	// Retrieve info needed for the update procedure
	txid := APIstub.GetTxID()

	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey, compositeErr := APIstub.CreateCompositeKey(compositeIndexName, []string{account, amountStr, op, "_", txid})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", account, compositeErr.Error()))
	}

	// Save the composite key index
	compositePutErr := APIstub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not operate supplyment of %s for %s in the ledger: %s", amountStr, account, compositePutErr.Error()))
	}

	return shim.Success([]byte(fmt.Sprintf("Successfully supplied %s to %s", amountStr, account)))
}

/**
 * Transfer amount of tokens from one account to another. The arguments
 * to give in the args array are as follows:
 *	- args[0] -> Account of transfer from
 *	- args[1] -> Account of transfer to
 *	- args[2] -> amount of tokens to edtransferred (float)
 *
 * @param APIstub The chaincode shim
 * @param args The arguments array for the transfer invocation
 *
 * @return A response structure indicating success or failure with a message
 */
func (t *TokenChaincode) transfer(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Check we have a valid number of args
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments, expecting 3")
	}

	// Extract the args
	from := args[0]
	to := args[1]
	amountStr := args[2]

	// Amount must be a number
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return shim.Error("Transfer amount was not a number")
	}

	// Make sure amount is larger than 0
	if amount <= 0 {
		return shim.Error("Transfer amount must larger than 0")
	}

	// Retrieve info needed for the update procedure
	txid := APIstub.GetTxID()

	// Create the composite keys that will allow us to query for all deltas on a particular variable
	// For from account
	op := "-"
	compositeKey, compositeErr := APIstub.CreateCompositeKey(compositeIndexName, []string{from, amountStr, op, to, txid})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s, %s, %s, %s: %s", from, amountStr, op, to, compositeErr.Error()))
	}

	// Save the composite key index
	compositePutErr := APIstub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not save composite key for %s, %s, %s, %s in the ledger: %s", from, amountStr, op, to, compositePutErr.Error()))
	}

	// For to account
	op = "+"
	compositeKey, compositeErr = APIstub.CreateCompositeKey(compositeIndexName, []string{to, amountStr, op, from, txid})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s, %s, %s, %s: %s", to, amountStr, op, from, compositeErr.Error()))
	}

	// Save the composite key index
	compositePutErr = APIstub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not save composite key for %s, %s, %s, %s in the ledger: %s", from, amountStr, op, to, compositePutErr.Error()))
	}

	return shim.Success([]byte(fmt.Sprintf("Successfully transfer %s from %s to %s", amountStr, from, to)))
}

/**
 * Retrieves the aggregate value of a variable in the ledger. Gets all delta rows for the variable
 * and computes the final value from all deltas. The args array for the invocation must contain the
 * following argument:
 *	- args[0] -> The name of account to get the value of
 *
 * @param APIstub The chaincode shim
 * @param args The arguments array for the get invocation
 *
 * @return A response structure indicating success or failure with a message
 */
func (t *TokenChaincode) balanceOf(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Check we have a valid number of args
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting 1")
	}

	account := strings.TrimSpace(args[0])
	if len(account) == 0 {
		return shim.Error("Account must not be empty")
	}

	// Get all deltas for the variable
	deltaResultsIterator, deltaErr := APIstub.GetStateByPartialCompositeKey(compositeIndexName, []string{account})
	if deltaErr != nil {
		return shim.Error(fmt.Sprintf("Could not retrieve value for %s: %s", account, deltaErr.Error()))
	}
	defer deltaResultsIterator.Close()

	// Check the variable existed
	if !deltaResultsIterator.HasNext() {
		return shim.Error(fmt.Sprintf("No variable by the account %s exists", account))
	}

	// Iterate through result set and compute final value
	var finalVal float64
	var i int
	for i = 0; deltaResultsIterator.HasNext(); i++ {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return shim.Error(nextErr.Error())
		}

		// Split the composite key into its component parts
		_, keyParts, splitKeyErr := APIstub.SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return shim.Error(splitKeyErr.Error())
		}

		// Retrieve the values
		amountStr := keyParts[1]
		op := keyParts[2]

		// Convert the amount string and perform the operation
		amount, convErr := strconv.ParseFloat(amountStr, 64)
		if convErr != nil {
			return shim.Error(convErr.Error())
		}

		switch op {
		case "+":
			finalVal += amount
		case "-":
			finalVal -= amount
		default:
			return shim.Error(fmt.Sprintf("operation was not recognized: %s", op))
		}
	}

	return shim.Success([]byte(strconv.FormatFloat(finalVal, 'f', -1, 64)))
}

/**
 * All functions below this are for testing traditional editing of a single row
 */
func (t *TokenChaincode) putStandard(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	name := args[0]
	valAddStr := args[1]
	fmt.Sprintln("valAddStr: " + valAddStr)
	valAdd, err := strconv.Atoi(valAddStr)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to convert value to int %s: %s", valAddStr, err.Error()))
	}

	// get current value
	valNowByte, getErr := APIstub.GetState(name)
	if getErr != nil {
		return shim.Error(fmt.Sprintf("Failed to retrieve the state of %s: %s", name, getErr.Error()))
	}
	fmt.Sprintln("valNowByte: " + string(valNowByte))

	// can't conver empty string to int
	valNowStr := string(valNowByte)
	if valNowStr == "" {
		valNowStr = "0"
	}

	// convert valNow to int
	valNow, err := strconv.Atoi(valNowStr)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to convert value to int %s: %s", valNowStr, err.Error()))
	}
	fmt.Sprintln("valNow: " + valNowStr)

	// add valNow with valAdd
	valNew := valNow + valAdd
	fmt.Sprintln("valNew: " + string(valNew))

	// upload new value
	putErr := APIstub.PutState(name, []byte(strconv.Itoa(valNew)))
	if putErr != nil {
		return shim.Error(fmt.Sprintf("Failed to put state: %s", putErr.Error()))
	}

	return shim.Success(nil)
}

func (t *TokenChaincode) getStandard(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	name := args[0]

	val, getErr := APIstub.GetState(name)
	if getErr != nil {
		return shim.Error(fmt.Sprintf("Failed to get state: %s", getErr.Error()))
	}

	return shim.Success(val)
}

/**
 * Updates the ledger to include a new delta for a particular variable. If this is the first time
 * this variable is being added to the ledger, then its initial value is assumed to be 0. The arguments
 * to give in the args array are as follows:
 *	- args[0] -> name of the variable
 *	- args[1] -> new delta (float)
 *	- args[2] -> operation (currently supported are addition "+" and subtraction "-")
 *
 * @param APIstub The chaincode shim
 * @param args The arguments array for the update invocation
 *
 * @return A response structure indicating success or failure with a message
 */
func (t *TokenChaincode) update(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	// Check we have a valid number of args
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments, expecting 3")
	}

	// Extract the args
	name := args[0]
	op := args[2]
	_, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return shim.Error("Provided value was not a number")
	}

	// Make sure a valid operator is provided
	if op != "+" && op != "-" {
		return shim.Error(fmt.Sprintf("Operator %s is unrecognized", op))
	}

	// Retrieve info needed for the update procedure
	txid := APIstub.GetTxID()
	compositeIndexName := "varName~op~value~txID"

	// Create the composite key that will allow us to query for all deltas on a particular variable
	compositeKey, compositeErr := APIstub.CreateCompositeKey(compositeIndexName, []string{name, op, args[1], txid})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", name, compositeErr.Error()))
	}

	// Save the composite key index
	compositePutErr := APIstub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return shim.Error(fmt.Sprintf("Could not put operation for %s in the ledger: %s", name, compositePutErr.Error()))
	}

	return shim.Success([]byte(fmt.Sprintf("Successfully added %s%s to %s", op, args[1], name)))
}
