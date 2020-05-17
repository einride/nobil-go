package nobil

type AttributeID string

type Attribute struct {
	ID          AttributeID `json:"attrtypeid"`
	Name        string      `json:"attrname"`
	ValueID     string      `json:"attrvalid"`
	Translation string      `json:"trans"`
	Value       interface{} `json:"attrval"`
}
