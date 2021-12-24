package accessor

type GribValuesType int8

type GribSection struct {
	owner    *GribAccessor
	aclength *GribAccessor
	block    *GribBlockOfAccessors
	length   uint32
	padding  uint32
}

type GribDependency struct {
	next     *GribDependency
	observed *GribAccessor
	observer *GribAccessor
	run      int
}

type GribValues struct {
	name         string
	Type         GribValuesType
	long_value   int64
	double_value int64
	string_value string
	err          error
	has_value    bool
	equal        bool
	next         *GribValues
}
