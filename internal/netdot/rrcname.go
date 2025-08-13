package netdot

type RRCnameQuery struct {
	Cname          *string `url:"cname"`
	RR             *int64  `url:"rr"`
	TTL            *int64  `url:"ttl"`
	SkipDeletingRR *int64  `url:"skip_deleting_rr,omitempty"`
}

func (q RRCnameQuery) Builder() RRCnameQueryBuilder {
	return RRCnameQueryBuilder{query: q}
}

type RRCnameQueryBuilder struct {
	query RRCnameQuery
}

func NewRRCnameQueryBuilder() *RRCnameQueryBuilder {
	return &RRCnameQueryBuilder{
		query: RRCnameQuery{},
	}
}

func (b *RRCnameQueryBuilder) SkipDeletingRR(skip bool) *RRCnameQueryBuilder {
	boolInt := boolToInt(skip)
	b.query.SkipDeletingRR = &boolInt
	return b
}

func (b *RRCnameQueryBuilder) Cname(cname string) *RRCnameQueryBuilder {
	b.query.Cname = &cname
	return b
}

func (b *RRCnameQueryBuilder) RR(rr int64) *RRCnameQueryBuilder {
	b.query.RR = &rr
	return b
}

func (b *RRCnameQueryBuilder) TTL(ttl int64) *RRCnameQueryBuilder {
	b.query.TTL = &ttl
	return b
}

func (b *RRCnameQueryBuilder) Build() RRCnameQuery {
	return b.query
}
