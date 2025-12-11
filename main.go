package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/cristalhq/acmd"
)

type Database struct {
	database *bolt.DB
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



func showTasks(db *bolt.DB) {

	db.View(func(tx *bolt.Tx) error {

	b := tx.Bucket([]byte("TasksBucket"))

	b.ForEach(func(k, v []byte) error {
		fmt.Printf("%v: %s\n", k[7], v[1 : len(v) - 1])
		return nil
	})

	return nil
})
}

func addKeyValue(args []string, db *bolt.DB) {

	db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("TasksBucket")) //open bucket

		id, _ := b.NextSequence() //generate ID 

		buff, err := json.Marshal(args[0]) //marshal task data into bytes

		if err != nil {
			return err
		}

		fmt.Printf("the id is %d", id)

		e := b.Put([]byte(itob(int(id))), buff) //persist bytes to users bucket
		
		return e
	})
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func initializeBucket(db *bolt.DB) {

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("TasksBucket"))

		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func deleteTask(taskId string, db *bolt.DB) {

	integerId, err := strconv.Atoi(taskId) //convert string to byte before deleting

	if err != nil {
		log.Fatalf("Error converting string to int %v", taskId)
	}

	db.Update(func(tx *bolt.Tx) error {
	
		b := tx.Bucket([]byte("TasksBucket"))
		err := b.Delete([]byte(itob(integerId)))

		return err
	})
}

func deleteBucket(db *bolt.DB) {


	db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("TasksBucket"))

		b.ForEach(func(k, v []byte) error {
		
		err := b.Delete([]byte(k))
		
		return err
	})

		return nil
	})
}


func main() {

	db := &Database{} 

	d, _ := bolt.Open("my.db", 0600, nil)

	db.database = d

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
			showTasks(db.database)
			return nil
		},
	},
	{
		Name: "init",
		Description: "initializes tasks bucket",
		ExecFunc: func(ctx context.Context, args []string) error {
			initializeBucket(db.database)
			return nil
		},
	},
	{
		Name: "add",
		Description: "adds a task to database",
		ExecFunc: func(ctx context.Context, args []string) error {
			addKeyValue(args, db.database)
			return nil
		},
	},
	{
		Name: "complete",
		Description: "completes a task (removes from db)",
		ExecFunc: func(ctx context.Context, args []string) error {
			deleteTask(args[0], db.database)
			return nil
		},
	},
	{
		Name: "reset",
		Description: "deletes the current bucket",
		ExecFunc: func(ctx context.Context, args []string) error {
			deleteBucket(db.database)
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






