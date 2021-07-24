package main

import (
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
	fmt.Printf("%d\n",a)
	w.Write([]byte("Hello"))
}
