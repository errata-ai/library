package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/errata-ai/library/internal/nlp"
	"github.com/errata-ai/library/pkg/data"
	"github.com/errata-ai/library/pkg/search"
	"github.com/go-resty/resty/v2"
	"github.com/urfave/cli/v2"
)

func printJSON(t interface{}) error {
	bf := bytes.NewBuffer([]byte{})

	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)

	err := jsonEncoder.Encode(t)
	if err != nil {
		panic(err)
	}

	fmt.Println(bf.String())
	return nil
}

func main() {
	app := &cli.App{
		Name:  "lois",
		Usage: "A local, offline search engine",
		Commands: []*cli.Command{
			{
				Name: "questions",
				Action: func(c *cli.Context) error {
					src, err := data.FromGHIssues([]string{c.Args().First()})
					if err != nil {
						return err
					}
					return printJSON(src)
				},
			},
			{
				Name: "read",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
					if len(args) != 2 {
						return errors.New("2 argument required")
					}

					src, err := data.FromLinkList(args[0])
					if err != nil {
						return err
					}

					issues, err := data.FromGHIssues([]string{"vale"})
					if err != nil {
						return err
					}
					src.Sources = append(src.Sources, issues.Sources...)

					_, err = search.NewEngineFromData(args[1], src)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name: "create",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
					if len(args) != 2 {
						return errors.New("2 argument required")
					}

					src, err := data.FromMapping(args[0])
					if err != nil {
						return err
					}

					_, err = search.NewEngineFromData(args[1], src)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name: "search",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
					if len(args) != 2 {
						return errors.New("2 arguments required")
					}

					engine, err := search.LoadEngine(args[0])
					if err != nil {
						return err
					}

					results, err := engine.Search(args[1])
					if err != nil {
						return err
					}

					return printJSON(results)
				},
			},
			{
				Name:  "lookup",
				Usage: "Finds the context for a given result ID",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
					if len(args) != 2 {
						return errors.New("2 arguments required")
					}

					m, err := data.ParseMap(args[1])
					if err != nil {
						return err
					}

					context, err := data.Context(args[0], m.Path)
					if err != nil {
						return err
					}
					fmt.Println(context)

					return nil
				},
			},
			{
				Name:  "parse",
				Usage: "",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
					if len(args) != 2 {
						return errors.New("2 arguments required")
					}

					m, err := data.ParseMap(args[1])
					if err != nil {
						return err
					}

					context, err := data.Context(args[0], m.Path)
					if err != nil {
						return err
					}

					tokens, err := nlp.TextToTokens(context)
					if err != nil {
						return err
					}

					return printJSON(tokens)
				},
			},
			{
				Name:  "post",
				Usage: "",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
					if len(args) != 2 {
						return errors.New("2 arguments required")
					}
					client := resty.New()

					resp, err := client.R().
						SetHeader("Content-Type", "application/x-www-form-urlencoded").
						SetQueryParam("text", args[1]).
						Post(args[0])

					if err != nil {
						return err
					}
					fmt.Println(resp)

					return nil
				},
			},
			{
				Name:  "fsm",
				Usage: "",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
					if len(args) != 2 {
						return errors.New("2 arguments required")
					}

					out, err := nlp.DoSectionize(args[0], args[1])
					if err != nil {
						return err
					}

					return printJSON(out)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
