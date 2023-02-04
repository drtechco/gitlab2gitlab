// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"drtech.co/gl2gl/orm/model"
)

func newErrorCode(db *gorm.DB) errorCode {
	_errorCode := errorCode{}

	_errorCode.errorCodeDo.UseDB(db)
	_errorCode.errorCodeDo.UseModel(&model.ErrorCode{})

	tableName := _errorCode.errorCodeDo.TableName()
	_errorCode.ALL = field.NewField(tableName, "*")
	_errorCode.Code = field.NewInt32(tableName, "code")
	_errorCode.I18nKey = field.NewString(tableName, "i18n_key")
	_errorCode.Memo = field.NewString(tableName, "memo")

	_errorCode.fillFieldMap()

	return _errorCode
}

type errorCode struct {
	errorCodeDo errorCodeDo

	ALL     field.Field
	Code    field.Int32
	I18nKey field.String
	Memo    field.String

	fieldMap map[string]field.Expr
}

func (e errorCode) Table(newTableName string) *errorCode {
	e.errorCodeDo.UseTable(newTableName)
	return e.updateTableName(newTableName)
}

func (e errorCode) As(alias string) *errorCode {
	e.errorCodeDo.DO = *(e.errorCodeDo.As(alias).(*gen.DO))
	return e.updateTableName(alias)
}

func (e *errorCode) updateTableName(table string) *errorCode {
	e.ALL = field.NewField(table, "*")
	e.Code = field.NewInt32(table, "code")
	e.I18nKey = field.NewString(table, "i18n_key")
	e.Memo = field.NewString(table, "memo")

	e.fillFieldMap()

	return e
}

func (e *errorCode) WithContext(ctx context.Context) *errorCodeDo {
	return e.errorCodeDo.WithContext(ctx)
}

func (e errorCode) TableName() string { return e.errorCodeDo.TableName() }

func (e errorCode) Alias() string { return e.errorCodeDo.Alias() }

func (e *errorCode) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := e.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (e *errorCode) fillFieldMap() {
	e.fieldMap = make(map[string]field.Expr, 3)
	e.fieldMap["code"] = e.Code
	e.fieldMap["i18n_key"] = e.I18nKey
	e.fieldMap["memo"] = e.Memo
}

func (e errorCode) clone(db *gorm.DB) errorCode {
	e.errorCodeDo.ReplaceDB(db)
	return e
}

type errorCodeDo struct{ gen.DO }

func (e errorCodeDo) Debug() *errorCodeDo {
	return e.withDO(e.DO.Debug())
}

func (e errorCodeDo) WithContext(ctx context.Context) *errorCodeDo {
	return e.withDO(e.DO.WithContext(ctx))
}

func (e errorCodeDo) ReadDB() *errorCodeDo {
	return e.Clauses(dbresolver.Read)
}

func (e errorCodeDo) WriteDB() *errorCodeDo {
	return e.Clauses(dbresolver.Write)
}

func (e errorCodeDo) Clauses(conds ...clause.Expression) *errorCodeDo {
	return e.withDO(e.DO.Clauses(conds...))
}

func (e errorCodeDo) Returning(value interface{}, columns ...string) *errorCodeDo {
	return e.withDO(e.DO.Returning(value, columns...))
}

func (e errorCodeDo) Not(conds ...gen.Condition) *errorCodeDo {
	return e.withDO(e.DO.Not(conds...))
}

func (e errorCodeDo) Or(conds ...gen.Condition) *errorCodeDo {
	return e.withDO(e.DO.Or(conds...))
}

func (e errorCodeDo) Select(conds ...field.Expr) *errorCodeDo {
	return e.withDO(e.DO.Select(conds...))
}

func (e errorCodeDo) Where(conds ...gen.Condition) *errorCodeDo {
	return e.withDO(e.DO.Where(conds...))
}

