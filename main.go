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
	fmt.Println("- Add Task")
	fmt.Println("- View Tasks")
	fmt.Println("- Mark Task as Completed")
	fmt.Println("- Run the program with the <help> argument to see how to use the above operations\n")
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
	t := Task{}

	db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte("TasksBucket")) //open bucket

		id, _ := b.NextSequence() //generate ID 
		t.Id = int(id)
		t.TaskName = args[0] //fill in task object

		buff, err := json.Marshal(t.TaskName) //marshal user data into bytes

		if err != nil {
			return err
		}

		e := b.Put([]byte(itob(t.Id)), buff) //persist bytes to users bucket

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
		Description: "    initialized tasks bucket\n\t    Usage: init\n",
		ExecFunc: func(ctx context.Context, args []string) error {
			initializeBucket()
			return nil
		},
	},
	{
		Name: "add",
		Description: "     adds a task to db\n\t     Usage: add <taskName>\n",
		ExecFunc: func(ctx context.Context, args []string) error {
	
			addKeyValue(args)
			return nil
		},
	},
	{
		Name: "complete",
		Description: "completes a task (removes from db)\n\tUsage: complete <taskID>\n",
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




