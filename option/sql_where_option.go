package option

import (
	"fmt"
)

type SQLWhereOption interface{
	ToSQL() (string, error)
}

type SQLEqualStringPhraseOption struct{
	Column string
	Value string
}

func (opt *SQLEqualStringPhraseOption) ToSQL() (string, error){
	return fmt.Sprintf("(%s = \"%s\")", opt.Column, opt.Value), nil
}

type SQLLikeStringPhraseOption struct{
	Column string
	Value string
}

func (opt *SQLLikeStringPhraseOption) ToSQL() (string, error){
	return fmt.Sprintf("(%s LIKE \"%%%s%%\")", opt.Column, opt.Value), nil
}

type SQLEqualIntPhraseOption struct{
	SQLWhereOption
	Column string
	Value int
}

func (opt *SQLEqualIntPhraseOption) ToSQL() (string, error){
	return fmt.Sprintf("(%s = %d)", opt.Column, opt.Value), nil
}

type SQLNotEqualStringPhraseOption struct{
	Column string
	Value string
}

func (opt *SQLNotEqualStringPhraseOption) ToSQL() (string, error){
	return fmt.Sprintf("(%s != \"%s\")", opt.Column, opt.Value), nil
}

type SQLNotEqualIntPhraseOption struct{
	Column string
	Value int
}

func (opt *SQLNotEqualIntPhraseOption) ToSQL() (string, error){
	return fmt.Sprintf("(%s != %d)", opt.Column, opt.Value), nil
}

type SQLGreaterThanIntPhraseOption struct{
	Column string
	Value int
}

func (opt *SQLGreaterThanIntPhraseOption) ToSQL() (string, error){
	return fmt.Sprintf("(%s > %d)", opt.Column, opt.Value), nil
}

type SQLGreaterThanOrEqualIntPhraseOption struct{
	Column string
	Value int
}

func (opt *SQLGreaterThanOrEqualIntPhraseOption) ToSQL() (string, error){
	return fmt.Sprintf("(%s >= %d)", opt.Column, opt.Value), nil
}

type SQLLessThanIntPhraseOption struct{
	Column string
	Value int
}

func (opt *SQLLessThanIntPhraseOption) ToSQL() (string, error){
	return fmt.Sprintf("(%s < %d)", opt.Column, opt.Value), nil
}

type SQLLessThanOrEqualIntPhraseOption struct{
	Column string
	Value int
}

func (opt *SQLLessThanOrEqualIntPhraseOption) ToSQL() (string, error){
	return fmt.Sprintf("(%s <= %d)", opt.Column, opt.Value), nil
}

type SQLAndPhraseOption struct{
	WherePhraseA SQLWhereOption
	WherePhraseB SQLWhereOption
}

func (opt *SQLAndPhraseOption) ToSQL() (string, error){
	phraseA, err := opt.WherePhraseA.ToSQL()
	if err != nil{
		return "", err
	}

	phraseB, err := opt.WherePhraseB.ToSQL()
	if err != nil{
		return "", err
	}
	return fmt.Sprintf("(%s AND %s)", phraseA, phraseB), nil
}

type SQLOrPhraseOption struct{
	WherePhraseA SQLWhereOption
	WherePhraseB SQLWhereOption
}

func (opt *SQLOrPhraseOption) ToSQL() (string, error){
	phraseA, err := opt.WherePhraseA.ToSQL()
	if err != nil{
		return "", err
	}

	phraseB, err := opt.WherePhraseB.ToSQL()
	if err != nil{
		return "", err
	}
	return fmt.Sprintf("(%s OR %s)", phraseA, phraseB), nil
}