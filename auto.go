package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

const index = "color~name"

// Define a new struct for the BookingConfirmation
type BookingConfirmation struct {
	DocType                       string `json:"docType"`
	BookingRequestID              string `json:"BookingRequestID"`
	TransportPlan                 string `json:"TransportPlan"`
	SpaceAllocation               string `json:"SpaceAllocation"`
	EmptyEquipmentRelease         string `json:"EmptyEquipmentRelease"`
	CarrierHaulageOrder           string `json:"CarrierHaulageOrder"`
	AdditionalServices            string `json:"AdditionalServices"`
	PricingConfirmation           string `json:"PricingConfirmation"`
	RestrictedPartyScreeningCheck string `json:"RestrictedPartyScreeningCheck"`
	ID                            string `json:"ID"`
}


type BookingRequest struct {
	DocType            string `json:"docType"`
	ShipperDetails     string `json:"ShipperDetails"`
	ConsigneeDetails   string `json:"ConsigneeDetails"`
	Destination        string `json:"Destination"`
	Commodity          string `json:"Commodity"`
	ContainerSize      string `json:"ContainerSize"`
	ContainerType      string `json:"ContainerType"`
	ContractPricing    string `json:"ContractPricing"`
	VesselIntermodal   string `json:"VesselIntermodal"`
	SpecialCargo       string `json:"SpecialCargo"`
	ID                 string `json:"ID"`
}

// InitLedger adds a base set of booking requests to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	bookingRequests := []BookingRequest{
		{
			ID:                "booking1",
			ShipperDetails:    "Shipper 1 Details",
			ConsigneeDetails:  "Consignee 1 Details",
			Destination:       "Destination 1",
			Commodity:         "Commodity 1",
			ContainerSize:     "Size 1",
			ContainerType:     "Type 1",
			ContractPricing:   "Pricing 1",
			VesselIntermodal:  "Vessel 1",
			SpecialCargo:      "Special Cargo 1",
		},
		// Add more booking requests here
	}

	for _, bookingRequest := range bookingRequests {
		bookingRequestJSON, err := json.Marshal(bookingRequest)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(bookingRequest.ID, bookingRequestJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)
		}
	}

	return nil
}

// CreateBookingRequest issues a new booking request to the world state with given details.
func (s *SmartContract) CreateBookingRequest(ctx contractapi.TransactionContextInterface, id string, shipperDetails string, consigneeDetails string, destination string, commodity string, containerSize string, containerType string, contractPricing string, vesselIntermodal string, specialCargo string) error {
	exists, err := s.BookingRequestExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the booking request %s already exists", id)
	}

	bookingRequest := BookingRequest{
		ID:               id,
		ShipperDetails:   shipperDetails,
		ConsigneeDetails: consigneeDetails,
		Destination:      destination,
		Commodity:        commodity,
		ContainerSize:    containerSize,
		ContainerType:    containerType,
		ContractPricing:  contractPricing,
		VesselIntermodal: vesselIntermodal,
		SpecialCargo:     specialCargo,
	}
	bookingRequestJSON, err := json.Marshal(bookingRequest)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, bookingRequestJSON)
}

// BookingRequestExists returns true when a booking request with the given ID exists in the world state
func (s *SmartContract) BookingRequestExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	bookingRequestJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return bookingRequestJSON != nil, nil
}

// ReadBookingRequest returns the booking request stored in the world state with the given ID.
func (s *SmartContract) ReadBookingRequest(ctx contractapi.TransactionContextInterface, id string) (*BookingRequest, error) {
	bookingRequestJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if bookingRequestJSON == nil {
		return nil, fmt.Errorf("the booking request %s does not exist", id)
	}

	var bookingRequest BookingRequest
	err = json.Unmarshal(bookingRequestJSON, &bookingRequest)
	if err != nil {
		return nil, err
	}

	return &bookingRequest, nil
}

// UpdateBookingRequest updates an existing booking request in the world state with provided parameters.
func (s *SmartContract) UpdateBookingRequest(ctx contractapi.TransactionContextInterface, id string, shipperDetails string, consigneeDetails string, destination string, commodity string, containerSize string, containerType string, contractPricing string, vesselIntermodal string, specialCargo string) error {
	exists, err := s.BookingRequestExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the booking request %s does not exist", id)
	}

	// Overwrite the original booking request with the new booking request
	bookingRequest := BookingRequest{
		ID:               id,
		ShipperDetails:   shipperDetails,
		ConsigneeDetails: consigneeDetails,
		Destination:      destination,
		Commodity:        commodity,
		ContainerSize:    containerSize,
		ContainerType:    containerType,
		ContractPricing:  contractPricing,
		VesselIntermodal: vesselIntermodal,
		SpecialCargo:     specialCargo,
	}
	bookingRequestJSON, err := json.Marshal(bookingRequest)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, bookingRequestJSON)
}

