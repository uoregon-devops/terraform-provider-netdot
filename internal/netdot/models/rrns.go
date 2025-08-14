package models

import "encoding/xml"

// <RRADDR id="544859" ipblock="184.171.15.38" ipblock_xlink="Ipblock/6577810709" rr="mollman-test-record.uoregon.edu" rr_xlink="RR/999289" ttl="800"/>

type RRNs struct {
	ID            int64  `xml:"id,attr"`
	NsDName       string `xml:"nsdname,attr"`
	RR            string `xml:"rr,attr"`
	RRXLinkString string `xml:"rr_xlink,attr"`
	RRXlink       Xlink
	TTL           int64 `xml:"ttl,attr"`
}

// XML Unmarshaler for RRAddr
func (r *RRNs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type bufferRRNs RRNs
	finalRRNs := bufferRRNs{}
	err := d.DecodeElement(&finalRRNs, &start)
	if err != nil {
		return err
	}

	rr_xlink, err := parseXlink(finalRRNs.RRXLinkString)
	if err != nil {
		return err
	}
	finalRRNs.RRXlink = rr_xlink

	*r = RRNs(finalRRNs)
	return nil
}
