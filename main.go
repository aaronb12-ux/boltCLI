package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"github.com/boltdb/bolt"
	"github.com/cristalhq/acmd"
)

type Task struct {
	Id int
	TaskName string
}

func home() {

	fmt.Println("\nWelcome to the TODO task manager!\n")
	fmt.Println("Your following options are:\n")
	fmt.Println("- Initialize Task Bucket: init")
	fmt.Println("- Add Task: add <taskToAdd>")
	fmt.Println("- View Tasks: show")
	fmt.Println("- Mark Task as Completed: complete <taskID>")
	fmt.Println("- Run the program with the <help> argument to see how to use the above operations\n")
}



func showTasks() {

	db, _ := bolt.Open("my.db", 0600, nil)

	db.View(func(tx *bolt.Tx) error {

	b := tx.Bucket([]byte("TasksBucket"))

	b.ForEach(func(k, v []byte) error {
		fmt.Printf("key=%s, value=%s\n", k, v)
		return nil
	})

	return nil
})
}

func addKeyValue(args []string) {

	db, _ := bolt.Open("my.db", 0600, nil)
	

	db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("TasksBucket")) //open bucket

		id, _ := b.NextSequence() //generate ID 

		buff, err := json.Marshal(args[0]) //marshal task data into bytes

		if err != nil {
			return err
		}

		e := b.Put([]byte(itob(int(id))), buff) //persist bytes to users bucket
		
		return e
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
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

func deleteBucket() {

	db, err := bolt.Open("my.db", 0600, nil)
	

	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TasksBucket"))

		b.ForEach(func(k, v []byte) error {
		
		err := b.Delete([]byte(k))
		
		return err
	})

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
		Description: "initializes tasks bucket",
		ExecFunc: func(ctx context.Context, args []string) error {
			initializeBucket()
			return nil
		},
	},
	{
		Name: "add",
		Description: "adds a task to database",
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
	{
		Name: "reset",
		Description: "deletes the current bucket",
		ExecFunc: func(ctx context.Context, args []string) error {
			deleteBucket()
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




