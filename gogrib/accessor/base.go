package accessor

import "gogrib/datatypes"

const MaxAccessorNames int = 20
const MaxAccessorAttributes int = 20
const AccessorArraySize int = 5000

type GribAccessor struct {
	name        string        // < name of the accessor
	name_space  string        // < namespace to which the accessor belongs
	length      int64         // < byte length of the accessor
	offset      int64         // < offset of the data in the buffer
	parent      *GribSection  //  < section to which the accessor is attached
	next        *GribAccessor //  < next accessor in list
	previous    *GribAccessor //  < next accessor in list
	flags       uint64        //  < Various flags
	sub_section *GribSection

	all_names       [MaxAccessorNames]string // < name of the accessor
	all_name_spaces [MaxAccessorNames]string // < namespace to which the accessor belongs
	dirty           bool

	same                *GribAccessor               // < accessors with the same name
	loop                int64                       // < used in lists
	bufr_subset_number  int64                       // < bufr subset (bufr data accessors belong to different subsets)
	bufr_group_number   int64                       // < used in bufr
	vvalue              *datatypes.GribVirtualValue // < virtual value used when transient flag on *
	set                 string
	attributes          *[MaxAccessorAttributes]GribAccessor // < attributes are accessors
	parent_as_attribute *GribAccessor
}

type GribBlockOfAccessors struct {
	first *GribAccessor
	last  *GribAccessor
}

type GribLoader struct {
	data             interface{}
	list_is_resized  int
	changing_edition int
}
