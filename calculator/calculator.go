package calculator

import (
	"fmt"
	"strconv"
	"strings"

	"errors"
	"math"

	hbot "github.com/whyrusleeping/hellabot"
)

var commands map[string]struct{} = map[string]struct{}{
	"!calc": struct{}{},
}

func CalcTrigger() *hbot.Trigger {
	return &hbot.Trigger{
		Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
			if m.Command != "PRIVMSG" || len(m.Trailing()) == 0 {
				return false
			}

			maybeCommand := strings.Split(strings.TrimSpace(m.Trailing()), " ")[0]
			_, ok := commands[maybeCommand]
			return ok
		},
		Action: func(bot *hbot.Bot, m *hbot.Message) bool {
			input, _ := strings.CutPrefix(m.Trailing(), "!calc")
			lexemes, err := LexTrailing(input)
			if err != nil {
				bot.Reply(m, err.Error())
				return false
			}
			ast, err := Parse(lexemes)
			if err != nil {
				bot.Reply(m, err.Error())
				return false
			}
			eval, err := ast.Eval()
			if err != nil {
				bot.Reply(m, err.Error())
				return false
			}
			bot.Reply(m, fmt.Sprintf("%s => %f", input, eval))
			return false
		},
	}
}

type Expr interface {
	Eval() (float64, error)
}

type CalcValue struct {
	value float64
}

func (cv CalcValue) Eval() (float64, error) {
	return cv.value, nil
}

type CalcExpr struct {
	op   string
	args []Expr
}

func (ce CalcExpr) Eval() (float64, error) {
	if len(ce.args) == 0 {
		return -3.14159, errors.New(ce.op + " called without arguments!")
	}
	result, err := (ce.args[0]).Eval()
	if err != nil {
		return -3.14159, err
	}
	for i := 1; i < len(ce.args); i++ {
		next, err := ce.args[i].Eval()
		if err != nil {
			return -3.14159, err
		}
		switch ce.op{
			case "+":
			result += next
			case "-":
			result -= next
			case "*":
			result *= next
			case "/":
			if next == 0.0 {
				return -3.14159, errors.New("Error, divide by ZERO!")
			}
			result /= next
			default:
			return -3.14159, errors.New(ce.op + ", unknown operation!")
		}
	}
	return result, nil
}

func Parse(lexemes []string) (Expr, error) {
	if len(lexemes) == 0 {
		return nil, errors.New("No expression provides!")
	}

	if len(lexemes) == 1 && isNumber(lexemes[0]){
		val, _ := strconv.ParseFloat(lexemes[0], 64)
		return CalcValue{value: val}, nil
	}

	return parseExpr(lexemes, 0)
}

func parseExpr(lexemes []string, index int) (Expr, error) {
	if lexemes[index] != "(" {
		return nil, errors.New("expected ( at starting expression")
	}
	if index + 1 == len(lexemes) - 1 {
		return CalcValue{value: math.NaN()}, nil
	}
	ce := CalcExpr{op: lexemes[index + 1]}
	for i := index + 2; i < len(lexemes); i++ {
		if lexemes[i] == "(" {
			subExpr, err := parseExpr(lexemes, i)
			if err != nil {
				return nil, err
			}
			ce.args = append(ce.args, subExpr)
			// a pretty bad O(n^2) right here
			for lexemes[i] != ")" {
				i++
			}
		} else if lexemes[i] == ")" {
			break
		} else if isNumber(lexemes[i]) {
			val, _ := strconv.ParseFloat(lexemes[i], 64)
			ce.args = append(ce.args, CalcValue{value: val})
		} else {
			return nil, errors.New(lexemes[i] + " lexeme is not a valid number")
		}
	}
	return ce, nil
}

func isNumber(lexeme string) bool {
	count := 0
	for i := 0; i < len(lexeme); i++ {
		if lexeme[i] == '.' {
			count++
			continue
		}
		if lexeme[i] < '0' || lexeme[i] > '9' {
			return false
		}
	}
	return count <= 1
}

func LexTrailing(input string) ([]string, error) {
	var lexemes []string = []string{}
	size := len(input)
	stack := 0
	for i := 0; i < size; i++ {
		if input[i] == '(' {
			stack++
			lexemes = append(lexemes, "(")
		} else if input[i] == ')' {
			lexemes = append(lexemes, ")")
			if stack == 0 {
				return nil, errors.New("unbalanced parentheses )")
			}
			stack--
		} else if input[i] != ' ' {
			var j int
			for j = i + 1; j < size; j++ {
				if input[j] == ' ' || input[j] == ')' || input[j] == '(' {
					break
				}
			}
			if i < j {
				lexemes = append(lexemes, input[i:j])
				if j < size && input[j] == ')' {
					i = j - 1
				} else {
					i = j
				}
			}
		}
	}
	if stack > 0 {
		return nil, errors.New("unbalanced parentheses (")
	}
	return lexemes, nil
}
