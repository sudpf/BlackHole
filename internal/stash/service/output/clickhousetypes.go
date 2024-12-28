package output

import (
	"errors"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

func toClickhouseType(value interface{}, valueType string) (interface{}, error) {
	switch valueType {
	case "Float32":
		return cast.ToFloat32E(value)
	case "Float64":
		return cast.ToFloat64E(value)
	case "Int8":
		return cast.ToInt8E(value)
	case "Int16":
		return cast.ToInt16E(value)
	case "Int32":
		return cast.ToInt32E(value)
	case "Int64":
		return cast.ToInt64E(value)
	case "UInt8":
		return cast.ToUint8E(value)
	case "UInt16":
		return cast.ToUint16E(value)
	case "UInt32":
		return cast.ToUint32E(value)
	case "UInt64":
		return cast.ToUint64E(value)
	case "IPv4", "IPv6":
		ip, _, err := net.ParseCIDR(value.(string))
		return ip, err
	case "Bool", "Boolean":
		return cast.ToBoolE(value)
	case "Date", "Date32", "DateTime":
		return cast.ToTimeE(value)
	case "DateTime64(3)":
		// 处理 DateTime64(3) 类型
		// 假设输入值为字符串或时间戳（秒数）
		switch v := value.(type) {
		case string:
			t, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, err
			}
			return t, nil
		case int64:
			// 处理时间戳（毫秒）
			return time.Unix(0, v*int64(time.Millisecond)).UTC(), nil
		case float64:
			// 处理时间戳（秒），假设浮点数代表秒
			return time.Unix(int64(v), 0).UTC(), nil
		default:
			return nil, errors.New("unsupported value type for DateTime64(3)")
		}
	case "UUID":
		return uuid.Parse(value.(string))
	case "String":
		return cast.ToStringE(value)
	case "varchar":
		return cast.ToStringE(value)
	case "int":
		return cast.ToInt32E(value)
	default:
		return "", errors.New("unsupported type:" + valueType)
	}
}