// GetAllBookingRequests returns all booking requests found in the world state
func (s *SmartContract) GetAllBookingRequests(ctx contractapi.TransactionContextInterface) ([]*BookingRequest, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var bookingRequests []*BookingRequest
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var bookingRequest BookingRequest
		err = json.Unmarshal(queryResponse.Value, &bookingRequest)
		if err != nil {
			return nil, err
		}
		bookingRequests = append(bookingRequests, &bookingRequest)
	}

	return bookingRequests, nil
}

// Function for the carrier to validate and confirm a booking request
func (s *SmartContract) ConfirmBookingRequest(ctx contractapi.TransactionContextInterface, bookingRequestID string, transportPlan string, spaceAllocation string, emptyEquipmentRelease string, carrierHaulageOrder string, additionalServices string, pricingConfirmation string, restrictedPartyScreeningCheck string) error {
	// Check if the booking request exists
	exists, err := s.BookingRequestExists(ctx, bookingRequestID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the booking request %s does not exist", bookingRequestID)
	}

    // You can perform compliance checks here (e.g., check against a compliance database)
    complianceCheckPassed := performComplianceCheck(ctx, bookingRequestID)
    if !complianceCheckPassed {
        return fmt.Errorf("compliance check failed for booking request %s", bookingRequestID)
    }

    // Perform a credit check (simulated with a simple function)
    creditCheckPassed := performCreditCheck(ctx, bookingRequestID)
    if !creditCheckPassed {
        return fmt.Errorf("credit check failed for booking request %s", bookingRequestID)
    }

    // Check for special cargo requirements
    specialCargoRequirementsMet := checkSpecialCargoRequirements(ctx, bookingRequestID)
    if !specialCargoRequirementsMet {
        return fmt.Errorf("special cargo requirements not met for booking request %s", bookingRequestID)
    }

    // Check and request missing or incorrect information from the customer
    missingInfo := checkAndRequestMissingInfo(ctx, bookingRequestID)
    if len(missingInfo) > 0 {
        // Notify the customer of missing information
        notifyCustomer(ctx, bookingRequestID, missingInfo)
        return fmt.Errorf("missing or incorrect information for booking request %s", bookingRequestID)
    }
	// Create a booking confirmation
	bookingConfirmation := BookingConfirmation{
		ID:                    bookingRequestID + "-confirmation",
		BookingRequestID:      bookingRequestID,
		TransportPlan:         transportPlan,
		SpaceAllocation:       spaceAllocation,
		EmptyEquipmentRelease: emptyEquipmentRelease,
		CarrierHaulageOrder:   carrierHaulageOrder,
		AdditionalServices:    additionalServices,
		PricingConfirmation:   pricingConfirmation,
		RestrictedPartyScreeningCheck: restrictedPartyScreeningCheck,
	}

	// Store the booking confirmation in the world state
	confirmationJSON, err := json.Marshal(bookingConfirmation)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(bookingConfirmation.ID, confirmationJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state: %v", err)
	}

	// You can perform additional business logic or checks here before confirming the booking

	return nil
}


// Simulated function to perform compliance check
func performComplianceCheck(ctx contractapi.TransactionContextInterface, bookingRequestID string) bool {
    // Implement your compliance check logic here
    // Return true if the check passes, false otherwise
    return true
}

// Simulated function to perform credit check
func performCreditCheck(ctx contractapi.TransactionContextInterface, bookingRequestID string) bool {
    // Implement your credit check logic here
    // Return true if the check passes, false otherwise
    return true
}

// Simulated function to check special cargo requirements
func checkSpecialCargoRequirements(ctx contractapi.TransactionContextInterface, bookingRequestID string) bool {
    // Implement your special cargo check logic here
    // Return true if requirements are met, false otherwise
    return true
}

// Simulated function to check and request missing information
func checkAndRequestMissingInfo(ctx contractapi.TransactionContextInterface, bookingRequestID string) []string {
    // Implement your logic to check for missing or incorrect information
    // If missing or incorrect information is found, request it from the customer
    missingInfo := []string{"Missing info 1", "Missing info 2"}
    return missingInfo
}

// Simulated function to notify the customer of missing information
func notifyCustomer(ctx contractapi.TransactionContextInterface, bookingRequestID string, missingInfo []string) {
    // Implement your logic to notify the customer (e.g., send an email)
    // Include details about the missing information in the notification
}

