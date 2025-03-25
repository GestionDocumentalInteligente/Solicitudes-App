package authconn

import "database/sql"

func getInt64FromNullInt64(ni sql.NullInt64) *int64 {
	if ni.Valid {
		value := ni.Int64
		return &value
	}
	return nil
}
