package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Fault struct {
	FaultType string
	Address   string
}

// reduces duplicate pages, returning a shortened copy
func reduceDuplicates(faults []Fault) []Fault {
	tempFaults := make([]Fault, 0)
	i := 0
	for i < len(faults)-1 {
		tempFaults = append(tempFaults, faults[i])
		j := i
		for faults[i].Address == faults[j].Address && faults[i].FaultType == faults[j].FaultType {
			j = j + 1
		}
		i = j
	}
	tempFaults = append(tempFaults, faults[len(faults)-1])
	return tempFaults
}

func main() {

	faults := make([]Fault, 0)
	pathPtr := flag.String("path", "default value", "file path to the data")
	flag.Parse()
	data, err := os.ReadFile(*pathPtr)

	if err != nil {
		panic(err)
	}

	content := string(data)
	r := csv.NewReader(strings.NewReader(content))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			//fmt.Println(err)
		}
		if record[1] != "" && record[0] != "" {
			faults = append(faults, Fault{FaultType: record[0], Address: record[1]})
		} else {
			faults = append(faults, Fault{record[0], "na"})
		}
	}

	fmt.Println(len(faults))
	b, err := json.MarshalIndent(reduceDuplicates(faults), "", " ")

	if err != nil {
		log.Fatal(err)
	}

	err2 := os.WriteFile("data.json", b, 0666)

	if err2 != nil {

	}

	fmt.Println("done")
	fmt.Println(faults[0])
}
