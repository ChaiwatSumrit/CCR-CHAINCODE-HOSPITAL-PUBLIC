package main

import (

	//fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"

	//"strconv"
	"os"
)

var logger = shim.NewLogger("idf-chaincode")

type IdfChaincode struct {
	serviceRequest ServiceRequest
	//usersAndOrgs   UsersAndOrgs
}

//=================================================================================================================================
//  Main - main - Starts up the chaincode
//=================================================================================================================================

func main() {

	// LogDebug, LogInfo, LogNotice, LogWarning, LogError, LogCritical (Default: LogDebug)
	logger.SetLevel(shim.LogDebug)

	logLevel, _ := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL"))
	shim.SetLoggingLevel(logLevel)

	err := shim.Start(new(IdfChaincode))
	if err != nil {
		//fmt.Printf("Error starting IdfChaincode: %s", err)
	}
}

//==============================================================================================================================
//  Init Function - Called when the user deploys the chaincode
//==============================================================================================================================

func (t *IdfChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	t.serviceRequest.Init(stub)
	//t.usersAndOrgs.Init(stub)
	return shim.Success(nil)
}

//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Passes the
//  		 initial arguments passed are passed on to the called function.
//==============================================================================================================================

func (t *IdfChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	logger.Infof("CCR is running " + function)

	if function == "init" {
		return t.Init(stub)
	} else if function == "create_ccr" {
	}

	return shim.Error("Received unknown ccr function name " + function)
}
