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
	"crypto/rand"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type AttestClaim struct{
 Claim map[string]map[string]string `json:"claim"`
}

type Credential struct {
   Token    string `json:"token"`
   ValidDay int `json:"validDay"`
}

type ID struct {
    Claims       map[string]string `json:"claims"`
    Infoshared   map[string]map[string]Credential `json:"infoshared"`
}

const REQUEST = "requestAttest_"
const ATTEST = "attester_"

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "Identity"
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
	} else if function == "requestAttestation" {
		return s.requestAttestation(APIstub, args)
	} else if function == "createAttestion" {
		return s.createAttestion(APIstub, args)
	} else if function == "shareinfo" {
		return s.shareinfo(APIstub, args)
	} else if function == "queryRequestAttestation" {
		return s.queryRequestAttestation(APIstub, args)
	} else if function == "queryAttestation" {
		return s.queryAttestation(APIstub, args)
	} else if function == "removeUser" {
		return s.removeUser(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * SHAREINFORMATION
 * args: 0 => (idClient), 1 => (attester), 2 => (ClaimName), 3 => (token), 4 => (validDays)
 */
func (s *SmartContract) shareinfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
	   return shim.Error("Incorrect number of arguments. Expecting 5")
  }

	//get state of AttesterRequest and validate if existe user or not
  idAsBytes, ok := APIstub.GetState(args[0])

	if ok {
		id := ID{}
	  json.Unmarshal(idAsBytes, &id)
    // if not exist then create the object
		if id.Infoshared == nil {
			id.Infoshared = make(map[string]map[string]Credential)
		}
    // if not exist then create the object key
		_, ok2 := id.Infoshared[args[1]]
		if !ok2 {
			id.Infoshared[args[1]] = make(map[string]Credential)
		}
    // parse to integer the validDays
    validDays, _ := strconv.Atoi(args[4])
		id.Infoshared[args[1]][args[2]] = Credential {Token: args[3], ValidDay: validDays}
    //parse to bytes and save state
		idAsBytes, _ := json.Marshal(id)
		APIstub.PutState(args[0], idAsBytes)

		return shim.Success(nil)
	} else {
		return shim.Error("User Not Found!")
	}

}

/*
 * QUERY ALL ATTESTATION
 * args: 0 => (idAttester)
 */
func (s *SmartContract) queryAttestation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
	  return shim.Error("Incorrect number of arguments. Expecting 1")
  }
	// index of state all ATTESTATION
  attestationsIndex  := ATTEST  + args[0]
  // get the state
	requestAsBytes, ok := APIstub.GetState(attestationsIndex)
	if ok {
      requestObj := AttestClaim{Claim: make(map[string]map[string]string)}
			json.Unmarshal(requestAsBytes, &requestObj)
			if requestObj.Claim == nil {
          badResult := []byte(`{"error": "There are not request Attestations"}`)
					return shim.Success(badResult)
			}
			// buffer is a JSON array containing QueryResults
			var buffer bytes.Buffer
			i := 0
			buffer.WriteString("[")
			for user, _ := range requestObj.Claim {
			   jsonData, _ := json.Marshal(requestObj.Claim[user])
         jsonResp := "{\"user\": \"" +  user + "\", \"claims\":" + string(jsonData) + "}"
				 buffer.WriteString(jsonResp)
         if i < len(requestObj.Claim)-1 {
					 buffer.WriteString(", ")
				 }
				 i++
		  }
			buffer.WriteString("]")
			fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

			return shim.Success(buffer.Bytes())
	} else {
		 badResult := []byte(`{"error": "There are not request Attestations"}`)
		 return shim.Success(badResult)
	}
}

/*
 * SAVE ATTESTATION
 * args: 0 => (idAttester), 1 => (idClient), 2 => (ClaimName), 3 => (hashClaim)
 */
func (s *SmartContract) createAttestion(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
  // index to save and search into state
  attestationsIndex  := ATTEST  + args[0]
	idAttesterRequest  := REQUEST + args[0]
	//get state of AttesterRequest
  requestAsBytes, ok := APIstub.GetState(idAttesterRequest)
  requestObj := AttestClaim{Claim: make(map[string]map[string]string)}
  //if exist key decode Json to Object
	if ok {
     json.Unmarshal(requestAsBytes, &requestObj)
	}else {
		return shim.Error("Request of Attestation not Found!")
	}
  // get state of attestations
  attestationsAsBytes, ok2 := APIstub.GetState(attestationsIndex)
  attestationsObj := AttestClaim{Claim: make(map[string]map[string]string)}
  //if exist key decode Json to Object
  if ok2 {
		json.Unmarshal(attestationsAsBytes, &attestationsObj)
	}
  // if no exist then create the mapping
	if attestationsObj.Claim == nil{
		attestationsObj.Claim = make(map[string]map[string]string)
	}
  // valid if existe the Claim mapping
	_, exist = attestationsObj.Claim[args[1]]
	if !exist {
		attestationsObj.Claim[args[1]] = make(map[string]string)
	}
	//set the hash of the attestation
	attestationsObj.Claim[args[1]][args[2]] = args[3]
  //remove the key of that attestation
	delete(requestObj.Claim[args[1]], args[2])
	requestAsBytesFinal, _ := json.Marshal(requestObj)
	APIstub.PutState(idAttesterRequest, requestAsBytesFinal)
  //save the new attestation
	attestationsAsBytesFinal, _ := json.Marshal(attestationsObj)
	APIstub.PutState(attestationsIndex, attestationsAsBytesFinal)

  return shim.Success(nil)
}

