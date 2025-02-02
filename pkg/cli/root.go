// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	"fmt"
	"github.com/onosproject/onos-cli/pkg/mlb"

	"github.com/onosproject/onos-cli/pkg/config"
	"github.com/onosproject/onos-cli/pkg/e2sub"
	"github.com/onosproject/onos-cli/pkg/e2t"
	"github.com/onosproject/onos-cli/pkg/kpimon"
	"github.com/onosproject/onos-cli/pkg/mho"
	"github.com/onosproject/onos-cli/pkg/modelregistry"
	"github.com/onosproject/onos-cli/pkg/pci"
	"github.com/onosproject/onos-cli/pkg/topo"
	"github.com/onosproject/onos-cli/pkg/uenib"

	// Needed to keep ran-sim happy for the mo
	_ "github.com/onosproject/onos-lib-go/pkg/cli"

	"os"

	"github.com/onosproject/onos-cli/pkg/ransim"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// Execute runs the root command and any sub-commands.
func Execute() {
	rootCmd := GetRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// GenerateCliDocs generate markdown files for onos-cli commands
func GenerateCliDocs() {
	cmd := GetRootCommand()
	err := doc.GenMarkdownTree(cmd, "docs/cli")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// GetRootCommand returns the root onos command
func GetRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                    "onos",
		Short:                  "ONOS command line client",
		BashCompletionFunction: getBashCompletions(),
		SilenceUsage:           true,
		SilenceErrors:          true,
	}
	cmd.AddCommand(topo.GetCommand())
	cmd.AddCommand(uenib.GetCommand())
	cmd.AddCommand(config.GetCommand())
	cmd.AddCommand(e2t.GetCommand())
	cmd.AddCommand(e2sub.GetCommand())
	cmd.AddCommand(ransim.GetCommand())
	cmd.AddCommand(kpimon.GetCommand())
	cmd.AddCommand(mho.GetCommand())
	cmd.AddCommand(pci.GetCommand())
	cmd.AddCommand(mlb.GetCommand())
	cmd.AddCommand(modelregistry.GetCommand())
	cmd.AddCommand()

	cmd.AddCommand(getCompletionCommand())

	return cmd
}
