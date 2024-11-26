package main

func main() {
	server := NewAPIserver(":3000")
	server.Run()

}
