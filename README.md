# NOBIL-Go

Go client to the [NOBIL API][nobil-api] for location and metadata about
charging stations in Scandinavia.

[nobil-api]: https://info.nobil.no/api

## Usage

### Get an API key

Visit [info.nobil.no/api][nobil-api] to ask for an API key for your
project.

### go get

```bash
$ go get -u go.einride.tech/nobil
```

## Examples

### Search charging stations in a rectangle

```go
package main

import (
	"context"
	"fmt"

	"github.com/einride/nobil-go"
)

func main() {
	ctx := context.Background()
	client := nobil.NewClient("YOUR_API_KEY_HERE")
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
```

### Search nearby charging stations

```go
import (
	"context"
	"fmt"
	"os"

	"github.com/einride/nobil-go"
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
```
