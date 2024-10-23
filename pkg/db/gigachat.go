package db

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type GigachatRepo struct {
	db      orm.DB
	filters map[string][]Filter
	sort    map[string][]SortField
	join    map[string][]string
}

// NewGigachatRepo returns new repository
func NewGigachatRepo(db orm.DB) GigachatRepo {
	return GigachatRepo{
		db:      db,
		filters: map[string][]Filter{},
		sort: map[string][]SortField{
			Tables.Gigachatmessage.Name:       {{Column: Columns.Gigachatmessage.ID, Direction: SortDesc}},
			Tables.StudupGigachatmessage.Name: {{Column: Columns.StudupGigachatmessage.ID, Direction: SortDesc}},
		},
		join: map[string][]string{
			Tables.Gigachatmessage.Name:       {TableColumns},
			Tables.StudupGigachatmessage.Name: {TableColumns},
		},
	}
}

// WithTransaction is a function that wraps GigachatRepo with pg.Tx transaction.
func (gr GigachatRepo) WithTransaction(tx *pg.Tx) GigachatRepo {
	gr.db = tx
	return gr
}

// WithEnabledOnly is a function that adds "statusId"=1 as base filter.
func (gr GigachatRepo) WithEnabledOnly() GigachatRepo {
	f := make(map[string][]Filter, len(gr.filters))
	for i := range gr.filters {
		f[i] = make([]Filter, len(gr.filters[i]))
		copy(f[i], gr.filters[i])
		f[i] = append(f[i], StatusEnabledFilter)
	}
	gr.filters = f

	return gr
}

/*** Gigachatmessage ***/

// FullGigachatmessage returns full joins with all columns
func (gr GigachatRepo) FullGigachatmessage() OpFunc {
	return WithColumns(gr.join[Tables.Gigachatmessage.Name]...)
}

// DefaultGigachatmessageSort returns default sort.
func (gr GigachatRepo) DefaultGigachatmessageSort() OpFunc {
	return WithSort(gr.sort[Tables.Gigachatmessage.Name]...)
}

// GigachatmessageByID is a function that returns Gigachatmessage by ID(s) or nil.
func (gr GigachatRepo) GigachatmessageByID(ctx context.Context, id int, ops ...OpFunc) (*Gigachatmessage, error) {
	return gr.OneGigachatmessage(ctx, &GigachatmessageSearch{ID: &id}, ops...)
}

