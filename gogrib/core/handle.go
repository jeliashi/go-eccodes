package core

import (
	"bytes"
	"fmt"
	"gogrib/accessor"
	"gogrib/datatypes"
	errs "gogrib/errors"
)

const MaxNumSections int = 12

type HeaderMode int8

const (
	replaceMe HeaderMode = iota
)

type ProductKind int8

const (
	ProductAny ProductKind = iota
	ProductGrib
	ProductBufr
	ProductMetar
	ProductGts
	ProductTaf
)

type GribHandle struct {
	buffer             *accessor.GribBuffer
	root               *accessor.GribSection
	asserts            *accessor.GribSection
	dependencies       *accessor.GribDependency
	main               *GribHandle
	kid                *GribHandle
	loader             *accessor.GribLoader
	values_stack       int
	values             [MaxSetValues]accessor.GribValues
	values_count       [MaxSetValues]uint32
	dont_trigger       bool
	partial            bool
	header_mode        HeaderMode
	gts_header         string
	gts_header_len     uint32
	use_trie           bool
	trie_invalid       bool
	accessors          [accessor.AccessorArraySize]accessor.GribAccessor
	section_offset     [MaxNumSections]uint32
	section_length     [MaxNumSections]uint32
	sections_count     int
	offset             uint64
	bufr_subset_number int64 /* bufr subset number */
	bufr_group_number  int64 /* used in bufr */
	/* grib_accessor* groups[MAX_NUM_GROUPS]; */
	missingValueLong    int64
	missingValueDouble  float64
	product_kind        ProductKind
	bufr_elements_table *datatypes.GribTrie
}

type GribMultiHandle struct {
	buffer bytes.Buffer
	offset uint32
	length uint32
}

func GribHandleCreate(gl *GribHandle, gc *GribContext, data *bytes.Buffer, buflen uint32) (*GribHandle, errs.Error) {
	var err errs.Error
	if gl == nil {
		// invalid use pattern
		return nil, nil
	}
	gl.use_trie = true
	gl.trie_invalid = false
	gl.buffer = accessor.NewGribBuffer(accessor.GribUserBuffer, data, buflen)
	// err = SetGribContextReader(gc)
	// if err != nil {
	// 	// this means shit went wrong with os?
	// 	return nil, nil
	// }
	gl.root, err = accessor.GribCreateRootSection(gl, gc)
	if err != nil {
		// unable to create handle
		return nil, nil
		/*if (!gl->context->grib_reader || !gl->context->grib_reader->first) {
		grib_context_log(c, GRIB_LOG_ERROR, "grib_handle_create: cannot create handle, no definitions found");
		grib_handle_delete(gl);
		return NULL;
		}*/
	}
	next = gc.grib_reader.first.root
	for next != nil {
		err = accessor.GribCreateAccessor(gl.root, next, nil)
		if err != nil {
			return nil, nil
		}
		next = next.next
	}
	err = accessor.GribSectionAdjustSize(gl.root, 0, 0)
	if err != nil {
		return nil, nil
	}
	err = GribSectionPostInit(gl.root)
	return gl
}

func GribHandleNewFromMessage(gc *GribContext, data *bytes.Buffer, buflen uint32) (*GribHandle, errs.Error) {
	var err error
	var gl GribHandle
	if gc == nil {
		*gc, err = GetContextDefault()
		if err != nil {
			return nil, errs.NewContextError(errs.WithMessage("Cannot determine default context"))
		}
	}
	gl.product_kind = ProductGrib
	h, herr := GribHandleCreate(&gl, gc, data, buflen) //h                = grib_handle_create(gl, c, data, buflen);
	if herr != nil {
		return nil, herr
	}
	product_kind, perr := DetermineProductKind(h) //(determine_product_kind(h, &product_kind) == GRIB_SUCCESS)
	if perr != nil {
		return nil, errs.NewInvalidProductKind(errs.WithMessage("Bad message"))
	}

	if h.product_kind == ProductGrib {
		if !GribIsDefined(h, "7777") {
			fmt.Errorf("unable to create handle from message: No final 7777 in message")
			return nil, errs.NewInavalidGribError(errs.WithMessage("No final 7777 in message"))
		}
	}
	return *h, nil

}
