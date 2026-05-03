package regex

import (
	"apigateway/services/product/internal/domain"
	"regexp"
	"strings"
	"unicode"
)

func ValidateName(str string) error {
	matched, _ := regexp.MatchString(`^[A-Z][a-zA-Z\d, -]+$`, str)
	if !matched {
		return domain.ErrInvalidName
	}

	var letterCount int
	for _, char := range str {
		if unicode.IsLetter(char) {
			letterCount++
		}
	}
	if letterCount < 2 {
		return domain.ErrInvalidName
	}

	var digitCount int
	strSlice := strings.Fields(str)
	for _, word := range strSlice {
		letterCount = 0
		digitCount = 0
		for _, ch := range word {
			if unicode.IsLetter(ch) {
				letterCount++
			}
			if unicode.IsDigit(ch) {
				digitCount++
			}
		}
		if letterCount < 2 && digitCount == 0 {
			return domain.ErrInvalidName
		}
	}

	return nil
}

func ValidateManufacturer(str string) error {
	matched, _ := regexp.MatchString(`^[a-zA-Z -]+$`, str)
	if !matched {
		return domain.ErrInvalidManufacturer
	}

	var letterCount int
	for _, char := range str {
		if unicode.IsLetter(char) {
			letterCount++
		}
	}
	if letterCount < 2 {
		return domain.ErrInvalidManufacturer
	}

	return nil
}

func ValidateCategory(str string) error {
	matched, _ := regexp.MatchString(`^[A-Z][a-zA-Z, -]+$`, str)
	if !matched {
		return domain.ErrInvalidCategory
	}

	var letterCount int
	for _, char := range str {
		if unicode.IsLetter(char) {
			letterCount++
		}
	}
	if letterCount < 4 {
		return domain.ErrInvalidCategory
	}

	return nil
}
