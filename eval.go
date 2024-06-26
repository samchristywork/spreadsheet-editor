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
			cells = append(cells, fmt.Sprintf("%s%d", getColumnName(col1), row))
		}
	} else if row1 == row2 {
		for col := col1; col <= col2; col++ {
			cells = append(cells, fmt.Sprintf("%s%d", getColumnName(col), row1))
		}
	} else {
		return nil, fmt.Errorf("Error creating range: Only supports ranges in a single row or column")
	}

	return cells, nil
}

func strlen(args ...interface{}) (interface{}, error) {
	length := len(args[0].(string))
	return (float64)(length), nil
}

func sum(args ...interface{}) (interface{}, error) {
	if len(args) == 2 {
		return args[0].(float64) + args[1].(float64), nil
	}

	if len(args) == 1 {
		s := strings.Split(args[0].(string), ":")
		sum := 0.0

		cells, err := getCellRange(column(s[0]), row(s[0]), column(s[1]), row(s[1]))
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}

		for _, cell := range cells {
			value, err := getCellValue(row(cell), column(cell))
			if err != nil {
				return nil, fmt.Errorf("Error finding cell: " + cell)
			}

			value = strings.TrimSpace(value)

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
}

func collectParameters(evaluableExpression *govaluate.EvaluableExpression) (map[string]interface{}, error) {
	parameters := make(map[string]interface{}, 8)
	vars := evaluableExpression.Vars()

	for _, v := range vars {
		if !isCellIdentifier(v) {
			return nil, fmt.Errorf("Error applying function: %s is not a valid cell identifier", v)
		}

		val, err := getCellValue(row(v), column(v))
		if err != nil {
			return nil, fmt.Errorf("Error finding cell: " + v)
		}

		val = strings.TrimSpace(val)

		if val == "" {
			return nil, fmt.Errorf("Error applying function: Cell %s is empty", v)
		} else if isNumeric(val) {
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, fmt.Errorf("%v", err)
			}
			parameters[v] = f
		} else {
			parameters[v] = val
		}
	}

	return parameters, nil
}

func isFloatingPointNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func eval(expression string) string {
	functions := map[string]govaluate.ExpressionFunction{
		"strlen": strlen,
		"sum":    sum,
	}

	evaluableExpression, err := govaluate.NewEvaluableExpressionWithFunctions(expression, functions)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	parameters, err := collectParameters(evaluableExpression)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	result, err := evaluableExpression.Evaluate(parameters)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}

	if isFloatingPointNumber(fmt.Sprintf("%v", result)) {
		s := fmt.Sprintf("%.4f", result)
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
		return s
	} else {
		return fmt.Sprintf("%v", result)
	}
}
