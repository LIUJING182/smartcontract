/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"sync"
	"github.com/Nik-U/pbc"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Car struct {
	Base   string `json:"Base"`
	Num  string `json:"Num"`
	Rc  string `json:"Rc"`
}

/*
 * The Init method is called when the Smart Contract "log" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "log"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryCar" {
		return s.queryCar(APIstub, args)
	}  else if function == "createCar" {
		return s.createCar(APIstub, args)
	} else if function == "queryAllCars" {
		return s.queryAllCars(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	carAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(carAsBytes)
}


func (s *SmartContract) createCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	// string to []byte
	intBase := []byte(args[1])
	fmt.Println(args[1])
	intNum := []byte(args[2])
	fmt.Println(intBase)
	fmt.Println(intNum)

	N := new(big.Int)
	_, err := fmt.Sscan(args[3], N)
	if err != nil {
		fmt.Println("error scanning value:", err)
	} else {
		fmt.Println(N)
	}
	params := pbc.GenerateA1(N)
	// create a new pairing with given params
	pairing := pbc.NewPairing(params)
	G1 := pairing.NewG1()
	gsk := G1.NewFieldElement()
	gsk.SetBytes(intBase)
	csk := G1.NewFieldElement()
	csk.SetBytes(intNum)
	fmt.Println("gsk")
	fmt.Println(gsk)
	fmt.Println("csk")
	fmt.Println(csk)

	var tableG1 sync.Map
	messageSpace := big.NewInt(10000)
	bound := int64(math.Ceil(math.Sqrt(float64(messageSpace.Int64()))))+1
	aux1 := gsk.NewFieldElement()
	aux1.Set(gsk)

	for j := int64(0); j <= bound; j++ {
		tableG1.Store(aux1.String(), j)
		fmt.Println("j")
		fmt.Println(j)
		fmt.Println(aux1)
		aux1.Mul(aux1, gsk)
	}

	aux := csk.NewFieldElement()
	gamma := gsk.NewFieldElement()

	aux.Set(csk)
	aux.Mul(aux, gamma)

	gamma.Set(gsk)
	gamma.MulBig(gamma, big.NewInt(bound))

	var val *big.Int
	var found bool
	var myStr string

	for i := int64(0); i <= bound; i++ {

		found = false
		val = big.NewInt(0)
		fmt.Println(i)
		value, hit := tableG1.Load(aux.String())
		if v, ok := value.(int64); ok {
			val = big.NewInt(v)
			found = hit
		}

		if found {
			dl := big.NewInt(i*bound + val.Int64() + 1)
			myStr = dl.String()
			fmt.Println(myStr)
			break
		}
		aux.Div(aux, gamma)
	}

	var car = Car{Base: args[1], Num: args[2], Rc: myStr}
	carAsBytes, _ := json.Marshal(car)
	APIstub.PutState(args[0], carAsBytes)
	return shim.Success(nil)

}

func (s *SmartContract) queryAllCars(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "CAR0"
	endKey := "CAR999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("{\"Key\":")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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

	fmt.Printf("- hello:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
