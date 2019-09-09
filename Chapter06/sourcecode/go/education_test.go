package main
import (
              "fmt"
              "testing"
              "encoding/json"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)
func TestEducation(t *testing.T) {
              fmt.Println("Inside TestEducation")
              stub := shim.NewMockStub("mockStub", new(EducationChaincode))
              if stub == nil {
                                           t.Fatalf("MockStub creation failed")
              }
              var key = "std1231"
              // Here we perform a "mock invoke" to invoke the function "insertReceiver" method with associated parameters
              // The first parameter is the function we are invoking
              result := stub.MockInvoke("001",
                                           [][]byte{[]byte("insertReceiver"),
                                                                        []byte(key),
                                                                        []byte("Anand Y"),
                                                                        []byte("Blockchain")})
              // We expect a shim.ok if all goes well
              if result.Status != shim.OK {
                                           t.Fatalf("Expected unauthorized user error to be returned")
              }
              fmt.Println("Result insertReceiver:")
              fmt.Println(result)
              result = stub.MockInvoke("002",
                                           [][]byte{[]byte("queryReceiverById"),
                                                                        []byte(key)})
              // We expect a shim.ok if all goes well
              if result.Status != shim.OK {
                                           t.Fatalf("Expected unauthorized user error to be returned")
              }
              fmt.Println("Result queryReceiverById:")
              receiver := &Receiver{}
              err := json.Unmarshal([]byte(result.Payload), &receiver)
              if err != nil {
                             t.Fatalf(err.Error())
              }
              fmt.Println(receiver)
              certid:="cert123"
              result = stub.MockInvoke("003",
                                           [][]byte{[]byte("insertCertificate"),
                                                                        []byte(certid),
                                                                        []byte("12345"),
                                                                        []byte("ORU Blockchain Certificate"),
                                                                        []byte(key),
                                                                        []byte("ORU"),
                                                                        []byte("IT"),
                                                                        []byte("06/04/2019"),
                                                                        []byte("06/04/2019"),
                                                                        []byte("Blockchain course completed"),
                                                                        []byte(""),
                                                                        []byte(""),
                                                                        []byte("Active")})
              fmt.Println(result)
              if result.Status != shim.OK {
                                           t.Fatalf("Expected unauthorized user error to be returned")
              }
              fmt.Println("Result insertCertificate:")
              fmt.Println(result)
              result = stub.MockInvoke("004",
                                           [][]byte{[]byte("queryCertificateById"),
                                                                        []byte(certid)})
              // We expect a shim.ok if all goes well
              if result.Status != shim.OK {
                                           t.Fatalf("Expected unauthorized user error to be returned")
              }
              fmt.Println("Result queryCertificateById:")
              certificate := &Certificate{}
              err = json.Unmarshal([]byte(result.Payload), &certificate)
              if err != nil {
                             t.Fatalf(err.Error())
              }
              fmt.Println(certificate)
              result = stub.MockInvoke("005",
                                           [][]byte{[]byte("approveCertificate"),
                                                                        []byte(certid),
                                                                        []byte("Approved"),
                                                                        []byte("06/04/2019 10:41:50")})
              fmt.Println(result)
              if result.Status != shim.OK {
                                           t.Fatalf("Expected unauthorized user error to be returned")
              }
              fmt.Println("Result approveCertificate:")
              fmt.Println(result)
              result = stub.MockInvoke("004",
              [][]byte{[]byte("queryCertificateById"),
                                           []byte(certid)})
              // We expect a shim.ok if all goes well
              if result.Status != shim.OK {
                             t.Fatalf("Expected unauthorized user error to be returned")
              }
              fmt.Println("Result queryCertificateById:")
                certificate = &Certificate{}
              err = json.Unmarshal([]byte(result.Payload), &certificate)
              if err != nil {
              t.Fatalf(err.Error())
              }
              fmt.Println(certificate) 
}