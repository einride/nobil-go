package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/einride/nobil-go"
)

func main() {
	ctx := context.Background()
	log.SetFlags(0)
	type subCommand struct {
		name        string
		description string
		fn          func(context.Context, *flag.FlagSet, []string)
	}
	cmds := []subCommand{
		{name: "search-rectangle", description: "search for chargers in a rectangle", fn: searchRectangle},
		{name: "search-near", description: "search for chargers near a coordinate", fn: searchNear},
		{name: "dump", description: "dump the full charger database", fn: dump},
		{name: "scan", description: "scan a JSON dump", fn: scan},
		{name: "stream", description: "stream real-time events", fn: stream},
	}
	var cmdName string
	var args []string
	if len(os.Args) > 1 {
		cmdName, args = os.Args[1], os.Args[2:]
	}
	for _, cmd := range cmds {
		if cmdName == cmd.name {
			fs := flag.NewFlagSet("nobilctl "+cmd.name, flag.ExitOnError)
			fs.SetOutput(os.Stderr)
			fs.Usage = func() {
				log.Println("usage: nobilctl", cmd.name)
				log.Println()
				fs.PrintDefaults()
				log.Println()
			}
			cmd.fn(ctx, fs, args)
			os.Exit(0)
		}
	}
	var longestName int
	for _, cmd := range cmds {
		if len(cmd.name) > longestName {
			longestName = len(cmd.name)
		}
	}
	log.Println("usage: nobilctl <command>")
	log.Println()
	for _, cmd := range cmds {
		log.Println("    ", cmd.name, strings.Repeat(" ", longestName-len(cmd.name)), cmd.description)
	}
	log.Println()
	os.Exit(1)
}

func searchRectangle(ctx context.Context, fs *flag.FlagSet, args []string) {
	var apiKey string
	fs.StringVar(&apiKey, "apiKey", "", "API key")
	var southWest latLngValue
	fs.Var(&southWest, "southWest", "south-west coordinate (lat, lng)")
	var northEast latLngValue
	fs.Var(&northEast, "northEast", "north-east coordinate (lat, lng)")
	var existingIDs commaSeparatedStringValue
	fs.Var(&existingIDs, "existingIDs", "existing IDs")
	if err := fs.Parse(args); err != nil {
		log.Panic(err)
	}
	if apiKey == "" || southWest == (latLngValue{}) || northEast == (latLngValue{}) {
		fs.Usage()
		os.Exit(1)
	}
	client := nobil.NewClient(apiKey)
	response, err := client.SearchRectangle(ctx, &nobil.SearchRectangleRequest{
		SouthWest:   southWest.LatLng,
		NorthEast:   northEast.LatLng,
		ExistingIDs: existingIDs,
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(response.Raw))
}

func searchNear(ctx context.Context, fs *flag.FlagSet, args []string) {
	var apiKey string
	fs.StringVar(&apiKey, "apiKey", "", "API key")
	var coordinate latLngValue
	fs.Var(&coordinate, "coordinate", "coordinate (lat, lng)")
	var distanceMetres int
	fs.IntVar(&distanceMetres, "distance", 0, "distance (metres)")
	var limit int
	fs.IntVar(&limit, "limit", 0, "limit")
	if err := fs.Parse(args); err != nil {
		log.Panic(err)
	}
	if apiKey == "" || coordinate == (latLngValue{}) || distanceMetres <= 0 || limit <= 0 {
		fs.Usage()
		os.Exit(1)
	}
	client := nobil.NewClient(apiKey)
	response, err := client.SearchNear(ctx, &nobil.SearchNearRequest{
		Coordinate:     coordinate.LatLng,
		DistanceMetres: distanceMetres,
		Limit:          limit,
	})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(response.Raw))
	for _, result := range response.Results {
		log.Printf("%+v", result.ChargingStation)
		for _, connection := range result.ChargingStation.Connections {
			log.Printf("%+v", connection)
		}
	}
}

func dump(ctx context.Context, fs *flag.FlagSet, args []string) {
	var apiKey string
	fs.StringVar(&apiKey, "apiKey", "", "API key")
	var fromDate dateValue
	fs.Var(&fromDate, "fromDate", "dump from date")
	var countryCode string
	fs.StringVar(&countryCode, "countryCode", "", "country code of chargers to dump (or all countries if unspecified)")
	if err := fs.Parse(args); err != nil {
		log.Panic(err)
	}
	if apiKey == "" {
		fs.Usage()
		os.Exit(1)
	}
	client := nobil.NewClient(apiKey)
	dump, err := client.Dump(ctx, &nobil.DumpRequest{
		FromDate:    fromDate.Date,
		CountryCode: countryCode,
		Format:      nobil.FormatJSON,
	})
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		_ = dump.Close()
	}()
	n, err := io.Copy(os.Stdout, dump)
	if err != nil {
		log.Panic(err)
	}
	log.Println()
	log.Println("dumped", n, "bytes")
}

func scan(_ context.Context, fs *flag.FlagSet, args []string) {
	var file string
	fs.StringVar(&file, "file", "", "JSON file to scane")
	if err := fs.Parse(args); err != nil {
		log.Panic(err)
	}
	if file == "" {
		fs.Usage()
		os.Exit(1)
	}
	f, err := os.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Panic(err)
		}
	}()
	sc := nobil.NewJSONScanner(f)
	csvOut := csv.NewWriter(os.Stdout)
	_ = csvOut.Write(append((&nobil.ChargingStation{}).CSVHeader(), (&nobil.Connection{}).CSVHeader()...))
	for sc.Scan() {
		chargingStation := sc.ChargingStation()
		if chargingStation.LandCode != "SWE" {
			continue
		}
		for _, connection := range chargingStation.Connections {
			if connection.Connector != nobil.Connector_CcsCombo {
				continue
			}
			_ = csvOut.Write(append(chargingStation.CSVRecord(), connection.CSVRecord()...))
		}
	}
	if sc.Err() != nil {
		log.Panic(sc.Err())
	}
	csvOut.Flush()
}

func stream(ctx context.Context, fs *flag.FlagSet, args []string) {
	// TODO: Implement streaming websocket connection
}

type latLngValue struct {
	nobil.LatLng
}

var _ = flag.Value(&latLngValue{})

func (l *latLngValue) Set(s string) error {
	return l.UnmarshalString(s)
}

type commaSeparatedStringValue []string

var _ = flag.Value(&commaSeparatedStringValue{})

func (c commaSeparatedStringValue) String() string {
	return strings.Join(c, ",")
}

func (c *commaSeparatedStringValue) Set(s string) error {
	*c = strings.Split(s, ",")
	return nil
}

type dateValue struct {
	nobil.Date
}

var _ = flag.Value(&dateValue{})

func (l *dateValue) Set(s string) error {
	return l.UnmarshalString(s)
}
