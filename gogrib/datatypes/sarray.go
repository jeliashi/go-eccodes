package datatypes

type GribStringList struct {
	Value string
	Count int
	Next  *GribStringList
}
