package json_utils

const (
	bom0 = 0xef
	bom1 = 0xbb
	bom2 = 0xbf
)

// Remove BOM from byte array.
func RemoveBom(b []byte) []byte {
	if len(b) > 3 && b[0] == bom0 && b[1] == bom1 && b[2] == bom2 {
		return b[3:]
	}
	return b
}