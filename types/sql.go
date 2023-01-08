package types

import "database/sql"

func NewNullString(s string) sql.NullString {
	res := sql.NullString{
		String: s,
		Valid:  true,
	}
	if s == "" {
		res.Valid = false
	}
	return res
}
