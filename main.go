package main

import (
	"fmt"
	"github.com/ssh-alive/utils"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"os/exec"
	"runtime"
)

var Version = "v0.0.1-devel"

func main() {
	var session *ssh.Session
	connected := false
	app := cli.NewApp()
	app.Name = "SSH-alive"
	app.Version = Version
	app.Author = "baojiweicn github.com:baojiweicn/ssh-alive.git"
	app.Copyright = "(c)baojiweicn"
	app.Usage = "SSH-alive for " + runtime.GOOS + "/" + runtime.GOARCH
	app.Description = ``

	app.Action = func(c *cli.Context) (err error) {
		user := c.String("user")
		password := c.String("password")
		ip := c.String("host")
		port := c.Int("port")
		ciphers := []string{}
		key := "/Users/baojiwei/.ssh/id_rsa" // 先固定的pub_key
		session, err = utils.Connect(user, password, ip, key, port, ciphers)
		if err != nil {
			fmt.Printf("警告: 读取历史命令文件错误, %s\n", err)
			return
		}
		err = session.Shell()
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		session.Stdin = os.Stdin
		if err != nil {
			fmt.Printf("警告: 读取历史命令文件错误, %s\n", err)
		}
		session.Wait()
		return
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "user, u",
			Value: "root",
			Usage: "ssh user",
		},
		cli.StringFlag{
			Name:  "host, H",
			Value: "localhost",
			Usage: "ssh host",
		},
		cli.IntFlag{
			Name:  "port, p",
			Value: 22,
			Usage: "ssh port",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:     "run",
			Usage:    "执行系统命令",
			Category: "其他",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					cli.ShowCommandHelp(c, c.Command.Name)
					return nil
				}

				cmd := exec.Command(c.Args().First(), c.Args().Tail()...)
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				cmd.Stderr = os.Stderr

				err := cmd.Run()
				if err != nil {
					fmt.Println(err)
				}

				return nil
			},
		},
		{
			Name:     "ssh",
			Usage:    "执行远程连接",
			Category: "ssh",
			Action: func(c *cli.Context) (err error) {
				if c.NArg() == 0 {
					cli.ShowCommandHelp(c, c.Command.Name)
					return nil
				}
				fmt.Println(c.Args())
				user := c.String("user")
				password := c.String("password")
				ip := c.String("host")
				port := c.Int("port")
				ciphers := []string{}
				key := "/Users/baojiwei/.ssh/id_rsa" // 先固定的pub_key
				session, err = utils.Connect(user, password, ip, key, port, ciphers)
				if err != nil {
					fmt.Printf("警告: 读取历史命令文件错误, %s\n", err)
					return
				}
				session.Stdout = os.Stdout
				session.Stderr = os.Stderr
				session.Stdin = os.Stdin
				err = session.Shell()
				if err != nil {
					fmt.Printf("警告: 读取历史命令文件错误, %s\n", err)
					return
				}
				connected = true
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "user, u",
					Value: "root",
					Usage: "ssh user",
				},
				cli.StringFlag{
					Name:  "host, H",
					Value: "localhost",
					Usage: "ssh host",
				},
				cli.IntFlag{
					Name:  "port, p",
					Value: 22,
					Usage: "ssh port",
				},
			},
		},
		{
			Name:    "quit",
			Aliases: []string{"exit"},
			Usage:   "退出程序",
			Action: func(c *cli.Context) error {
				return cli.NewExitError("", 0)
			},
			Hidden:   true,
			HideHelp: true,
		},
	}

	// if c.Int("port") == 8000 {
	// 	return cli.NewExitError("invalid port", 88)
	// }

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}


