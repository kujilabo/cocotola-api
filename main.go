package main

import (
	"errors"
	"fmt"
	"net/http"

	f "github.com/kujilabo/cocotola-api/pkg/func"
)

func init() {
	http.HandleFunc("/", HealthcheckHandler)
}

func main() {
	port := "8080"
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	a:=f.Abc(2)
	validate(a)
	fmt.Printf("%d\n",a)
	w.Write([]byte("Hello"))
}

func validate(age int) error {
    if age < 20 {
        return errors.New("age should be 20 or more")
    }
    fmt.Println("ok~")
    return nil
}

