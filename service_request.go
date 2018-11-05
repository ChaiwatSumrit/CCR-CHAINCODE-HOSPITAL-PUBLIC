package main

import (
	"encoding/json"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type InvoiceDoubleFinance struct {
}

// Const for create
const TxtCCR = "CCR|"

// Const for request
const TxtCCRR = "CCRR|"

// CCRcertificateReceipt containing all main fields for organization request
type CCRcertificateReceipt struct {
	Asset
	// wait ICE!!!
}

// EndorseInvoiceDoubleFinance containing all main fields for organization request
type EndorseInvoiceDoubleFinanceModel struct {
	Asset
	InvIdentity         string    `json:"InvIdentity" valid:"required"` //Key
	EndorsementIdentity string    `json:"EndorsementIdentit" valid:"required"`
	PublicKey           string    `json:"PublicKey" valid:"required"`
	TimeStamp           time.Time `json:"TimeStamp" valid:"optional"`
	FinanceIdentity     []string  `json:"FinanceIdentity" valid:"optional"`
}

// FinanceInvoiceDoubleFinance containing all main fields for organization request
type FinanceInvoiceDoubleFinanceModel struct {
	Asset
	InvIdentity         string    `json:"FinanceIdentity" valid:"required"`
	EndorsementIdentity string    `json:"EndorsementIdentity" valid:"required"`
	FinanceIdentity     string    `json:"FinanceIdentity" valid:"required"` //Key
	FinanceRunningNo    string    `json:"FinanceRunningNo" valid:"required"`
	InvAmountUsed       float64   `json:"InvAmountUsed" valid:"required"`
	FinanceTime         time.Time `json:"FinanceTime" valid:"required"`
}

// EndorsementIdentity for CCRcertificateReceipt
type EndorsementIdentity struct {
	InvEndorsementIdentity string `json:"InvEndorsementIdentity" valid:"required"`
}

//go:generate enumer -type=State -json
type State int

const (
	NEW State = 1 + iota
	ENDORSED
	FINANCED
)

//Init initializes the InvoiceDoubleFinance model/smart contract
func (t *InvoiceDoubleFinance) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success(nil)

}

//CreateCertificateReceipt
func (t *InvoiceDoubleFinance) CreateCertificateReceipt(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	methodName := "[CreateCertificateReceipt]"
	// Check number of args.
	if len(args) != 2 { //N wait ICE !!
		logger.Errorf(methodName + " Incorrect number of arguments. Expecting 2")
		return shim.Error(methodName + " Incorrect number of arguments. Expecting 2")
		//todo: create custom event
	}
	//Check devorgId Tcert attribute
	devorgId, err := getOrgId(stub)
	if devorgId == "" || err != nil {
		//error getting devorg Id , return error
		logger.Errorf(methodName + " error retrieving devorgId attribute for user")
		// return error
		return shim.Error(methodName + "error retrieving devorgId attribute for user, please check if user was registered with a devorgId or invoke request has attributes")
	}

	var sr CCRcertificateReceipt

	sr = CCRcertificateReceipt{}

	// Validate Invoice Existing
	keyCheck := TxtInv + args[0]
	InvoiceDoubleFinanceJSONBytes, _ := stub.GetState(keyCheck)
	if CCRcertificateReceiptJSONBytes == nil {
		key := TxtInv + sr.InvIdentity
		srJSONBytes, _ := json.Marshal(sr)
		// Write service request to world state]
		err = stub.PutState(key, srJSONBytes)
		if err != nil {
			logger.Errorf(methodName + " PutState error: " + err.Error())
		}
	} else {
		logger.Errorf(methodName + "Invoice has Already Exist")
		return shim.Error(methodName + "Invoice has Already Exist")
	}

	return shim.Success(nil)
} //End of CreateIdfInvoiceDoubleFinance function


