package gosmarty

import (
	"fmt"
	"testing"
	"time"

	"github.com/szks-repo/gosmarty/modifier"
	"github.com/szks-repo/gosmarty/object"
)

// 変数置換の基本テスト
func TestVariableEvaluation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		env   *Environment
		want  string
	}{
		{
			input: `Hello, {$name}!`,
			env: Must(NewEnvironment(
				WithVariable("name", "Smarty"),
			)),
			want: "Hello, Smarty!",
		},
		{
			input: `Hello, {$first_name} {$given_name}!`,
			env: Must(NewEnvironment(
				WithVariable("first_name", "Go"),
				WithVariable("given_name", "Smarty"),
			)),
			want: "Hello, Go Smarty!",
		},
		// Numbers
		{
			input: `This is number test: {$num}.`,
			env: Must(NewEnvironment(
				WithVariable("num", 777),
			)),
			want: "This is number test: 777.",
		},
		{
			input: `This is number test: {$num}.`,
			env: Must(NewEnvironment(
				WithVariable("num", -777),
			)),
			want: "This is number test: -777.",
		},
		// Field access
		{
			input: `id:{$user.id} name:{$user.name}`,
			env: Must(NewEnvironment(
				WithVariable("user", map[string]any{
					"id":   "1",
					"name": "Tanaka",
				}),
			)),
			want: "id:1 name:Tanaka",
		},
		{
			input: `id:{$user.id} name:{$user.name} rank:{$user.rank} balance:{$user.balance | number_format} created_at:{$user.metadata.created_at}`,
			env: Must(NewEnvironment(
				WithVariable("user", map[string]any{
					"id":      "1",
					"name":    "Tanaka",
					"rank":    3,
					"balance": 1500000,
					"metadata": map[string]any{
						"created_at": time.Date(2024, 1, 3, 15, 4, 6, 0, time.UTC),
					},
				}),
			)),
			want: "id:1 name:Tanaka rank:3 balance:1,500,000 created_at:2024-01-03T15:04:06Z",
		},
		// Index access
		{
			input: `id[0]:{$ids[0]} id[1]:{$ids[1]} id[2]:{$ids[2]} id[3]:{$ids[3]}`,
			env: Must(NewEnvironment(
				WithVariable("ids", []string{"1", "2", "3", "4"}),
			)),
			want: "id[0]:1 id[1]:2 id[2]:3 id[3]:4",
		},
		{
			input: `id:{$user.id} name:{$user.name} comment:{$user.comments[1] | upper}`,
			env: Must(NewEnvironment(
				WithVariable("user", map[string]any{
					"id":       "1",
					"name":     "Tanaka",
					"comments": []string{"aaa", "bbb", "ccc"},
				}),
			)),
			want: "id:1 name:Tanaka comment:BBB",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", i+1), func(t *testing.T) {
			gsm := New()
			tmpl, err := gsm.Parse(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			evaled := tmpl.Execute(tt.env)
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

func TestModifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input     string
		env       *Environment
		want      string
		modifiers map[string]modifier.Modifier
	}{
		{
			input: `{$contents | nl2br}`,
			env: Must(NewEnvironment(
				WithVariable("contents", "Hello1\nHello2\nHello3"),
			)),
			want: "Hello1<br />Hello2<br />Hello3",
		},
		{
			input: `{$name | pipe_test | pipe_test | pipe_test} 1|2|3|4`,
			env: Must(NewEnvironment(
				WithVariable("name", "Smarty"),
			)),
			modifiers: map[string]modifier.Modifier{
				"pipe_test": func(input object.Object, args ...any) object.Object {
					if input.Type() == object.StringType {
						return object.NewString(input.Inspect() + "_test1")
					}
					return object.NewNull()
				},
			},
			want: "Smarty_test1_test1_test1 1|2|3|4",
		},
		{
			input: `This is number test: {$num | number_format}.`,
			env: Must(NewEnvironment(
				WithVariable("num", 777777777),
			)),
			want: "This is number test: 777,777,777.",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", i+1), func(t *testing.T) {
			gsm := New()
			for modName, mod := range tt.modifiers {
				RegisterModifier(modName, mod)
			}

			tmpl, err := gsm.Parse(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			evaled := tmpl.Execute(tt.env)
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
		env   *Environment
		want  string
	}{
		{
			input: `Hello,{* Comment *} {$name}!`,
			env: Must(NewEnvironment(
				WithVariable("name", "Smarty"),
			)),
			want: "Hello, Smarty!",
		},
		{
			input: `Hello,{* Comment *}{$name}!`,
			env: Must(NewEnvironment(
				WithVariable("name", "Smarty"),
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
			env: Must(NewEnvironment(
				WithVariable("name", "Smarty"),
			)),
			want: "\n<span>Hello</span>",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", i+1), func(t *testing.T) {
			gsm := New()
			tmpl, err := gsm.Parse(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			evaled := tmpl.Execute(tt.env)
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
		env      *Environment
		expected string
	}{
		{
			input: `{if $is_logged_in}Welcome, {$name}!{else}Hello, Guest.{/if}`,
			env: Must(NewEnvironment(
				WithVariable("is_logged_in", true),
				WithVariable("name", "Suzuki"),
			)),
			expected: "Welcome, Suzuki!",
		},
		{
			input: `{if $is_logged_in}Welcome, {$name}!{else}Hello, Guest.{/if}`,
			env: Must(NewEnvironment(
				WithVariable("is_logged_in", false),
				WithVariable("name", "Suzuki"),
			)),
			expected: "Hello, Guest.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(NewEnvironment(
				WithVariable("item_count", "exists"), // 空文字以外はtrue
			)),
			expected: "Your item is available.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(NewEnvironment(
				WithVariable("item_count", ""), // 空文字はfalse
			)),
			expected: "Your item is sold out.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(NewEnvironment(
				WithVariable("item_count", nil), // nilはfalse
			)),
			expected: "Your item is sold out.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(NewEnvironment(
				WithVariable("item_count", 0), // 0はfalse
			)),
			expected: "Your item is sold out.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(NewEnvironment(
				WithVariable("item_count", 1), // 0以外true
			)),
			expected: "Your item is available.",
		},
		{
			input: `{if $show_block}This block is shown.{/if}`,
			env: Must(NewEnvironment(
				WithVariable("show_block", true),
			)),
			expected: "This block is shown.",
		},
		{
			input: `{if $show_block}This block is shown.{/if}`,
			env: Must(NewEnvironment(
				WithVariable("show_block", false),
			)),
			expected: "",
		},
	}

	for _, tt := range tests {
		gsm := New()
		tmpl, err := gsm.Parse(tt.input)
		if err != nil {
			t.Fatal(err)
		}

		evaled := tmpl.Execute(tt.env)
		result, ok := evaled.(*object.String)
		if !ok {
			t.Error("isn't object.String")
		}

		if result.Value != tt.expected {
			t.Errorf("wrong result for input %q.\nexpected=%q\ngot     =%q", tt.input, tt.expected, result.Value)
		}
	}
}
