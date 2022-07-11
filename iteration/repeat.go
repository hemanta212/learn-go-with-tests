package main

func main() {
	Repeat("z", 5)
}

func Repeat(charectar string, n int) string {
	var result string
	for i := 0; i < n; i++ {
		result += charectar
	}
	return result
}
