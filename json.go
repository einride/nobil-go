package nobil

import (
	"html"
	"strconv"
	"strings"
	"time"
)

type jsonChargingStation struct {
	jsonMetadata   `json:"csmd"`
	jsonAttributes `json:"attr"`
}

func (j *jsonChargingStation) ChargingStation() *ChargingStation {
	c := &ChargingStation{
		ID:                          j.ID,
		Name:                        strings.TrimSpace(html.UnescapeString(j.Name)),
		Street:                      strings.TrimSpace(html.UnescapeString(j.Street)),
		HouseNumber:                 strings.TrimSpace(html.UnescapeString(j.HouseNumber)),
		ZipCode:                     strings.TrimSpace(html.UnescapeString(j.ZipCode)),
		City:                        strings.TrimSpace(html.UnescapeString(j.City)),
		MunicipalityID:              strings.TrimSpace(html.UnescapeString(j.MunicipalityID)),
		Municipality:                strings.TrimSpace(html.UnescapeString(j.Municipality)),
		CountyID:                    strings.TrimSpace(html.UnescapeString(j.CountyID)),
		County:                      strings.TrimSpace(html.UnescapeString(j.County)),
		Description:                 strings.TrimSpace(html.UnescapeString(j.Description)),
		Owner:                       strings.TrimSpace(html.UnescapeString(j.Owner)),
		ChargingPointCount:          j.ChargingPointCount,
		AvailableChargingPointCount: j.AvailableChargingPointCount,
		Position:                    j.Position,
		Image:                       strings.TrimSpace(html.UnescapeString(j.Image)),
		UserComment:                 strings.TrimSpace(html.UnescapeString(j.UserComment)),
		ContactInfo:                 strings.TrimSpace(html.UnescapeString(j.ContactInfo)),
		CreateTime:                  time.Time(j.CreateTime),
		UpdateTime:                  time.Time(j.UpdateTime),
		StationStatus:               j.StationStatus,
		LandCode:                    strings.TrimSpace(html.UnescapeString(j.LandCode)),
		InternationalID:             strings.TrimSpace(html.UnescapeString(j.InternationalID)),
	}
	c.unmarshalAttributes(j.StationAttributes)
	c.Connections = make([]*Connection, 0, len(j.ConnectionAttributes))
	for connectionID, connectionAttributes := range j.ConnectionAttributes {
		connection := &Connection{
			ID:         connectionID,
			Attributes: connectionAttributes,
		}
		connection.unmarshalAttributes(connectionAttributes)
		c.Connections = append(c.Connections, connection)
	}
	return c
}

type jsonAttributes struct {
	StationAttributes    map[AttributeID]*Attribute            `json:"st"`
	ConnectionAttributes map[string]map[AttributeID]*Attribute `json:"conn"`
}

type jsonMetadata struct {
	ID                          int           `json:"id"`
	Name                        string        `json:"name"`
	Street                      string        `json:"Street"`
	HouseNumber                 string        `json:"House_number"`
	ZipCode                     string        `json:"zipcode"`
	City                        string        `json:"city"`
	MunicipalityID              string        `json:"Municipality_ID"`
	Municipality                string        `json:"Municipality"`
	CountyID                    string        `json:"County_ID"`
	County                      string        `json:"County"`
	Description                 string        `json:"Description_of_location"`
	Owner                       string        `json:"Owned_by"`
	ChargingPointCount          int           `json:"Number_charging_points"`
	AvailableChargingPointCount int           `json:"Available_charging_points"`
	Position                    LatLng        `json:"Position"`
	Image                       string        `json:"Image"`
	UserComment                 string        `json:"User_comment"`
	ContactInfo                 string        `json:"Contact_info"`
	CreateTime                  jsonTimestamp `json:"Created"`
	UpdateTime                  jsonTimestamp `json:"Updated"`
	StationStatus               int           `json:"Station_status"`
	LandCode                    string        `json:"Land_code"`
	InternationalID             string        `json:"International_id"`
}

type jsonTimestamp time.Time

func (t *jsonTimestamp) UnmarshalJSON(bytes []byte) error {
	s, err := strconv.Unquote(string(bytes))
	if err != nil {
		return err
	}
	ts, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*t = jsonTimestamp(ts.UTC())
	return nil
}
