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

var (
	_ = syscall.Unmount
)

func run(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("unable to determine path")
	}
	path := args[0]

	for {
		longestBelow, err := mounts.LongestMountUnder(path)
		if err != nil {
			if os.IsNotExist(err) {
				break
			}
			return err
		}
		fmt.Printf("Unmounting: %q\n", longestBelow)
		// err := syscall.Unmount(longestBelow, 0)
		if err != nil {
			return err
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
	root.Execute()
}
