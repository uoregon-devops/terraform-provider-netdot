package netdot

type RRQuery struct {
	Active     *int64  `url:"active"`
	AutoUpdate *int64  `url:"auto_update"`
	Expiration *string `url:"expiration"`
	Info       *string `url:"info"`
	Name       *string `url:"name"`
	ZoneID     *int64  `url:"zone"`
}

type RRQueryBuilder struct {
	query RRQuery
}

func NewRRQueryBuilder() *RRQueryBuilder {
	return &RRQueryBuilder{
		query: RRQuery{},
	}
}

func (b *RRQueryBuilder) Active(active bool) *RRQueryBuilder {
	intBool := boolToInt(active)
	b.query.Active = &intBool
	return b
}

func (b *RRQueryBuilder) AutoUpdate(autoUpdate bool) *RRQueryBuilder {
	intBool := boolToInt(autoUpdate)
	b.query.AutoUpdate = &intBool
	return b
}

func (b *RRQueryBuilder) Expiration(expiration string) *RRQueryBuilder {
	b.query.Expiration = &expiration
	return b
}

func (b *RRQueryBuilder) Info(info string) *RRQueryBuilder {
	b.query.Info = &info
	return b
}

func (b *RRQueryBuilder) Name(name string) *RRQueryBuilder {
	b.query.Name = &name
	return b
}

func (b *RRQueryBuilder) ZoneID(zoneID int64) *RRQueryBuilder {
	b.query.ZoneID = &zoneID
	return b
}

func (b *RRQueryBuilder) Build() RRQuery {
	return b.query
}
