package main

import (
	"context"
	
	"fmt"


	"github.com/cristalhq/acmd"
)

func home() {

	fmt.Println("\nWelcome to the TODO task manager!\n")
	fmt.Println("Your following options are:\n")
	fmt.Println("Add Task - <command>")
	fmt.Println("View Tasks - <command>")
	fmt.Println("Mark Task as Completed")
	fmt.Println("Type 'help' for a basic rundown of the tool")
}

func main() {

	cmds := []acmd.Command{
	{
		Name:        "todo",
		Description: "shows home screen",
		ExecFunc: func(ctx context.Context, args []string) error {
			home()
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
