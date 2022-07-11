package main

import "fmt"

const spanish = "Spanish"
const french = "French"
const nepali = "Nepali"
const englishHelloPrefix = "Hello, "
const spanishHelloPrefix = "Hola, "
const frenchHelloPrefix = "Bonjour, "
const nepaliHelloPrefix = "Namaste, "

func Hello(name string, language string) string {
	if name == "" {
		name = "World!"
	}
	return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
	switch language {
	case spanish:
		prefix = spanishHelloPrefix
	case french:
		prefix = frenchHelloPrefix
	case nepali:
		prefix = nepaliHelloPrefix
	default:
		prefix = englishHelloPrefix
	}
	return
}

func main() {
	fmt.Println(Hello("", ""))
}
