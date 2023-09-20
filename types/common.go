package types

type QueryFilter struct {
	Key   string
	Value string
}

type QueryFilterMap map[string]string

type Count struct {
	Cnt uint64 `json:"cnt" db:"cnt"`
}
