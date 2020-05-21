package lexer

import "testing"

func TestParsingNumber(t *testing.T) {
	var caseMap = map[string]TokenType{
		".11": ERR,
		//"11.":   ERR,
		"11a":   ERR,
		"11.a0": ERR,
		"11a0":  ERR,
		//"-11":   NUMBER,
		"a11": IDENTIFIER,
		"1.0": NUMBER,
	}
	evalCase(t, caseMap)
}

func TestParsingString(t *testing.T) {
	var caseMap = map[string]TokenType{
		"\"dame\"": STRING,
		"\"dame":   ERR,
		//"'dame":    ERR,
	}
	evalCase(t, caseMap)
}

func TestParsingExpression(t *testing.T) {
	caseMap := map[string][]TokenType{
		"(+ 1 2)":                  []TokenType{OP, PLUS, NUMBER, NUMBER, CP},
		"(+ 1.2 3)":                []TokenType{OP, PLUS, NUMBER, NUMBER, CP},
		"(+ (- 3 4) 3)":            []TokenType{OP, PLUS, OP, MINUS, NUMBER, NUMBER, CP, NUMBER, CP},
		"(_ 1.2 3)":                []TokenType{OP, PLACEHOLDER, NUMBER, NUMBER, CP},
		"(_ a b)":                  []TokenType{OP, PLACEHOLDER, IDENTIFIER, IDENTIFIER, CP},
		"(declare a 10)":           []TokenType{OP, DECLARE, IDENTIFIER, NUMBER, CP},
		"(== a 10)":                []TokenType{OP, EQUAL_EQUAL, IDENTIFIER, NUMBER, CP},
		"(>= a 10)":                []TokenType{OP, GREATER_EQUAL, IDENTIFIER, NUMBER, CP},
		"(<= a 10)":                []TokenType{OP, LESS_EQUAL, IDENTIFIER, NUMBER, CP},
		"if":                       []TokenType{IF},
		"(== false true)":          []TokenType{OP, EQUAL_EQUAL, BOOL_FALSE, BOOL_TRUE, CP},
		"(a 10) # (if equal true)": []TokenType{OP, IDENTIFIER, NUMBER, CP, COMMENT},
	}
	evalExpr(t, caseMap)
}

func evalExpr(t *testing.T, caseMap map[string][]TokenType) {
	for inputExp, expectExp := range caseMap {
		grinder := startGrinding(inputExp)
		for _, expectType := range expectExp {
			tokenType := grinder.nextToken().tokenType
			if expectType != tokenType {
				t.Errorf("tokenType %s is %s", inputExp, reverseKeys[tokenType])
			}
		}
	}
}

func evalCase(t *testing.T, caseMap map[string]TokenType) {
	for inputExp, expectExp := range caseMap {
		tokenType := startGrinding(inputExp).nextToken().tokenType
		if tokenType != expectExp {
			t.Errorf("tokenType %s is %s", inputExp, reverseKeys[tokenType])
		}
	}
}
