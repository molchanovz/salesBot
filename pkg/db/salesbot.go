package db

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type SalesbotRepo struct {
	db      orm.DB
	filters map[string][]Filter
	sort    map[string][]SortField
	join    map[string][]string
}

// NewSalesbotRepo returns new repository
func NewSalesbotRepo(db orm.DB) SalesbotRepo {
	return SalesbotRepo{
		db:      db,
		filters: map[string][]Filter{},
		sort: map[string][]SortField{
			Tables.GigachatMessage.Name: {{Column: Columns.GigachatMessage.ID, Direction: SortDesc}},
		},
		join: map[string][]string{
			Tables.GigachatMessage.Name: {TableColumns},
		},
	}
}

// WithTransaction is a function that wraps SalesbotRepo with pg.Tx transaction.
func (sr SalesbotRepo) WithTransaction(tx *pg.Tx) SalesbotRepo {
	sr.db = tx
	return sr
}

// WithEnabledOnly is a function that adds "statusId"=1 as base filter.
func (sr SalesbotRepo) WithEnabledOnly() SalesbotRepo {
	f := make(map[string][]Filter, len(sr.filters))
	for i := range sr.filters {
		f[i] = make([]Filter, len(sr.filters[i]))
		copy(f[i], sr.filters[i])
		f[i] = append(f[i], StatusEnabledFilter)
	}
	sr.filters = f

	return sr
}

/*** GigachatMessage ***/

// FullGigachatMessage returns full joins with all columns
func (sr SalesbotRepo) FullGigachatMessage() OpFunc {
	return WithColumns(sr.join[Tables.GigachatMessage.Name]...)
}

// DefaultGigachatMessageSort returns default sort.
func (sr SalesbotRepo) DefaultGigachatMessageSort() OpFunc {
	return WithSort(sr.sort[Tables.GigachatMessage.Name]...)
}

// GigachatMessageByID is a function that returns GigachatMessage by ID(s) or nil.
func (sr SalesbotRepo) GigachatMessageByID(ctx context.Context, id int, ops ...OpFunc) (*GigachatMessage, error) {
	return sr.OneGigachatMessage(ctx, &GigachatMessageSearch{ID: &id}, ops...)
}

// OneGigachatMessage is a function that returns one GigachatMessage by filters. It could return pg.ErrMultiRows.
func (sr SalesbotRepo) OneGigachatMessage(ctx context.Context, search *GigachatMessageSearch, ops ...OpFunc) (*GigachatMessage, error) {
	obj := &GigachatMessage{}
	err := buildQuery(ctx, sr.db, obj, search, sr.filters[Tables.GigachatMessage.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// GigachatMessagesByFilters returns GigachatMessage list.
func (sr SalesbotRepo) GigachatMessagesByFilters(ctx context.Context, search *GigachatMessageSearch, pager Pager, ops ...OpFunc) (gigachatMessages []GigachatMessage, err error) {
	err = buildQuery(ctx, sr.db, &gigachatMessages, search, sr.filters[Tables.GigachatMessage.Name], pager, ops...).Select()
	return
}

// CountGigachatMessages returns count
func (sr SalesbotRepo) CountGigachatMessages(ctx context.Context, search *GigachatMessageSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, sr.db, &GigachatMessage{}, search, sr.filters[Tables.GigachatMessage.Name], PagerOne, ops...).Count()
}

// AddGigachatMessage adds GigachatMessage to DB.
func (sr SalesbotRepo) AddGigachatMessage(ctx context.Context, gigachatMessage *GigachatMessage, ops ...OpFunc) (*GigachatMessage, error) {
	q := sr.db.ModelContext(ctx, gigachatMessage)
	applyOps(q, ops...)
	_, err := q.Insert()

	return gigachatMessage, err
}

// UpdateGigachatMessage updates GigachatMessage in DB.
func (sr SalesbotRepo) UpdateGigachatMessage(ctx context.Context, gigachatMessage *GigachatMessage, ops ...OpFunc) (bool, error) {
	q := sr.db.ModelContext(ctx, gigachatMessage).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.GigachatMessage.ID)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteGigachatMessage deletes GigachatMessage from DB.
func (sr SalesbotRepo) DeleteGigachatMessage(ctx context.Context, id int) (deleted bool, err error) {
	gigachatMessage := &GigachatMessage{ID: id}

	res, err := sr.db.ModelContext(ctx, gigachatMessage).WherePK().Delete()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}
