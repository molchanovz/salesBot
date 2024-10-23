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
			Tables.Gigachatmessage.Name: {{Column: Columns.Gigachatmessage.ID, Direction: SortDesc}},
		},
		join: map[string][]string{
			Tables.Gigachatmessage.Name: {TableColumns},
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

/*** Gigachatmessage ***/

// FullGigachatmessage returns full joins with all columns
func (sr SalesbotRepo) FullGigachatmessage() OpFunc {
	return WithColumns(sr.join[Tables.Gigachatmessage.Name]...)
}

// DefaultGigachatmessageSort returns default sort.
func (sr SalesbotRepo) DefaultGigachatmessageSort() OpFunc {
	return WithSort(sr.sort[Tables.Gigachatmessage.Name]...)
}

// GigachatmessageByID is a function that returns Gigachatmessage by ID(s) or nil.
func (sr SalesbotRepo) GigachatmessageByID(ctx context.Context, id int, ops ...OpFunc) (*Gigachatmessage, error) {
	return sr.OneGigachatmessage(ctx, &GigachatmessageSearch{ID: &id}, ops...)
}

// OneGigachatmessage is a function that returns one Gigachatmessage by filters. It could return pg.ErrMultiRows.
func (sr SalesbotRepo) OneGigachatmessage(ctx context.Context, search *GigachatmessageSearch, ops ...OpFunc) (*Gigachatmessage, error) {
	obj := &Gigachatmessage{}
	err := buildQuery(ctx, sr.db, obj, search, sr.filters[Tables.Gigachatmessage.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// GigachatmessagesByFilters returns Gigachatmessage list.
func (sr SalesbotRepo) GigachatmessagesByFilters(ctx context.Context, search *GigachatmessageSearch, pager Pager, ops ...OpFunc) (gigachatmessages []Gigachatmessage, err error) {
	err = buildQuery(ctx, sr.db, &gigachatmessages, search, sr.filters[Tables.Gigachatmessage.Name], pager, ops...).Select()
	return
}

// CountGigachatmessages returns count
func (sr SalesbotRepo) CountGigachatmessages(ctx context.Context, search *GigachatmessageSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, sr.db, &Gigachatmessage{}, search, sr.filters[Tables.Gigachatmessage.Name], PagerOne, ops...).Count()
}

// AddGigachatmessage adds Gigachatmessage to DB.
func (sr SalesbotRepo) AddGigachatmessage(ctx context.Context, gigachatmessage *Gigachatmessage, ops ...OpFunc) (*Gigachatmessage, error) {
	q := sr.db.ModelContext(ctx, gigachatmessage)
	applyOps(q, ops...)
	_, err := q.Insert()

	return gigachatmessage, err
}

// UpdateGigachatmessage updates Gigachatmessage in DB.
func (sr SalesbotRepo) UpdateGigachatmessage(ctx context.Context, gigachatmessage *Gigachatmessage, ops ...OpFunc) (bool, error) {
	q := sr.db.ModelContext(ctx, gigachatmessage).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.Gigachatmessage.ID)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteGigachatmessage deletes Gigachatmessage from DB.
func (sr SalesbotRepo) DeleteGigachatmessage(ctx context.Context, id int) (deleted bool, err error) {
	gigachatmessage := &Gigachatmessage{ID: id}

	res, err := sr.db.ModelContext(ctx, gigachatmessage).WherePK().Delete()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}
