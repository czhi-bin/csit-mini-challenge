package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flight struct {
	ID              primitive.ObjectID `bson:"_id"`
	Airline         string             `bson:"airline"`
	AirlineID       int32              `bson:"airlineid"`
	SrcAirport      string             `bson:"srcairport"`
	SrcAirportID    int32              `bson:"srcairportid"`
	DestAirport     string             `bson:"destairport"`
	DestAirportID   int32              `bson:"destairportid"`
	CodeShare       string             `bson:"codeshare"`
	Stops           int32              `bson:"stops"`
	Equipment       string             `bson:"eq"`
	AirlineName     string             `bson:"airlinename"`
	SrcAirportName  string             `bson:"srcairportname"`
	SrcCity         string             `bson:"srccity"`
	SrcCountry      string             `bson:"srccountry"`
	DestAirportName string             `bson:"destairportname"`
	DestCity        string             `bson:"destcity"`
	DestCountry     string             `bson:"destcountry"`
	Price           int32              `bson:"price"`
	Date            primitive.DateTime `bson:"date"`
}

type ReturnFlight struct {
	ID              primitive.ObjectID `bson:"_id"`
	Airline         string             `bson:"airline"`
	AirlineID       int32              `bson:"airlineid"`
	SrcAirport      string             `bson:"srcairport"`
	SrcAirportID    int32              `bson:"srcairportid"`
	DestAirport     string             `bson:"destairport"`
	DestAirportID   int32              `bson:"destairportid"`
	CodeShare       string             `bson:"codeshare"`
	Stops           int32              `bson:"stops"`
	Equipment       string             `bson:"eq"`
	AirlineName     string             `bson:"airlinename"`
	SrcAirportName  string             `bson:"srcairportname"`
	SrcCity         string             `bson:"srccity"`
	SrcCountry      string             `bson:"srccountry"`
	DestAirportName string             `bson:"destairportname"`
	DestCity        string             `bson:"destcity"`
	DestCountry     string             `bson:"destcountry"`
	Price           int32              `bson:"price"`
	Date            primitive.DateTime `bson:"date"`
	ReturnFlight    Flight             `bson:"return_flights"`
	TotalPrice      int32              `bson:"totalprice"`
}

type FlightResponse struct {
	City             string `json:"City"`
	DepartureDate    string `json:"Departure Date"`
	DepartureAirline string `json:"Departure Airline"`
	DeparturePrice   int32  `json:"Departure Price"`
	ReturnDate       string `json:"Return Date"`
	ReturnAirline    string `json:"Return Airline"`
	ReturnPrice      int32  `json:"Return Price"`
}

type HotelID struct {
	HotelName string `bson:"hotelname"`
	City      string `bson:"city"`
}

type HotelBooking struct {
	ID           HotelID            `bson:"_id"`
	Price        int32              `bson:"totalprice"`
	CheckInDate  primitive.DateTime `bson:"checkindate"`
	CheckOutDate primitive.DateTime `bson:"checkoutdate"`
}

type HotelResponse struct {
	City         string `json:"City"`
	CheckInDate  string `json:"Check In Date"`
	CheckOutDate string `json:"Check Out Date"`
	Hotel        string `json:"Hotel"`
	Price        int32  `json:"Price"`
}
