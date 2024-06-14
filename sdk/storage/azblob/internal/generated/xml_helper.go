package generated

import (
	"encoding/xml"
)

//type additionalProperties map[string]*string

// MarshalXML implements the xml.Marshaler interface for additionalProperties.
func (ap additionalProperties) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}
	for k, v := range ap {
		err := e.EncodeElement(v, xml.StartElement{Name: xml.Name{Local: k}})
		if err != nil {
			return err
		}
	}
	err = e.EncodeToken(xml.EndElement{Name: start.Name})
	return err
}
