package models

import "encoding/xml"

// <RRPTR
// 	id="448896"
// 	ipblock="184.171.2.2"
// 	ipblock_xlink="Ipblock/5405116861"
// 	ptrdname="mollman-test-record.uoregon.edu"
// 	rr="2.2.171.184.in-addr.arpa"
// 	rr_xlink="RR/999290"
// 	ttl="86400"
// />

type RRPtr struct {
	ID                 int64  `xml:"id,attr"`
	IpBlock            string `xml:"ipblock,attr"`
	IpBlockXLinkString string `xml:"ipblock_xlink,attr"`
	PTRdname           string `xml:"ptrdname,attr"`
	RR                 string `xml:"rr,attr"`
	RRXLinkString      string `xml:"rr_xlink,attr"`
	TTL                int64  `xml:"ttl,attr"`
	IpBlockXlink       Xlink
	RRXlink            Xlink
}

// XML Unmarshaler for RRPtr
func (r *RRPtr) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type bufferRRPtr RRPtr
	finalRRPtr := bufferRRPtr{}
	err := d.DecodeElement(&finalRRPtr, &start)
	if err != nil {
		return err
	}

	ipblock_xlink, err := parseXlink(finalRRPtr.IpBlockXLinkString)
	if err != nil {
		return err
	}
	finalRRPtr.IpBlockXlink = ipblock_xlink

	rr_xlink, err := parseXlink(finalRRPtr.RRXLinkString)
	if err != nil {
		return err
	}
	finalRRPtr.RRXlink = rr_xlink

	*r = RRPtr(finalRRPtr)
	return nil
}
