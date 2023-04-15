package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"flag"
	"log"
	"encoding/csv"
	"encoding/json"
)

type Fault struct {

	FaultType string
	Address string

}

func main() {

	faults := make([]Fault, 5)
	pathPtr := flag.String("path", "default value", "file path to the data")

	flag.Parse()
	data, err := os.ReadFile(*pathPtr)
	if err != nil {
		panic(err)
	}
	content := string(data)
	r := csv.NewReader(strings.NewReader(content))
	for {
		record, err := r.Read();
		if err == io.EOF {
			break
		}
		if err != nil {
			//fmt.Println(err)
		}

		if (record[1] != "" && record[0] != "") {
			faults = append(faults,Fault{FaultType: record[0],Address: record[1]})

		} else {
			faults = append(faults,Fault{record[0],"na"})
			
		}
		
	}
		b, err := json.MarshalIndent(faults,""," ")

		if err != nil{
			log.Fatal(err)
		}

		err2 := os.WriteFile("data.json",b,0666)
		if err2 != nil {

		}

		fmt.Println("done")
		fmt.Println(faults[0])
}