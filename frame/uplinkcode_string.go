// Code generated by "stringer -type=UplinkCode"; DO NOT EDIT.

package frame

import "strconv"

const (
	_UplinkCode_name_0 = "Device"
	_UplinkCode_name_1 = "Network"
	_UplinkCode_name_2 = "KeepaliveResponse"
	_UplinkCode_name_3 = "Data"
)

var (
	_UplinkCode_index_2 = [...]uint8{0, 9, 17}
)

func (i UplinkCode) String() string {
	switch {
	case i == 16:
		return _UplinkCode_name_0
	case i == 32:
		return _UplinkCode_name_1
	case 48 <= i && i <= 49:
		i -= 48
		return _UplinkCode_name_2[_UplinkCode_index_2[i]:_UplinkCode_index_2[i+1]]
	case i == 64:
		return _UplinkCode_name_3
	default:
		return "UplinkCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
