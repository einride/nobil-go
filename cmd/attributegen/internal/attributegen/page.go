package attributegen

import (
	"fmt"
	"regexp"
	"strconv"
)

type AttributesPage struct {
	Attributes []*Attribute
}

func (a *AttributesPage) UnmarshalHTML(text string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("unmarshal attributes page HTML: %w", err)
		}
	}()
	tableMatches := regexp.MustCompile(
		`(?s)(<table.*?width: 980px.*?\n</table>)`,
	).FindAllStringSubmatch(text, -1)
	if len(tableMatches) != 3 {
		return fmt.Errorf("unexpected table match count: %d", len(tableMatches))
	}
	if err := a.unmarshalAttributeTableHTML(tableMatches[0][0]); err != nil {
		return err
	}
	if err := a.unmarshalAttributeValueTableHTML(tableMatches[1][0]); err != nil {
		return err
	}
	if err := a.unmarshalAttributeValueTableHTML(tableMatches[2][0]); err != nil {
		return err
	}
	return nil
}

func (a *AttributesPage) unmarshalAttributeTableHTML(text string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("unmarshal attribute table HTML: %w", err)
		}
	}()
	attributeMatches := regexp.MustCompile(
		`<td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td>`,
	).FindAllStringSubmatch(text, -1)
	a.Attributes = a.Attributes[:0]
	for _, attributeMatch := range attributeMatches {
		if len(attributeMatch) != 6 {
			return fmt.Errorf("unexpected match count: %d", len(attributeMatch))
		}
		attributeID, err := strconv.Atoi(attributeMatch[1])
		if err != nil {
			return fmt.Errorf("attribute ID: %w", err)
		}
		attributeName := attributeMatch[2]
		if attributeName == "" {
			return fmt.Errorf("attribute name: empty")
		}
		translationID, err := strconv.Atoi(attributeMatch[3])
		if err != nil {
			return fmt.Errorf("translation ID: %w", err)
		}
		a.Attributes = append(a.Attributes, &Attribute{
			ID:             attributeID,
			Name:           attributeName,
			TranslationID:  translationID,
			UserSelectable: attributeMatch[4] == "1",
			Required:       attributeMatch[5] == "1",
		})
	}
	return nil
}

func (a *AttributesPage) unmarshalAttributeValueTableHTML(text string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("unmarshal attribute table HTML: %w", err)
		}
	}()
	attributeMatches := regexp.MustCompile(
		`(?s)(<tr><td><strong>.*?</td>\n\s*?</tr>)`,
	).FindAllStringSubmatch(text, -1)
	for _, attributeMatch := range attributeMatches {
		if len(attributeMatch) != 2 {
			return fmt.Errorf("unexpected attribute match count: %d", len(attributeMatch))
		}
		attributeIDMatches := regexp.MustCompile(
			`<tr><td><strong>(.*?)</strong></td>`,
		).FindStringSubmatch(attributeMatch[0])
		if len(attributeIDMatches) != 2 {
			return fmt.Errorf("unexpected attribute ID match count: %d", len(attributeIDMatches))
		}
		attributeID, err := strconv.Atoi(attributeIDMatches[1])
		if err != nil {
			return fmt.Errorf("attribute ID: %w", err)
		}
		var attribute *Attribute
		for _, candidate := range a.Attributes {
			if candidate.ID == attributeID {
				attribute = candidate
				break
			}
		}
		if attribute == nil {
			return fmt.Errorf("no attribute for ID: %d", attributeID)
		}
		attributeValueMatches := regexp.MustCompile(
			`<td>(.*?)</td><td>(.*?)</td><td>(.*?)</td><td>(.*?)</td></tr>`,
		).FindAllStringSubmatch(attributeMatch[0], -1)
		for _, attributeValueMatch := range attributeValueMatches {
			if len(attributeValueMatch) != 5 {
				return fmt.Errorf("unexpected attribute value match count: %d", len(attributeValueMatch))
			}
			valueID, err := strconv.Atoi(attributeValueMatch[1])
			if err != nil {
				return fmt.Errorf("value ID: %w", err)
			}
			valueName := attributeValueMatch[2]
			if valueName == "" {
				return fmt.Errorf("value name: empty")
			}
			attribute.Values = append(attribute.Values, &AttributeValue{
				ID:         valueID,
				Name:       valueName,
				Translated: attributeValueMatch[3],
				Key:        attributeValueMatch[4],
			})
		}
	}
	return nil
}
