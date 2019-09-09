package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

//  Chaincode implementation
type EducationChaincode struct {

}

//  receiver/student struct
type Receiver struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Receiver_id string `json:"receiver_id"`
	Receiver_name string `json:"receiver_name"`
	Upload_org string `json:"upload_org"`
}

//  certificate data struct
type Certificate struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Cert_id       string  `json:"cert_id"`
	Cert_no       string  `json:"cert_no"`
	Cert_name          string `json:"cert_name"`
	Cert_receiver     string `json:"cert_receiver"`      // student name
	Cert_receiver_id   string `json:"cert_receiver_id"`  // student id
	Cert_issuer        string `json:"cert_issuer"`       // org name
	Cert_industry     string `json:"cert_industry"`
	Cert_create_time	      string `json:"cert_create_time"`
	Cert_update_time       string `json:"cert_update_time"`
	Cert_remark         string `json:"cert_remark"`
	Cert_url_image      string `json:"cert_url_image"`
	Cert_status string `json:"cert_status"`
}



// main - Start execution

func main() {
	err := shim.Start(new(EducationChaincode))
	if err != nil {
		fmt.Printf("Error starting Xebest Trace chaincode: %s", err)
	}
}

// Init initializes chaincode
func (t *EducationChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}



// Invoke - Invoking user transactions
func (t *EducationChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "insertReceiver" { //create a new Receiver or student
		return t.insertReceiver(stub, args)
	} else if function == "queryReceiverById" { // query a receiver by id, stupid name - -!
		return t.queryReceiverById(stub,args)
	} else if function == "insertCertificate" { //insert a cert
		return t.insertCertificate(stub, args)
	} else if function == "queryCertificateById" { // query a certificate
		return t.queryCertificateById(stub, args)
	} else if function == "getRecordHistory"{ //query hisitory of one key for the record
		return t.getRecordHistory(stub,args)
	} else if function == "queryAllCertificates"{ // query all of all students
		return t.queryAllCertificates(stub,args)
	} else if function == "approveCertificate" { // change status
		return t.approveCertificate(stub,args) 
	}else if function == "deleteRecord" { // delete student or certificate
		return t.deleteRecord(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// initReceiver - insert a new Receiver into chaincode state
func (t *EducationChaincode) insertReceiver(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	fmt.Println("start insert receiver")
	
	receiver_id := args[0]
	receiver_name := args[1]
	upload_org := args[2]

	// Check if the receiver already exists with the id
	receiverAsBytes, err := stub.GetState(receiver_id)
	if err != nil {
		return shim.Error("Failed to get receiver: " + err.Error())
	} else if receiverAsBytes != nil {
		fmt.Println("This receiver already exists: " + receiver_id)
		return shim.Error("This receiver already exists: " + receiver_id)
	}

	// Create receiver object and marshal to JSON 
	objectType := "receiver"
	receiver := &Receiver{objectType, receiver_id, receiver_name,upload_org}
	receiverJSONasBytes, err := json.Marshal(receiver)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	fmt.Println("receiver: ")
	fmt.Println(receiver)
	// Save the receiver to ledger state
	err = stub.PutState(receiver_id, receiverJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// receiver saved and indexed. Return success
	fmt.Println("End init receiver")
	return shim.Success(nil)

}



// queryReceiverById - read data for the given receiver from the chaincode state
func (t *EducationChaincode) queryReceiverById(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var recev_id, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting receiver_id to query")
	}

	recev_id = args[0]
	//Read the Receiver from the chaincode state
	valAsbytes, err := stub.GetState(recev_id) 
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + recev_id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"receiver does not exist: " + recev_id + "\"}"
		return shim.Error(jsonResp)
	}	
	return shim.Success(valAsbytes)
}

// insertCertificate - insert a new certificate information into the ledger state
func (t *EducationChaincode) insertCertificate(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	
	if len(args) != 11 {
		return shim.Error("Incorrect number of arguments. expecting 11 args")
	}

	cert_id := args[0]
	cert_no := args[1]
	cert_name := args[2]	
	cert_receiver_id := args[3]
	cert_issuer := args[4]
	cert_industry := args[5]
	cert_create_time := args[6]
	cert_update_time := args[7]
	cert_remark := args[8]
	cert_url_image := args[9]
	cert_status := args[10]


	// check if receiver exists
	ReceAsBytes, err := stub.GetState(cert_receiver_id)
	if err != nil {
		return shim.Error("Failed to get Receiver:" + cert_receiver_id + "," + err.Error())
	} else if ReceAsBytes == nil {
	  fmt.Println("Receiver does not exist with id: " + cert_receiver_id )
		return shim.Error("Receiver does not exist with id: " + cert_receiver_id )
	}

	//Fetch receiver name from the state
	receiver := &Receiver{}
	err = json.Unmarshal([]byte(ReceAsBytes), &receiver)
	if err != nil {
		return shim.Error(err.Error())
	}
	cert_receiver :=receiver.Receiver_name;
	fmt.Println("cert_receiver: "+cert_receiver)

	objectType := "certificate"
	certificate := &Certificate{objectType,cert_id,cert_no,cert_name,cert_receiver,cert_receiver_id,cert_issuer,cert_industry,cert_create_time,cert_update_time,cert_remark,cert_url_image,cert_status}
	certificateJSONasBytes, err := json.Marshal(certificate)
	if err != nil {
		return shim.Error(err.Error())
	}

	// insert the certificate into the ledger
	err = stub.PutState(cert_id, certificateJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// certificate saved - Return success	
	return shim.Success(nil)
}


// queryCertificateById - read a certificate by given id from the ledger state
func (t *EducationChaincode) queryCertificateById(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var cert_id, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting id of the certificate to query")
	}

	cert_id = args[0]
	//Read the certificate from chaincode state
	valAsbytes, err := stub.GetState(cert_id) 
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + cert_id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"certificate does not exist: " + cert_id + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// approveCertificate - approve the certificate by authority
func (t *EducationChaincode) approveCertificate(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	var err error
	// check args
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}

	cert_id := args[0]
	status := args[1]
	update_time := args[2]

	//Read certificate details from the ledger
	valAsbytes, err := stub.GetState(cert_id) 
		if err != nil {
			return shim.Error(err.Error())
		} else if valAsbytes == nil {
			return shim.Error("certificate not exist")
		}
	
	certificate := &Certificate{}
	err = json.Unmarshal([]byte(valAsbytes), &certificate)
	if err != nil {
		return shim.Error(err.Error())
	}
	certificate.Cert_status = status
	certificate.Cert_update_time = update_time

	valAsbytes, err = json.Marshal(certificate)
	if err != nil {
		return shim.Error(err.Error())
	}
	//Update the certificate in the ledger
	err = stub.PutState(cert_id, valAsbytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//Publishing a custom event
	var testEventValue []byte
	testEventValue=[]byte("Certificate "+cert_id+" status is changed to "+status)
	stub.SetEvent("testEvent",testEventValue)
	
	return shim.Success(nil)

}

// queryAllCertificates - Query all certificates from the ledger state
func (t *EducationChaincode) queryAllCertificates(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	queryString := "{\"selector\":{\"docType\":\"certificate\"}}"

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}


// getRecordHistory - Fetches the historical state transitions for a given key of a record
func (t *EducationChaincode) getRecordHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting an id of Receiver or Certificate")
	}

	recordKey := args[0]

	fmt.Printf("Fetching history for record: %s\n", recordKey)

	resultsIterator, err := stub.GetHistoryForKey(recordKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the key/value pair
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON goods)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("Result of getHistoryForRecord :\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}



// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("getQueryResultForQueryString queryString:\n%s\n", queryString)

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
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}


// deleteRecord - Mark the record deleted by given key

func (t *EducationChaincode) deleteRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	
	if len(args) != 1{
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	id := args[0]
	err := stub.DelState(id)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	return shim.Success(nil)
}
