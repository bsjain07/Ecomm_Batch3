package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type buyer_seller struct {
}

type Product_Asset struct {
	Product_Id          string `json:"Product_Id"`
	Product_Name        string `json:"Product_Name"`
	Product_Description string `json:"Product_Description"`
	Product_Price       int    `json:"Product_Price"`
	Product_Quantity    int    `json:"Product_Quantity"`
}

type Order_Asset struct {
	Order_Id           string `json:"Order_Id"`
	Order_Receiver     string `json:"Order_Receiver"`
	Order_ProductId    string `json:"Order_ProductId"`
	Order_Quantity     int    `json:"Order_Quantity"`
	Order_Status       string `json:"Order_Status"`
	Order_CreationDate string `json:"Order_CreationDate"`
}

type CounterNO struct {
	Counter int `json:"counter"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(buyer_seller))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}

}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *buyer_seller) Init(APIstub shim.ChaincodeStubInterface) pb.Response {

	// Initializing Product Counter
	ProductCounterBytes, _ := APIstub.GetState("ProductCounterNO")
	if ProductCounterBytes == nil {
		var ProductCounter = CounterNO{Counter: 0}
		ProductCounterBytes, _ := json.Marshal(ProductCounter)
		err := APIstub.PutState("ProductCounterNO", ProductCounterBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to Intitate Product Counter"))
		}
	}

	// Initializing Order Counter
	OrderCounterBytes, _ := APIstub.GetState("OrderCounterNO")
	if OrderCounterBytes == nil {
		var OrderCounter = CounterNO{Counter: 0}
		OrderCounterBytes, _ := json.Marshal(OrderCounter)
		err := APIstub.PutState("OrderCounterNO", OrderCounterBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("Failed to Intitate Order Counter"))
		}
	}

	return shim.Success(nil)

}

///Start of Private Function

//getCounter to the latest value of the counter based on the Asset Type provided as input parameter
func getCounter(APIstub shim.ChaincodeStubInterface, AssetType string) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	fmt.Sprintf("Counter Current Value %d of Asset Type %s", counterAsset.Counter, AssetType)

	return counterAsset.Counter
}

//incrementCounter to the increase value of the counter based on the Asset Type provided as input parameter by 1
func incrementCounter(APIstub shim.ChaincodeStubInterface, AssetType string) int {
	counterAsBytes, _ := APIstub.GetState(AssetType)
	counterAsset := CounterNO{}

	json.Unmarshal(counterAsBytes, &counterAsset)
	counterAsset.Counter++
	counterAsBytes, _ = json.Marshal(counterAsset)

	err := APIstub.PutState(AssetType, counterAsBytes)
	if err != nil {

		fmt.Sprintf("Failed to Increment Counter")

	}
	return counterAsset.Counter
}

///End of Private Function

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *buyer_seller) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("function is ==> :" + function)
	//action := args[0]
	//fmt.Println(" action is ==> :" + action)
	fmt.Println(args)

	if function == "queryAsset" {
		return t.queryAsset(stub, args)
	} else if function == "queryAllAsset" {
		return t.queryAllAsset(stub, args)
	} else if function == "getHistoryForRecord" {
		return t.getHistoryForRecord(stub, args)
	} else if function == "createProduct" {
		return t.createProduct(stub, args)
	} else if function == "updateProduct" {
		return t.updateProduct(stub, args)
	} else if function == "createOrder" {
		return t.createOrder(stub, args)
	} else if function == "updateOrderStatus" {
		return t.updateOrderStatus(stub, args)
	} else if function == "query" {
		return t.Query(stub, args)
	}
	//else if action == "history" {
	// 	return t.GetHistory(stub, args)
	// }

	fmt.Println("invoke did not find func: " + function) //error

	return shim.Error("Received unknown function")
}

// ===== Example: Ad hoc rich query ========================================================
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *buyer_seller) Query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	//fmt.Println("GetQueryResultForQueryString() : getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()
	fmt.Println(resultsIterator)
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
		queryResponseStr := string(queryResponse.Value)
		fmt.Println(queryResponseStr)
		buffer.WriteString(queryResponseStr)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	//fmt.Println("GetQueryResultForQueryString(): getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

//E-com

// GetTxTimestampChannel Function gets the Transaction time when the chain code was executed it remains same on all the peers where chaincode executes
func (t *buyer_seller) GetTxTimestampChannel(APIstub shim.ChaincodeStubInterface) (string, error) {
	txTimeAsPtr, err := APIstub.GetTxTimestamp()
	if err != nil {
		fmt.Printf("Returning error in TimeStamp \n")
		return "Error", err
	}
	fmt.Printf("\t returned value from APIstub: %v\n", txTimeAsPtr)
	timeStr := time.Unix(txTimeAsPtr.Seconds, int64(txTimeAsPtr.Nanos)).String()

	return timeStr, nil
}

// queryAsset Function gets the assets based on Id provided as input
func (t *buyer_seller) queryAsset(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, Required 1")
	}

	fmt.Println("In Query Asset")

	AssetAsBytes, _ := APIstub.GetState(args[0])
	fmt.Println(AssetAsBytes, "In Query Asset")
	if AssetAsBytes == nil {
		return shim.Error("Could not locate Asset")

	}

	return shim.Success(AssetAsBytes)
}

// create Product
func (t *buyer_seller) createProduct(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	//To check number of arguments are 4
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments, Required 4 arguments")
	}

	//To check each argument is not null
	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			return shim.Error(string(i+1) + "st argument must be a non-empty string")
		}
	}

	i1, errPrice := strconv.Atoi(args[2])
	if errPrice != nil {
		return shim.Error(fmt.Sprintf("Failed to Convert Price: %s", errPrice))
		//fmt.Println(i1)
	}

	i2, errQuantity := strconv.Atoi(args[3])
	if errQuantity != nil {
		return shim.Error(fmt.Sprintf("Failed to Convert Quantity: %s", errQuantity))
		//fmt.Println(i2)
	}

	productCounter := getCounter(APIstub, "ProductCounterNO")
	productCounter++

	var comAsset = Product_Asset{Product_Id: "Product" + strconv.Itoa(productCounter), Product_Name: args[0], Product_Description: args[1], Product_Price: i1, Product_Quantity: i2}

	comAssetAsBytes, errMarshal := json.Marshal(comAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in Product: %s", errMarshal))
	}

	errPut := APIstub.PutState(comAsset.Product_Id, comAssetAsBytes)

	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to create Product Asset: %s", comAsset.Product_Id))
	}

	//TO Increment the Product Counter
	incrementCounter(APIstub, "ProductCounterNO")

	fmt.Println("Success in creating Product Asset %v", comAsset)

	return shim.Success(nil)

}

// update Product Attributes
func (t *buyer_seller) updateProduct(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments, Required 6")
	}

	if len(args[1]) == 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	productBytes, _ := APIstub.GetState(args[1])

	if productBytes == nil {
		return shim.Error("Cannot Find Product Asset ")
	}

	i1, errPrice := strconv.Atoi(args[4])
	if errPrice != nil {
		return shim.Error(fmt.Sprintf("Failed to Convert Price: %s", errPrice))
		//fmt.Println(i1)
	}

	i2, errQuantity := strconv.Atoi(args[5])
	if errQuantity != nil {
		return shim.Error(fmt.Sprintf("Failed to Convert Quantity: %s", errQuantity))
		//fmt.Println(i2)
	}

	var comAsset = Product_Asset{Product_Id: args[1], Product_Name: args[2], Product_Description: args[3], Product_Price: i1, Product_Quantity: i2}

	comAssetAsBytes, errMarshal := json.Marshal(comAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error: %s", errMarshal))
	}

	errPut := APIstub.PutState(comAsset.Product_Id, comAssetAsBytes)

	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to Update Product: %s", comAsset.Product_Id))
	}

	fmt.Println("Success in updating Product Asset %v ", comAsset)

	return shim.Success(nil)
}

// create Order
func (t *buyer_seller) createOrder(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	//To check number of arguments are 3
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments, Required 3 arguments")
	}

	//To check each argument is not null
	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			return shim.Error(string(i+1) + "st argument must be a non-empty string")
		}
	}

	i1, errQuantity := strconv.Atoi(args[2])
	if errQuantity != nil {
		return shim.Error(fmt.Sprintf("Failed to Convert Quantity: %s", errQuantity))
		//fmt.Println(i2)
	}

	// TO Get Product details
	productBytes, _ := APIstub.GetState(args[1])

	if productBytes == nil {
		return shim.Error("Cannot Find Product Asset ")
	}

	productAsset := Product_Asset{}

	json.Unmarshal(productBytes, &productAsset)

	// To check if the quantity is available
	if i1 > productAsset.Product_Quantity {
		return shim.Error("Quantity Requested is more than the Available Product Quantity")
	}

	productAsset.Product_Quantity = productAsset.Product_Quantity - i1

	updateProductBytes, errMarshal := json.Marshal(productAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error: %s", errMarshal))
	}

	errPutProduct := APIstub.PutState(productAsset.Product_Id, updateProductBytes)

	if errPutProduct != nil {
		return shim.Error(fmt.Sprintf("Failed to Update Product: %s", productAsset.Product_Id))
	}

	orderCounter := getCounter(APIstub, "OrderCounterNO")
	orderCounter++

	//To Get the transaction TimeStamp from the Channel Header
	txTimeAsPtr, errTx := t.GetTxTimestampChannel(APIstub)
	if errTx != nil {
		return shim.Error("Returning error in Transaction TimeStamp")
	}

	var comAsset = Order_Asset{Order_Id: "Order" + strconv.Itoa(orderCounter), Order_Receiver: args[0], Order_ProductId: args[1], Order_Quantity: i1, Order_Status: "PLACED", Order_CreationDate: txTimeAsPtr}

	comAssetAsBytes, errMarshal := json.Marshal(comAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error in Order: %s", errMarshal))
	}

	errPut := APIstub.PutState(comAsset.Order_Id, comAssetAsBytes)

	if errPut != nil {
		return shim.Error(fmt.Sprintf("Failed to create Order Asset: %s", comAsset.Order_Id))
	}

	//TO Increment the Product Counter
	incrementCounter(APIstub, "OrderCounterNO")

	fmt.Println("Success in creating Order Asset %v", comAsset)

	return shim.Success(nil)

}

// update Order Status
func (t *buyer_seller) updateOrderStatus(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments, Required 2")
	}

	if len(args[0]) == 0 {
		return shim.Error("1st argument must be a non-empty string")
	}

	if len(args[1]) == 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}

	orderBytes, _ := APIstub.GetState(args[0])

	if orderBytes == nil {
		return shim.Error("Cannot Find Order Asset ")
	}

	orderAsset := Order_Asset{}

	json.Unmarshal(orderBytes, &orderAsset)

	orderAsset.Order_Status = args[1]

	updateOrderBytes, errMarshal := json.Marshal(orderAsset)

	if errMarshal != nil {
		return shim.Error(fmt.Sprintf("Marshal Error: %s", errMarshal))
	}

	errPutOrder := APIstub.PutState(orderAsset.Order_Id, updateOrderBytes)

	if errPutOrder != nil {
		return shim.Error(fmt.Sprintf("Failed to Update Order: %s", orderAsset.Order_Id))
	}

	return shim.Success(nil)
}

// query all assets
func (t *buyer_seller) queryAllAsset(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	startKey := ""

	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)

	if err != nil {

		return shim.Error(err.Error())

	}

	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults

	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()
		// respValue := string(queryResponse.Value)
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

	fmt.Printf("- queryAllAssets:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())

}

// get History For Record
func (t *buyer_seller) getHistoryForRecord(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	recordKey := args[0]

	fmt.Printf("- start getHistoryForRecord: %s\n", recordKey)

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
		//as-is (as the Value itself a JSON vehiclePart)
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

	fmt.Printf("- getHistoryForRecord returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}