package netdot

type RRNsQuery struct {
	NsDName *string `url:"nsdname"`
	RR      *int64  `url:"rr"`
	TTL     *int64  `url:"ttl"`
}

func (q RRNsQuery) Builder() RRNsQueryBuilder {
	return RRNsQueryBuilder{query: q}
}

type RRNsQueryBuilder struct {
	query RRNsQuery
}

func NewRRNsQueryBuilder() *RRNsQueryBuilder {
	return &RRNsQueryBuilder{
		query: RRNsQuery{},
	}
}

func (b *RRNsQueryBuilder) Ns(nameserver_domain_name string) *RRNsQueryBuilder {
	b.query.NsDName = &nameserver_domain_name
	return b
}

func (b *RRNsQueryBuilder) RR(rr int64) *RRNsQueryBuilder {
	b.query.RR = &rr
	return b
}

func (b *RRNsQueryBuilder) TTL(ttl int64) *RRNsQueryBuilder {
	b.query.TTL = &ttl
	return b
}

func (b *RRNsQueryBuilder) Build() RRNsQuery {
	return b.query
}
