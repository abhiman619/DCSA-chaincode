package main

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
    "github.com/google/uuid"
    "reflect"
    "strings"
    "time"
)

// BookingChaincode represents the smart contract.
type BookingChaincode struct {
    contractapi.Contract
}

// BookingRequest represents the booking request data.
type BookingRequest struct {
    ID                             string       `json:"ID"`
    CarrierBookingRequestReference string       `json:"carrierBookingRequestReference"`
    DocumentStatus                 string       `json:"documentStatus"`
    BookingRequestCreatedDateTime  string       `json:"bookingRequestCreatedDateTime"`
    BookingRequestUpdatedDateTime  string       `json:"bookingRequestUpdatedDateTime"`
    TransportDocumentIssuer        Carrier      `json:"carrierDetails"`
    ShipperEntity                  Shipper      `json:"shipperEntity"`
    Consignee                      Consignee    `json:"consignee"`
    PlaceOfReceipt                 PlaceOfReceipt `json:"placeOfReceipt"`
    PlaceOfDelivery                string       `json:"placeOfDelivery"`
    ServiceType                    string       `json:"serviceType"`
    CargoMovementOrigin            string       `json:"cargoMovementSrc"`
    CargoMovementDestination       string       `json:"cargoMovementDst"`
    Commodity                      Commodity    `json:"commodity"`
    CargoGrossWeight               float64      `json:"cargoGrossWeight"`
    ContainerTypeSize              string       `json:"containerType"`
    CreditAmount                   int          `json:"creditAmount"`
}

// Carrier represents the carrier details.
type Carrier struct {
    IssuerID    string `json:"IssuerID"`
    EntityName  string `json:"EntityName"`
    Address     string `json:"Address"`
    Phone       string `json:"phone"`
    EmailOrFax  string `json:"emailOrFax"`
}

// Shipper represents the shipper details.
type Shipper struct {
    CompanyName     string `json:"companyName"`
    PhysicalAddress string `json:"physicalAddress"`
    ContactName     string `json:"contactName"`
    EmailOrFax      string `json:"emailOrFax"`
    Phone           string `json:"phone"`
    LEIOrTaxID      string `json:"leiOrTaxID"`
}

// Consignee represents the consignee details.
type Consignee struct {
    CompanyName       string `json:"companyName"`
    PhysicalAddress   string `json:"physicalAddress"`
    ContactName       string `json:"contactName"`
    EmailOrFax        string `json:"emailOrFax"`
    Phone             string `json:"phone"`
    LEIOrTaxID        string `json:"leiOrTaxID"`
    ToOrderIdentifier string `json:"toOrderIdentifier"`
}

// PlaceOfReceipt represents the place of receipt details.
type PlaceOfReceipt struct {
    CompanyName       string `json:"companyName"`
    PhysicalAddress   string `json:"physicalAddress"`
    ContactName       string `json:"contactName"`
    EmailOrFax        string `json:"emailOrFax"`
    Phone             string `json:"phone"`
    LEIOrTaxID        string `json:"leiOrTaxID"`
    ToOrderIdentifier string `json:"toOrderIdentifier"`
}

// Commodity represents commodity details.
type Commodity struct {
    CommodityName string `json:"commodityName"`
    Description   string `json:"description"`
    Quantity      int    `json:"quantity"`
}

