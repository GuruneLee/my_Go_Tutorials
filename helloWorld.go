package main
import (
	"fmt"
	"encoding/json"
)

//Member -
type Member struct {
	Name string
	Age int
	Active bool
}

func main() {
	//테스트용 JSON 데이터
	jsonBytes, err := json.Marshal(Member{"Tim", 1, true})

	//JSON 디코딩
	var mem Member
	err = json.Unmarshal(jsonBytes, &mem)
	if err != nil {
		panic(err)
	}

	fmt.Println(mem.Name, mem.Age, mem.Active)
}
