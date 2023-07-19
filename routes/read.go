package route

import (
	getcollection "csit-mini-challenge/Collection"
	database "csit-mini-challenge/databases"
	model "csit-mini-challenge/model"

	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCheapestFlights(c *gin.Context) {
	db := database.ConnectDB()
	ctx := context.TODO()
	defer func() {
		if err := db.Disconnect(ctx); err != nil {
		  	log.Fatal(err)
		}
	}()
	collection := getcollection.GetCollection(db, "flights")

	departureDate, returnDate, destination, success := getFlightQueryParams(c)
	if !success {
		return
	}

	pipeline := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"srccity", "Singapore"},
					{"destcity", bson.D{{"$regex", primitive.Regex{Pattern: destination, Options: "i"}}}},
					{"date", departureDate},
				},
			},
		},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "flights"},
					{"localField", "destcity"},
					{"foreignField", "srccity"},
					{"as", "return_flights"},
				},
			},
		},
		bson.D{{"$unwind", bson.D{{"path", "$return_flights"}}}},
		bson.D{
			{"$match",
				bson.D{
					{"return_flights.destcity", "Singapore"},
					{"return_flights.date", returnDate},
				},
			},
		},
		bson.D{
			{"$addFields",
				bson.D{
					{"totalprice",
						bson.D{
							{"$add",
								bson.A{
									"$price",
									"$return_flights.price",
								},
							},
						},
					},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"totalprice", 1}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	var results []model.ReturnFlight
	for cursor.Next(ctx) {
		var returnFlight model.ReturnFlight
		err := cursor.Decode(&returnFlight)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}
		if len(results) == 0  || returnFlight.TotalPrice == results[0].TotalPrice{
			results = append(results, returnFlight)
		} else {
			break
		}
	}

	arr := []model.FlightResponse{}
	for _, result := range results {
		response := model.FlightResponse{
			City: result.DestCity,
			DepartureDate: result.Date.Time().Format("2006-01-02"),
			DepartureAirline: result.AirlineName,
			DeparturePrice: result.Price,
			ReturnDate: result.ReturnFlight.Date.Time().Format("2006-01-02"),
			ReturnAirline: result.ReturnFlight.AirlineName,
			ReturnPrice: result.ReturnFlight.Price,
		}

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		arr = append(arr, response)
	}
	c.IndentedJSON(http.StatusOK, arr)
}

func getFlightQueryParams(c *gin.Context) (primitive.DateTime, primitive.DateTime, string, bool) {
	emptyDate := primitive.NewDateTimeFromTime(time.Time{})
	departureDate, exists := c.GetQuery("departureDate")
	if !exists {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "departureDate not provided"})
		return emptyDate, emptyDate, "", false
	}
	returnDate, exists := c.GetQuery("returnDate")
	if !exists {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "returnDate not provided"})
		return emptyDate, emptyDate, "", false
	}
	destination, exists := c.GetQuery("destination")
	if !exists {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "destination not provided"})
		return emptyDate, emptyDate, "", false
	}

	parsedDepartureDate, err := parseDate(departureDate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid departureDate"})
		return emptyDate, emptyDate, "", false
	}
	parsedReturnDate, err := parseDate(returnDate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid returnDate"})
		return emptyDate, emptyDate, "", false
	}

	return 	primitive.NewDateTimeFromTime(parsedDepartureDate), 
			primitive.NewDateTimeFromTime(parsedReturnDate), 
			destination, true
}

func GetCheapestHotels(c *gin.Context) {
	db := database.ConnectDB()
	ctx := context.TODO()
	defer func() {
		if err := db.Disconnect(ctx); err != nil {
		  	log.Fatal(err)
		}
	}()
	collection := getcollection.GetCollection(db, "hotels")

	checkInDate, checkOutDate, destination, success := getHotelQueryParams(c)
	if !success {
		return
	}

	pipeline := bson.A{
		bson.D{
			{"$match",
				bson.D{
					{"city", bson.D{{"$regex", primitive.Regex{Pattern: destination, Options: "i"}}}},
					{"date",
						bson.D{
							{"$gte", checkInDate},
							{"$lte", checkOutDate},
						},
					},
				},
			},
		},
		bson.D{
			{"$group",
				bson.D{
					{"_id",
						bson.D{
							{"hotelName", "$hotelName"},
							{"city", "$city"},
						},
					},
					{"totalPrice", bson.D{{"$sum", "$price"}}},
					{"checkindate", bson.D{{"$min", "$date"}}},
					{"checkoutdate", bson.D{{"$max", "$date"}}},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"totalPrice", 1}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	var results []model.HotelBooking
	for cursor.Next(ctx) {
		var hotelBooking model.HotelBooking
		err := cursor.Decode(&hotelBooking)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}
		if len(results) == 0  || hotelBooking.Price == results[0].Price{
			results = append(results, hotelBooking)
		} else {
			break
		}
	}

	arr := []model.HotelResponse{}
	for _, result := range results {
		response := model.HotelResponse{
			City: result.ID.City,
			CheckInDate: result.CheckInDate.Time().Format("2006-01-02"),
			CheckOutDate: result.CheckOutDate.Time().Format("2006-01-02"),
			Hotel: result.ID.HotelName,
			Price: result.Price,
		}

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		arr = append(arr, response)
	}
	c.IndentedJSON(http.StatusOK, arr)
}

func getHotelQueryParams(c *gin.Context) (primitive.DateTime, primitive.DateTime, string, bool) {
	emptyDate := primitive.NewDateTimeFromTime(time.Time{})
	checkInDate, exists := c.GetQuery("checkInDate")
	if !exists {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "checkInDate not provided"})
		return emptyDate, emptyDate, "", false
	}
	checkOutDate, exists := c.GetQuery("checkOutDate")
	if !exists {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "checkOutDate not provided"})
		return emptyDate, emptyDate, "", false
	}
	destination, exists := c.GetQuery("destination")
	if !exists {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "destination not provided"})
		return emptyDate, emptyDate, "", false
	}

	parsedCheckInDate, err := parseDate(checkInDate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid departureDate"})
		return emptyDate, emptyDate, "", false
	}
	parsedCheckOutDate, err := parseDate(checkOutDate)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid returnDate"})
		return emptyDate, emptyDate, "", false
	}

	return 	primitive.NewDateTimeFromTime(parsedCheckInDate), 
			primitive.NewDateTimeFromTime(parsedCheckOutDate), 
			destination, true
}

func parseDate(date string) (time.Time, error) {
	dt, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, err
	}
	return dt, nil
}