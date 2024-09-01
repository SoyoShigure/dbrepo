package dbrepo

import(
	"github.com/soyoshigure/dbrepo/option"
	"reflect"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	repositories map[reflect.Type]*option.RepositoryOption
)

func RegisterRepository[T any](opt option.DatabaseOption, table string){
	if repositories == nil{
		repositories =  map[reflect.Type]*option.RepositoryOption{}
	}

	modelType := reflect.TypeOf((*T)(nil)).Elem()

	repoOpt := &option.RepositoryOption{
		DBOpt: opt,
		Table: table,
	}

	repositories[modelType] = repoOpt
}

func Do[T any](ctx context.Context, fn func(ctx context.Context, repo Repository[T]) error) error{

	//TODO:　NoSuchElement
	modelType := reflect.TypeOf((*T)(nil)).Elem()

	repoOpt := repositories[modelType]

	//sql.DBのインスタンス化
	dbOpt := repoOpt.DBOpt

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", dbOpt.User, dbOpt.Password, dbOpt.Host, dbOpt.Port, dbOpt.Name)

	db, err := sql.Open("mysql", dsn)

	if(err != nil){
		//TODO: 例外追加
		return nil
	}

	tx, err := db.BeginTx(ctx, nil)

	if(err != nil){
		//TODO: 例外追加
		return nil
	}

	repo := newRepository[T](tx, repoOpt.Table)

	if err = fn(ctx, repo); err != nil{
		if txErr := tx.Rollback(); txErr != nil{
			return txErr
		}
		return err
	}

	if err = tx.Commit(); err != nil{
		return err
	}

	return nil
}