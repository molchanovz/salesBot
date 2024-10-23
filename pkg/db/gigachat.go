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
			Tables.GigachatMessage.Name: {{Column: Columns.GigachatMessage.ID, Direction: SortDesc}},
		},
		join: map[string][]string{
			Tables.GigachatMessage.Name: {TableColumns},
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

/*** GigachatMessage ***/

// FullGigachatMessage returns full joins with all columns
func (gr GigachatRepo) FullGigachatMessage() OpFunc {
	return WithColumns(gr.join[Tables.GigachatMessage.Name]...)
}

// DefaultGigachatMessageSort returns default sort.
func (gr GigachatRepo) DefaultGigachatMessageSort() OpFunc {
	return WithSort(gr.sort[Tables.GigachatMessage.Name]...)
}

// GigachatMessageByID is a function that returns GigachatMessage by ID(s) or nil.
func (gr GigachatRepo) GigachatMessageByID(ctx context.Context, id int, ops ...OpFunc) (*GigachatMessage, error) {
	return gr.OneGigachatMessage(ctx, &GigachatMessageSearch{ID: &id}, ops...)
}

// OneGigachatMessage is a function that returns one GigachatMessage by filters. It could return pg.ErrMultiRows.
func (gr GigachatRepo) OneGigachatMessage(ctx context.Context, search *GigachatMessageSearch, ops ...OpFunc) (*GigachatMessage, error) {
	obj := &GigachatMessage{}
	err := buildQuery(ctx, gr.db, obj, search, gr.filters[Tables.GigachatMessage.Name], PagerTwo, ops...).Select()

	if errors.Is(err, pg.ErrMultiRows) {
		return nil, err
	} else if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}

	return obj, err
}

// GigachatMessagesByFilters returns GigachatMessage list.
func (gr GigachatRepo) GigachatMessagesByFilters(ctx context.Context, search *GigachatMessageSearch, pager Pager, ops ...OpFunc) (gigachatMessages []GigachatMessage, err error) {
	err = buildQuery(ctx, gr.db, &gigachatMessages, search, gr.filters[Tables.GigachatMessage.Name], pager, ops...).Select()
	return
}

// CountGigachatMessages returns count
func (gr GigachatRepo) CountGigachatMessages(ctx context.Context, search *GigachatMessageSearch, ops ...OpFunc) (int, error) {
	return buildQuery(ctx, gr.db, &GigachatMessage{}, search, gr.filters[Tables.GigachatMessage.Name], PagerOne, ops...).Count()
}

// AddGigachatMessage adds GigachatMessage to DB.
func (gr GigachatRepo) AddGigachatMessage(ctx context.Context, gigachatMessage *GigachatMessage, ops ...OpFunc) (*GigachatMessage, error) {
	q := gr.db.ModelContext(ctx, gigachatMessage)
	applyOps(q, ops...)
	_, err := q.Insert()

	return gigachatMessage, err
}

// UpdateGigachatMessage updates GigachatMessage in DB.
func (gr GigachatRepo) UpdateGigachatMessage(ctx context.Context, gigachatMessage *GigachatMessage, ops ...OpFunc) (bool, error) {
	q := gr.db.ModelContext(ctx, gigachatMessage).WherePK()
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
func (gr GigachatRepo) DeleteGigachatMessage(ctx context.Context, id int) (deleted bool, err error) {
	gigachatMessage := &GigachatMessage{ID: id}

	res, err := gr.db.ModelContext(ctx, gigachatMessage).WherePK().Delete()
	if err != nil {
		return false, err
	}

	return res.RowsAffected() > 0, err
}
