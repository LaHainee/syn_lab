package model

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"unicode"
)

type Phone struct {
	number int64 // Хранит телефон в формате 79151596781
}

// NewPhone - телефон в формате +7 (915) 159-67-81
func NewPhone(number string) (Phone, error) {
	compiler, err := regexp.Compile("^(\\+7)+ \\(\\d{3}\\) \\d{3}-\\d{2}-\\d{2}$")
	if err != nil {
		return Phone{}, fmt.Errorf("failed to compile regex: %w", err)
	}

	if !compiler.MatchString(number) {
		return Phone{}, errors.New("invalid format, example +7 (915) 159-67-81")
	}

	filtered := ""
	for _, ch := range number {
		if unicode.IsDigit(ch) {
			filtered += string(ch)
		}
	}

	value, err := strconv.ParseInt(filtered, 10, 64)
	if err != nil {
		return Phone{}, fmt.Errorf("failed to parse int: %w", err)
	}

	return Phone{
		number: value,
	}, nil
}

func NewPhoneFromInt64(number int64) Phone {
	return Phone{
		number: number,
	}
}

func (p Phone) Number() int64 {
	return p.number
}
