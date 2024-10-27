/*******************************************************************************
*
* Copyright 2019 SAP SE
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package cmd

import (
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/romnn/helm-outdated/pkg/helm"
	"github.com/spf13/cobra"
)

type updateCmd struct {
	chartPath               string
	maxColumnWidth          uint
	indent                  int
	isIncrementChartVersion bool
	dependencyFilter        *helm.Filter

	// **Experimental**
	// isAutoUpdate updates the dependencies, increments version of the chart with the dependency and (git) commits the changes.
	isAutoUpdate,
	isOnlyPullRequest bool
	authorName,
	authorEmail string
}

var updateLongUsage = `
Update outdated dependencies of a given chart to their latest version.

Examples:
  # Update dependencies of the given chart.
  $ helm outdated update <chartPath>

	# Only update specific dependencies of the given chart.
	$ helm outdated update <chartPath> --dependencies kube-state-metrics,prometheus-operator
`

func newUpdateOutdatedDependenciesCmd() *cobra.Command {
	u := &updateCmd{
		dependencyFilter: &helm.Filter{},
		maxColumnWidth:   60,
	}

	cmd := &cobra.Command{
		Use:          "update",
		Long:         updateLongUsage,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if debug, err := cmd.Flags().GetBool("debug"); err == nil {
				if debug == true {
					log.SetLevel(log.DebugLevel)
				} else {
					log.SetLevel(log.InfoLevel)
				}
			}

			if maxColumnWidth, err := cmd.Flags().GetUint("max-column-width"); err == nil {
				u.maxColumnWidth = maxColumnWidth
			}

			if repositories, err := cmd.Flags().GetStringSlice("repositories"); err == nil {
				u.dependencyFilter.Repositories = repositories
			}

			if deps, err := cmd.Flags().GetStringSlice("dependencies"); err == nil {
				u.dependencyFilter.DependencyNames = deps
			}

			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			u.chartPath = path

			return fmt.Errorf("update is not yet implemented")
		},
	}

	addCommonFlags(cmd)
	cmd.Flags().
		BoolVarP(&u.isIncrementChartVersion, "increment-chart-version", "", false, "Increment the version of the Helm chart if requirements are updated.")
	cmd.Flags().
		IntVarP(&u.indent, "indent", "", 4, "Indent to use when writing the requirements.yaml .")

	// **Experimental** Update dependencies of the given chart, commit and push to upstream using git.
	cmd.Flags().
		BoolVar(&u.isAutoUpdate, "auto-update", false, "**Experimental** Update dependencies of the given chart, commit and push to upstream using git.")
	cmd.Flags().
		StringVar(&u.authorName, "author-name", "", "The name of the author and committer to be used when auto update is enabled.")
	cmd.Flags().
		StringVar(&u.authorEmail, "author-email", "", "The email of the author and committer to be used when auto update is enabled.")
	cmd.Flags().
		BoolVar(&u.isOnlyPullRequest, "only-pull-requests", false, "Only use pull requests. Do not commit minor changes to master branch.")

	return cmd
}
