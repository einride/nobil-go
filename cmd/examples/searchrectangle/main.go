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
	response, err := client.SearchRectangle(ctx, &nobil.SearchRectangleRequest{
		NorthEast: nobil.LatLng{
			Latitude:  59.943921193288915,
			Longitude: 10.826683044433594,
		},
		SouthWest: nobil.LatLng{
			Latitude:  59.883683240905256,
			Longitude: 10.650901794433594,
		},
		ExistingIDs: []string{"189", "195", "199", "89", "48"},
	})
	if err != nil {
		panic(err)
	}
	for _, chargingStation := range response.ChargingStations {
		fmt.Printf("%+v\n", chargingStation)
	}
}
