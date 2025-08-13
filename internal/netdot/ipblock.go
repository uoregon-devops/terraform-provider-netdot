package netdot

import (
	"fmt"
	"net/http"
	"net/netip"
	"terraform-provider-netdot/internal/netdot/models"
)

// address, prefix, parent, version, status, info, description

type IpBlockQuery struct {
	Address                *string `url:"address"`
	Prefix                 *int64  `url:"prefix"`
	Version                *int64  `url:"version"`
	Status                 *string `url:"status"`
	ASNID                  *int64  `url:"asn"`
	Description            *string `url:"description"`
	Info                   *string `url:"info"`
	InterfaceID            *int64  `url:"interface"`
	Monitored              *int64  `url:"monitored"`
	OwnerID                *int64  `url:"owner"`
	ParentID               *int64  `url:"parent"`
	RIR                    *string `url:"rir"`
	UseNetworkBroadcast    *int64  `url:"use_network_broadcast"`
	UsedByID               *int64  `url:"used_by"`
	VLANID                 *int64  `url:"vlan"`
	SkipInheritParentOwner *int64  `url:"skip_inherit_parent_owner,omitempty"`
	SkipReserveFirstN      *int64  `url:"skip_reserve_first_n,omitempty"`
	SkipUpdateTree         *int64  `url:"no_update_tree,omitempty"`
	Validate               *int64  `url:"validate,omitempty"`
}

func (q IpBlockQuery) Builder() IpBlockQueryBuilder {
	return IpBlockQueryBuilder{query: q}
}

type IpBlockQueryBuilder struct {
	query IpBlockQuery
}

func NewIpBlockQueryBuilder() IpBlockQueryBuilder {
	return IpBlockQueryBuilder{query: IpBlockQuery{}}
}

// build function for every field

func (b *IpBlockQueryBuilder) SkipInheritParentOwner(skip bool) *IpBlockQueryBuilder {
	boolInt := boolToInt(skip)
	b.query.SkipInheritParentOwner = &boolInt
	return b
}

func (b *IpBlockQueryBuilder) SkipReserveFirstN(skip bool) *IpBlockQueryBuilder {
	boolInt := boolToInt(skip)
	b.query.SkipReserveFirstN = &boolInt
	return b
}

func (b *IpBlockQueryBuilder) SkipUpdateTree(skip bool) *IpBlockQueryBuilder {
	boolInt := boolToInt(skip)
	b.query.SkipUpdateTree = &boolInt
	return b
}

func (b *IpBlockQueryBuilder) Validate(validate bool) *IpBlockQueryBuilder {
	boolInt := boolToInt(validate)
	b.query.Validate = &boolInt
	return b
}

func (b *IpBlockQueryBuilder) Address(address string) *IpBlockQueryBuilder {
	b.query.Address = &address
	return b
}

func (b *IpBlockQueryBuilder) Prefix(prefix int64) *IpBlockQueryBuilder {
	b.query.Prefix = &prefix
	return b
}

func (b *IpBlockQueryBuilder) Version(version int64) *IpBlockQueryBuilder {
	b.query.Version = &version
	return b
}

func (b *IpBlockQueryBuilder) Status(status string) *IpBlockQueryBuilder {
	b.query.Status = &status
	return b
}

func (b *IpBlockQueryBuilder) ASNID(asnID int64) *IpBlockQueryBuilder {
	b.query.ASNID = &asnID
	return b
}

func (b *IpBlockQueryBuilder) Description(description string) *IpBlockQueryBuilder {
	b.query.Description = &description
	return b
}

func (b *IpBlockQueryBuilder) Info(info string) *IpBlockQueryBuilder {
	b.query.Info = &info
	return b
}

func (b *IpBlockQueryBuilder) InterfaceID(interfaceID int64) *IpBlockQueryBuilder {
	b.query.InterfaceID = &interfaceID
	return b
}

