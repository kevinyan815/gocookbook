package main

import (
    "errors"
    "flag"
    "fmt"
    "os"
)

func NewGreetCommand() *GreetCommand {
    gc := &GreetCommand{
        fs: flag.NewFlagSet("greet", flag.ContinueOnError),
    }

    gc.fs.StringVar(&gc.name, "name", "World", "name of the person to be greeted")

    return gc
}

type GreetCommand struct {
    fs *flag.FlagSet

    name string
}

func (g *GreetCommand) Name() string {
    return g.fs.Name()
}

func (g *GreetCommand) Init(args []string) error {
    return g.fs.Parse(args)
}

func (g *GreetCommand) Run() error {
    fmt.Println("Hello", g.name, "!")
    return nil
}

type Runner interface {
    Init([]string) error
    Run() error
    Name() string
}

func root(args []string) error {
    if len(args) < 1 {
        return errors.New("You must pass a sub-command")
    }

    cmds := []Runner{
        NewGreetCommand(),
    }

    subcommand := os.Args[1]

    for _, cmd := range cmds {
        if cmd.Name() == subcommand {
            cmd.Init(os.Args[2:])
            return cmd.Run()
        }
    }

    return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func main() {
    if err := root(os.Args[1:]); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
