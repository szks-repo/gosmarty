package object

import (
	"reflect"
	"testing"
)

func Test_newObjectFromAny(t *testing.T) {
	t.Parallel()

	type (
		StringUnderlying string
		BoolUnderlying   bool
		IntUnderlying    int
		Int64Underlying  int64
	)
	tests := []struct {
		anyVal  any
		want    Object
		wantErr bool
	}{
		{
			anyVal: "string",
			want:   &String{Value: "string"},
		},
		{
			anyVal: StringUnderlying("userId"),
			want:   &String{Value: "userId"},
		},
		{
			anyVal: true,
			want:   &Boolean{Value: true},
		},
		{
			anyVal: false,
			want:   &Boolean{Value: false},
		},
		{
			anyVal: BoolUnderlying(true),
			want:   &Boolean{Value: true},
		},
		{
			anyVal: nil,
			want:   &Null{},
		},
		{
			anyVal: int(100),
			want:   &Number{Value: float64(100)},
		},
		{
			anyVal: IntUnderlying(100),
			want:   &Number{Value: float64(100)},
		},
		{
			anyVal: Int64Underlying(100),
			want:   &Number{Value: float64(100)},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := newObjectFromAny(tt.anyVal)
			if err != nil && !tt.wantErr {
				t.Error(err)
			} else {
				if tt.wantErr {
					t.Error("wantErr=true, but err is nil")
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("want=%#v, bug got=%#v", tt.want, got)
				}
			}
		})
	}
}
