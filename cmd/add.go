// Copyright 2022 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new malware family to the malware souk database",
	Long: `Generates markdown source code for a new malware family for
saferwall's malware souk database`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add command")
	},
}
