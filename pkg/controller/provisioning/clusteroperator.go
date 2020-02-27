package provisioning

import (
	"context"

	configv1 "github.com/openshift/api/config/v1"
	operatorv1helpers "github.com/openshift/library-go/pkg/operator/v1helpers"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func syncClusterOperator(c client.Client, targetNamespace string, version string, done, disabled bool) error {
	name := types.NamespacedName{Name: "baremetal"}

	co := &configv1.ClusterOperator{}
	err := c.Get(context.TODO(), name, co)
	if err != nil {
		// Not found - create the resource
		if errors.IsNotFound(err) {
			co = &configv1.ClusterOperator{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterOperator",
					APIVersion: "config.openshift.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: name.Name,
				},
				Status: configv1.ClusterOperatorStatus{
					RelatedObjects: []configv1.ObjectReference{
						{
							Group:    "",
							Resource: "namespaces",
							Name:     targetNamespace,
						},
					},
				},
			}
			err = c.Create(context.TODO(), co)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	prevConditions := co.Status.Conditions
	co.Status.Conditions = updateConditions(prevConditions, done, disabled)

	operatorv1helpers.SetOperandVersion(&co.Status.Versions, configv1.OperandVersion{Name: "operator", Version: version})

	if conditionsEquals(co.Status.Conditions, prevConditions) {
		return nil
	}

	err = c.Update(context.TODO(), co)
	if err != nil {
		return err
	}

	return nil
}

// OperatorDisabled reports when the primary function of the operator has been disabled.
const OperatorDisabled configv1.ClusterStatusConditionType = "Disabled"

func updateConditions(conditions []configv1.ClusterOperatorStatusCondition, done, disabled bool) []configv1.ClusterOperatorStatusCondition {
	// FIXME: actually implement the expected semantics of these conditions
	conditions = []configv1.ClusterOperatorStatusCondition{
		{
			Type:   configv1.OperatorAvailable,
			Status: configv1.ConditionFalse,
		}, {
			Type:   configv1.OperatorProgressing,
			Status: configv1.ConditionFalse,
		}, {
			Type:   configv1.OperatorDegraded,
			Status: configv1.ConditionFalse,
		}, {
			Type:   OperatorDisabled,
			Status: configv1.ConditionFalse,
		},
	}
	if done {
		conditions[0].Status = configv1.ConditionTrue
	} else {
		conditions[1].Status = configv1.ConditionTrue
	}
	if disabled {
		conditions[3].Status = configv1.ConditionTrue
	}
	return conditions
}

func conditionsEquals(conditions []configv1.ClusterOperatorStatusCondition, prev []configv1.ClusterOperatorStatusCondition) bool {
	// FIXME: only update when something has changed
	return false
}
