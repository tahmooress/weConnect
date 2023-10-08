package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tahmooress/weConnect-task/cmd"
	"github.com/tahmooress/weConnect-task/internal/entity"
	"github.com/tahmooress/weConnect-task/internal/reader"
	"github.com/tahmooress/weConnect-task/internal/repository"
	"github.com/tahmooress/weConnect-task/internal/repository/mongodb"
	"github.com/tahmooress/weConnect-task/internal/workerpool"
)

var InvalidDataFormat = errors.New("invalid data format")

func main() {
	path := flag.String("path", "", "path of data source file")
	workers := flag.Int("workers", 10, "number of workers")

	flag.Parse()

	if *path == "" {
		log.Fatal("set data source file address with path flag")
	}

	ctx, cancel := context.WithCancel(context.Background())

	repo, err := mongodb.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go cmd.Shutdown(ctx, cancel)

	if err := run(ctx, repo, *path, *workers); err != nil {
		cancel()
		log.Fatal(err)
	}

	cancel()
	os.Exit(0)
}

func run(ctx context.Context, repo repository.Repository, filePath string, workers int) error {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	r := reader.NewReader(csv.NewReader(f))
	records := r.ReadAll()

	wp := workerpool.NewPool(workers)

	// defer are called in reverse order so first close
	// the incomming stream then stop the workerPool.
	defer wp.Stop()
	defer r.Stop()

	taskQueue := make(chan func() error)
	go func() {
		defer close(taskQueue)
		for record := range records {
			taskQueue <- func() error {
				statistics, err := statisticsBuilder(record)
				if err != nil {
					// if data format is invalid just threw way
					if errors.Is(err, InvalidDataFormat) {
						fmt.Printf("invalid record format: %s\n", printSlice(record))
						return nil
					}

					return err
				}
				id, err := repo.Insert(ctx, statistics)
				if err != nil {
					return err
				}

				fmt.Println("record inserted with id: ", id)

				return nil
			}
		}
	}()

	errc := wp.Run(taskQueue)

	for {
		select {
		case err, ok := <-errc:
			if !ok || err != nil {
				return err
			}
		case err := <-r.Err():
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func statisticsBuilder(record []string) (*entity.Statistics, error) {
	if len(record) < 14 {
		return nil, InvalidDataFormat
	}

	dataValue, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		if record[2] == "" {
			dataValue = 0.0
		} else {
			return nil, InvalidDataFormat
		}
	}

	magnitude, err := strconv.Atoi(record[6])
	if err != nil {
		if record[6] == "" {
			magnitude = 0
		} else {
			return nil, InvalidDataFormat
		}
	}

	return &entity.Statistics{
		SeriesReference: record[0],
		Period:          record[1],
		DataValue:       dataValue,
		Suppressed:      record[3],
		Status:          record[4],
		Units:           record[5],
		Magnitude:       magnitude,
		Subject:         record[7],
		Group:           record[8],
		SeriesTitle1:    record[9],
		SeriesTitle2:    record[10],
		SeriesTitle3:    record[11],
		SeriesTitle4:    record[12],
		SeriesTitle5:    record[13],
	}, nil
}

func printSlice(strSlice []string) string {
	var res string
	for _, str := range strSlice {
		res += fmt.Sprintf("%s, ", str)
	}

	return strings.TrimRight(res, ",")
}
