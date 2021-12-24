package datatypes

type GribVirtualValue struct {
	lval    int64
	dval    float64
	cval    byte
	missing int
	length  uint32
	Type    int
}
