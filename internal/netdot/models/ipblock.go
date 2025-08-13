package models

import "encoding/xml"

// <opt id="6577810709" address="184.171.15.38" asn="0" description="" first_seen="2025-03-24 13:08:20" info="" interface="0" last_seen="2025-03-24 13:08:20" monitored="0" owner="0" parent="184.171.0.0/20" parent_xlink="Ipblock/2711826724" prefix="32" rir="" status="Static" status_xlink="IpblockStatus/3" use_network_broadcast="0" used_by="0" version="4" vlan="0"/>

type IpBlock struct {
	FirstSeen            string `xml:"first_seen,attr"`
	LastSeen             string `xml:"last_seen,attr"`
	ID                   int64  `xml:"id,attr"`
	Address              string `xml:"address,attr"`
	Prefix               int64  `xml:"prefix,attr"`
	Version              int64  `xml:"version,attr"`
	Status               string `xml:"status,attr"`
	StatusXLinkString    string `xml:"status_xlink,attr"`
	ASN                  int64  `xml:"asn,attr"`
	ASNXLinkString       string `xml:"asn_xlink,attr"`
	Description          string `xml:"description,attr"`
	Info                 string `xml:"info,attr"`
	Interface            string `xml:"interface,attr"`
	InterfaceXLinkString string `xml:"interface_xlink,attr"`
	Monitored            bool   `xml:"monitored,attr"`
	Owner                string `xml:"owner,attr"`
	OwnerXlinkString     string `xml:"owner_xlink,attr"`
	Parent               string `xml:"parent,attr"`
	ParentXLinkString    string `xml:"parent_xlink,attr"`
	RIR                  string `xml:"rir,attr"`
	UseNetworkBroadcast  bool   `xml:"use_network_broadcast,attr"`
	UsedBy               string `xml:"used_by,attr"`
	UsedByXLinkString    string `xml:"used_by_xlink,attr"`
	VLAN                 int64  `xml:"vlan,attr"`
	VLANXLinkString      string `xml:"vlan_xlink,attr"`
	AsnXlink             Xlink
	InterfaceXlink       Xlink
	OwnerXlink           Xlink
	ParentXlink          Xlink
	StatusXlink          Xlink
	UsedByXlink          Xlink
	VLANXlink            Xlink
}

// XML Unmarshaler for IpBlock
func (r *IpBlock) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type bufferIpBlock IpBlock
	finalIpBlock := bufferIpBlock{}
	err := d.DecodeElement(&finalIpBlock, &start)
	if err != nil {
		return err
	}

	asn_xlink, err := parseXlink(finalIpBlock.ASNXLinkString)
	if err != nil {
		return err
	}
	finalIpBlock.AsnXlink = asn_xlink

	interface_xlink, err := parseXlink(finalIpBlock.InterfaceXLinkString)
	if err != nil {
		return err
	}
	finalIpBlock.InterfaceXlink = interface_xlink

	parent_xlink, err := parseXlink(finalIpBlock.ParentXLinkString)
	if err != nil {
		return err
	}
	finalIpBlock.ParentXlink = parent_xlink

	status_xlink, err := parseXlink(finalIpBlock.StatusXLinkString)
	if err != nil {
		return err
	}
	finalIpBlock.StatusXlink = status_xlink

	owner_xlink, err := parseXlink(finalIpBlock.OwnerXlinkString)
	if err != nil {
		return err
	}
	finalIpBlock.OwnerXlink = owner_xlink

	used_by_xlink, err := parseXlink(finalIpBlock.UsedByXLinkString)
	if err != nil {
		return err
	}
	finalIpBlock.UsedByXlink = used_by_xlink

	vlan_xlink, err := parseXlink(finalIpBlock.VLANXLinkString)
	if err != nil {
		return err
	}
	finalIpBlock.VLANXlink = vlan_xlink

	*r = IpBlock(finalIpBlock)
	return nil
}