// Function to get booking confirmation details
func (s *SmartContract) GetBookingConfirmation(ctx contractapi.TransactionContextInterface, confirmationID string) (*BookingConfirmation, error) {
	confirmationJSON, err := ctx.GetStub().GetState(confirmationID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if confirmationJSON == nil {
		return nil, fmt.Errorf("the booking confirmation %s does not exist", confirmationID)
	}

	var bookingConfirmation BookingConfirmation
	err = json.Unmarshal(confirmationJSON, &bookingConfirmation)
	if err != nil {
		return nil, err
	}

	return &bookingConfirmation, nil
}

// Add other functions as needed for querying and history tracking.

// // Asset describes basic details of what makes up a simple asset
// // Insert struct field in alphabetic order => to achieve determinism across languages
// // golang keeps the order when marshal to json but doesn't order automatically
// type Asset struct {
// 	DocType            string `json:"docType"`
// 	AppraisedValue     int    `json:"AppraisedValue"`
// 	Color              string `json:"Color"`
// 	ID                 string `json:"ID"`
// 	Owner              string `json:"Owner"`
// 	Size               int    `json:"Size"`
// }

// // InitLedger adds a base set of assets to the ledger
// func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
// 	assets := []Asset{
// 		{ID: "asset1", Color: "blue", Size: 5, Owner: "Tomoko", AppraisedValue: 300},
// 		{ID: "asset2", Color: "red", Size: 5, Owner: "Brad", AppraisedValue: 400},
// 		{ID: "asset3", Color: "green", Size: 10, Owner: "Jin Soo", AppraisedValue: 500},
// 		{ID: "asset4", Color: "yellow", Size: 10, Owner: "Max", AppraisedValue: 600},
// 		{ID: "asset5", Color: "black", Size: 15, Owner: "Adriana", AppraisedValue: 700},
// 		{ID: "asset6", Color: "white", Size: 15, Owner: "Michel", AppraisedValue: 800},
// 	}

// 	for _, asset := range assets {
// 		assetJSON, err := json.Marshal(asset)
// 		if err != nil {
// 			return err
// 		}

// 		err = ctx.GetStub().PutState(asset.ID, assetJSON)
// 		if err != nil {
// 			return fmt.Errorf("failed to put to world state. %v", err)
// 		}
// 	}

// 	return nil
// }

// // CreateAsset issues a new asset to the world state with given details.
// func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
// 	exists, err := s.AssetExists(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	if exists {
// 		return fmt.Errorf("the asset %s already exists", id)
// 	}

// 	asset := Asset{
// 		ID:             id,
// 		Color:          color,
// 		Size:           size,
// 		Owner:          owner,
// 		AppraisedValue: appraisedValue,
// 	}
// 	assetJSON, err := json.Marshal(asset)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(id, assetJSON)
// }


// // AssetExists returns true when asset with given ID exists in world state
// func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
// 	assetJSON, err := ctx.GetStub().GetState(id)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to read from world state: %v", err)
// 	}

// 	return assetJSON != nil, nil
// }


// // ReadAsset returns the asset stored in the world state with given id.
// func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
// 	assetJSON, err := ctx.GetStub().GetState(id)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read from world state: %v", err)
// 	}
// 	if assetJSON == nil {
// 		return nil, fmt.Errorf("the asset %s does not exist", id)
// 	}

// 	var asset Asset
// 	err = json.Unmarshal(assetJSON, &asset)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &asset, nil
// }

// // UpdateAsset updates an existing asset in the world state with provided parameters.
// func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
// 	exists, err := s.AssetExists(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	if !exists {
// 		return fmt.Errorf("the asset %s does not exist", id)
// 	}

// 	// overwriting original asset with new asset
// 	asset := Asset{
// 		ID:             id,
// 		Color:          color,
// 		Size:           size,
// 		Owner:          owner,
// 		AppraisedValue: appraisedValue,
// 	}
// 	assetJSON, err := json.Marshal(asset)
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.GetStub().PutState(id, assetJSON)
// }

// // GetAllAssets returns all assets found in world state
// func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
// 	// range query with empty string for startKey and endKey does an
// 	// open-ended query of all assets in the chaincode namespace.
// 	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	var assets []*Asset
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}

// 		var asset Asset
// 		err = json.Unmarshal(queryResponse.Value, &asset)
// 		if err != nil {
// 			return nil, err
// 		}
// 		assets = append(assets, &asset)
// 	}

// 	return assets, nil
// }


// // getQueryResultForQueryStringWithPagination executes the passed in query string with
// // pagination info. The result set is built and returned as a byte array containing the JSON results.
// func getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

// 	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	assets, err := constructQueryResponseFromIterator(resultsIterator)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &PaginatedQueryResult{
// 		Records:             assets,
// 		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
// 		Bookmark:            responseMetadata.Bookmark,
// 	}, nil
// }

// func (t *SimpleChaincode) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]HistoryQueryResult, error) {
// 	log.Printf("GetAssetHistory: ID %v", assetID)

// 	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	var records []HistoryQueryResult
// 	for resultsIterator.HasNext() {
// 		response, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}

// 		var asset Asset
// 		if len(response.Value) > 0 {
// 			err = json.Unmarshal(response.Value, &asset)
// 			if err != nil {
// 				return nil, err
// 			}
// 		} else {
// 			asset = Asset{
// 				ID: assetID,
// 			}
// 		}

// 		timestamp, err := ptypes.Timestamp(response.Timestamp)
// 		if err != nil {
// 			return nil, err
// 		}

// 		record := HistoryQueryResult{
// 			TxId:      response.TxId,
// 			Timestamp: timestamp,
// 			Record:    &asset,
// 			IsDelete:  response.IsDelete,
// 		}
// 		records = append(records, record)
// 	}

// 	return records, nil
// }