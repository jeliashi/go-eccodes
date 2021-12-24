package accessor

const MaxSmartTableColumns int = 20

type CodeTableEntry struct {
	abbreviation string
	title        string
	units        string
}

type GribCodetable struct {
	filename        [2]string
	recomposed_name [2]string
	next            *GribCodetable
	size            uint32
	entries         *CodeTableEntry
}

type GribSmartTableEntry struct {
	abbreviation string
	column       [MaxSmartTableColumns]string
}

type GribSmartTable struct {
	filename        [3]string
	recomposed_name [3]string
	next            *GribSmartTable
	numberOfEntries uint32
	entries         GribSmartTableEntry
}
