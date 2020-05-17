package nobil

// Connection represents a single connection at a charging station.
type Connection struct {
	ID                    string
	VehicleType           VehicleType
	Accessibility         Accessibility
	ChargingCapacity      ChargingCapacity
	Connector             Connector
	FixedCable            FixedCable
	ChargeMode            ChargeMode
	PaymentMethod         PaymentMethod
	Reservable            Reservable
	ConnectorSensorStatus ConnectorSensorStatus
	ConnectorErrorStatus  ConnectorErrorStatus
	ConnectorStatus       ConnectorStatus
	EnergyCarrier         EnergyCarrier
	EVSEID                string
	Manufacturer          string
	Attributes            map[AttributeID]*Attribute
}

func (c *Connection) unmarshalAttributes(attrs map[AttributeID]*Attribute) {
	c.VehicleType.unmarshalAttributes(attrs)
	c.Accessibility.unmarshalAttributes(attrs)
	c.ChargingCapacity.unmarshalAttributes(attrs)
	c.Connector.unmarshalAttributes(attrs)
	c.FixedCable.unmarshalAttributes(attrs)
	c.ChargeMode.unmarshalAttributes(attrs)
	c.PaymentMethod.unmarshalAttributes(attrs)
	c.Reservable.unmarshalAttributes(attrs)
	c.ConnectorSensorStatus.unmarshalAttributes(attrs)
	c.ConnectorErrorStatus.unmarshalAttributes(attrs)
	c.ConnectorStatus.unmarshalAttributes(attrs)
	c.EnergyCarrier.unmarshalAttributes(attrs)
	if attr, ok := attrs[AttributeID_Evseid]; ok {
		if val, ok := attr.Value.(string); ok {
			c.EVSEID = val
		}
	}
	if attr, ok := attrs[AttributeID_Manufacturer]; ok {
		if val, ok := attr.Value.(string); ok {
			c.Manufacturer = val
		}
	}
	c.Attributes = attrs
}

func (c *Connection) CSVHeader() []string {
	return []string{
		"connection_id",
		"vehicle_type",
		"accessibility",
		"charging_capacity",
		"connector",
		"fixed_cable",
		"charge_mode",
		"payment_method",
		"reservable",
		"energy_carrier",
		"evseid",
		"manufacturer",
	}
}

func (c *Connection) CSVRecord() []string {
	return []string{
		c.ID,
		c.VehicleType.String(),
		c.Accessibility.String(),
		c.ChargingCapacity.String(),
		c.Connector.String(),
		c.FixedCable.String(),
		c.ChargeMode.String(),
		c.PaymentMethod.String(),
		c.Reservable.String(),
		c.EnergyCarrier.String(),
		c.EVSEID,
		c.Manufacturer,
	}
}