//GetInvoiceById - Get a InvoiceDoubleFinance by ID
func (t *InvoiceDoubleFinance) GetInvoiceById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	methodName := "[GetInvoiceById]"
	if len(args) != 1 {
		logger.Errorf(methodName + " Incorrect number of arguments. Expecting 1")
		return shim.Error(methodName + " Incorrect number of arguments. Expecting 1")
	}
	err := validateFirstArgument(args)
	if err != nil {
		logger.Errorf(methodName + " Incorrect format of Invoice ID")
		return shim.Error(err.Error())
	}
	//Check devorgId
	devorgId, err := getOrgId(stub)
	if devorgId == "" || err != nil {
		//error getting devorg Id , return error
		logger.Errorf(methodName + " error retrieving devorgId attribute for user")
		// return error
		return shim.Error(methodName + "error retrieving devorgId attribute for user, please check if user was registered with a devorgId or invoke request has attributes")
	}

	id := args[0]
	key := id

	InvoiceDoubleFinanceJSONBytes, err2 := stub.GetState(key)

	if err2 != nil {
		logger.Errorf(methodName + " Get state error: " + err2.Error())
		return shim.Error(err2.Error())
	}
	if InvoiceDoubleFinanceJSONBytes == nil {
		logger.Errorf(methodName + " No Invoice with " + id + " is found")
		return shim.Error(methodName + " Error: No Invoice found")
	}

	return shim.Success(InvoiceDoubleFinanceJSONBytes)
}

// func (t *InvoiceDoubleFinance) GetAllInvoiceDoubleFinance(stub shim.ChaincodeStubInterface) pb.Response {
// 	methodName := "[GetAllInvoiceDoubleFinance]"

// 	//Check devorgId
// 	devorgId, err := getOrgId(stub)
// 	if devorgId == "" || err != nil {
// 		//error getting devorg Id , return error
// 		logger.Errorf(methodName + " error retrieving devorgId attribute for user")
// 		// return error
// 		return shim.Error(methodName + "error retrieving devorgId attribute for user, please check if user was registered with a devorgId or invoke request has attributes")
// 	}
// 	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"InvoiceDoubleFinance\",\"read\": {\"$all\": [\"%s\"]}}}", devorgId)

// 	queryResults, err := getQueryResultForQueryString(stub, queryString)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	return shim.Success(queryResults)
// }

func validateCreateUpdateInvoiceDoubleFinance(args []string) (InvCreateUpdateInvoiceDoubleFinanceTransaction, error) {
	methodName := "[validateCreateUpdateInvoiceDoubleFinance]"
	var InvoiceDoubleFinanceAsStruct InvCreateUpdateInvoiceDoubleFinanceTransaction

	invAmount, _ := strconv.ParseFloat(args[1], 64)
	invRemainingAmount := invAmount

	InvoiceDoubleFinanceAsStruct = InvCreateUpdateInvoiceDoubleFinanceTransaction{
		InvIdentity:        args[0],
		InvAmount:          invAmount,
		InvRemainingAmount: invRemainingAmount}

	//attempt to validate the whole service request struct. if there are errors, immediately abort
	result, err := govalidator.ValidateStruct(InvoiceDoubleFinanceAsStruct)
	if !result && err != nil {
		logger.Errorf(methodName + " Validation error for Service Request Struct: " + err.Error())
		return InvoiceDoubleFinanceAsStruct, err
	}

	return InvoiceDoubleFinanceAsStruct, nil
}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	methodName := "[getQueryResultForQueryString]"
	logger.Infof(methodName+" getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")
	count := 0
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, errors.New(methodName + " " + err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
		count++
	}
	buffer.WriteString("]")
	// data, _ := govalidator.ToJSON(buffer.String())
	message := `{"count": ` + strconv.Itoa(count) + `, "data":` + buffer.String() + ` }`
	logger.Infof("- getQueryResultForQueryString queryResult:\n%s\n", message)
	return []byte(message), nil
}


// //CreateIdfInvoiceDoubleFinance
// func (t *InvoiceDoubleFinance) CreateIdfInvoiceDoubleFinance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	methodName := "[CreateIdfInvoiceDoubleFinance]"
// 	// Check number of args.
// 	if len(args) != 2 {
// 		logger.Errorf(methodName + " Incorrect number of arguments. Expecting 4")
// 		return shim.Error(methodName + " Incorrect number of arguments. Expecting 4")
// 		//todo: create custom event
// 	}
// 	//Check devorgId Tcert attribute
// 	devorgId, err := getOrgId(stub)
// 	if devorgId == "" || err != nil {
// 		//error getting devorg Id , return error
// 		logger.Errorf(methodName + " error retrieving devorgId attribute for user")
// 		// return error
// 		return shim.Error(methodName + "error retrieving devorgId attribute for user, please check if user was registered with a devorgId or invoke request has attributes")
// 	}

