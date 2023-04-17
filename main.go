package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Fault struct {
	FaultType string
	Address   int64
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

// returns ordred array of faults.
// NOTE: our data seems to be in referenetial order already but this might not always be the case
func getUnique(faults []Fault) []int64 {
	uniqueAdresses := make([]int64, 0)
	//starts at 2 because of s at start
	for i := 2; i < len(faults); i++ {
		key := faults[i]
		j := i - 1

		for j >= 0 && key.Address < faults[j].Address {
			faults[j+1] = faults[j]
			j--
		}
		faults[j+1] = key
	}

	for i := 0; i < len(faults); i++ {
		uniqueAdresses = append(uniqueAdresses, faults[i].Address)
	}

	return uniqueAdresses
}

func main() {

	faults := make([]Fault, 0)
	pathPtr := flag.String("path", "default value", "file path to the data")
	flag.Parse()
	data, err := os.Open(*pathPtr)

	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(data)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			//fmt.Println(err)
		}
		if record[1] != "" && record[0] != "" {
			conv, err := strconv.ParseInt(record[1], 16, 0)
			if err != nil {
				log.Fatal(err)
			}
			faults = append(faults, Fault{FaultType: record[0], Address: conv})
		} else {
			faults = append(faults, Fault{record[0], -1})
		}
	}
	data.Close()

	fmt.Println(len(faults))
	b, err := json.MarshalIndent(reduceDuplicates(faults), "", " ")

	if err != nil {
		log.Fatal(err)
	}

	err2 := os.WriteFile("data.json", b, 0666)

	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println(faults[0])

	uf, err3 := os.Create("unique.txt")
	if err3 != nil {
		log.Fatal(err3)
	}

	uniqueAdresses := getUnique(reduceDuplicates(faults))
	for i := 0; i < len(uniqueAdresses); i++ {
		line := strconv.FormatInt(uniqueAdresses[i], 10) + "\n"
		_, err := uf.Write([]byte(line))
		if err != nil {
			log.Fatal((err))
		}

	}
	uf.Close()
	fmt.Println("done")
}
