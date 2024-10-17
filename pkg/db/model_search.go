// Code generated by mfd-generator v0.4.0; DO NOT EDIT.

//nolint:all
//lint:file-ignore U1000 ignore unused code, it's generated
package db

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

const condition = "?.? = ?"

// base filters
type applier func(query *orm.Query) (*orm.Query, error)

type search struct {
	appliers []applier
}

func (s *search) apply(query *orm.Query) {
	for _, applier := range s.appliers {
		query.Apply(applier)
	}
}

func (s *search) where(query *orm.Query, table, field string, value interface{}) {
	query.Where(condition, pg.Ident(table), pg.Ident(field), value)
}

func (s *search) WithApply(a applier) {
	if s.appliers == nil {
		s.appliers = []applier{}
	}
	s.appliers = append(s.appliers, a)
}

func (s *search) With(condition string, params ...interface{}) {
	s.WithApply(func(query *orm.Query) (*orm.Query, error) {
		return query.Where(condition, params...), nil
	})
}

// Searcher is interface for every generated filter
type Searcher interface {
	Apply(query *orm.Query) *orm.Query
	Q() applier

	With(condition string, params ...interface{})
	WithApply(a applier)
}

type UserSearch struct {
	search

	ID                 *int
	CreatedAt          *time.Time
	Login              *string
	Password           *string
	AuthKey            *string
	LastActivityAt     *time.Time
	StatusID           *int
	IDs                []int
	NotID              *int
	LoginILike         *string
	PasswordILike      *string
	AuthKeyILike       *string
	LastActivityAtFrom *time.Time
	LastActivityAtTo   *time.Time
}

func (us *UserSearch) Apply(query *orm.Query) *orm.Query {
	if us == nil {
		return query
	}
	if us.ID != nil {
		us.where(query, Tables.User.Alias, Columns.User.ID, us.ID)
	}
	if us.CreatedAt != nil {
		us.where(query, Tables.User.Alias, Columns.User.CreatedAt, us.CreatedAt)
	}
	if us.Login != nil {
		us.where(query, Tables.User.Alias, Columns.User.Login, us.Login)
	}
	if us.Password != nil {
		us.where(query, Tables.User.Alias, Columns.User.Password, us.Password)
	}
	if us.AuthKey != nil {
		us.where(query, Tables.User.Alias, Columns.User.AuthKey, us.AuthKey)
	}
	if us.LastActivityAt != nil {
		us.where(query, Tables.User.Alias, Columns.User.LastActivityAt, us.LastActivityAt)
	}
	if us.StatusID != nil {
		us.where(query, Tables.User.Alias, Columns.User.StatusID, us.StatusID)
	}
	if len(us.IDs) > 0 {
		Filter{Columns.User.ID, us.IDs, SearchTypeArray, false}.Apply(query)
	}
	if us.NotID != nil {
		Filter{Columns.User.ID, *us.NotID, SearchTypeEquals, true}.Apply(query)
	}
	if us.LoginILike != nil {
		Filter{Columns.User.Login, *us.LoginILike, SearchTypeILike, false}.Apply(query)
	}
	if us.PasswordILike != nil {
		Filter{Columns.User.Password, *us.PasswordILike, SearchTypeILike, false}.Apply(query)
	}
	if us.AuthKeyILike != nil {
		Filter{Columns.User.AuthKey, *us.AuthKeyILike, SearchTypeILike, false}.Apply(query)
	}
	if us.LastActivityAtFrom != nil {
		Filter{Columns.User.LastActivityAt, *us.LastActivityAtFrom, SearchTypeGE, false}.Apply(query)
	}
	if us.LastActivityAtTo != nil {
		Filter{Columns.User.LastActivityAt, *us.LastActivityAtTo, SearchTypeLE, false}.Apply(query)
	}

	us.apply(query)

	return query
}

func (us *UserSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		if us == nil {
			return query, nil
		}
		return us.Apply(query), nil
	}
}

type VfsFileSearch struct {
	search

	ID            *int
	FolderID      *int
	Title         *string
	Path          *string
	Params        *string
	IsFavorite    *bool
	MimeType      *string
	FileSize      *int
	FileExists    *bool
	CreatedAt     *time.Time
	StatusID      *int
	IDs           []int
	TitleILike    *string
	PathILike     *string
	ParamsILike   *string
	MimeTypeILike *string
}

