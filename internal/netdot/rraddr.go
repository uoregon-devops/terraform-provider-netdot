package netdot

type RRAddrQuery struct {
	IpBlock        *int64 `url:"ipblock"`
	RR             *int64 `url:"rr"`
	TTL            *int64 `url:"ttl"`
	NoChangeStatus *int64 `url:"no_change_status,omitempty"`
	SkipDeletingRR *int64 `url:"skip_deleting_rr,omitempty"`
}

func (q RRAddrQuery) Builder() RRAddrQueryBuilder {
	return RRAddrQueryBuilder{query: q}
}

type RRAddrQueryBuilder struct {
	query RRAddrQuery
}

func NewRRAddrQueryBuilder() *RRAddrQueryBuilder {
	return &RRAddrQueryBuilder{
		query: RRAddrQuery{},
	}
}

func (b *RRAddrQueryBuilder) NoChangeStatus(skip bool) *RRAddrQueryBuilder {
	boolInt := boolToInt(skip)
	b.query.NoChangeStatus = &boolInt
	return b
}

func (b *RRAddrQueryBuilder) SkipDeletingRR(skip bool) *RRAddrQueryBuilder {
	boolInt := boolToInt(skip)
	b.query.SkipDeletingRR = &boolInt
	return b
}

func (b *RRAddrQueryBuilder) IpBlock(ipBlock int64) *RRAddrQueryBuilder {
	b.query.IpBlock = &ipBlock
	return b
}

func (b *RRAddrQueryBuilder) RR(rr int64) *RRAddrQueryBuilder {
	b.query.RR = &rr
	return b
}

func (b *RRAddrQueryBuilder) TTL(ttl int64) *RRAddrQueryBuilder {
	b.query.TTL = &ttl
	return b
}

func (b *RRAddrQueryBuilder) Build() RRAddrQuery {
	return b.query
}
