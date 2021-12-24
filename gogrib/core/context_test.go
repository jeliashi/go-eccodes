package core

import (
	"reflect"
	"testing"
)

func TestDefault(t *testing.T) {
	gc, err := GetContextDefault()
	if err != nil {
		t.Errorf("not nil err")
	}
	if gc.debug {
		t.Errorf("bad")
	}
}

func TestGetContextDefault(t *testing.T) {
	tests := []struct {
		name    string
		want    GribContext
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContextDefault()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContextDefault() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetContextDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
