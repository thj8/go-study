package engine

import "fmt"

func doWork(r Request) Respone {
	fmt.Println(r.Url)
	return Respone{Back: "back " + r.Url}
}
