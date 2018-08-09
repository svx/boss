package main

import (
	"fmt"
	"os"

	"github.com/crosbymichael/boss/config"
	"github.com/crosbymichael/boss/system"
	"github.com/crosbymichael/boss/systemd"
	"github.com/urfave/cli"
)

var cfg *system.Config

func main() {
	app := cli.NewApp()
	app.Name = "boss"
	app.Version = "7"
	app.Usage = "container tools for my setup"
	app.Description = `

                    ___           ___           ___     
     _____         /\  \         /\__\         /\__\    
    /::\  \       /::\  \       /:/ _/_       /:/ _/_   
   /:/\:\  \     /:/\:\  \     /:/ /\  \     /:/ /\  \  
  /:/ /::\__\   /:/  \:\  \   /:/ /::\  \   /:/ /::\  \ 
 /:/_/:/\:|__| /:/__/ \:\__\ /:/_/:/\:\__\ /:/_/:/\:\__\
 \:\/:/ /:/  / \:\  \ /:/  / \:\/:/ /:/  / \:\/:/ /:/  /
  \::/_/:/  /   \:\  /:/  /   \::/ /:/  /   \::/ /:/  / 
   \:\/:/  /     \:\/:/  /     \/_/:/  /     \/_/:/  /  
    \::/  /       \::/  /        /:/  /        /:/  /   
     \/__/         \/__/         \/__/         \/__/    

run containers like a boss`
	app.Commands = []cli.Command{
		buildCommand,
		createCommand,
		deleteCommand,
		initCommand,
		killCommand,
		listCommand,
		rollbackCommand,
		startCommand,
		stopCommand,
		systemdCommand,
		upgradeCommand,
	}
	app.Before = func(clix *cli.Context) error {
		c, err := system.Load("/etc/boss/boss.toml")
		if err != nil {
			return err
		}
		cfg = c
		if os.Geteuid() != 0 {
			return nil
		}
		if err := os.MkdirAll(config.Root, 0711); err != nil {
			return err
		}
		return systemd.Install()
	}
	app.After = func(clix *cli.Context) error {
		if cfg == nil {
			return nil
		}
		return cfg.Close()
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func ReadyBefore(clix *cli.Context) error {
	return system.Ready(cfg)
}
