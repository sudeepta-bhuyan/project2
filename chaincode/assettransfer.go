package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "strings"

  "github.com/hyperledger/fabric/core/chaincode/shim"
  pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleAssetTransferChaincode struct {
}

type asset struct {
  Name       string `json:"name"`
  Owner      string `json:"owner"`
}

func main() {
  err := shim.Start(new(SimpleAssetTransferChaincode))
  if err != nil {
    fmt.Printf("Error starting SimpleAssetTransfer chaincode: %s", err)
  }
}

func (t *SimpleAssetTransferChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
  return shim.Success(nil)
}

func (t *SimpleAssetTransferChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
  function, args := stub.GetFunctionAndParameters()
  fmt.Println("Invoking " + function)

  if function == "createAsset" {
    return t.createAsset(stub, args)
  } else if function == "deleteAsset" {
    return t.deleteAsset(stub, args)
  } else if function == "transferAsset" {
    return t.transferAsset(stub, args)
  } else if function == "getAsset" {
    return t.getAsset(stub, args)
  }
  fmt.Println("Could did not find function: " + function)
  return shim.Error("Unknown Function")
}

func (t *SimpleAssetTransferChaincode) createAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  var err error

  if len(args) != 2 {
    return shim.Error("Incorrect number of arguments. Expecting 2")
  }

  if len(args[0]) <= 0 {
    return shim.Error("1st argument must be a non-empty string")
  }
  if len(args[1]) <= 0 {
    return shim.Error("2nd argument must be a non-empty string")
  }

  assetName := args[0]
  assetOwner := strings.ToLower(args[1])

  assetAsBytes, err := stub.GetState(assetName)
  if err != nil {
    return shim.Error("Failed to get asset: " + err.Error())
  } else if assetAsBytes != nil {
    fmt.Println("This asset already exists: " + assetName)
    return shim.Error("This asset already exists: " + assetName)
  }

  asset := &asset{assetName, color, size, owner}
  assetJSONasBytes, err := json.Marshal(asset)
  if err != nil {
    return shim.Error(err.Error())
  }

  err = stub.PutState(assetName, assetJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }

  return shim.Success(nil)
}

func (t *SimpleAssetTransferChaincode) deleteAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Expecting 1")
  }
  assetName := args[0]

  assetAsbytes, err := stub.GetState(assetName)
  if err != nil {
    return shim.Error("Failed to get state for " + assetName)
  } else if assetAsbytes == nil {
    return shim.Error(assetName + " does not exist")
  }

  err = stub.DelState(assetName)
  if err != nil {
    return shim.Error("Failed to delete state:" + err.Error())
  }
  return shim.Success(nil)
}

func (t *SimpleAssetTransferChaincode) transferAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  if len(args) < 2 {
    return shim.Error("Incorrect number of arguments. Expecting 2")
  }

  assetName := args[0]
  newOwner := strings.ToLower(args[1])

  assetAsBytes, err := stub.GetState(assetName)
  if err != nil {
    return shim.Error("Failed to get asset:" + err.Error())
  } else if assetAsBytes == nil {
    return shim.Error("Asset does not exist")
  }

  assetToTransfer := asset{}
  err = json.Unmarshal(assetAsBytes, &assetToTransfer)
  if err != nil {
    return shim.Error(err.Error())
  }
  assetToTransfer.Owner = newOwner

  assetJSONasBytes, _ := json.Marshal(assetToTransfer)
  err = stub.PutState(assetName, assetJSONasBytes)
  if err != nil {
    return shim.Error(err.Error())
  }
  return shim.Success(nil)
}

func (t *SimpleAssetTransferChaincode) getAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
  var err error

  if len(args) != 1 {
    return shim.Error("Incorrect number of arguments. Expecting name of the asset to query")
  }

  name = args[0]
  assetAsbytes, err := stub.GetState(name)
  if err != nil {
    return shim.Error("Failed to get asset:" + err.Error())
  } else if assetAsbytes == nil {
    return shim.Error("Asset does not exist")
  }

  return shim.Success(assetAsbytes)
}