func (b *IpBlockQueryBuilder) Monitored(monitor bool) *IpBlockQueryBuilder {
	intbool := boolToInt(monitor)
	b.query.Monitored = &intbool
	return b
}

func (b *IpBlockQueryBuilder) OwnerID(ownerID int64) *IpBlockQueryBuilder {
	b.query.OwnerID = &ownerID
	return b
}

func (b *IpBlockQueryBuilder) ParentID(parentID int64) *IpBlockQueryBuilder {
	b.query.ParentID = &parentID
	return b
}

func (b *IpBlockQueryBuilder) RIR(rir string) *IpBlockQueryBuilder {
	b.query.RIR = &rir
	return b
}

func (b *IpBlockQueryBuilder) UseNetworkBroadcast(useNetworkBroadcast bool) *IpBlockQueryBuilder {
	intbool := boolToInt(useNetworkBroadcast)
	b.query.UseNetworkBroadcast = &intbool
	return b
}

func (b *IpBlockQueryBuilder) UsedByID(usedByID int64) *IpBlockQueryBuilder {
	b.query.UsedByID = &usedByID
	return b
}

func (b *IpBlockQueryBuilder) VLANID(vlanID int64) *IpBlockQueryBuilder {
	b.query.VLANID = &vlanID
	return b
}
func (b *IpBlockQueryBuilder) Build() IpBlockQuery {
	return b.query
}

type ipAllocationStrategy int

const (
	IPAllocationStrategyFirstFree ipAllocationStrategy = iota
	// IPAllocationStrategyLastFree
)

// get next available IP in Subnet IPblock, this could be
func (c *Client) GetNextAvailableIP(subnetID int64, strategy ipAllocationStrategy) (*int64, string, error) {
	if subnetID <= 0 {
		return nil, "", fmt.Errorf("invalid subnetID")
	}

	switch strategy {
	case IPAllocationStrategyFirstFree:
		return c.getFirstFreeIP(subnetID)
	// case IPAllocationStrategyLastFree:
	// 	return getLastFreeIP(subnetID)
	default:
		return nil, "", fmt.Errorf("invalid strategy")
	}
}

func (c *Client) getFirstFreeIP(subnetID int64) (*int64, string, error) {
	if subnetID <= 0 {
		return nil, "", fmt.Errorf("invalid subnetID")
	}

	var subnetIpBlock models.IpBlock
	statusCode, err := c.GetResourceByID("ipblock", subnetID, &subnetIpBlock)
	if err != nil {
		return nil, "", fmt.Errorf("Error reading IP block", err.Error())
	}

	var existingIPs struct {
		IpBlocks []models.IpBlock `xml:"Ipblock"`
	}

	statusCode, err = c.Get(fmt.Sprintf("/rest/ipblock?parent=%d", subnetID), &existingIPs)
	if err != nil {
		if statusCode != nil && *statusCode != http.StatusNotFound {
			return nil, "", fmt.Errorf("Error reading IP block", err.Error())
		}
	}

	parentIP, err := netip.ParseAddr(subnetIpBlock.Address)
	if err != nil {
		return nil, "", fmt.Errorf("Error parsing parent IP address", err.Error())
	}

	parentPrefix := netip.PrefixFrom(parentIP, int(subnetIpBlock.Prefix))

	// start one address from gateway, just in case
	currentAddress := parentPrefix.Addr().Next().Next()

	for parentPrefix.Contains(currentAddress.Next()) {
		var existingIP *models.IpBlock = nil
		for _, ip := range existingIPs.IpBlocks {
			if ip.Address == currentAddress.String() {
				existingIP = &ip
				break
			}
		}
		if existingIP == nil {
			return nil, currentAddress.String(), nil
		}

		if *&existingIP.Status == "Available" {
			return &existingIP.ID, existingIP.Address, nil
		}
		currentAddress = currentAddress.Next()
	}

	return nil, "", fmt.Errorf("No available IP address found")
}
