# NOBIL Go

Go client to the [NOBIL.no][nobil-no] public charger [API][nobil-api].

[nobil-no]: https://info.nobil.no
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

```bash
$ NOBIL_API_KEY=YOUR_API_KEY go run go.einride.tech/cmd/examples/searchrectangle
```

```go
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
			Latitude:  59.94392,
			Longitude: 10.82668,
		},
		SouthWest: nobil.LatLng{
			Latitude:  59.88368,
			Longitude: 10.65090,
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

```bash
$ NOBIL_API_KEY=YOUR_API_KEY go run go.einride.tech/cmd/examples/searchnear
```

```go
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
```
