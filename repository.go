package dbrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt" 
	"reflect"

	//"strings"

	"github.com/soyoshigure/dbrepo/option"
	//date "google.golang.org/genproto/googleapis/type/date"
)

type Repository[T any] interface {
	Select(ctx context.Context, opt *option.SQLSelectOption) (*T, error)

	SelectAll(ctx context.Context, opt *option.SQLSelectOption) ([]*T, error)

	Insert(ctx context.Context, value *T) (*T, error)

	Update(ctx context.Context, value *T) (*T, error)

	Delete(ctx context.Context, value *T) (error)
}

func newRepository[T any](tx *sql.Tx, table string) Repository[T]{
	return &repository[T]{
		tx: tx,
		table: table,
	}
}

type repository[T any] struct {
	tx *sql.Tx
	table string
}

func (*repository[T]) getColumnInfos(withoutIndex bool, withOutReadOnly bool) ([]*columnInfo){
	modelType := reflect.TypeOf((*T)(nil)).Elem()

	var columns []*columnInfo

	for i := 0; i < modelType.NumField(); i++{

		column := &columnInfo{

		}

		field := modelType.Field(i)

		

		column.Field = field.Name

		column.FieldType = field.Type

		if t := field.Tag.Get("column"); t != ""{
			column.Name = t
		}else{
			continue
		}

		if t := field.Tag.Get("type"); t != ""{
			column.Type = t
		}else{
			continue
		}

		if t := field.Tag.Get("defVal"); t != ""{
			column.DefaultValue = t
		}

		if t := field.Tag.Get("index"); t == "true"{
			column.IsIndex = true
		}else{
			column.IsIndex = false
		}

		if t := field.Tag.Get("readOnly"); t == "true"{
			column.IsReadOnly = true
		}else{
			column.IsReadOnly = false
		}

		if withoutIndex && column.IsIndex{
			continue
		}

		if withOutReadOnly && column.IsReadOnly{
			continue
		}
		
		columns = append(columns, column)
	}

	return columns
}

func (*repository[T]) getIndexColumn() *columnInfo{
	modelType := reflect.TypeOf((*T)(nil)).Elem()

	for i := 0; i < modelType.NumField(); i++{

		column := &columnInfo{

		}

		field := modelType.Field(i)
	

		column.Field = field.Name

		column.FieldType = field.Type

		if t := field.Tag.Get("column"); t != ""{
			column.Name = t
		}else{
			continue
		}

		if t := field.Tag.Get("type"); t != ""{
			column.Type = t
		}else{
			continue
		}

		if t := field.Tag.Get("defVal"); t != ""{
			column.DefaultValue = t
		}

		if t := field.Tag.Get("index"); t == "true"{
			column.IsIndex = true
			return column
		}else{
			column.IsIndex = false
		}
		
	}

	return nil
}

