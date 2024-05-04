package option

type SQLSelectOption struct{
	WherePhrase SQLWhereOption
	OrderBy *SQLOrderByOption
	Limit int
	Offset *int
}