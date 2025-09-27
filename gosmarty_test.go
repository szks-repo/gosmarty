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
		{
			input: `Hello, {$name}! {$bio}{$bio2}`,
			env: Must(NewEnvironment(
				WithVariable("name", "Smarty"),
				WithVariable("bio", Ptr("I am Gopher!")),
				WithVariable("bio2", new(string)),
			)),
			want: "Hello, Smarty! I am Gopher!",
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

func TestForeach(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		env   *Environment
		want  string
	}{
		{
			name:  "array iteration",
			input: `{foreach from=$items item=item}{$item}{/foreach}`,
			env: Must(NewEnvironment(
				WithVariable("items", []string{"apple", "banana", "cherry"}),
			)),
			want: "applebananacherry",
		},
		{
			name:  "array with key",
			input: `{foreach from=$items item=item key=index}{$index}:{$item};{/foreach}`,
			env: Must(NewEnvironment(
				WithVariable("items", []string{"apple", "banana"}),
			)),
			want: "0:apple;1:banana;",
		},
		{
			name:  "map iteration sorted by key",
			input: `{foreach from=$pairs item=value key=key}{$key}={$value};{/foreach}`,
			env: Must(NewEnvironment(
				WithVariable("pairs", map[string]any{
					"b": "Beta",
					"a": "Alpha",
				}),
			)),
			want: "a=Alpha;b=Beta;",
		},
		{
			name:  "map array iteration",
			input: `{foreach from=$users item=value}{$value.userId}:{$value.name},{/foreach}`,
			env: Must(NewEnvironment(
				WithVariable("users", []map[string]any{
					{
						"userId": 1,
						"name":   "Tom",
					},
					{
						"userId": 2,
						"name":   "Bob",
					},
				}),
			)),
			want: "1:Tom,2:Bob,",
		},
		{
			name:  "first last flags",
			input: `{foreach from=$items item=item name=list}{$smarty.foreach.list.first}:{$smarty.foreach.list.last}:{$item};{/foreach}`,
			env: Must(NewEnvironment(
				WithVariable("items", []string{"A", "B", "C"}),
			)),
			want: "true:false:A;false:false:B;false:true:C;",
		},
		{
			name:  "first flag in if",
			input: `{foreach from=$items item=item name=list}{if $smarty.foreach.list.first}First!{/if}{$item};{/foreach}`,
			env: Must(NewEnvironment(
				WithVariable("items", []string{"A", "B"}),
			)),
			want: "First!A;B;",
		},
		{
			name:  "last flag in if",
			input: `{foreach from=$items item=item name=list}{if $smarty.foreach.list.last}Last!{/if}{$item};{/foreach}`,
			env: Must(NewEnvironment(
				WithVariable("items", []string{"A", "B"}),
			)),
			want: "A;Last!B;",
		},
		{
			name:  "foreachelse fallback",
			input: `{foreach from=$items item=item}{$item}{foreachelse}empty{/foreach}`,
			env: Must(NewEnvironment(
				WithVariable("items", []string{}),
			)),
			want: "empty",
		},
		{
			name:  "foreachelse fallback",
			input: `{foreach from=$items item=item}{$item}{foreachelse}empty{/foreach}`,
			env: Must(NewEnvironment(
				WithVariable("items", 123),
			)),
			want: "empty",
		},
		{
			name:  "restore outer variable",
			input: `start:{$item} {foreach from=$items item=item}{$item}{/foreach} end:{$item}`,
			env: Must(NewEnvironment(
				WithVariable("item", "outside"),
				WithVariable("items", []string{"inside"}),
			)),
			want: "start:outside inside end:outside",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gsm := New()
			tmpl, err := gsm.Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error: %v", err)
			}

			evaled := tmpl.Execute(tt.env)
			result, ok := evaled.(*object.String)
			if !ok {
				t.Fatal("result is not *object.String")
			}

			if result.Value != tt.want {
				t.Errorf("unexpected output. got=%q, want=%q", result.Value, tt.want)
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
		input string
		env   *Environment
		want  string
	}{
		{
			input: `{if $is_logged_in}Welcome, {$name}!{else}Hello, Guest.{/if}`,
			env: Must(NewEnvironment(
				WithVariable("is_logged_in", true),
				WithVariable("name", "Suzuki"),
			)),
			want: "Welcome, Suzuki!",
		},
		{
			input: `{if $is_logged_in}Welcome, {$name}!{else}Hello, Guest.{/if}`,
			env: Must(NewEnvironment(
				WithVariable("is_logged_in", false),
				WithVariable("name", "Suzuki"),
			)),
			want: "Hello, Guest.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(NewEnvironment(
				WithVariable("item_count", "exists"), // 空文字以外はtrue
			)),
			want: "Your item is available.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(NewEnvironment(
				WithVariable("item_count", ""), // 空文字はfalse
			)),
			want: "Your item is sold out.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(NewEnvironment(
				WithVariable("item_count", nil), // nilはfalse
			)),
			want: "Your item is sold out.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(NewEnvironment(
				WithVariable("item_count", 0), // 0はfalse
			)),
			want: "Your item is sold out.",
		},
		{
			input: `Your item is {if $item_count}available{else}sold out{/if}.`,
			env: Must(NewEnvironment(
				WithVariable("item_count", 1), // 0以外true
			)),
			want: "Your item is available.",
		},
		{
			input: `{if $show_block}This block is shown.{/if}`,
			env: Must(NewEnvironment(
				WithVariable("show_block", true),
			)),
			want: "This block is shown.",
		},
		{
			input: `{if $show_block}This block is shown.{/if}`,
			env: Must(NewEnvironment(
				WithVariable("show_block", false),
			)),
			want: "",
		},
		{
			input: `{if $primary}Primary{elseif $secondary}Secondary{else}Fallback{/if}`,
			env: Must(NewEnvironment(
				WithVariable("primary", false),
				WithVariable("secondary", true),
			)),
			want: "Secondary",
		},
		{
			input: `{if $primary}Primary{elseif $secondary}Secondary{elseif $tertiary}Tertiary{else}Fallback{/if}`,
			env: Must(NewEnvironment(
				WithVariable("primary", false),
				WithVariable("secondary", false),
				WithVariable("tertiary", true),
			)),
			want: "Tertiary",
		},
		{
			input: `{if $primary}Primary{elseif $secondary}Secondary{else}Fallback{/if}`,
			env: Must(NewEnvironment(
				WithVariable("primary", false),
				WithVariable("secondary", false),
			)),
			want: "Fallback",
		},
		{
			input: `{if $num > 50 and $inSellingPeriod}在庫あり{else}在庫なし{/if}`,
			env: Must(NewEnvironment(
				WithVariable("num", 60),
				WithVariable("inSellingPeriod", true),
			)),
			want: "在庫あり",
		},
		{
			input: `{if $num > 50 and $inSellingPeriod}在庫あり{else}在庫なし{/if}`,
			env: Must(NewEnvironment(
				WithVariable("num", 60),
				WithVariable("inSellingPeriod", false),
			)),
			want: "在庫なし",
		},
		{
			input: `{if $num > 50 and $inSellingPeriod}在庫あり{else}在庫なし{/if}`,
			env: Must(NewEnvironment(
				WithVariable("num", 40),
				WithVariable("inSellingPeriod", true),
			)),
			want: "在庫なし",
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

		if result.Value != tt.want {
			t.Errorf("wrong result for input %q.\nwant=%q\ngot     =%q", tt.input, tt.want, result.Value)
		}
	}
}
