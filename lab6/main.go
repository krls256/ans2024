package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"math/rand/v2"
	"os"
	"strings"
)

type Record struct {
	ID      string   `json:"id"`
	Friends []string `json:"friends"`
}

func readCSV(filename string) []Record {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	records := []Record{}

	i := 0
	for _, rowBts := range bytes.Split(content, []byte("\n"))[1:] {
		row := string(rowBts)
		i++

		if i%1000 == 0 {
			fmt.Println(i)
		}

		if len(row) == 0 {
			continue
		}

		firstColumnIndex := strings.Index(row, ",")
		from := strings.LastIndex(row, "[")
		to := strings.LastIndex(row, "]")

		friends := row[from+1 : to]

		records = append(records, Record{
			ID: strings.ReplaceAll(strings.TrimSpace(row[:firstColumnIndex]), `"`, ""),
			Friends: lo.Map(strings.Split(friends, ","), func(item string, index int) string {
				return strings.ReplaceAll(strings.TrimSpace(item), `"`, "")
			}),
		})
	}

	return records
}

func basicCleanup(records []Record) []Record {
	recordsMap := map[string]struct{}{}
	friendsMap := map[string]struct{}{}

	for _, r := range records {
		recordsMap[r.ID] = struct{}{}

		for _, f := range r.Friends {
			friendsMap[f] = struct{}{}
		}
	}

	recordsFiltered := []Record{}

	for i, r := range records {
		realFriends := lo.Filter(r.Friends, func(item string, index int) bool {
			_, ok := recordsMap[item]

			return ok
		})

		if i%1000 == 0 {
			fmt.Println(len(recordsMap), len(r.Friends), len(realFriends))
		}

		_, hasFriend := friendsMap[r.ID]
		if len(realFriends) == 0 && !hasFriend {
			continue
		}

		recordsFiltered = append(recordsFiltered, Record{
			ID:      r.ID,
			Friends: realFriends,
		})
	}

	return recordsFiltered
}

func filterSomeLowFriendsNodes(records []Record, prob float64, lowLimit int) []Record {
	friendsMap := map[string]int{}

	for _, r := range records {

		for _, f := range r.Friends {
			friendsMap[f]++
		}
	}

	var newRecords []Record

	for _, r := range records {
		if friendsMap[r.ID] < lowLimit && rand.Float64() < prob {
			continue
		}

		newRecords = append(newRecords, r)
	}

	return basicCleanup(newRecords)
}

func addSomeConnections(records []Record) []Record {
	count := len(records) / 2

	for range count {
		from, to := rand.IntN(len(records)), rand.IntN(len(records))

		if from != to {
			records[from].Friends = lo.Uniq(append(records[from].Friends, records[to].ID))
		}
	}

	return records
}

func main() {
	records := readCSV("./data.csv")
	recordsFiltered := basicCleanup(records)
	smallGraph := filterSomeLowFriendsNodes(recordsFiltered, 0.97, 30)
	recordsWithNewEdges := addSomeConnections(smallGraph)

	bytes, err := json.Marshal(recordsWithNewEdges)

	err = os.WriteFile("output.json", bytes, 0777)
	if err != nil {
		panic(err)
	}
}
