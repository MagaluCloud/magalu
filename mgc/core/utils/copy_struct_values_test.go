package utils

import (
	"reflect"
	"testing"
)

type Dst struct {
	FieldString     string
	ExclusiveField  string
	unexportedField string
}

type Src struct {
	FieldString     string
	FiedlInt        int
	unexportedField string
}

func TestCopyStructValues(t *testing.T) {
	src := Src{
		FieldString:     "src FieldString",
		FiedlInt:        10,
		unexportedField: "src unexportedField",
	}

	dst := Dst{
		FieldString:     "dst FieldString",
		ExclusiveField:  "dst ExclusiveField",
		unexportedField: "dst unexportedField",
	}

	dstExpected := Dst{
		FieldString:     src.FieldString,
		ExclusiveField:  dst.ExclusiveField,
		unexportedField: dst.unexportedField,
	}

	if err := CopyStructValues(&src, &dst); err != nil {
		t.Error("error:", err)
	}

	if !reflect.DeepEqual(dst, dstExpected) {
		t.Errorf("dst struct\nexpected: %+v \nfounded:  %+v", dstExpected, dst)
	}
}