// 	saveDraftSRTx, err := validateCreateUpdateInvoiceDoubleFinance(args)
// 	if err != nil {
// 		return shim.Error(methodName + err.Error())
// 		//todo: create custom event
// 	}

// 	var sr CCRcertificateReceipt

// 	sr = CCRcertificateReceipt{
// 		InvIdentity:        saveDraftSRTx.InvIdentity,
// 		InvAmount:          saveDraftSRTx.InvAmount,
// 		InvState:           NEW,
// 		InvRemainingAmount: saveDraftSRTx.InvRemainingAmount}

// 	// Validate Invoice Existing
// 	keyCheck := TxtInv + args[0]
// 	InvoiceDoubleFinanceJSONBytes, _ := stub.GetState(keyCheck)
// 	if InvoiceDoubleFinanceJSONBytes == nil {
// 		key := TxtInv + sr.InvIdentity
// 		srJSONBytes, _ := json.Marshal(sr)
// 		// Write service request to world state]
// 		err = stub.PutState(key, srJSONBytes)
// 		if err != nil {
// 			logger.Errorf(methodName + " PutState error: " + err.Error())
// 		}
// 	} else {
// 		logger.Errorf(methodName + "Invoice has Already Exist")
// 		return shim.Error(methodName + "Invoice has Already Exist")
// 	}

// 	return shim.Success(nil)
// } //End of CreateIdfInvoiceDoubleFinance function

// //EndorseInvoiceDoubleFinance
// func (t *InvoiceDoubleFinance) EndorseInvoiceDoubleFinance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	methodName := "[EndorseInvoiceDoubleFinance]"

// 	// Check number of args.
// 	if len(args) != 4 {
// 		logger.Errorf(methodName + " Incorrect number of arguments. Expecting 4")
// 		return shim.Error(methodName + " Incorrect number of arguments. Expecting 4")
// 	}

// 	timeStamp, err5 := time.Parse(time.RFC3339, args[3])
// 	if err5 != nil {
// 		logger.Errorf(methodName+" Failed in parsing Date: %s", err5.Error())

// 	}
// 	var sr EndorseInvoiceDoubleFinanceModel

// 	sr = EndorseInvoiceDoubleFinanceModel{
// 		InvIdentity:         args[0],
// 		EndorsementIdentity: args[1],
// 		PublicKey:           args[2],
// 		TimeStamp:           timeStamp}

// 	endorseVerify, endorseVerifyBool := validateDigitalSignature(args[2], args[1], args[0])
// 	if endorseVerifyBool == false {
// 		return shim.Error(methodName + endorseVerify)
// 	}

// 	// Get data from Invoice
// 	keyCheckInv := TxtInv + args[0]
// 	InvInvoiceDoubleFinanceJSONBytes, _ := stub.GetState(keyCheckInv)
// 	InvSr := CCRcertificateReceipt{}
// 	err3 := json.Unmarshal(InvInvoiceDoubleFinanceJSONBytes, &InvSr)
// 	if err3 != nil {
// 		return shim.Error(methodName + " Error unmarshaling InvoiceDoubleFinance:" + err3.Error())
// 	}
// 	InvSr.InvState = ENDORSED
// 	InvSr.EndorsementIdentity = append(InvSr.EndorsementIdentity, args[1])

// 	//Get data from Endorse
// 	keyCheck := TxtEndorse + args[1]
// 	InvoiceDoubleFinanceJSONBytes, _ := stub.GetState(keyCheck)
// 	//EndSr := EndorseInvoiceDoubleFinanceModel{}
// 	// err4 := json.Unmarshal(InvoiceDoubleFinanceJSONBytes, &EndSr)

// 	// keyCheck2 := TxtFinance + EndSr.FinanceIdentity[0]
// 	// FinInvoiceDoubleFinanceJSONBytes, _ := stub.GetState(keyCheck2)
// 	// FinSr := FinanceInvoiceDoubleFinanceModel{}
// 	// err6 := json.Unmarshal(FinInvoiceDoubleFinanceJSONBytes, &FinSr)

