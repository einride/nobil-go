package nobil

import (
	"strconv"
	"time"
)

// ChargingStation represents a charging station.
type ChargingStation struct {
	ID                          int
	Name                        string
	Street                      string
	HouseNumber                 string
	ZipCode                     string
	City                        string
	MunicipalityID              string
	Municipality                string
	CountyID                    string
	County                      string
	Description                 string
	Owner                       string
	ChargingPointCount          int
	AvailableChargingPointCount int
	Position                    LatLng
	Image                       string
	UserComment                 string
	ContactInfo                 string
	CreateTime                  time.Time
	UpdateTime                  time.Time
	StationStatus               int
	LandCode                    string
	InternationalID             string
	Location                    Location
	Availability                Availability
	Open24h                     Open24h
	ParkingFee                  ParkingFee
	TimeLimit                   TimeLimit
	RealTimeInformation         RealTimeInformation
	PublicFunding               PublicFunding
	Connections                 []*Connection
	Attributes                  map[AttributeID]*Attribute
}

func (c *ChargingStation) unmarshalAttributes(attrs map[AttributeID]*Attribute) {
	c.Location.unmarshalAttributes(attrs)
	c.Availability.unmarshalAttributes(attrs)
	c.Open24h.unmarshalAttributes(attrs)
	c.ParkingFee.unmarshalAttributes(attrs)
	c.TimeLimit.unmarshalAttributes(attrs)
	c.RealTimeInformation.unmarshalAttributes(attrs)
	c.PublicFunding.unmarshalAttributes(attrs)
	c.Attributes = attrs
}

func (c *ChargingStation) CSVHeader() []string {
	return []string{
		"station_id",
		"name",
		"street",
		"house_number",
		"zip_code",
		"city",
		"municipality_id",
		"municipality",
		"county_id",
		"county",
		"description",
		"owner",
		"charging_point_count",
		"latitude",
		"longitude",
		"create_time",
		"update_time",
		"country_code",
		"international_id",
		"location",
		"availability",
		"open_24h",
		"parking_fee",
		"time_limit",
		"public_funding",
	}
}

func (c *ChargingStation) CSVRecord() []string {
	return []string{
		strconv.Itoa(c.ID),
		c.Name,
		c.Street,
		c.HouseNumber,
		c.ZipCode,
		c.City,
		c.MunicipalityID,
		c.Municipality,
		c.CountyID,
		c.County,
		c.Description,
		c.Owner,
		strconv.Itoa(c.ChargingPointCount),
		strconv.FormatFloat(c.Position.Latitude, 'f', -1, 64),
		strconv.FormatFloat(c.Position.Longitude, 'f', -1, 64),
		c.CreateTime.String(),
		c.UpdateTime.String(),
		c.LandCode,
		c.InternationalID,
		c.Location.String(),
		c.Availability.String(),
		c.Open24h.String(),
		c.ParkingFee.String(),
		c.TimeLimit.String(),
		c.PublicFunding.String(),
	}
}
