package gosmarty

import (
	"fmt"
	"testing"

	"github.com/szks-repo/gosmarty/object"
)

// 変数置換の基本テスト
func TestVariableEvaluation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input      string
		envFactory func() *object.Environment
		want       string
	}{
		{
			input: `Hello, {$name}!`,
			envFactory: func() *object.Environment {
				env := object.NewEnvironment()
				env.Set("name", &object.String{Value: "Smarty"})
				return env
			},
			want: "Hello, Smarty!",
		},
		{
			input: `Hello, {$first_name} {$given_name}!`,
			envFactory: func() *object.Environment {
				env := object.NewEnvironment()
				env.Set("first_name", &object.String{Value: "Go"})
				env.Set("given_name", &object.String{Value: "Smarty"})
				return env
			},
			want: "Hello, Go Smarty!",
		},
		{
			input: `{$contents | nl2br}`,
			envFactory: func() *object.Environment {
				env := object.NewEnvironment()
				env.Set("contents", &object.String{Value: "Hello1\nHello2\nHello3"})
				return env
			},
			want: "Hello1<br>Hello2<br>Hello3",
		},
		{
			input: `{$name | devtest1 | devtest1 | devtest1} 1|2|3|4`,
			envFactory: func() *object.Environment {
				env := object.NewEnvironment()
				env.Set("name", &object.String{Value: "Smarty"})
				return env
			},
			want: "Smarty_test1_test1_test1 1|2|3|4",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", i+1), func(t *testing.T) {
			gsm, err := New("test").Parse(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			evaled := gsm.Exec(tt.envFactory())
			result, ok := evaled.(*object.String)
			if !ok {
				t.Error("isn't object.String")
			}

			if result.Value != tt.want {
				t.Errorf("result has wrong value. got=%q, want=%q", result.Value, tt.want)
			}
		})
	}
}

// if文のテスト
func TestIfStatements(t *testing.T) {
	t.Parallel()

	// テストケースを定義
	tests := []struct {
		input    string
		envSetup map[string]object.Object
		expected string
	}{
		{
			input: `{if $is_logged_in}Welcome, {$name}!{else}Hello, Guest.{/if}`,
			envSetup: map[string]object.Object{
				"is_logged_in": object.TRUE,
				"name":         &object.String{Value: "Suzuki"},
			},
			expected: "Welcome, Suzuki!",
		},
		{
			input: `{if $is_logged_in}Welcome, {$name}!{else}Hello, Guest.{/if}`,
			envSetup: map[string]object.Object{
				"is_logged_in": object.FALSE,
				"name":         &object.String{Value: "Suzuki"},
			},
			expected: "Hello, Guest.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			envSetup: map[string]object.Object{
				"item_count": &object.String{Value: "exists"}, // 空文字以外はtrue
			},
			expected: "Your item is available.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			envSetup: map[string]object.Object{
				"item_count": &object.String{Value: ""}, // 空文字はfalse
			},
			expected: "Your item is sold out.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			envSetup: map[string]object.Object{
				"item_count": &object.Null{}, // nullはfalse
			},
			expected: "Your item is sold out.",
		},
		{
			input: `{if $show_block}This block is shown.{/if}`,
			envSetup: map[string]object.Object{
				"show_block": object.TRUE,
			},
			expected: "This block is shown.",
		},
		{
			input: `{if $show_block}This block is shown.{/if}`,
			envSetup: map[string]object.Object{
				"show_block": object.FALSE,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		gsm, err := New("").Parse(tt.input)
		if err != nil {
			t.Error(err)
			continue
		}

		env := object.NewEnvironment()
		for key, val := range tt.envSetup {
			env.Set(key, val)
		}
		evaled := gsm.Exec(env)

		// 結果を検証
		result, ok := evaled.(*object.String)
		if !ok {
			t.Fatalf("evaluation result is not String. got=%T (%+v)", evaled, evaled)
		}

		if result.Value != tt.expected {
			t.Errorf("wrong result for input %q.\nexpected=%q\ngot     =%q", tt.input, tt.expected, result.Value)
		}
	}
}
