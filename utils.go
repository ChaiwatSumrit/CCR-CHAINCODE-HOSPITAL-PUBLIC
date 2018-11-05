package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/xeipuuv/gojsonschema"
)

type pkcs1PublicKey struct {
	N *big.Int
	E int
}

func validateFirstArgument(args []string) error {
	methodName := "[validateFirstArgument]"
	// GetLockerById id := args[0] //id generated from node.js bridge TODO ensure format is uuid v4
	// GetLockersByOwnerId ownerId := args[0]
	// GetLockersByUserId userId := args[0]
	// Check for nullity
	if &args[0] == nil {
		logger.Infof(methodName + " Null field detected. Field index: " + strconv.Itoa(0))
		return errors.New(methodName + " Null field detected. Field index: " + strconv.Itoa(0))
	}
	// Check for empty
	if args[0] == "" {
		logger.Infof(methodName + " Empty field detected. Field index: " + strconv.Itoa(0))
		return errors.New(methodName + " Empty field detected. Field index: " + strconv.Itoa(0))
	}
	return nil
}

func convertStubTimeToRFC3339(nowSeconds int64, nowNano int32) (string, time.Time) {
	methodName := "[convertStubTimeToRFC3339]"
	tnan := int64(nowNano)
	currentTime := time.Unix(nowSeconds, tnan)
	currentTimeRFC3339 := currentTime.Format(time.RFC3339)
	logger.Debugf(methodName + "Converted Shim Time to: " + currentTimeRFC3339)
	currentTimeRFC3339Time, _ := time.Parse(time.RFC3339, currentTimeRFC3339)
	return currentTimeRFC3339, currentTimeRFC3339Time
}

//check on cert
func getOrgId(stub shim.ChaincodeStubInterface) (string, error) {
	// //Check devorgId Tcert attribute
	// // GetCreator returns marshaled serialized identity of the client
	serializedID, _ := stub.GetCreator()

	// /* CODE TO RETRIEVE CREATOR IDENTITY IN BYTES AND WRITE TO FILE
	// logger.Debugf("----------DEBUG---------")
	// logger.Debugf(string(serializedID))
	// err0 := ioutil.WriteFile("/tmp/dat1", serializedID, 0644)
	// if err0 != nil {
	// 	return "", fmt.Errorf("Error writing file %s", err0)
	// }
	// */

	sID := &msp.SerializedIdentity{}
	err := proto.Unmarshal(serializedID, sID)
	if err != nil {
		return "", fmt.Errorf("Could not deserialize a SerializedIdentity, err %s", err)
	}
	bl, _ := pem.Decode(sID.IdBytes)
	if bl == nil {
		return "", fmt.Errorf("Failed to decode PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return "", fmt.Errorf("Unable to parse certificate %s", err)
	}
	fullCommonName := strings.Split(cert.Subject.CommonName, "-")
	devorgID := fullCommonName[0]
	logger.Debugf("devorgID: %s", devorgID)
	return string(devorgID), err
}

//retrieve enrollment id of user
func getEnrollmentID(stub shim.ChaincodeStubInterface) (string, error) {
	// // GetCreator returns marshaled serialized identity of the client
	serializedID, _ := stub.GetCreator()
	sID := &msp.SerializedIdentity{}
	err := proto.Unmarshal(serializedID, sID)
	if err != nil {
		return "", fmt.Errorf("Could not deserialize a SerializedIdentity, err %s", err)
	}
	bl, _ := pem.Decode(sID.IdBytes)
	if bl == nil {
		return "", fmt.Errorf("Failed to decode PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return "", fmt.Errorf("Unable to parse certificate %s", err)
	}
	enrollmentID := string(cert.Subject.CommonName)

	logger.Debugf("enrollmentId: %s", enrollmentID)
	return enrollmentID, err
}

//validateJSONSchema : Validate JSON Schema
func validateJSONSchema(schema string, jsonObj string) error {
	methodName := "[validateJSONSchema]"
	loaderSchema := gojsonschema.NewStringLoader(schema)
	loaderObj := gojsonschema.NewStringLoader(jsonObj)
	result, _ := gojsonschema.Validate(loaderSchema, loaderObj)
	if !result.Valid() {
		logger.Errorf(methodName + " The document is not valid. see errors :")
		errMsg := ""
		for _, err := range result.Errors() {
			errMsg = errMsg + " " + err.String()
		}
		return errors.New(methodName + " Can't be Validate JSON Schema. Error: " + errMsg)
	}
	return nil
}

func resolveUnmarshalErr(data []byte, err error) string {
	if e, ok := err.(*json.UnmarshalTypeError); ok {
		// grab stuff ahead of the error
		var i int
		for i = int(e.Offset) - 1; i != -1 && data[i] != '\n' && data[i] != ','; i-- {
		}
		info := strings.TrimSpace(string(data[i+1 : int(e.Offset)]))
		s := fmt.Sprintf("%s - at: %s", e.Error(), info)
		return s
	}
	if e, ok := err.(*json.UnmarshalFieldError); ok {
		return e.Error()
	}
	if e, ok := err.(*json.InvalidUnmarshalError); ok {
		return e.Error()
	}
	return err.Error()
}

// copy from go version 1.10 because go 1.9.x don't have tihs function
func ParsePKCS1PublicKey(der []byte) (*rsa.PublicKey, error) {
	var pub pkcs1PublicKey
	rest, err := asn1.Unmarshal(der, &pub)
	if err != nil {
		return nil, err
	}
	if len(rest) > 0 {
		return nil, asn1.SyntaxError{Msg: "trailing data"}
	}

	if pub.N.Sign() <= 0 || pub.E <= 0 {
		return nil, errors.New("x509: public key contains zero or negative value")
	}
	if pub.E > 1<<31-1 {
		return nil, errors.New("x509: public key contains large public exponent")
	}

	return &rsa.PublicKey{
		E: pub.E,
		N: pub.N,
	}, nil
}

func validateDigitalSignature(publicKeyArg string, dataSignatureArg string, infoDataArg string) (string, bool) {
	//methodName := "[validateDigitalSignature]"

	fmt.Println("Arg1 :", publicKeyArg)
	fmt.Println("Arg2 :", dataSignatureArg)
	fmt.Println("Arg3 :", infoDataArg)
	var msg string
	var msgBool bool
	pemString := publicKeyArg
	block, _ := pem.Decode([]byte(pemString))
	publicKey, _ := ParsePKCS1PublicKey(block.Bytes)

	fmt.Println("**********************************")
	fmt.Println("*********  Public Key  ***********")
	fmt.Println("**********************************")
	fmt.Println("Public Key  :", publicKey)

	dataSignature := dataSignatureArg
	stringSignature := string(dataSignature)
	Signature, err := base64.StdEncoding.DecodeString(stringSignature)
	fmt.Println("**********************************")
	fmt.Println("******  Signature Detail  ********")
	fmt.Println("**********************************")
	fmt.Printf("Signature =  %x\nError = %s\n", Signature, err)

	infoData := infoDataArg
	//infoData := args[2]
	hashed := sha256.Sum256([]byte(infoData))
	errv := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], Signature)
	fmt.Println("**********************************")
	fmt.Println("******  Data Detail  ********")
	fmt.Println("**********************************")
	fmt.Println("Data: ", hashed)
	//errv := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, []byte(infoData[:]), Signature)
	if errv != nil {
		msg = "Verify False"
		msgBool = false
	} else {
		msg = "Verify True"
		msgBool = true
	}
	fmt.Println(msg)
	return msg, msgBool
}
