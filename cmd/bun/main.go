package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/uptrace/bun-starter-kit/bunapp"
	"github.com/uptrace/bun-starter-kit/cmd/bun/migrations"
	_ "github.com/uptrace/bun-starter-kit/example"
	"github.com/uptrace/bun-starter-kit/httputil"
	"github.com/uptrace/bun/migrate"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "bun",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "dev",
				Usage: "environment",
			},
		},
		Commands: []*cli.Command{
			serverCommand,
			newDBCommand(migrations.Migrations),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var serverCommand = &cli.Command{
	Name:  "runserver",
	Usage: "start API server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Value: "localhost:8000",
			Usage: "serve address",
		},
	},
	Action: func(c *cli.Context) error {
		ctx, app, err := bunapp.Start(c.Context, "api", c.String("env"))
		if err != nil {
			return err
		}
		defer app.Stop()

		var handler http.Handler
		handler = app.Router()
		handler = httputil.ExitOnPanicHandler{Next: handler}

		srv := &http.Server{
			Addr:         c.String("addr"),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
			Handler:      handler,
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil && !isServerClosed(err) {
				log.Printf("ListenAndServe failed: %s", err)
			}
		}()

		fmt.Printf("listening on http://%s\n", srv.Addr)
		fmt.Println(bunapp.WaitExitSignal())

		return srv.Shutdown(ctx)
	},
}

func newDBCommand(migrations *migrate.Migrations) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					ctx, app, err := bunapp.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Init(ctx)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					ctx, app, err := bunapp.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					group, err := migrator.Migrate(ctx)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no new migrations to run\n")
						return nil
					}

					fmt.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					ctx, app, err := bunapp.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					group, err := migrator.Rollback(ctx)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no groups to roll back\n")
						return nil
					}

					fmt.Printf("rolled back %s\n", group)
					return nil
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := bunapp.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Lock(ctx)
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := bunapp.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Unlock(ctx)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					ctx, app, err := bunapp.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					mf, err := migrator.CreateGo(ctx, c.Args().Get(0))
					if err != nil {
						return err
					}
					fmt.Printf("created migration %s (%s)\n", mf.FileName, mf.FilePath)
					return nil
				},
			},
			{
				Name:  "create_sql",
				Usage: "create SQL migration",
				Action: func(c *cli.Context) error {
					ctx, app, err := bunapp.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					mf, err := migrator.CreateSQL(ctx, c.Args().Get(0))
					if err != nil {
						return err
					}
					fmt.Printf("created migration %s (%s)\n", mf.FileName, mf.FilePath)
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					ctx, app, err := bunapp.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					status, err := migrator.Status(ctx)
					if err != nil {
						return err
					}

					fmt.Printf("migrations: %s\n", status.Migrations)
					fmt.Printf("new migrations: %s\n", status.NewMigrations)
					fmt.Printf("last group: %s\n", status.LastGroup)
					return nil
				},
			},
			{
				Name:  "mark_completed",
				Usage: "mark migrations as completed without actually running them",
				Action: func(c *cli.Context) error {
					ctx, app, err := bunapp.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					group, err := migrator.MarkCompleted(ctx)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no new migrations to mark as completed\n")
						return nil
					}

					fmt.Printf("marked as completed %s\n", group)
					return nil
				},
			},
		},
	}
}

func isServerClosed(err error) bool {
	return err.Error() == "http: Server closed"
}
