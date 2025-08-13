package models

import "encoding/xml"

// <RRADDR id="544859" ipblock="184.171.15.38" ipblock_xlink="Ipblock/6577810709" rr="mollman-test-record.uoregon.edu" rr_xlink="RR/999289" ttl="800"/>

type RRCname struct {
	ID            int64  `xml:"id,attr"`
	Cname         string `xml:"cname,attr"`
	RR            string `xml:"rr,attr"`
	RRXLinkString string `xml:"rr_xlink,attr"`
	RRXlink       Xlink
	TTL           int64 `xml:"ttl,attr"`
}

// XML Unmarshaler for RRAddr
func (r *RRCname) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type bufferRRCname RRCname
	finalRRCname := bufferRRCname{}
	err := d.DecodeElement(&finalRRCname, &start)
	if err != nil {
		return err
	}

	rr_xlink, err := parseXlink(finalRRCname.RRXLinkString)
	if err != nil {
		return err
	}
	finalRRCname.RRXlink = rr_xlink

	*r = RRCname(finalRRCname)
	return nil
}