func (vfs *VfsFileSearch) Apply(query *orm.Query) *orm.Query {
	if vfs == nil {
		return query
	}
	if vfs.ID != nil {
		vfs.where(query, Tables.VfsFile.Alias, Columns.VfsFile.ID, vfs.ID)
	}
	if vfs.FolderID != nil {
		vfs.where(query, Tables.VfsFile.Alias, Columns.VfsFile.FolderID, vfs.FolderID)
	}
	if vfs.Title != nil {
		vfs.where(query, Tables.VfsFile.Alias, Columns.VfsFile.Title, vfs.Title)
	}
	if vfs.Path != nil {
		vfs.where(query, Tables.VfsFile.Alias, Columns.VfsFile.Path, vfs.Path)
	}
	if vfs.Params != nil {
		vfs.where(query, Tables.VfsFile.Alias, Columns.VfsFile.Params, vfs.Params)
	}
	if vfs.IsFavorite != nil {
		vfs.where(query, Tables.VfsFile.Alias, Columns.VfsFile.IsFavorite, vfs.IsFavorite)
	}
	if vfs.MimeType != nil {
		vfs.where(query, Tables.VfsFile.Alias, Columns.VfsFile.MimeType, vfs.MimeType)
	}
	if vfs.FileSize != nil {
		vfs.where(query, Tables.VfsFile.Alias, Columns.VfsFile.FileSize, vfs.FileSize)
	}
	if vfs.FileExists != nil {
		vfs.where(query, Tables.VfsFile.Alias, Columns.VfsFile.FileExists, vfs.FileExists)
	}
	if vfs.CreatedAt != nil {
		vfs.where(query, Tables.VfsFile.Alias, Columns.VfsFile.CreatedAt, vfs.CreatedAt)
	}
	if vfs.StatusID != nil {
		vfs.where(query, Tables.VfsFile.Alias, Columns.VfsFile.StatusID, vfs.StatusID)
	}
	if len(vfs.IDs) > 0 {
		Filter{Columns.VfsFile.ID, vfs.IDs, SearchTypeArray, false}.Apply(query)
	}
	if vfs.TitleILike != nil {
		Filter{Columns.VfsFile.Title, *vfs.TitleILike, SearchTypeILike, false}.Apply(query)
	}
	if vfs.PathILike != nil {
		Filter{Columns.VfsFile.Path, *vfs.PathILike, SearchTypeILike, false}.Apply(query)
	}
	if vfs.ParamsILike != nil {
		Filter{Columns.VfsFile.Params, *vfs.ParamsILike, SearchTypeILike, false}.Apply(query)
	}
	if vfs.MimeTypeILike != nil {
		Filter{Columns.VfsFile.MimeType, *vfs.MimeTypeILike, SearchTypeILike, false}.Apply(query)
	}

	vfs.apply(query)

	return query
}

func (vfs *VfsFileSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		if vfs == nil {
			return query, nil
		}
		return vfs.Apply(query), nil
	}
}

type VfsFolderSearch struct {
	search

	ID             *int
	ParentFolderID *int
	Title          *string
	IsFavorite     *bool
	CreatedAt      *time.Time
	StatusID       *int
	IDs            []int
	TitleILike     *string
}

func (vfs *VfsFolderSearch) Apply(query *orm.Query) *orm.Query {
	if vfs == nil {
		return query
	}
	if vfs.ID != nil {
		vfs.where(query, Tables.VfsFolder.Alias, Columns.VfsFolder.ID, vfs.ID)
	}
	if vfs.ParentFolderID != nil {
		vfs.where(query, Tables.VfsFolder.Alias, Columns.VfsFolder.ParentFolderID, vfs.ParentFolderID)
	}
	if vfs.Title != nil {
		vfs.where(query, Tables.VfsFolder.Alias, Columns.VfsFolder.Title, vfs.Title)
	}
	if vfs.IsFavorite != nil {
		vfs.where(query, Tables.VfsFolder.Alias, Columns.VfsFolder.IsFavorite, vfs.IsFavorite)
	}
	if vfs.CreatedAt != nil {
		vfs.where(query, Tables.VfsFolder.Alias, Columns.VfsFolder.CreatedAt, vfs.CreatedAt)
	}
	if vfs.StatusID != nil {
		vfs.where(query, Tables.VfsFolder.Alias, Columns.VfsFolder.StatusID, vfs.StatusID)
	}
	if len(vfs.IDs) > 0 {
		Filter{Columns.VfsFolder.ID, vfs.IDs, SearchTypeArray, false}.Apply(query)
	}
	if vfs.TitleILike != nil {
		Filter{Columns.VfsFolder.Title, *vfs.TitleILike, SearchTypeILike, false}.Apply(query)
	}

	vfs.apply(query)

	return query
}

func (vfs *VfsFolderSearch) Q() applier {
	return func(query *orm.Query) (*orm.Query, error) {
		if vfs == nil {
			return query, nil
		}
		return vfs.Apply(query), nil
	}
}
