package main

import (
	"errors"
	"fmt"

	"sis/internal/service"
	"sis/internal/storage"
	"sis/internal/task"
)

func main() {
	svc := service.New(storage.NewMemoryStore())

	a := mustCreate(svc, "Send email", "user@mail.com")
	b := mustCreate(svc, "Parse report", "report.csv")
	c := mustCreate(svc, "Backup DB", "prod")

	must(svc.Start(a.ID))
	must(svc.Complete(a.ID, "email sent"))

	must(svc.Start(b.ID))
	must(svc.Fail(b.ID, "file is corrupted"))

	must(svc.Cancel(c.ID))

	fmt.Println("All tasks")
	for _, t := range svc.List() {
		fmt.Printf("#%d %-13s %-10s result=%q error=%q\n",
			t.ID, t.Title, t.Status, t.Result, t.Error)
	}

	fmt.Println("\nError handling")

	_, err := svc.Create("   ", "x")
	fmt.Println("empty title:", err, "| is ErrEmptyTitle:", errors.Is(err, task.ErrEmptyTitle))

	_, err = svc.Start(999)
	fmt.Println("missing task:", err, "| is ErrNotFound:", errors.Is(err, storage.ErrNotFound))

	_, err = svc.Start(a.ID) // a уже done
	var te *task.TransitionError
	if errors.As(err, &te) {
		fmt.Printf("bad transition: %v (from=%s, to=%s)\n", err, te.From, te.To)
	}
}

func mustCreate(svc *service.Service, title, payload string) task.Task {
	t, err := svc.Create(title, payload)
	if err != nil {
		panic(err)
	}
	return t
}

func must(_ task.Task, err error) {
	if err != nil {
		panic(err)
	}
}