// InitLedger initializes the ledger with sample data.
func (s *BookingChaincode) InitLedger(ctx contractapi.TransactionContextInterface) error {
    bookingRequests := []BookingRequest{
        {
            ID: uuid.New().String(),
            CarrierBookingRequestReference: uuid.New().String(),
            DocumentStatus:                 "RECE",
            BookingRequestCreatedDateTime:  time.Now().Format(time.RFC3339),
            BookingRequestUpdatedDateTime:  time.Now().Format(time.RFC3339),
            TransportDocumentIssuer: Carrier{
                IssuerID:    "carrier123",
                EntityName:  "Carrier Corp",
                Address:     "123 Main St",
                Phone:       "+1-123-456-7890",
                EmailOrFax:  "carrier@company.com",
            },
            ShipperEntity: Shipper{
                CompanyName:     "Shipper Co",
                PhysicalAddress: "456 Elm St",
                ContactName:     "John Doe",
                EmailOrFax:      "shipper@company.com",
                Phone:           "+1-987-654-3210",
                LEIOrTaxID:      "SHIP12345",
            },
            Consignee: Consignee{
                CompanyName:       "Consignee Inc",
                PhysicalAddress:   "789 Oak St",
                ContactName:       "Jane Smith",
                EmailOrFax:        "consignee@company.com",
                Phone:             "+1-555-555-5555",
                LEIOrTaxID:        "CONSGN6789",
                ToOrderIdentifier: "TO123",
            },
            PlaceOfReceipt: PlaceOfReceipt{
                CompanyName:       "Receipt Place Ltd",
                PhysicalAddress:   "101 Pine St",
                ContactName:       "Receptionist",
                EmailOrFax:        "receipt@company.com",
                Phone:             "+1-777-777-7777",
                LEIOrTaxID:        "RCPT5555",
                ToOrderIdentifier: "TO456",
            },
            PlaceOfDelivery:         "Destination City",
            ServiceType:             "Express",
            CargoMovementOrigin:     "Source Location",
            CargoMovementDestination: "Destination Location",
            Commodity: Commodity{
                CommodityName: "Electronics",
                Description:   "Consumer electronics",
                Quantity:      100,
            },
            CargoGrossWeight: 1500.75,
            ContainerTypeSize: "20 ft Container",
            CreditAmount:      75000,
        },
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

// CreateBookingRequest creates a new booking request.
func (s *BookingChaincode) CreateBookingRequest(ctx contractapi.TransactionContextInterface, bookingRequest BookingRequest) (map[string]string, error) {
    response := make(map[string]string)

    // Automatically generate carrierBookingRequestReference
    bookingRequest.CarrierBookingRequestReference = uuid.New().String()
    // Set bookingRequestCreatedDateTime and bookingRequestUpdatedDateTime
    currentTime := time.Now().Format(time.RFC3339)
    bookingRequest.BookingRequestCreatedDateTime = currentTime
    bookingRequest.BookingRequestUpdatedDateTime = currentTime

    // Check if the booking request already exists
    exists, err := s.BookingRequestExists(ctx, bookingRequest.ID)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, fmt.Errorf("the asset %s already exists", bookingRequest.ID)
    }

    bookingrequest := BookingRequest{
                {
                    ID: uuid.New().String(),
                    CarrierBookingRequestReference: uuid.New().String(),
                    DocumentStatus:                 "RECE",
                    BookingRequestCreatedDateTime:  time.Now().Format(time.RFC3339),
                    BookingRequestUpdatedDateTime:  time.Now().Format(time.RFC3339),
                    TransportDocumentIssuer: Carrier{
                        IssuerID: carrier123,
                        EntityName: Carrier Corp,
                        Address: 123 Main St,
                        phone: +1-123-456-7890,
                        emailOrFax: carrier@company.com
                    },
                    shipperEntity: {
                        companyName: Shipper Co,
                        physicalAddress: 456 Elm St,
                        contactName: John Doe,
                        emailOrFax: shipper@company.com,
                        phone: +1-987-654-3210,
                        leiOrTaxID: SHIP12345
                    },
                    consignee: {
                        companyName: Consignee Inc,
                        physicalAddress: 789 Oak St,
                        contactName: Jane Smith,
                        emailOrFax: consignee@company.com,
                        phone: +1-555-555-5555,
                        leiOrTaxID: CONSGN6789,
                        toOrderIdentifier: TO123
                    },
                    placeOfReceipt: {
                        companyName: Receipt Place Ltd,
                        physicalAddress: 101 Pine St,
                        contactName: Receptionist,
                        emailOrFax: receipt@company.com,
                        phone: +1-777-777-7777,
                        leiOrTaxID: RCPT5555,
                        toOrderIdentifier: TO456
                    },
                    placeOfDelivery: Destination City,
                    serviceType: Express,
                    cargoMovementSrc: Source Location,
                    cargoMovementDst: Destination Location,
                    commodity: {
                        commodityName: Electronics,
                        description: Consumer electronics,
                        quantity: 100
                    },
                    cargoGrossWeight: 1500.75,
                    containerType: 20 ft Container,
                    creditAmount: 75000
                }
    }

    // Perform additional validation logic here for other fields
    if BookingRequest.CreditAmount < 50000 {
        response["carrierBookingRequestReference"] = BookingRequest.CarrierBookingRequestReference
        response["documentStatus"] = "CANC"
        response["bookingRequestCreatedDateTime"] = BookingRequest.BookingRequestCreatedDateTime
        response["bookingRequestUpdatedDateTime"] = BookingRequest.BookingRequestUpdatedDateTime
        return response, fmt.Errorf("Credit amount is less than 50000. Booking request canceled.")
    
    
    // Store the booking request in the blockchain ledger
    bookingRequestJSON, err := json.Marshal(bookingRequest)
    if err != nil {
        return nil, err
    }
    err = ctx.GetStub().PutState(bookingRequest.ID, bookingRequestJSON)
    return response, nil
    if err != nil {
        return nil, fmt.Errorf("failed to put to world state: %v", err)
    }

    // Set the response data for a successful booking
    response["carrierBookingRequestReference"] = bookingRequest.CarrierBookingRequestReference
    response["documentStatus"] = "RECE"
    response["bookingRequestCreatedDateTime"] = bookingRequest.BookingRequestCreatedDateTime
    response["bookingRequestUpdatedDateTime"] = bookingRequest.BookingRequestUpdatedDateTime

    return response, nil
}

// BookingRequestExists checks if a booking request with a given ID exists.
func (s *BookingChaincode) BookingRequestExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
    bookingRequestJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
        return false, fmt.Errorf("failed to read from world state: %v", err)
    }
    return bookingRequestJSON != nil, nil
}