func (repo *repository[T]) Select(ctx context.Context, opt *option.SQLSelectOption) (*T, error){
	sql := "SELECT "

	columns := repo.getColumnInfos(false, false)

	for i, column := range columns{
		if i == 0{
			sql += fmt.Sprintf("%s", column.Name)
		}else{
			sql += fmt.Sprintf(", %s", column.Name)
		}
		
	}

	sql += fmt.Sprintf(" FROM %s", repo.table)

	if opt.WherePhrase != nil{
		where, err := opt.WherePhrase.ToSQL()
		if err != nil{
			return nil, err
		}
		sql += fmt.Sprintf(" WHERE %s", where)
	}

	if opt.OrderBy != nil {
		if opt.OrderBy.IsASC {
			sql += fmt.Sprintf(" ORDER BY %s ASC", opt.OrderBy.Column)
		} else {
			sql += fmt.Sprintf(" ORDER BY %s DESC", opt.OrderBy.Column)
		}
	}

	sql += " LIMIT 1"

	if opt.Offset != nil {
		sql += fmt.Sprintf(" OFFSET %d", opt.Offset)
	}

	row := repo.tx.QueryRowContext(ctx, sql)

	ptrs := make([]any, len(columns))
	vals := make([]reflect.Value, len(columns))

	for i, column := range columns{  
		if column.Type == "Json" || column.Type == "json"{
			/*if column.FieldType.Kind() == reflect.Pointer{
				vals[i] = reflect.ValueOf([]byte{})
				print(vals[i].IsNil())
			}else{
				vals[i] = reflect.New(column.FieldType).Elem()
			}*/
			var data string
			vals[i] = reflect.New(reflect.TypeOf(data)).Elem()
			ptrs[i] = vals[i].Addr().Interface()
		}else{
		if column.FieldType.Kind() == reflect.Pointer{
			vals[i] = reflect.New(column.FieldType.Elem()).Elem()
			ptrs[i] = vals[i].Interface()
		}else{
			vals[i] = reflect.New(column.FieldType).Elem()
			ptrs[i] = vals[i].Addr().Interface()
		}
	}
		
	}

	if err := row.Scan(ptrs...); err != nil{
		return nil, err
	}

	modelType := reflect.TypeOf((*T)(nil)).Elem()
	modelValue := reflect.New(modelType).Elem()

	for i, column := range columns{
		//modelValue.FieldByName(column.Field).Set(vals[i])
		if column.Type == "Json" || column.Type == "json"{
			if column.FieldType.Kind() == reflect.Pointer{
				//print(*ptrs[i].(*json.RawMessage))
				
				d := reflect.New(column.FieldType.Elem()).Interface()
				//err  := json.NewDecoder(bytes.NewReader(vals[i].Interface().(json.RawMessage))).Decode(d)
				print("---")
				print(vals[i].Interface().(string))
				print("---")
				err := json.Unmarshal([]byte(vals[i].Interface().(string)), d)

				if err != nil{
					return nil, err
				}
				reflect.Indirect(modelValue).FieldByName(column.Field).Set(reflect.ValueOf(d))
			}else{
				err := json.Unmarshal(*ptrs[i].(*json.RawMessage), vals[i].Addr().Interface())
				if err != nil{
					return nil, err
				}
				reflect.Indirect(modelValue).FieldByName(column.Field).Set(vals[i])
			}
			ptrs[i] = &json.RawMessage{}

			
		}else{
			if column.FieldType.Kind() == reflect.Pointer{
			}else{
				reflect.Indirect(modelValue).FieldByName(column.Field).Set(vals[i])
			}
		}
	}

	return modelValue.Addr().Interface().(*T), nil
}

func (repo *repository[T]) SelectAll(ctx context.Context, opt *option.SQLSelectOption) ([]*T, error){
	sql := "SELECT "

	columns := repo.getColumnInfos(false,false)

	for _, column := range columns{
		sql += fmt.Sprintf("%s, ", column.Name)
	}

	sql += fmt.Sprintf("FROM %s ", repo.table)

	if opt.WherePhrase != nil{
		where, err := opt.WherePhrase.ToSQL()
		if err != nil{
			return nil, err
		}
		sql += fmt.Sprintf(" WHERE %s", where)
	}

	if opt.OrderBy != nil {
		if opt.OrderBy.IsASC {
			sql += fmt.Sprintf(" ORDER BY %s ASC", opt.OrderBy.Column)
		} else {
			sql += fmt.Sprintf(" ORDER BY %s DESC", opt.OrderBy.Column)
		}
	}

	if opt.Limit != 0 {
		sql += fmt.Sprintf(" LIMIT %d", opt.Limit)
	}

	if opt.Offset != nil {
		sql += fmt.Sprintf(" OFFSET %d", opt.Offset)
	}

	rows, _ := repo.tx.QueryContext(ctx, sql)

	var results []*T

	for rows.Next(){
		ptrs := make([]any, len(columns))
		vals := make([]reflect.Value, len(columns))

		for i, column := range columns{
			vals[i] = reflect.New(column.FieldType)
			ptrs[i] = vals[i].Interface()
		}

		if err := rows.Scan(ptrs...); err != nil{
			return nil, err
		}

		modelType := reflect.TypeOf((*T)(nil)).Elem()
		modelValue := reflect.New(modelType).Elem()

		for i, column := range columns{
			modelValue.FieldByName(column.Field).Set(vals[i])
		}

		results = append(results, modelValue.Interface().(*T))
	
		
	}

	return results, nil
}

