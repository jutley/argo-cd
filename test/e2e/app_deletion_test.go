package e2e

import (
	"testing"

	. "github.com/argoproj/argo-cd/errors"
	. "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	. "github.com/argoproj/argo-cd/test/e2e/fixture"
	. "github.com/argoproj/argo-cd/test/e2e/fixture/app"
	"github.com/argoproj/argo-cd/test/fixture/test"
)

// when a app gets stuck in sync, and we try to delete it, it won't delete, instead we must then terminate it
// and deletion will then just happen
func TestDeletingAppStuckInSync(t *testing.T) {

	test.Flaky(t)

	Given(t).
		Path("hook").
		When().
		PatchFile("hook.yaml", `[{"op": "replace", "path": "/spec/containers/0/command", "value": ["sleep", "999"]}]`).
		Create().
		Sync().
		Then().
		// stuck in running state
		Expect(OperationPhaseIs(OperationRunning)).
		When().
		Delete(true).
		Then().
		// delete is ignored, still stuck in running state
		Expect(OperationPhaseIs(OperationRunning)).
		When().
		And(func() {
			// force delete the resource
			FailOnErr(Run("", "kubectl", "-n", DeploymentNamespace(), "delete", "pod", "hook", "--force", "--grace-period", "0"))
		}).
		TerminateOp().
		Then().
		// delete is successful
		Expect(DoesNotExist())
}
