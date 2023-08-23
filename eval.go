package main

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"strconv"
	"strings"
)

func getCellRange(col1 int, row1 int, col2 int, row2 int) ([]string, error) {
	cells := make([]string, 0, 8)

	if col1 == col2 {
		for row := row1; row <= row2; row++ {
			cells = append(cells, fmt.Sprintf("%v%v", getColumnName(col1), row))
		}
	} else if row1 == row2 {
		for col := col1; col <= col2; col++ {
			cells = append(cells, fmt.Sprintf("%v%v", getColumnName(col), row1))
		}
	} else {
		return nil, fmt.Errorf("Error creating range: Only supports ranges in a single row or column")
	}

	return cells, nil
}

func eval(expression string) string {
	functions := map[string]govaluate.ExpressionFunction{
		"strlen": func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))
			return (float64)(length), nil
		},
		"sum": func(args ...interface{}) (interface{}, error) {
			if len(args) == 2 {
				return args[0].(float64) + args[1].(float64), nil
			}

			if len(args) == 1 {
				s := strings.Split(args[0].(string), ":")
				sum := 0.0

				cells := getCellRange(row(s[0]), column(s[0]), row(s[1]), column(s[1]))

				for _, cell := range cells {
					value, err := getCellValue(row(cell), column(cell))
					if err != nil {
						return nil, fmt.Errorf("Error finding cell: " + cell)
					}

					if isNumeric(value) {
						f, err := strconv.ParseFloat(value, 64)
						if err != nil {
							return nil, fmt.Errorf("%v", err)
						}
						sum += f
					} else {
						return nil, fmt.Errorf("Error applying sum: Can only be used with numeric values")
					}
				}

				return sum, nil
			}

			return nil, fmt.Errorf("Error applying sum: Function takes 1 or 2 arguments")
		},
	}

	evaluableExpression, err := govaluate.NewEvaluableExpressionWithFunctions(expression, functions)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	parameters := make(map[string]interface{}, 8)
	vars := evaluableExpression.Vars()
	for _, v := range vars {
		val, err := getCellValue(row(v), column(v))
		if err != nil {
			return fmt.Sprintf("%v", err)
		}

		if isNumeric(val) {
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return fmt.Sprintf("%v", err)
			}
			parameters[v] = f
		} else {
			parameters[v] = val
		}
	}

	result, err := evaluableExpression.Evaluate(parameters)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	return fmt.Sprintf("%v", result)
}
