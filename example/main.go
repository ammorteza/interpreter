package main

import (
	"log"
	"os"

	"github.com/ammorteza/interpreter"
	"github.com/urfave/cli"
)

func main() {
	inter := interpreter.New()
	app := &cli.App{
		Name:  "Brainfuck Interpreter",
		Usage: "A Brainfuck cli interpreter",
		Action: func(c *cli.Context) error {
			if len(c.Args()) > 0 {
				file, err := os.Open(c.Args().Get(0))
				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()
				inter.AddOpt('*', func(p *interpreter.Payload) {
					p.Ram[*p.Ptr] *= p.Ram[*p.Ptr]
				})
				inter.Interpret(file)
				res, err := inter.Execute()
				if err != nil {
					log.Println(err)
				}
				log.Println(string(res))
			} else {
				log.Fatal("Fatal error: No input file\n")
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
