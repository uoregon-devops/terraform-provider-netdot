package models

// <opt id="6577810709" address="184.171.15.38" asn="0" description="" first_seen="2025-03-24 13:08:20" info="" interface="0" last_seen="2025-03-24 13:08:20" monitored="0" owner="0" parent="184.171.0.0/20" parent_xlink="Ipblock/2711826724" prefix="32" rir="" status="Static" status_xlink="IpblockStatus/3" use_network_broadcast="0" used_by="0" version="4" vlan="0"/>

type IpBlockStatus struct {
	// computed
	ID   int64  `xml:"id,attr"`
	Name string `xml:"name,attr"`
}
