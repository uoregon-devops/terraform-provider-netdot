package models

import "encoding/xml"

// <RRADDR id="544859" ipblock="184.171.15.38" ipblock_xlink="Ipblock/6577810709" rr="mollman-test-record.uoregon.edu" rr_xlink="RR/999289" ttl="800"/>

type RRAddr struct {
	ID                 int64  `xml:"id,attr"`
	IpBlock            string `xml:"ipblock,attr"`
	IpBlockXLinkString string `xml:"ipblock_xlink,attr"`
	IpBlockXlink       Xlink
	RR                 string `xml:"rr,attr"`
	RRXLinkString      string `xml:"rr_xlink,attr"`
	RRXlink            Xlink
	TTL                int64 `xml:"ttl,attr"`
}

// XML Unmarshaler for RRAddr
func (r *RRAddr) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type bufferRRAddr RRAddr
	finalRRAddr := bufferRRAddr{}
	err := d.DecodeElement(&finalRRAddr, &start)
	if err != nil {
		return err
	}

	ipblock_xlink, err := parseXlink(finalRRAddr.IpBlockXLinkString)
	if err != nil {
		return err
	}
	finalRRAddr.IpBlockXlink = ipblock_xlink

	rr_xlink, err := parseXlink(finalRRAddr.RRXLinkString)
	if err != nil {
		return err
	}
	finalRRAddr.RRXlink = rr_xlink

	*r = RRAddr(finalRRAddr)
	return nil
}
