/*
Go version of eccodes `src/grib_context.c`
*/
package core

import (
	"fmt"
	"gogrib/accessor"
	"gogrib/datatypes"
	errs "gogrib/errors"
	"gogrib/expression"
	"io"
	"os"
	"path"
	"reflect"
)

// TODO: determine if pthread or openmp logic should be added here
// TODO: determine if manage_mem logic should be added here
const MaxSetValues int = 10
const MaxNumConcepts int = 2000
const MaxNumHashArray int = 2000
const ECC_PATH_MAXLEN int = 8192

var DEFAULT_FILE_POOL_MAX_OPENED_FILES = 0

type RealMode int8

const (
	GribRealMode8 RealMode = iota
)

type BufrdcMode int8

const (
	BufrdcMode8 BufrdcMode = iota
)

type GribMultiSupport struct {
	offset                uint32
	message               []byte
	message_length        uint32
	sections              [8]string
	bitmap_section        string
	bitmap_section_length uint32
	sections_length       [9]uint32
	section_number        uint8
	next                  *GribMultiSupport
}

type GribConceptCondition struct {
	name       string
	expression *expression.GribExpression
	iarray     *datatypes.GribIarray
	next       *GribConceptCondition
}

type GribConceptValue struct {
	name       string
	conditions *GribConceptCondition
	index      *datatypes.GribTrie
	next       *GribConceptValue
}

type GribContext struct {
	inited                     bool
	debug                      bool
	write_on_fail              bool
	no_abort                   bool
	io_buffer_size             int
	no_big_group_split         bool
	no_spd                     bool
	keep_matrix                bool
	grib_definition_files_path string
	grib_samples_path          string
	grib_concept_path          string
	grib_reader                []expression.GribActionFileList
	user_data                  interface{}
	real_mode                  RealMode

	codetable                           *accessor.GribCodetable
	smart_table                         *accessor.GribSmartTable
	outfilename                         string
	mult_support_on                     bool
	multi_support                       *GribMultiSupport
	grib_definition_files_dir           *datatypes.GribStringList
	handle_file_count                   int
	handle_total_count                  int
	message_file_offset                 int64
	no_fail_on_wrong_length             bool
	gts_header_on                       bool
	gribex_mode_on                      bool
	large_constant_fields               bool
	keys                                []datatypes.GribItrie
	keys_count                          int
	concepts_index                      []datatypes.GribItrie
	concepts_count                      int
	concepts                            [MaxNumConcepts]GribConceptValue
	hash_array_index                    datatypes.GribItrie
	hash_array_count                    int
	hash_array                          [MaxNumHashArray]datatypes.GribHashArrayValue
	def_files                           *datatypes.GribTrie
	blocklist                           *datatypes.GribStringList
	ieee_packing                        int
	bufrdc_mode                         BufrdcMode
	bufr_set_to_missing_if_out_of_range int
	bufr_multi_element_constant_arrays  int
	grib_data_quality_checks            int
	log_stream                          io.Writer
	classes                             *datatypes.GribTrie
	lists                               *datatypes.GribTrie
	expanded_descriptors                *datatypes.GribTrie
	file_pool_max_opened_files          int
}
type GribContextInterface interface {
	// free_mem(data interface)
	// alloc_mem(size uint32)
	// realloc_mem(data interface, size uint32)
	// free_persistent_mem(data interface)
	// alloc_persistent_mem(size uint32)
	// free_buffer_mem(data interface)
	// alloc_buffer_mem(size uint32)
	// realloc_buffer_mem(data interface, size uint32)

	// kinda feel like these should have methods contained within io
	// read(ptr interface{}, size uint32, stream interface{}) uint32
	// write(ptr interface{}, size uint32, stream interface{}) uint32
	// tell(stream interface{}) uint32
	// seek(offset uint32, whence int, stream interface{})
	// eof(stream interface{}) int

	log(level int, msg string)
	print(descriptor interface{}, msg string)
}

func GetContextDefault() (GribContext, error) {
	definitions_path, ok := os.LookupEnv("ECCODES_DEFINITION_PATH")
	if !ok {
		definitions_path = "/usr/local/Cellar/eccodes/2.24.0/share/eccodes/definitions"
	}
	gc := GribContext{real_mode: GribRealMode8, grib_definition_files_path: definitions_path}
	return gc, nil
}

func (gc *GribContext) initDefinitionFilesDir() errs.Error {
	if gc.grib_definition_files_dir != nil {
		return nil
	}
	if len(gc.grib_definition_files_path) == 0 {
		return errs.NewContextError(errs.WithMessage("Grib no Definition"))
	}
	defPathList := splitEccodeString(gc.grib_definition_files_path)
	if len(defPathList) == 1 {
		gc.grib_definition_files_dir = &datatypes.GribStringList{Value: path.Clean(gc.grib_definition_files_path)}
	} else {
		var next *datatypes.GribStringList
		for _, dir := range defPathList {
			if next != nil {
				next = next.Next
			} else {
				next = gc.grib_definition_files_dir
			}
			next.Value = path.Clean(dir)
		}
	}
	return nil

}

func (gc *GribContext) GribContextFullDefsPath(s string) (string, errs.Error) {
	if s == string('/') || s == string('.') {
		return s, nil
	}
	// TODO: figure out how this works...
	// This is referenced from here:
	// https://github.com/ecmwf/eccodes/blob/5407ecd6d19afad0db729a73043b2a31843140c9/src/grib_context.c#L697
	fullpath := gc.def_files.GribTrieGet(s)
	if fullpath != nil {
		t := reflect.TypeOf(fullpath)
		v, has := t.FieldByName("value")
		if has {
			return v.Type.String(), nil
		}
		return string(""), errs.NewContextError(errs.WithMessage("unable to ascertain value from GribTrie"))
	}

	if gc.grib_definition_files_dir == nil {
		err := gc.initDefinitionFilesDir()
		if err != nil {
			// wrap above error?
			return s, errs.NewContextError(errs.WithMessage("unable to find definition files directory"))
		}
	}
	var dir *datatypes.GribStringList
	dir = gc.grib_definition_files_dir
	for dir != nil {
		full := fmt.Sprintf("%s/%s", (*dir).Value, s)
		isReadable := checkReadAccess(full)
		if isReadable {
			fullpath := &datatypes.GribStringList{Value: string(full)}
			gc.def_files.Insert(s, fullpath)
			fmt.Printf("Found def file %s", full)
			return fullpath.Value, nil
		}
		dir = dir.Next
	}
	gc.def_files.Insert(s, &datatypes.GribStringList{})
	return s, errs.NewContextError(errs.WithMessage("grib file not found"))
}
func (gc *GribContext) GribParseFile(s string) {

}
func (gc *GribContext) SetGribContextReader() errs.Error {
	if gc.grib_reader != nil {
		return nil
	}
	fpath, err := gc.GribContextFullDefsPath("boot.def")
	if err != nil {
		// wrap above error?
		return errs.NewContextError(
			errs.WithMessage(
				fmt.Sprintf("Unable to find boot.def. Context path=%s\nCheck ECCOES_DEFINITION_PATH", gc.grib_definition_files_path),
			),
		)
	}
	expression.GribParseFile(gc.grib_reader, fpath)

	return nil
}
