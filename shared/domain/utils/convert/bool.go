package convert

import "strconv"

func AsBool(src interface{}) (bool, error) {
	switch v := src.(type) {
	case bool:
		return v, nil
	case int:
		return v > 0, nil
	case int8:
		return v > 0, nil
	case int16:
		return v > 0, nil
	case int32:
		return v > 0, nil
	case int64:
		return v > 0, nil
	case uint:
		return v > 0, nil
	case uint8:
		return v > 0, nil
	case uint16:
		return v > 0, nil
	case uint32:
		return v > 0, nil
	case uint64:
		return v > 0, nil
	case float32:
		return v > 0, nil
	case float64:
		return v > 0, nil
	case []byte:
		if len(v) == 0 {
			return false, nil
		}

		if v[0] == 0x00 {
			return false, nil
		} else if v[0] == 0x01 {
			return true, nil
		}

		return strconv.ParseBool(string(v))
	case string:
		return strconv.ParseBool(v)
	}

	return false, ErrValueIncompatibleType
}