func (repo *repository[T]) Insert(ctx context.Context, value *T) (*T, error){
	columns := repo.getColumnInfos(true, true)

	modelValue := reflect.ValueOf(value)

	sql := fmt.Sprintf("INSERT INTO %s (", repo.table)

	sqlValues := ""
	for i, column := range columns{
		if i != 0 {
			sql += fmt.Sprintf(", %s", column.Name)
			sqlValues += ", ?"
		} else {
			sql += fmt.Sprintf("%s", column.Name)
			sqlValues += "?"
		}
	}

	sql += fmt.Sprintf(") VALUES (%s)", sqlValues)

	print(sql)

	vals := make([]any, len(columns))

	for i, column := range columns{
		print(column)
		//vals[i] = modelValue.FieldByName(column.Field).Interface()
		//Date型の場合かつgRPC用の型を使用している場合　いらない
		/*if (column.Type == "date" || column.Type == "Date") && column.FieldType == reflect.TypeOf(&date.Date{}){
			d := reflect.Indirect(modelValue).FieldByName(column.Field) .Interface().(*date.Date)

			ds, err := fmt.Printf("%d-%d-%d", d.Year, d.Month, d.Day)
			if err != nil{
				return nil, err
			}
			vals[i] = ds
		}else */if(column.Type == "json" || column.Type == "Json"){

			var js []byte
			var err error
			if column.FieldType.Kind() == reflect.Pointer{
				js, err = json.Marshal(reflect.Indirect(reflect.Indirect(modelValue).FieldByName(column.Field)).Interface())
			}else{
				js, err = json.Marshal(reflect.Indirect(modelValue).FieldByName(column.Field).Interface())
			}

			if err != nil{
				return nil, err
			}

			vals[i] = js
			
		}else{
			if column.FieldType.Kind() == reflect.Pointer{
				vals[i] = reflect.Indirect(reflect.Indirect(modelValue).FieldByName(column.Field)).Interface()
			}else{
				vals[i] = reflect.Indirect(modelValue).FieldByName(column.Field).Interface()
			}
			
		}
	}

	result, err := repo.tx.ExecContext(ctx, sql, *&vals...)

	if err != nil{
		return nil, err
	}

	index := repo.getIndexColumn()

	lastId, err := result.LastInsertId()

	if err != nil{
		return nil, err
	}

	opt := &option.SQLSelectOption{
		WherePhrase: &option.SQLEqualIntPhraseOption{
			Column: index.Name,
			Value: int(lastId),
		},
	}

	return repo.Select(ctx, opt)
}

func (repo *repository[T]) Update(ctx context.Context, value *T) (*T, error){
	columns := repo.getColumnInfos(false, true)
	index := repo.getIndexColumn()

	modelValue := reflect.ValueOf(value)

	sql := fmt.Sprintf("UPDATE %s SET ", repo.table)

	for i, column := range columns{
		if i != 0 {
			sql += fmt.Sprintf(", %s = ?", column.Name)
		} else {
			sql += fmt.Sprintf("%s = ?", column.Name)
		}
	}

	opt := &option.SQLEqualIntPhraseOption{
		Column: index.Name,
		Value: int(modelValue.FieldByName(index.Field).Int()),
	}

	where, _ := opt.ToSQL()

	sql += fmt.Sprintf(" WHERE %s", where)

	vals := make([]any, len(columns))

	for i, column := range columns{

		vals[i] = modelValue.FieldByName(column.Field).Interface()
		
	}

	_, _ = repo.tx.ExecContext(ctx, sql, *&vals...)

	selectOpt := &option.SQLSelectOption{
		WherePhrase: opt,
	}

	return repo.Select(ctx, selectOpt)
}

func (repo *repository[T]) Delete(ctx context.Context, value *T) (error){
	index := repo.getIndexColumn()

	modelValue := reflect.ValueOf(value)

	sql := fmt.Sprintf("DELETE FROM %s", repo.table)

	opt := &option.SQLEqualIntPhraseOption{
		Column: index.Name,
		Value: int(modelValue.FieldByName(index.Field).Int()),
	}

	where, _ := opt.ToSQL()

	sql += fmt.Sprintf(" WHERE %s", where)

	_, err := repo.tx.ExecContext(ctx, sql)

	return err

}

type columnInfo struct{
	Field string
	FieldType reflect.Type
	Name string
	Type string
	IsIndex bool
	IsReadOnly bool
	DefaultValue string
}