package main

import (
	"context"
	"fmt"
	"github.com/cristalhq/acmd"
	"github.com/boltdb/bolt"
	"log"
)

func home() {

	fmt.Println("\nWelcome to the TODO task manager!\n")
	fmt.Println("Your following options are:\n")
	fmt.Println("Add Task")
	fmt.Println("View Tasks")
	fmt.Println("Mark Task as Completed")
}

func openDataBase() {

	db, err := bolt.Open("my.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

}

func showTasks() {

	db, _ := bolt.Open("my.db", 0600, nil)

	db.View(func(tx *bolt.Tx) error {
	// Assume bucket exists and has keys
	b := tx.Bucket([]byte("TasksBucket"))

	c := b.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		fmt.Printf("%s. %s\n", k, v)
	}

	return nil
})
}

func addKeyValue(args []string) {

	db, _ := bolt.Open("my.db", 0600, nil)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TasksBucket"))

		err := b.Put([]byte(args[0]), []byte(args[1]))

		return err
	})
}

func initializeBucket() {

	db, err := bolt.Open("my.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("TasksBucket"))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func deleteTask(taskId string) {

	db, err := bolt.Open("my.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TasksBucket"))

		err := b.Delete([]byte(taskId))

		return err
	})
}



func main() {

	cmds := []acmd.Command{
	{
		Name: "todo",
		Description: "shows home screen",
		ExecFunc: func(ctx context.Context, args []string) error {
			home()
			openDataBase()
			return nil
		},
	},
	{
		Name: "show",
		Description: "shows all tasks",
		ExecFunc: func(ctx context.Context, args []string) error {
			showTasks()
			return nil
		},
	},
	{
		Name: "init",
		Description: "initialized tasks bucket",
		ExecFunc: func(ctx context.Context, args []string) error {
			initializeBucket()
			return nil
		},
	},
	{
		Name: "add",
		Description: "adds a task to db",
		ExecFunc: func(ctx context.Context, args []string) error {
			addKeyValue(args)
			return nil
		},
	},
	{
		Name: "complete",
		Description: "completes a task (removes from db)",
		ExecFunc: func(ctx context.Context, args []string) error {
			deleteTask(args[0])
			return nil
		},
	},
}

r := acmd.RunnerOf(cmds, acmd.Config{
})

if err := r.Run(); err != nil {
	fmt.Println(err)
	r.Exit(err)
}
}


//we have a bucket (collection) called 'tasks' that we query from
//we add key/value pairs to this tasks bucket and read and delete from it


