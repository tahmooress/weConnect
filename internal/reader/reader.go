package reader

import (
	"encoding/csv"
	"errors"
	"io"
	"sync"
)

// Reader takes a csv.Reader and read it concurrenly
// and send each records through a retruned channel
// is up to users to call Err() to takes any occuring
// error during the proccess.
type Reader struct {
	reader     *csv.Reader
	out        chan []string
	errc       chan error
	onceCloser *sync.Once
	done       chan struct{}
}

func NewReader(reader *csv.Reader) *Reader {
	return &Reader{
		reader:     reader,
		out:        make(chan []string),
		errc:       make(chan error, 1),
		onceCloser: &sync.Once{},
		done:       make(chan struct{}),
	}
}

func (r *Reader) ReadAll() <-chan []string {
	go func() {
		defer close(r.out)

		for {
			select {
			case <-r.done:
				return
			default:
				records, err := r.reader.Read()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					r.errc <- err
					return
				}
				r.out <- records
			}
		}
	}()

	return r.out
}

// Err return any encountering error during
// reading the csv file. caller should always call
// and handle the possible error if exsit any.
func (r *Reader) Err() <-chan error {
	return r.errc
}

// Stop will stop the reading proccess of the underling
// reader error and close the returned channel
// from ReadAll function call.
func (r *Reader) Stop() {
	r.onceCloser.Do(func() {
		close(r.done)
	})
}
