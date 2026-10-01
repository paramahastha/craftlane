package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/craftlane/craftlane/internal/app"
	"github.com/craftlane/craftlane/internal/diagnostics"
	"github.com/craftlane/craftlane/internal/project"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "craftlane:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	command := "start"
	if len(args) > 0 {
		command, args = args[0], args[1:]
	}
	ctx := context.Background()
	switch command {
	case "start":
		flags := flag.NewFlagSet("start", flag.ContinueOnError)
		noOpen := flags.Bool("no-open", false, "do not open browser")
		headless := flags.Bool("headless", false, "do not open browser")
		port := flags.Int("port", 7331, "loopback port")
		if err := flags.Parse(args); err != nil {
			return err
		}
		explicitPort := false
		flags.Visit(func(item *flag.Flag) { explicitPort = explicitPort || item.Name == "port" })
		record, err := app.StartDetached(ctx, ".", *port, explicitPort, !*noOpen && !*headless)
		if err != nil {
			return err
		}
		fmt.Printf("Craftlane ready: %s\n", record.Origin)
		return nil
	case "serve":
		flags := flag.NewFlagSet("serve", flag.ContinueOnError)
		root := flags.String("root", ".", "project root")
		port := flags.Int("port", 7331, "loopback port")
		runtimeID := flags.String("runtime-id", "", "runtime identity")
		if err := flags.Parse(args); err != nil {
			return err
		}
		return app.Serve(*root, *port, *runtimeID)
	case "status":
		p, err := project.Detect(ctx, ".")
		if err != nil {
			return err
		}
		record, live, err := app.RuntimeStatus(ctx, p)
		if err != nil {
			return err
		}
		fmt.Printf("Project: %s\nState: %s\n", p.Name, record.State)
		if live {
			fmt.Printf("Origin: %s\n", record.Origin)
		}
		checks, _ := app.Doctor(ctx, p)
		fmt.Print(diagnostics.Render(checks))
		return nil
	case "doctor":
		p, err := project.Detect(ctx, ".")
		if err != nil {
			return err
		}
		checks, ok := app.Doctor(ctx, p)
		fmt.Print(diagnostics.Render(checks))
		if !ok {
			return fmt.Errorf("one or more required checks failed")
		}
		return nil
	case "open":
		p, err := project.Detect(ctx, ".")
		if err != nil {
			return err
		}
		return app.OpenProject(ctx, p)
	case "stop":
		p, err := project.Detect(ctx, ".")
		if err != nil {
			return err
		}
		if err = app.StopProject(ctx, p); err != nil {
			return err
		}
		fmt.Println("Craftlane stopped")
		return nil
	case "restart":
		p, err := project.Detect(ctx, ".")
		if err != nil {
			return err
		}
		if err = app.StopProject(ctx, p); err != nil {
			return err
		}
		record, err := app.StartDetached(ctx, ".", 7331, false, false)
		if err != nil {
			return err
		}
		fmt.Println("Craftlane restarted at", record.Origin)
		return nil
	default:
		return fmt.Errorf("unknown command %q", command)
	}
}
