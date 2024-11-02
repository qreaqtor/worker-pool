package apiconsole

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type workersPool interface {
	Delete(id []int)
	Add(n int)
	Alive() []int
	Work(jobs []string)
	GetJobs() []string
}

type ConsoleWorkersAPI struct {
	workers workersPool
}

func NewConsoleWorkersAPI(workers workersPool) *ConsoleWorkersAPI {
	return &ConsoleWorkersAPI{
		workers: workers,
	}
}

const (
	cmdAdd    = "add"
	cmdDelete = "delete"
	cmdAlive  = "alive"
	cmdWork   = "work"
	cmdJobs   = "jobs"
)

func (w *ConsoleWorkersAPI) Register(reader io.Reader, closeCtx chan<- os.Signal) {
	scanner := bufio.NewScanner(reader)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			tokens := strings.Split(line, " ")

			switch tokens[0] {
			case cmdAdd:
				w.add(tokens)
			case cmdDelete:
				w.delete(tokens)
			case cmdAlive:
				w.alive()
			case cmdWork:
				w.work(tokens)
			case cmdJobs:
				w.jobs()
			}
		}

		closeCtx <- nil // пишу в канал, который создал в main для отмены контекста если получил EOF
	}()
}

func (w *ConsoleWorkersAPI) delete(tokens []string) {
	if len(tokens) < 2 {
		return
	}

	ids := make([]int, 0, len(tokens)-1)
	for i := 1; i < len(tokens); i++ {
		id, err := strconv.Atoi(tokens[i])
		if err != nil {
			slog.Debug(fmt.Sprint("can't parse to int:", tokens[i]))
			continue
		}

		ids = append(ids, id)
	}

	w.workers.Delete(ids)
}

func (w *ConsoleWorkersAPI) add(tokens []string) {
	if len(tokens) != 2 {
		return
	}

	count, err := strconv.Atoi(tokens[1])
	if err != nil {
		slog.Debug(fmt.Sprint("can't parse to int:", tokens[1]))
		return
	}

	w.workers.Add(count)
}

func (w *ConsoleWorkersAPI) alive() {
	slog.Info(fmt.Sprint(w.workers.Alive()))
}

func (w *ConsoleWorkersAPI) work(tokens []string) {
	if len(tokens) < 2 {
		return
	}
	w.workers.Work(tokens[1:])
}

func (w *ConsoleWorkersAPI) jobs() {
	slog.Info(fmt.Sprint(w.workers.GetJobs()))
}
