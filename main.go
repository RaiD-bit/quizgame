package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type Record struct {
	Problem string
	Answer  string
}

func readQuestionCsv(path string) ([]Record, error) {
	file, err := os.Open(path)
	if err != nil {
		return []Record{}, err
	}
	defer file.Close()
	var res []Record
	records, err := csv.NewReader(file).ReadAll()
	for _, item := range records {
		ques := item[0]
		ans := item[1]
		res = append(res, Record{Answer: ans, Problem: ques})
	}
	return res, nil
}

func process(ctx context.Context, questionCsv []Record, correct *int, incorrect *int) {
	for _, record := range questionCsv {
		fmt.Printf("Problem: %s\n", record.Problem)
		var userinput string
		fmt.Scanln(&userinput)
		if userinput == record.Answer {
			*correct += 1
		} else {
			*incorrect += 1
		}
		fmt.Println("--------------------")
	}

	for {
		select {
		case <-time.After(1 * time.Second):
			// do nothing
		case <-ctx.Done():
			fmt.Println("times up !! exiting")
			return
		}
	}
}

func main() {
	correct := 0
	incorrect := 0
	total := 0
	// get path to a file
	totalTime := flag.Int64("time", 0, "quiz timing in minutes")
	filepath := flag.String("csv", "", "provide the path for the csv")
	flag.Parse()
	fmt.Println(*totalTime)
	questionCsv, err := readQuestionCsv(*filepath)
	if err != nil {
		return
	}
	total = len(questionCsv)

	ctx, cancel := context.WithCancel(context.Background())

	go process(ctx, questionCsv, &correct, &incorrect)

	time.Sleep(time.Duration(*totalTime) * time.Minute)

	cancel()

	fmt.Printf("Correct : %d Incorrect : %d Total : %d", correct, incorrect, total)
}
