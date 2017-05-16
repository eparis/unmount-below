/*
Copyright 2017 Eric Paris

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/eparis/unmount-below/mounts"

	"github.com/spf13/cobra"
)

const (
	dryRunVar = "dry-run"
)

var (
	_      = syscall.Unmount
	dryRun = true
)

func run(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("unable to determine path")
	}
	path := args[0]

	mnts, err := mounts.MountsUnder(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No mounts found under: %s\n", path)
			return nil
		}
		return err
	}
	for _, mnt := range mnts {
		t := mnt.Target()
		if !dryRun {
			fmt.Printf("Unmounting: %s\n", t)
			//err = syscall.Unmount(longestBelow, 0)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("Would unmount, but --dry-run=true: %s\n", t)
		}
	}
	return nil
}

func main() {
	root := &cobra.Command{
		Use:   fmt.Sprintf("%s [pathname]", filepath.Base(os.Args[0])),
		Short: "A program to convert blunderbuss.yaml",
		RunE:  run,
	}
	root.Flags().BoolVar(&dryRun, "dry-run", dryRun, "If we should do an unmount of print what we would unmount")
	root.Execute()
}
