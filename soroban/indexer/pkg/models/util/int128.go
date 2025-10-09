package util

import (
	"database/sql/driver"
	"fmt"
	"math/big"
)

type Int128 struct {
	big.Int
}

// Value implements the driver.Valuer interface for database serialization
func (i Int128) Value() (driver.Value, error) {
	return i.String(), nil
}

// Scan implements the sql.Scanner interface for database deserialization
func (i *Int128) Scan(value interface{}) error {
	var strValue string
	switch t := value.(type) {
	case []byte:
		strValue = string(t)
	case string:
		strValue = t
	default:
		return fmt.Errorf("unsupported type for Int128: %T", value)
	}

	int128 := new(big.Int)
	_, ok := int128.SetString(strValue, 10)
	if !ok {
		return fmt.Errorf("failed to parse Int128 value: %s", strValue)
	}

	i.Int = *int128
	return nil
}