func (s *BookingChaincode) ReadBookingRequest(ctx contractapi.TransactionContextInterface, id string) (*BookingRequest, error) {
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

func (s *BookingChaincode) UpdateBookingRequest(ctx contractapi.TransactionContextInterface, id string, carrierBookingRequestReference string, documentStatus string, bookingRequestCreatedDateTime string, bookingRequestUpdatedDateTime string, transportDocumentIssuer Carrier, shipperEntity Shipper, consignee Consignee, placeOfReceipt PlaceOfReceipt, placeOfDelivery string, serviceType string, cargoMovementOrigin string, cargoMovementDestination string, commodity commodity, cargoGrossWeight float64, containerTypeSize string, creditAmount int) error {
    exists, err := s.BookingRequestExists(ctx, id)
    if err != nil {
        return err
    }
    if !exists {
        return fmt.Errorf("the booking request %s does not exist", id)
    }

    // Overwriting the original booking request with the updated booking request
    bookingRequest := BookingRequest{
            CarrierBookingRequestReference: uuid.New().String(),
            TransportDocumentIssuer: Carrier{
                IssuerID:    "carrier123",
                EntityName:  "Carrier Corp",
                Address:     "123 Main St",
                Phone:       "+1-123-456-7890",
                EmailOrFax:  "carrier@company.com",
            },
            ShipperEntity: Shipper{
                CompanyName:     "Shipper Co",
                PhysicalAddress: "456 Elm St",
                ContactName:     "John Doe",
                EmailOrFax:      "shipper@company.com",
                Phone:           "+1-987-654-3210",
                LEIOrTaxID:      "SHIP12345",
            },
            Consignee: Consignee{
                CompanyName:       "Consignee Inc",
                PhysicalAddress:   "789 Oak St",
                ContactName:       "Jane Smith",
                EmailOrFax:        "consignee@company.com",
                Phone:             "+1-555-555-5555",
                LEIOrTaxID:        "CONSGN6789",
                ToOrderIdentifier: "TO123",
            },
            PlaceOfReceipt: PlaceOfReceipt{
                CompanyName:       "Receipt Place Ltd",
                PhysicalAddress:   "101 Pine St",
                ContactName:       "Receptionist",
                EmailOrFax:        "receipt@company.com",
                Phone:             "+1-777-777-7777",
                LEIOrTaxID:        "RCPT5555",
                ToOrderIdentifier: "TO456",
            },
            PlaceOfDelivery:         "Destination City",
            ServiceType:             "Express",
            CargoMovementOrigin:     "Source Location",
            CargoMovementDestination: "Destination Location",
            Commodity: Commodity{
                CommodityName: "Electronics",
                Description:   "Consumer electronics",
                Quantity:      100,
            },
            CargoGrossWeight: 1500.75,
            ContainerTypeSize: "20 ft Container",
            CreditAmount:      75000,
        },
    }
    bookingRequestJSON, err := json.Marshal(bookingRequest)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(id, bookingRequestJSON)
}


func (s *BookingChaincode) GetAllBookingRequests(ctx contractapi.TransactionContextInterface) ([]*BookingRequest, error) {
    // Range query with empty string for startKey and endKey does an
    // open-ended query of all booking requests in the chaincode namespace.
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
