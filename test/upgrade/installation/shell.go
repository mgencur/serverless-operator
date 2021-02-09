/*
Copyright 2020 The Knative Authors

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

package installation

import (
	"knative.dev/hack/shell"
	pkgupgrade "knative.dev/pkg/test/upgrade"
)

// UpgradeServerless upgrades the product to the latest version/CSV.
func UpgradeServerless() pkgupgrade.Operation {
	return install("UpgradeServerless", "upgrade_serverless")
}

// UpgradeCluster upgrades the OCP cluster.
func UpgradeCluster() pkgupgrade.Operation {
	return install("UpgradeCluster", "upgrade_ocp_cluster")
}

func install(installName, shellFunc string) pkgupgrade.Operation {
	return pkgupgrade.NewOperation(installName, func(c pkgupgrade.Context) {
		c.Log.Info("Running shell function: ", shellFunc)
		if err := callShellFunction(shellFunc); err != nil {
			c.T.Error(err)
		}
	})
}

func callShellFunction(funcName string) error {
	loc, err := shell.NewProjectLocation("../../..")
	if err != nil {
		return err
	}
	exec := shell.NewExecutor(shell.ExecutorConfig{
		ProjectLocation: loc,
	})
	fn := shell.Function{
		Script: shell.Script{
			Label:      funcName,
			ScriptPath: "test/lib.bash",
		},
		FunctionName: funcName,
	}
	return exec.RunFunction(fn)
}