// 	// if err4 != nil {
// 	// 	return shim.Error(methodName + " Error unmarshaling InvoiceDoubleFinance:" + err4.Error())
// 	// }

// 	// if err6 != nil {
// 	// 	return shim.Error(methodName + " Error unmarshaling InvoiceDoubleFinance:" + err6.Error())
// 	// }

// 	//Validate EndorseIdentity
// 	if InvoiceDoubleFinanceJSONBytes != nil {
// 		logger.Errorf(methodName + "EndorseIdentity has Already Exist")
// 		return shim.Error(methodName + "EndorseIdentity has Already Exist")
// 	}

// 	// Validate State || FinSr.InvState == FINANCED
// 	// if !(InvSr.InvState == NEW) {
// 	// 	return shim.Error(methodName + " Invalid State")
// 	// }

// 	key := TxtEndorse + sr.EndorsementIdentity
// 	srJSONBytes, _ := json.Marshal(sr)
// 	// Write service request to world state]
// 	err := stub.PutState(key, srJSONBytes)
// 	if err != nil {
// 		logger.Errorf(methodName + " PutState error: ")
// 	}
// 	InvSrJSONBytes, _ := json.Marshal(InvSr)
// 	err1 := stub.PutState(keyCheckInv, InvSrJSONBytes)
// 	if err1 != nil {
// 		logger.Errorf(methodName + " PutState error: ")
// 	}

// 	return shim.Success(nil)
// } //End of EndorseInvoiceDoubleFinance function

// //FinanceIdfInvoiceDoubleFinance
// func (t *InvoiceDoubleFinance) FinanceIdfInvoiceDoubleFinance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	methodName := "[FinanceIdfInvoiceDoubleFinance]"
// 	// Check number of args.
// 	if len(args) != 6 { //args[]
// 		logger.Errorf(methodName + " Incorrect number of arguments. Expecting 6")
// 		return shim.Error(methodName + " Incorrect number of arguments. Expecting 6")
// 		//todo: create custom event
// 	}

// 	//Check devorgId Tcert attribute
// 	devorgId, err := getOrgId(stub)
// 	if devorgId == "" || err != nil {
// 		//error getting devorg Id , return error
// 		logger.Errorf(methodName + " error retrieving devorgId attribute for user")
// 		// return error
// 		return shim.Error(methodName + "error retrieving devorgId attribute for user, please check if user was registered with a devorgId or invoke request has attributes")
// 	}

// 	//Parse VALUES
// 	argsInvAmountUsed, _ := strconv.ParseFloat(args[5], 64)
// 	var sr FinanceInvoiceDoubleFinanceModel

// 	Date, err := time.Parse(time.RFC3339, args[4])
// 	if err != nil {
// 		return shim.Error(methodName + " Failed in parsing Date: " + err.Error())
// 	}
// 	//PUT args
// 	sr = FinanceInvoiceDoubleFinanceModel{
// 		InvIdentity:         args[0],
// 		EndorsementIdentity: args[1],
// 		FinanceIdentity:     args[2],
// 		FinanceRunningNo:    args[3],
// 		FinanceTime:         Date,
// 		InvAmountUsed:       argsInvAmountUsed}

// 	//@1 Use endorsement_identity to get public_key. => SR.public_key
// 	// GET SR1
// 	EnkeyCheck := TxtEndorse + sr.EndorsementIdentity
// 	EnInvoiceDoubleFinanceJSONBytes, _ := stub.GetState(EnkeyCheck)
// 	EnSR := EndorseInvoiceDoubleFinanceModel{}
// 	err1 := json.Unmarshal(EnInvoiceDoubleFinanceJSONBytes, &EnSR)
// 	if err1 != nil {
// 		//error unmarshaling
// 		return shim.Error(methodName + "Error unmarshaling InvoiceDoubleFinance:" + err1.Error())
// 	}
// 	Endorsement_Running := sr.EndorsementIdentity + "|" + sr.FinanceRunningNo

// 	//@2 Verrify
// 	// GET EnSR PublicKey for Endorsement
// 	_, msgBool := validateDigitalSignature(EnSR.PublicKey, sr.FinanceIdentity, Endorsement_Running)
// 	if msgBool != true {
// 		return shim.Error(`ERROR Verify False`)
// 	}

