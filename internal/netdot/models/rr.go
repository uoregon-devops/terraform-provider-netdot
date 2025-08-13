package models

import "encoding/xml"

type RR struct {
	ID              int64  `xml:"id,attr"`
	Active          bool   `xml:"active,attr"`
	AutoUpdate      bool   `xml:"auto_update,attr"`
	Expiration      string `xml:"expiration,attr"`
	Info            string `xml:"info,attr"`
	Created         string `xml:"created,attr"`
	Modified        string `xml:"modified,attr"`
	Name            string `xml:"name,attr"`
	Zone            string `xml:"zone,attr"`
	ZoneXlinkString string `xml:"zone_xlink,attr"`
	ZoneXlink       Xlink
}

// XML Unmarshaler for RR
func (r *RR) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type bufferRR RR
	finalRR := bufferRR{}
	err := d.DecodeElement(&finalRR, &start)
	if err != nil {
		return err
	}

	zone_xlink, err := parseXlink(finalRR.ZoneXlinkString)
	if err != nil {
		return err
	}
	finalRR.ZoneXlink = zone_xlink

	*r = RR(finalRR)
	return nil
}