func (e errorCodeDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *errorCodeDo {
	return e.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (e errorCodeDo) Order(conds ...field.Expr) *errorCodeDo {
	return e.withDO(e.DO.Order(conds...))
}

func (e errorCodeDo) Distinct(cols ...field.Expr) *errorCodeDo {
	return e.withDO(e.DO.Distinct(cols...))
}

func (e errorCodeDo) Omit(cols ...field.Expr) *errorCodeDo {
	return e.withDO(e.DO.Omit(cols...))
}

func (e errorCodeDo) Join(table schema.Tabler, on ...field.Expr) *errorCodeDo {
	return e.withDO(e.DO.Join(table, on...))
}

func (e errorCodeDo) LeftJoin(table schema.Tabler, on ...field.Expr) *errorCodeDo {
	return e.withDO(e.DO.LeftJoin(table, on...))
}

func (e errorCodeDo) RightJoin(table schema.Tabler, on ...field.Expr) *errorCodeDo {
	return e.withDO(e.DO.RightJoin(table, on...))
}

func (e errorCodeDo) Group(cols ...field.Expr) *errorCodeDo {
	return e.withDO(e.DO.Group(cols...))
}

func (e errorCodeDo) Having(conds ...gen.Condition) *errorCodeDo {
	return e.withDO(e.DO.Having(conds...))
}

func (e errorCodeDo) Limit(limit int) *errorCodeDo {
	return e.withDO(e.DO.Limit(limit))
}

func (e errorCodeDo) Offset(offset int) *errorCodeDo {
	return e.withDO(e.DO.Offset(offset))
}

func (e errorCodeDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *errorCodeDo {
	return e.withDO(e.DO.Scopes(funcs...))
}

func (e errorCodeDo) Unscoped() *errorCodeDo {
	return e.withDO(e.DO.Unscoped())
}

func (e errorCodeDo) Create(values ...*model.ErrorCode) error {
	if len(values) == 0 {
		return nil
	}
	return e.DO.Create(values)
}

func (e errorCodeDo) CreateInBatches(values []*model.ErrorCode, batchSize int) error {
	return e.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (e errorCodeDo) Save(values ...*model.ErrorCode) error {
	if len(values) == 0 {
		return nil
	}
	return e.DO.Save(values)
}

func (e errorCodeDo) First() (*model.ErrorCode, error) {
	if result, err := e.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.ErrorCode), nil
	}
}

func (e errorCodeDo) Take() (*model.ErrorCode, error) {
	if result, err := e.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.ErrorCode), nil
	}
}

func (e errorCodeDo) Last() (*model.ErrorCode, error) {
	if result, err := e.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.ErrorCode), nil
	}
}

func (e errorCodeDo) Find() ([]*model.ErrorCode, error) {
	result, err := e.DO.Find()
	return result.([]*model.ErrorCode), err
}

func (e errorCodeDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ErrorCode, err error) {
	buf := make([]*model.ErrorCode, 0, batchSize)
	err = e.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (e errorCodeDo) FindInBatches(result *[]*model.ErrorCode, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return e.DO.FindInBatches(result, batchSize, fc)
}

func (e errorCodeDo) Attrs(attrs ...field.AssignExpr) *errorCodeDo {
	return e.withDO(e.DO.Attrs(attrs...))
}

func (e errorCodeDo) Assign(attrs ...field.AssignExpr) *errorCodeDo {
	return e.withDO(e.DO.Assign(attrs...))
}

func (e errorCodeDo) Joins(fields ...field.RelationField) *errorCodeDo {
	for _, _f := range fields {
		e = *e.withDO(e.DO.Joins(_f))
	}
	return &e
}

func (e errorCodeDo) Preload(fields ...field.RelationField) *errorCodeDo {
	for _, _f := range fields {
		e = *e.withDO(e.DO.Preload(_f))
	}
	return &e
}

func (e errorCodeDo) FirstOrInit() (*model.ErrorCode, error) {
	if result, err := e.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.ErrorCode), nil
	}
}

func (e errorCodeDo) FirstOrCreate() (*model.ErrorCode, error) {
	if result, err := e.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.ErrorCode), nil
	}
}

func (e errorCodeDo) FindByPage(offset int, limit int) (result []*model.ErrorCode, count int64, err error) {
	result, err = e.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = e.Offset(-1).Limit(-1).Count()
	return
}

func (e errorCodeDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = e.Count()
	if err != nil {
		return
	}

	err = e.Offset(offset).Limit(limit).Scan(result)
	return
}

func (e errorCodeDo) Scan(result interface{}) (err error) {
	return e.DO.Scan(result)
}

func (e *errorCodeDo) withDO(do gen.Dao) *errorCodeDo {
	e.DO = *do.(*gen.DO)
	return e
}