// OneGigachatmessage is a function that returns one Gigachatmessage by filters. It could return pg.ErrMultiRows.
func (gr GigachatRepo) OneGigachatmessage(ctx context.Context, search *GigachatmessageSearch, ops ...OpFunc) (*Gigachatmessage, error) {
	obj := &Gigachatmessage{}
	err := buildQuery(ctx, gr.db, obj, search, gr.filters[Tables.Gigachatmessage.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// GigachatmessagesByFilters returns Gigachatmessage list.
func (gr GigachatRepo) GigachatmessagesByFilters(ctx context.Context, search *GigachatmessageSearch, pager Pager, ops ...OpFunc) (gigachatmessages []Gigachatmessage, err error) {
	err = buildQuery(ctx, gr.db, &gigachatmessages, search, gr.filters[Tables.Gigachatmessage.Name], pager, ops...).Select()
	return
}

// CountGigachatmessages returns count
func (gr GigachatRepo) CountGigachatmessages(ctx context.Context, search *GigachatmessageSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, gr.db, &Gigachatmessage{}, search, gr.filters[Tables.Gigachatmessage.Name], PagerOne, ops...).Count()
}

// AddGigachatmessage adds Gigachatmessage to DB.
func (gr GigachatRepo) AddGigachatmessage(ctx context.Context, gigachatmessage *Gigachatmessage, ops ...OpFunc) (*Gigachatmessage, error) {
	q := gr.db.ModelContext(ctx, gigachatmessage)
	applyOps(q, ops...)
	_, err := q.Insert()

	return gigachatmessage, err
}

// UpdateGigachatmessage updates Gigachatmessage in DB.
func (gr GigachatRepo) UpdateGigachatmessage(ctx context.Context, gigachatmessage *Gigachatmessage, ops ...OpFunc) (bool, error) {
	q := gr.db.ModelContext(ctx, gigachatmessage).WherePK()
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
func (gr GigachatRepo) DeleteGigachatmessage(ctx context.Context, id int) (deleted bool, err error) {
	gigachatmessage := &Gigachatmessage{ID: id}

	res, err := gr.db.ModelContext(ctx, gigachatmessage).WherePK().Delete()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

/*** StudupGigachatmessage ***/

// FullStudupGigachatmessage returns full joins with all columns
func (gr GigachatRepo) FullStudupGigachatmessage() OpFunc {
	return WithColumns(gr.join[Tables.StudupGigachatmessage.Name]...)
}

// DefaultStudupGigachatmessageSort returns default sort.
func (gr GigachatRepo) DefaultStudupGigachatmessageSort() OpFunc {
	return WithSort(gr.sort[Tables.StudupGigachatmessage.Name]...)
}

// StudupGigachatmessageByID is a function that returns StudupGigachatmessage by ID(s) or nil.
func (gr GigachatRepo) StudupGigachatmessageByID(ctx context.Context, id int, ops ...OpFunc) (*StudupGigachatmessage, error) {
	return gr.OneStudupGigachatmessage(ctx, &StudupGigachatmessageSearch{ID: &id}, ops...)
}

// OneStudupGigachatmessage is a function that returns one StudupGigachatmessage by filters. It could return pg.ErrMultiRows.
func (gr GigachatRepo) OneStudupGigachatmessage(ctx context.Context, search *StudupGigachatmessageSearch, ops ...OpFunc) (*StudupGigachatmessage, error) {
	obj := &StudupGigachatmessage{}
	err := buildQuery(ctx, gr.db, obj, search, gr.filters[Tables.StudupGigachatmessage.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// StudupGigachatmessagesByFilters returns StudupGigachatmessage list.
func (gr GigachatRepo) StudupGigachatmessagesByFilters(ctx context.Context, search *StudupGigachatmessageSearch, pager Pager, ops ...OpFunc) (studupGigachatmessages []StudupGigachatmessage, err error) {
	err = buildQuery(ctx, gr.db, &studupGigachatmessages, search, gr.filters[Tables.StudupGigachatmessage.Name], pager, ops...).Select()
	return
}

// CountStudupGigachatmessages returns count
func (gr GigachatRepo) CountStudupGigachatmessages(ctx context.Context, search *StudupGigachatmessageSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, gr.db, &StudupGigachatmessage{}, search, gr.filters[Tables.StudupGigachatmessage.Name], PagerOne, ops...).Count()
}

// AddStudupGigachatmessage adds StudupGigachatmessage to DB.
func (gr GigachatRepo) AddStudupGigachatmessage(ctx context.Context, studupGigachatmessage *StudupGigachatmessage, ops ...OpFunc) (*StudupGigachatmessage, error) {
	q := gr.db.ModelContext(ctx, studupGigachatmessage)
	applyOps(q, ops...)
	_, err := q.Insert()

	return studupGigachatmessage, err
}

// UpdateStudupGigachatmessage updates StudupGigachatmessage in DB.
func (gr GigachatRepo) UpdateStudupGigachatmessage(ctx context.Context, studupGigachatmessage *StudupGigachatmessage, ops ...OpFunc) (bool, error) {
	q := gr.db.ModelContext(ctx, studupGigachatmessage).WherePK()
	if len(ops) == 0 {
		q = q.ExcludeColumn(Columns.StudupGigachatmessage.ID)
	}
	applyOps(q, ops...)
	res, err := q.Update()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}

// DeleteStudupGigachatmessage deletes StudupGigachatmessage from DB.
func (gr GigachatRepo) DeleteStudupGigachatmessage(ctx context.Context, id int) (deleted bool, err error) {
	studupGigachatmessage := &StudupGigachatmessage{ID: id}

	res, err := gr.db.ModelContext(ctx, studupGigachatmessage).WherePK().Delete()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}
