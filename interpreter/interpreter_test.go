package interpreter

import (
	"testing"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		name     string
		tokens   []string
		expected string
	}{
		{
			name: "set and math",
			tokens: []string{
				"do=set", "arg=x", "arg=10",
				"do=add", "var=x", "arg=5",
			},
			expected: "15",
		},
		{
			name: "eval substitution",
			tokens: []string{
				"do=set", "arg=y", "eval=begin", "do=sub", "arg=10", "arg=3", "eval=end",
				"do=mul", "var=y", "arg=2",
			},
			expected: "14",
		},
		{
			name: "if true",
			tokens: []string{
				"do=set", "arg=x", "arg=1",
				"do=if", "block=begin", "do=eq", "var=x", "arg=1", "block=end",
				"block=begin", "do=set", "arg=result", "arg=yes", "block=end",
				"do=set", "arg=dummy", "var=result",
			},
			expected: "yes",
		},
		{
			name: "ifelse false",
			tokens: []string{
				"do=ifelse", 
				"block=begin", "do=eq", "arg=1", "arg=2", "block=end",
				"block=begin", "do=set", "arg=res", "arg=true", "block=end",
				"block=begin", "do=set", "arg=res", "arg=false", "block=end",
				"do=set", "arg=dummy", "var=res",
			},
			expected: "false",
		},
		{
			name: "while loop",
			tokens: []string{
				"do=set", "arg=i", "arg=0",
				"do=while",
				"block=begin", "do=lt", "var=i", "arg=3", "block=end",
				"block=begin", "do=set", "arg=i", "eval=begin", "do=add", "var=i", "arg=1", "eval=end", "block=end",
				"do=set", "arg=dummy", "var=i",
			},
			expected: "3",
		},
		{
			name: "def custom function",
			tokens: []string{
				"do=def", "arg=double", "block=begin", "arg=x", "block=end",
				"block=begin", "do=mul", "var=x", "arg=2", "block=end",
				"do=double", "arg=5",
			},
			expected: "10",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			i := NewInterpreter()
			RegisterBuiltins(i)
			i.Tokens = tc.tokens

			res, err := i.Eval()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, res)
			}
		})
	}
}
