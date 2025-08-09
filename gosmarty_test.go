package gosmarty

import (
	"testing"

	"github.com/szks-repo/gosmarty/object"
)

// 変数置換の基本テスト
func TestVariableEvaluation(t *testing.T) {
	input := `Hello, {$name}!`
	expected := "Hello, Smarty!"

	gsm, err := New("test").Parse(input)
	if err != nil {
		t.Fatal(err)
	}

	env := object.NewEnvironment()
	env.Set("name", &object.String{Value: "Smarty"})
	evaled := gsm.Exec(env)

	// 3. 結果を検証
	result, ok := evaled.(*object.String)
	if !ok {
		t.Error("isn't object.String")
	}

	if result.Value != expected {
		t.Errorf("result has wrong value. got=%q, want=%q", result.Value, expected)
	}
}

// if文のテスト
func TestIfStatements(t *testing.T) {
	// テストケースを定義
	tests := []struct {
		input    string
		envSetup map[string]object.Object
		expected string
	}{
		{
			input: `{if $is_logged_in}Welcome, {$name}!{else}Hello, Guest.{/if}`,
			envSetup: map[string]object.Object{
				"is_logged_in": object.TRUE, // 真
				"name":         &object.String{Value: "Suzuki"},
			},
			expected: "Welcome, Suzuki!",
		},
		{
			input: `{if $is_logged_in}Welcome, {$name}!{else}Hello, Guest.{/if}`,
			envSetup: map[string]object.Object{
				"is_logged_in": object.FALSE, // 偽
				"name":         &object.String{Value: "Suzuki"},
			},
			expected: "Hello, Guest.",
		},
		{
			input:    `Your item is {if $item_count}available{else}sold out{/if}.`,
			envSetup: map[string]object.Object{"item_count": &object.String{Value: "exists"}}, // 空文字以外は真
			expected: "Your item is available.",
		},
		{
			input:    `Your item is {if $item_count}available{else}sold out{/if}.`,
			envSetup: map[string]object.Object{"item_count": &object.String{Value: ""}}, // 空文字は偽
			expected: "Your item is sold out.",
		},
		{
			input:    `{if $show_block}This block is shown.{/if}`,
			envSetup: map[string]object.Object{"show_block": object.TRUE},
			expected: "This block is shown.",
		},
		{
			input:    `{if $show_block}This block is shown.{/if}`, // else節がない場合
			envSetup: map[string]object.Object{"show_block": object.FALSE},
			expected: "", // 何も出力されない
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