// 	//@4 Check if InvRemainingAmount < inv_amount_use.
// 	//If true, return ERROR (case trying to finance more than remain)
// 	InvkeyCheck := TxtInv + sr.InvIdentity
// 	InvInvoiceDoubleFinanceJSONBytes, _ := stub.GetState(InvkeyCheck)
// 	// GET InvSR
// 	InvSR := CCRcertificateReceipt{}
// 	err2 := json.Unmarshal(InvInvoiceDoubleFinanceJSONBytes, &InvSR)
// 	if err2 != nil {
// 		//error unmarshaling
// 		return shim.Error(methodName + " Error unmarshaling InvoiceDoubleFinance:" + err2.Error())
// 	}

// 	//@3 Check state of invoice_identity, return ERROR if not ok.
// 	// GET EnSR INV STATE for Endorsement
// 	if !(InvSR.InvState == ENDORSED || InvSR.InvState == FINANCED) {
// 		logger.Errorf(methodName + "Invalid State")
// 		return shim.Error(methodName + "Invalid State")
// 	}

// 	// GET InvSR InvRemainingAmount for INV
// 	if InvSR.InvRemainingAmount < sr.InvAmountUsed {
// 		logger.Errorf(methodName + "ERROR case trying to finance more than remaining amount")
// 		return shim.Error("ERROR case trying to finance more than remaining amount")
// 	}

// 	/*@5 Blockchain txn
// 	Add finance _identity and its related fields into blockchain.
// 	Add FinanceIdentity to endorsement data in blockchain.
// 	Update InvRemainingAmount into invoice table.
// 	*/

// 	EnSR.FinanceIdentity = append(EnSR.FinanceIdentity, sr.FinanceIdentity)
// 	InvSR.InvState = FINANCED
// 	InvSR.InvRemainingAmount = InvSR.InvRemainingAmount - sr.InvAmountUsed

// 	FinkeyCheck := TxtFinance + sr.FinanceIdentity
// 	FinInvoiceDoubleFinanceJSONBytes, _ := stub.GetState(FinkeyCheck)
// 	// TxtFinance + sr.FinanceIdentity == nil
// 	if FinInvoiceDoubleFinanceJSONBytes == nil {
// 		srJSONBytes, _ := json.Marshal(sr)
// 		// Write service request to world state]
// 		//@5.1PutState Finance
// 		err := stub.PutState(FinkeyCheck, srJSONBytes)
// 		if err != nil {
// 			logger.Errorf(methodName + " PutState error: ")
// 		}
// 		//@5.2PutState Endorsement
// 		sr2JSONBytes, _ := json.Marshal(EnSR)
// 		err1 := stub.PutState(EnkeyCheck, sr2JSONBytes)
// 		if err1 != nil {
// 			logger.Errorf(methodName + " PutState error: ")
// 		}
// 		//@5.3Update InvRemainingAmount into invoice table.
// 		sr3JSONBytes, _ := json.Marshal(InvSR)
// 		err2 := stub.PutState(InvkeyCheck, sr3JSONBytes)
// 		if err2 != nil {
// 			logger.Errorf(methodName + " PutState error: ")
// 		}

// 	} else {
// 		logger.Errorf(methodName + "Already Exist")
// 		return shim.Error(methodName + "Already Exist")
// 	}
// 	//*end
// 	return shim.Success(nil)
// } //End of FinanceIdfInvoiceDoubleFinance function

// //GetInvoiceById - Get a InvoiceDoubleFinance by ID
// func (t *InvoiceDoubleFinance) GetInvoiceById(stub shim.ChaincodeStubInterface, args []string) pb.Response {
// 	methodName := "[GetInvoiceById]"
// 	if len(args) != 1 {
// 		logger.Errorf(methodName + " Incorrect number of arguments. Expecting 1")
// 		return shim.Error(methodName + " Incorrect number of arguments. Expecting 1")
// 	}
// 	err := validateFirstArgument(args)
// 	if err != nil {
// 		logger.Errorf(methodName + " Incorrect format of Invoice ID")
// 		return shim.Error(err.Error())
// 	}
// 	//Check devorgId
// 	devorgId, err := getOrgId(stub)
// 	if devorgId == "" || err != nil {
// 		//error getting devorg Id , return error
// 		logger.Errorf(methodName + " error retrieving devorgId attribute for user")
// 		// return error
// 		return shim.Error(methodName + "error retrieving devorgId attribute for user, please check if user was registered with a devorgId or invoke request has attributes")
// 	}