/*
 * QUERY ALL REQUEST ATTESTATION
 * args: 0 => (idAttester)
 */
func (s *SmartContract) queryRequestAttestation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
	  return shim.Error("Incorrect number of arguments. Expecting 1")
  }
	// index of state REQUEST ATTESTATION
  idAttesterRequest := REQUEST + args[0]
  // get the state
	requestAsBytes, ok := APIstub.GetState(idAttesterRequest)
	if ok {
      requestObj := AttestClaim{Claim: make(map[string]map[string]string)}
			json.Unmarshal(requestAsBytes, &requestObj)
			if requestObj.Claim == nil {
          badResult := []byte(`{"error": "There are not request Attestations"}`)
					return shim.Success(badResult)
			}
			// buffer is a JSON array containing QueryResults
			var buffer bytes.Buffer
			i := 0
			buffer.WriteString("[")
			for user, _ := range requestObj.Claim {
			   jsonData, _ := json.Marshal(requestObj.Claim[user])
         jsonResp := "{\"user\": \"" +  user + "\", \"claims\":" + string(jsonData) + "}"
				 buffer.WriteString(jsonResp)
         if i < len(requestObj.Claim)-1 {
					 buffer.WriteString(", ")
				 }
				 i++
		  }
			buffer.WriteString("]")
			fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

			return shim.Success(buffer.Bytes())
	} else {
		 badResult := []byte(`{"error": "There are not request Attestations"}`)
		 return shim.Success(badResult)
	}
}

/*
 * REQUEST ATTESTATION
 * args: 0 => (idAttester), 1 => (idClient), 2 => (ClaimName), 3 => (ClaimUrl)
 */
func (s *SmartContract) requestAttestation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
  // index to save and search into state
  idAttesterRequest := REQUEST + args[0]
	requestAsBytes, ok := APIstub.GetState(idAttesterRequest)
  requestObj := AttestClaim{Claim: make(map[string]map[string]string)}

	if ok {
     json.Unmarshal(requestAsBytes, &requestObj)
     if requestObj.Claim == nil{
			 requestObj.Claim = make(map[string]map[string]string)
		 }

		 _, exist = requestObj.Claim[args[1]]
		 if !exist {
			 requestObj.Claim[args[1]] = make(map[string]string)
		 }

		 requestObj.Claim[args[1]][args[2]] = args[3]
	}else{
     requestObj.Claim[args[1]] = make(map[string]string)
		 requestObj.Claim[args[1]][args[2]] = args[3]
	}

  requestAsBytesFinal, _ := json.Marshal(requestObj)
	APIstub.PutState(idAttesterRequest, requestAsBytesFinal)

	return shim.Success(nil)
}

/*
 * Remove user from State
 * Args: 0 => "userid or hashId"
 */
func (s *SmartContract) removeUser(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
  // search the user by id
	idAsBytes, ok := APIstub.GetState(args[0])
  // if it's found then remove it else return error
	if ok {
     APIstub.DelState(args[0])
		 return shim.Error("User not exist! :(")
	} else {
     return shim.Error("User not exist! :(")
	}
}

/*
 *  create a User Identity
 *  Args: 0 => "userid or hashId", 1 => "fullname", 2 => "docid"
 */
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

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
  for i := 1; i < 10; i++ {
		 u := pseudo_uuid()
     id: ID {Claims: map[string]string{"fullname": "name"+strconv.Itoa(i), "docid": u}, Infoshared: make(map[string]map[string]Credential)}
		 id.Infoshared["fullname"] = make(map[string]Credential)
     id.Infoshared["fullname"]["GOOGLE"]    = Credential{Token: "token1", ValidDay: 30}
     id.Infoshared["fullname"]["FACEBOOK"]  = Credential{Token: "token2", ValidDay: 60}

		 idAsBytes, _ := json.Marshal(id)
	 	 APIstub.PutState("ID"+strconv.Itoa(i), idAsBytes)
	}
  return shim.Success(nil)
}

// generate random id to test
func pseudo_uuid() (uuid string) {

    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        fmt.Println("Error: ", err)
        return
    }

    uuid = fmt.Sprintf("%X%X%X%X%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

    return
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
