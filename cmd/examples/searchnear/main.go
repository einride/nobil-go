package main

import (
	"context"
	"fmt"
	"os"

	"go.einride.tech/nobil"
)

func main() {
	ctx := context.Background()
	client := nobil.NewClient(os.Getenv("NOBIL_API_KEY"))
	response, err := client.SearchNear(ctx, &nobil.SearchNearRequest{
		Coordinate: nobil.LatLng{
			Latitude:  59.91673,
			Longitude: 10.74782,
		},
		DistanceMetres: 2_000,
		Limit:          10,
	})
	if err != nil {
		panic(err)
	}
	for _, result := range response.Results {
		fmt.Printf("%dm: %+v\n", result.DistanceMetres, result.ChargingStation)
	}
}
