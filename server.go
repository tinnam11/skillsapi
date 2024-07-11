package main

import "fmt"
import "os"

func main() {
	fmt.Println("Hello, World!")

	url := os.Getenv("DATABASE_URL")

	fmt.Println("url:", url)

	fmt.Println("done")
}
