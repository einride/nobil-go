package nobil

import (
	"strconv"
	"strings"
)

type commaSeparatedStrings []string

func (c *commaSeparatedStrings) UnmarshalString(s string) error {
	*c = strings.Split(s, ",")
	return nil
}

func (c commaSeparatedStrings) String() string {
	return strings.Join(c, ",")
}

func (c commaSeparatedStrings) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(c.String())), nil
}
