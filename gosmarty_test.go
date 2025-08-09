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
		input string
		env   *object.Environment
		want  string
	}{
		{
			input: `Hello, {$name}!`,
			env: Must(object.NewEnvironment(
				object.WithVariable("name", "Smarty"),
			)),
			want: "Hello, Smarty!",
		},
		{
			input: `Hello, {$first_name} {$given_name}!`,
			env: Must(object.NewEnvironment(
				object.WithVariable("first_name", "Go"),
				object.WithVariable("given_name", "Smarty"),
			)),
			want: "Hello, Go Smarty!",
		},
		{
			input: `{$contents | nl2br}`,
			env: Must(object.NewEnvironment(
				object.WithVariable("contents", "Hello1\nHello2\nHello3"),
			)),
			want: "Hello1<br>Hello2<br>Hello3",
		},
		{
			input: `{$name | devtest1 | devtest1 | devtest1} 1|2|3|4`,
			env: Must(object.NewEnvironment(
				object.WithVariable("name", "Smarty"),
			)),
			want: "Smarty_test1_test1_test1 1|2|3|4",
		},
		// Numbers
		{
			input: `This is number test: {$num}.`,
			env: Must(object.NewEnvironment(
				object.WithVariable("num", 777),
			)),
			want: "This is number test: 777.",
		},
		{
			input: `This is number test: {$num}.`,
			env: Must(object.NewEnvironment(
				object.WithVariable("num", -777),
			)),
			want: "This is number test: -777.",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", i+1), func(t *testing.T) {
			gsm, err := New("test").Parse(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			evaled := gsm.Exec(tt.env)
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

func TestComment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		env   *object.Environment
		want  string
	}{
		{
			input: `Hello,{* Comment *} {$name}!`,
			env: Must(object.NewEnvironment(
				object.WithVariable("name", "Smarty"),
			)),
			want: "Hello, Smarty!",
		},
		{
			input: `Hello,{* Comment *}{$name}!`,
			env: Must(object.NewEnvironment(
				object.WithVariable("name", "Smarty"),
			)),
			want: "Hello,Smarty!",
		},
		{
			input: `{*
Note:
  - One
  - Two
  - Three			
*}
<span>Hello</span>`,
			env: Must(object.NewEnvironment(
				object.WithVariable("name", "Smarty"),
			)),
			want: "\n<span>Hello</span>",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", i+1), func(t *testing.T) {
			gsm, err := New("test").Parse(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			evaled := gsm.Exec(tt.env)
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
		env      *object.Environment
		expected string
	}{
		{
			input: `{if $is_logged_in}Welcome, {$name}!{else}Hello, Guest.{/if}`,
			env: Must(object.NewEnvironment(
				object.WithVariable("is_logged_in", true),
				object.WithVariable("name", "Suzuki"),
			)),
			expected: "Welcome, Suzuki!",
		},
		{
			input: `{if $is_logged_in}Welcome, {$name}!{else}Hello, Guest.{/if}`,
			env: Must(object.NewEnvironment(
				object.WithVariable("is_logged_in", false),
				object.WithVariable("name", "Suzuki"),
			)),
			expected: "Hello, Guest.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(object.NewEnvironment(
				object.WithVariable("item_count", "exists"), // 空文字以外はtrue
			)),
			expected: "Your item is available.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(object.NewEnvironment(
				object.WithVariable("item_count", ""), // 空文字はfalse
			)),
			expected: "Your item is sold out.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(object.NewEnvironment(
				object.WithVariable("item_count", nil), // nilはfalse
			)),
			expected: "Your item is sold out.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(object.NewEnvironment(
				object.WithVariable("item_count", 0), // 0はfalse
			)),
			expected: "Your item is sold out.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(object.NewEnvironment(
				object.WithVariable("item_count", 1), // 0以外true
			)),
			expected: "Your item is available.",
		},
		{
			input: `{if $show_block}This block is shown.{/if}`,
			env: Must(object.NewEnvironment(
				object.WithVariable("show_block", true),
			)),
			expected: "This block is shown.",
		},
		{
			input: `{if $show_block}This block is shown.{/if}`,
			env: Must(object.NewEnvironment(
				object.WithVariable("show_block", false),
			)),
			expected: "",
		},
	}

	for _, tt := range tests {
		gsm, err := New("").Parse(tt.input)
		if err != nil {
			t.Error(err)
			continue
		}

		evaled := gsm.Exec(tt.env)
		result, ok := evaled.(*object.String)
		if !ok {
			t.Fatalf("evaluation result is not String. got=%T (%+v)", evaled, evaled)
		}

		if result.Value != tt.expected {
			t.Errorf("wrong result for input %q.\nexpected=%q\ngot     =%q", tt.input, tt.expected, result.Value)
		}
	}
}