// 	id := args[0]
// 	key := id

// 	InvoiceDoubleFinanceJSONBytes, err2 := stub.GetState(key)

// 	if err2 != nil {
// 		logger.Errorf(methodName + " Get state error: " + err2.Error())
// 		return shim.Error(err2.Error())
// 	}
// 	if InvoiceDoubleFinanceJSONBytes == nil {
// 		logger.Errorf(methodName + " No Invoice with " + id + " is found")
// 		return shim.Error(methodName + " Error: No Invoice found")
// 	}

// 	return shim.Success(InvoiceDoubleFinanceJSONBytes)
// }

// // func (t *InvoiceDoubleFinance) GetAllInvoiceDoubleFinance(stub shim.ChaincodeStubInterface) pb.Response {
// // 	methodName := "[GetAllInvoiceDoubleFinance]"

// // 	//Check devorgId
// // 	devorgId, err := getOrgId(stub)
// // 	if devorgId == "" || err != nil {
// // 		//error getting devorg Id , return error
// // 		logger.Errorf(methodName + " error retrieving devorgId attribute for user")
// // 		// return error
// // 		return shim.Error(methodName + "error retrieving devorgId attribute for user, please check if user was registered with a devorgId or invoke request has attributes")
// // 	}
// // 	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"InvoiceDoubleFinance\",\"read\": {\"$all\": [\"%s\"]}}}", devorgId)

// // 	queryResults, err := getQueryResultForQueryString(stub, queryString)
// // 	if err != nil {
// // 		return shim.Error(err.Error())
// // 	}
// // 	return shim.Success(queryResults)
// // }

// func validateCreateUpdateInvoiceDoubleFinance(args []string) (InvCreateUpdateInvoiceDoubleFinanceTransaction, error) {
// 	methodName := "[validateCreateUpdateInvoiceDoubleFinance]"
// 	var InvoiceDoubleFinanceAsStruct InvCreateUpdateInvoiceDoubleFinanceTransaction

// 	invAmount, _ := strconv.ParseFloat(args[1], 64)
// 	invRemainingAmount := invAmount

// 	InvoiceDoubleFinanceAsStruct = InvCreateUpdateInvoiceDoubleFinanceTransaction{
// 		InvIdentity:        args[0],
// 		InvAmount:          invAmount,
// 		InvRemainingAmount: invRemainingAmount}

// 	//attempt to validate the whole service request struct. if there are errors, immediately abort
// 	result, err := govalidator.ValidateStruct(InvoiceDoubleFinanceAsStruct)
// 	if !result && err != nil {
// 		logger.Errorf(methodName + " Validation error for Service Request Struct: " + err.Error())
// 		return InvoiceDoubleFinanceAsStruct, err
// 	}

// 	return InvoiceDoubleFinanceAsStruct, nil
// }

// func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
// 	methodName := "[getQueryResultForQueryString]"
// 	logger.Infof(methodName+" getQueryResultForQueryString queryString:\n%s\n", queryString)

// 	resultsIterator, err := stub.GetQueryResult(queryString)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	// buffer is a JSON array containing QueryRecords
// 	var buffer bytes.Buffer
// 	buffer.WriteString("[")
// 	count := 0
// 	bArrayMemberAlreadyWritten := false
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, errors.New(methodName + " " + err.Error())
// 		}
// 		// Add a comma before array members, suppress it for the first array member
// 		if bArrayMemberAlreadyWritten == true {
// 			buffer.WriteString(",")
// 		}
// 		// Record is a JSON object, so we write as-is
// 		buffer.WriteString(string(queryResponse.Value))
// 		bArrayMemberAlreadyWritten = true
// 		count++
// 	}
// 	buffer.WriteString("]")
// 	// data, _ := govalidator.ToJSON(buffer.String())
// 	message := `{"count": ` + strconv.Itoa(count) + `, "data":` + buffer.String() + ` }`
// 	logger.Infof("- getQueryResultForQueryString queryResult:\n%s\n", message)
// 	return []byte(message), nil
// }
