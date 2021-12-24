package datatypes

type GribHashArrayValue struct {
	next   *GribHashArrayValue
	name   string
	Type   int
	iarray *GribIarray
	darray *GribDarray
	index  *GribTrie
}
