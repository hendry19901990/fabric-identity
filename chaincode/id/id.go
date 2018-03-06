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
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type AttestClaim struct{
 claim map[string]map[string]string `json:"claim"`
}

type Credential struct {
   Token    string `json:"token"`
   ValidDay int `json:"validDay"`
}

type ID struct {
    Claims       map[string]string `json:"claims"`
    Infoshared   map[string]map[string]Credential `json:"infoshared"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryClaimsById" {
		return s.queryClaimsById(APIstub, args)
	} else if function == "createId" {
		return s.createId(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) createId(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	id := ID {Claims: make(map[string]string), Infoshared: make(map[string]map[string]Credential)}
	id.Claims["fullname"] = args[1]
  id.Claims["docid"]    = args[2]

	idAsBytes, _ := json.Marshal(id)
	APIstub.PutState(args[0], idAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryClaimsById(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

  id := ID{}
	idAsBytes, ok := APIstub.GetState(args[0])

	if ok {
		  json.Unmarshal(idAsBytes, &id)
			jsonData, _ := json.Marshal(id.Claims)
      jsonResp := "{\"user\": \"" +  args[0] + "\", \"claims\":" + string(jsonData) + "}"

			return shim.Success([]byte(jsonResp))
	}else{
      return shim.Error("User not exist! :(")
	}
}



// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
