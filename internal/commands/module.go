package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/charmbracelet/ssh"
	"github.com/erandelax/go-dossh/internal/configuration"
	"github.com/erandelax/go-dossh/internal/utils"
	"github.com/spf13/cobra"
)

func NewRoot(userCfg configuration.UserConfig, user string, s ssh.Session) *cobra.Command {
	cmd := &cobra.Command{}

	cmd.AddCommand(&cobra.Command{
		Use:   "ps",
		Short: "List available containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			containers, err := utils.GetRunningContainerNames()
			if err != nil {
				cmd.PrintErrln(err)
			}
			availableRows := []string{}
			notAvailableRows := []string{}
			for _, containerName := range containers {
				userContainerCfg, ok := userCfg.Containers[containerName]
				if ok {
					availableRows = append(availableRows, fmt.Sprintf("- %s %s", containerName, strings.Join(userContainerCfg, "|")))
				} else {
					notAvailableRows = append(notAvailableRows, fmt.Sprintf("- %s", containerName))
				}
			}
			if len(notAvailableRows) > 0 {
				cmd.Println("\nNon-accessible containers:")
				for _, row := range notAvailableRows {
					cmd.Println(row)
				}
			}
			if len(availableRows) > 0 {
				cmd.Println("\nAccessible containers:")
				for _, row := range availableRows {
					cmd.Println(row)
				}
			}
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "restart [container]",
		Args:  cobra.ExactArgs(1),
		Short: "Runs 'docker restart [container]'",
		RunE: func(cmd *cobra.Command, args []string) error {
			containerName := args[0]
			execIt := "restart"
			userContainerCfg, ok := userCfg.Containers[containerName]
			if !ok {
				cmd.Println("You don't have access to this container")
				return nil
			}
			if !utils.SliceContainsString(userContainerCfg, execIt) {
				cmd.Println("You are not allowed to execute", execIt, "for this container")
				return nil
			}

			parts := []string{}
			parts = append(parts, "docker", execIt, containerName)

			var c *exec.Cmd
			c = exec.Command(parts[0], parts[1:]...)
			c.Env = os.Environ()
			out, err := c.Output()
			if err != nil {
				cmd.Println(err.Error())
			}
			cmd.Println(strings.Trim(string(out), " \t\r\n"))

			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "logs [container]",
		Args:  cobra.ExactArgs(1),
		Short: "Runs 'docker logs [container] --follow'",
		RunE: func(cmd *cobra.Command, args []string) error {
			containerName := args[0]
			execIt := "logs"
			userContainerCfg, ok := userCfg.Containers[containerName]
			if !ok {
				cmd.Println("You don't have access to this container")
				return nil
			}
			if !utils.SliceContainsString(userContainerCfg, execIt) {
				cmd.Println("You are not allowed to execute", execIt, "for this container")
				return nil
			}

			parts := []string{}
			if runtime.GOOS == "windows" {
				parts = append(parts, "winpty", "-Xallow-non-tty")
			}
			parts = append(parts, "docker", execIt, containerName, "--follow")

			var c *exec.Cmd
			c = exec.Command(parts[0], parts[1:]...)
			c.Env = os.Environ()
			c.Stdin = cmd.InOrStdin()
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.ErrOrStderr()
			err := c.Run()
			if err != nil {
				cmd.Println(err.Error())
			}

			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "exec [container_name] [command=sh]",
		Short: "Runs 'docker exec -it [container_name] [command=sh]'",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			containerName := args[0]
			execIt := args[1]

			if execIt == "restart" || execIt == "logs" {
				cmd.Println("You are not allowed to execute", execIt, "-inside- of this container")
				return nil
			}
			userContainerCfg, ok := userCfg.Containers[containerName]
			if !ok {
				cmd.Println("You don't have access to this container")
				return nil
			}
			if !utils.SliceContainsString(userContainerCfg, execIt) {
				cmd.Println("You are not allowed to execute", execIt, "for this container")
				return nil
			}

			execAs := ""
			containerCfg, ok := configuration.Get().Containers[containerName]
			if ok && nil != containerCfg.As {
				execAs = *containerCfg.As
			}

			parts := []string{}
			if runtime.GOOS == "windows" {
				parts = append(parts, "winpty", "-Xallow-non-tty")
			}
			parts = append(parts, "docker", "exec")
			if execAs != "" {
				parts = append(parts, "-u", execAs)
			}
			parts = append(parts, "-it", containerName, execIt)

			var c *exec.Cmd
			c = exec.Command(parts[0], parts[1:]...)
			c.Env = os.Environ()
			c.Stdin = cmd.InOrStdin()
			c.Stdout = cmd.OutOrStdout()
			c.Stderr = cmd.ErrOrStderr()
			err := c.Run()
			if err != nil {
				cmd.Println(err.Error())
			}

			return nil
		},
	})

	return cmd
}
