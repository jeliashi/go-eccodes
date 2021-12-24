package datatypes

import (
	errs "gogrib/errors"
	"unicode"
)

const SIZE int = 39

func getMapping(c rune) int {
	if unicode.IsDigit(c) {
		return int(c-'0') + 1
	}
	if unicode.IsLetter(c) {
		if unicode.IsLower(c) {
			return int(c - 'a' + 10)
		}
		return int(c - 'A' + 10)
	}
	if c == '#' {
		return 38
	}
	if c == '_' {
		return 37
	}

	return 0

}

type GribTrie struct {
	next  [SIZE]*GribTrie
	first int
	last  int
	data  interface{}
}

func (t *GribTrie) GribTrieGet(key string) interface{} {
	// I honestly have no idea how this function works
	// or to be honest, how a trie data structure works
	// please see:
	// https://github.com/ecmwf/eccodes/blob/5407ecd6d19afad0db729a73043b2a31843140c9/src/grib_trie.c#L474

	for _, _k := range key {
		t = t.next[getMapping(_k)]
	}
	if t != nil && t.data != nil {
		return t.data
	}
	return nil
}

func (t *GribTrie) Insert(s string, data interface{}) (interface{}, errs.Error) {

	var last *GribTrie
	var j int = 0
	var k rune
	if t == nil {
		return nil, errs.NewContextError(errs.WithMessage("nil GribTrie"))
	}
	for _j, _k := range s {
		last = t
		t = t.next[getMapping(_k)]
		if t == nil {
			break
		}
		j = _j
		k = _k
	}
	if k != '0' {
		t = last
		for _, _k := range s[j:] {
			i := getMapping(_k)
			if i < t.first {
				t.first = i
			}
			if i > t.last {
				t.last = i
			}
			t.next[i] = &GribTrie{}
			t = t.next[i]
		}

	}
	old := t.data
	t.data = data
	// equiv of return data == old ? NULL : old; I think?
	if old == data {
		return nil, nil
	} else {
		return data, nil
	}
}

//     if (*k == 0) {
//         old     = t->data;
//         t->data = data;
//     }
//     else {
//         t = last;
//         while (*k) {
//             int j = 0;
//             DebugCheckBounds((int)*k, key);
//             j = mapping[(int)*k++];
//             if (j < t->first)
//                 t->first = j;
//             if (j > t->last)
//                 t->last = j;
//             t = t->next[j] = grib_trie_new(t->context);
//         }
//         old     = t->data;
//         t->data = data;
//     }
//     GRIB_MUTEX_UNLOCK(&mutex);
//     return data == old ? NULL : old;
// }
