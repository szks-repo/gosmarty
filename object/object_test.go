package object

import (
	"reflect"
	"testing"
	"time"
)

func TestNewObjectFromAny(t *testing.T) {
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
		{
			anyVal: []string{"1", "2", "3"},
			want: &Array{
				Value: []Object{
					NewString("1"),
					NewString("2"),
					NewString("3"),
				},
			},
		},
		{
			anyVal: struct {
				Id   int
				Name string
			}{
				Id:   123,
				Name: "Alice",
			},
			want: &Map{
				Value: map[string]Object{
					"Id":   &Number{Value: 123},
					"Name": &String{Value: "Alice"},
				},
			},
		},
		{
			anyVal: struct {
				Id       int
				Name     string
				Metadata *struct {
					CreatedAt time.Time
				}
			}{
				Id:   123,
				Name: "Alice",
				Metadata: &struct {
					CreatedAt time.Time
				}{
					CreatedAt: time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC),
				},
			},
			want: &Map{
				Value: map[string]Object{
					"Id":   &Number{Value: 123},
					"Name": &String{Value: "Alice"},
					"Metadata": &Map{
						Value: map[string]Object{
							"CreatedAt": &Time{Value: time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := NewObjectFromAny(tt.anyVal)
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
