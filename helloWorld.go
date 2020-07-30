package main
import (
	"fmt"
)

func main() {
		//m := make(map[string]int)
    clothes := [][]string{{"A","a"}, {"B","b"}}
    for _, vk := range clothes {
        fmt.Println(vk[0])
    }
}
