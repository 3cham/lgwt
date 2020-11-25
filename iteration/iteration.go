package iteration

func Repeat(char string) (value string) {
	for i := 0; i < 5; i++ {
		value += char
	}
	return
}
