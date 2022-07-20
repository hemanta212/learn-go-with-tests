package main

import (
	"strings"
)

type RomanNumeral struct {
	Value  uint16
	Symbol string
}

var allRomanNumerals = []RomanNumeral{
	{1000, "M"},
	{900, "CM"},
	{500, "D"},
	{400, "CD"},
	{100, "C"},
	{90, "XC"},
	{50, "L"},
	{40, "XL"},
	{10, "X"},
	{9, "IX"},
	{5, "V"},
	{4, "IV"},
	{1, "I"},
}

func ConvertToRoman(arabic uint16) string {
	var result strings.Builder
	for _, numeral := range allRomanNumerals {
		for arabic >= numeral.Value {
			result.WriteString(numeral.Symbol)
			arabic -= numeral.Value
		}
	}
	return result.String()
}

func ConvertToArabic(roman string) uint16 {
	total := uint16(0)
	for i := 0; i < len(roman); {
		roman_part := parseNumeral(i, roman)
		for _, numeral := range allRomanNumerals {
			if numeral.Symbol == roman_part {
				total += numeral.Value
			}
		}
		i += len(roman_part)
	}

	return total
}

func parseNumeral(i int, roman string) string {
	if i >= len(roman)-1 {
		return roman[i : i+1]
	}
	biNumeral, isBinumeral := roman[i:i+2], false
	for _, numeral := range allRomanNumerals {
		if numeral.Symbol == biNumeral {
			isBinumeral = true
		}
	}
	if !isBinumeral {
		return biNumeral[0:1]
	}
	return biNumeral
}
