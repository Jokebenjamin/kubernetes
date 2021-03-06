/*
Copyright 2016 The Kubernetes Authors.

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

package v1beta2

import (
	"fmt"
	"strconv"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/api"
	k8s_api_v1 "k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func addConversionFuncs(scheme *runtime.Scheme) error {
	// Add non-generated conversion functions to handle the *int32 -> int32
	// conversion. A pointer is useful in the versioned type so we can default
	// it, but a plain int32 is more convenient in the internal type. These
	// functions are the same as the autogenerated ones in every other way.
	err := scheme.AddConversionFuncs(
		Convert_v1beta2_StatefulSetSpec_To_apps_StatefulSetSpec,
		Convert_apps_StatefulSetSpec_To_v1beta2_StatefulSetSpec,
		Convert_v1beta2_StatefulSetUpdateStrategy_To_apps_StatefulSetUpdateStrategy,
		Convert_apps_StatefulSetUpdateStrategy_To_v1beta2_StatefulSetUpdateStrategy,
		Convert_extensions_RollingUpdateDaemonSet_To_v1beta2_RollingUpdateDaemonSet,
		Convert_v1beta2_RollingUpdateDaemonSet_To_extensions_RollingUpdateDaemonSet,
		Convert_v1beta2_StatefulSetStatus_To_apps_StatefulSetStatus,
		Convert_apps_StatefulSetStatus_To_v1beta2_StatefulSetStatus,
		Convert_v1beta2_Deployment_To_extensions_Deployment,
		Convert_extensions_Deployment_To_v1beta2_Deployment,
		// extensions
		// TODO: below conversions should be dropped in favor of auto-generated
		// ones, see https://github.com/kubernetes/kubernetes/issues/39865
		Convert_v1beta2_ScaleStatus_To_extensions_ScaleStatus,
		Convert_extensions_ScaleStatus_To_v1beta2_ScaleStatus,
		Convert_v1beta2_DeploymentSpec_To_extensions_DeploymentSpec,
		Convert_extensions_DeploymentSpec_To_v1beta2_DeploymentSpec,
		Convert_v1beta2_DeploymentStrategy_To_extensions_DeploymentStrategy,
		Convert_extensions_DeploymentStrategy_To_v1beta2_DeploymentStrategy,
		Convert_v1beta2_RollingUpdateDeployment_To_extensions_RollingUpdateDeployment,
		Convert_extensions_RollingUpdateDeployment_To_v1beta2_RollingUpdateDeployment,
		Convert_extensions_ReplicaSetSpec_To_v1beta2_ReplicaSetSpec,
		Convert_v1beta2_ReplicaSetSpec_To_extensions_ReplicaSetSpec,
	)
	if err != nil {
		return err
	}

	// Add field label conversions for kinds having selectable nothing but ObjectMeta fields.
	err = scheme.AddFieldLabelConversionFunc("apps/v1beta2", "StatefulSet",
		func(label, value string) (string, string, error) {
			switch label {
			case "metadata.name", "metadata.namespace", "status.successful":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported for appsv1beta2.StatefulSet: %s", label)
			}
		})
	if err != nil {
		return err
	}
	err = scheme.AddFieldLabelConversionFunc("apps/v1beta2", "Deployment",
		func(label, value string) (string, string, error) {
			switch label {
			case "metadata.name", "metadata.namespace":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label %q not supported for appsv1beta2.Deployment", label)
			}
		})
	if err != nil {
		return err
	}
	err = scheme.AddFieldLabelConversionFunc("apps/v1beta2", "ReplicaSet",
		func(label, value string) (string, string, error) {
			switch label {
			case "metadata.name", "metadata.namespace":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label %q not supported for appsv1beta2.ReplicaSet", label)
			}
		})
	if err != nil {
		return err
	}

	return nil
}

func Convert_extensions_RollingUpdateDaemonSet_To_v1beta2_RollingUpdateDaemonSet(in *extensions.RollingUpdateDaemonSet, out *appsv1beta2.RollingUpdateDaemonSet, s conversion.Scope) error {
	if out.MaxUnavailable == nil {
		out.MaxUnavailable = &intstr.IntOrString{}
	}
	if err := s.Convert(&in.MaxUnavailable, out.MaxUnavailable, 0); err != nil {
		return err
	}
	return nil
}

func Convert_v1beta2_RollingUpdateDaemonSet_To_extensions_RollingUpdateDaemonSet(in *appsv1beta2.RollingUpdateDaemonSet, out *extensions.RollingUpdateDaemonSet, s conversion.Scope) error {
	if err := s.Convert(in.MaxUnavailable, &out.MaxUnavailable, 0); err != nil {
		return err
	}
	return nil
}

func Convert_v1beta2_StatefulSetSpec_To_apps_StatefulSetSpec(in *appsv1beta2.StatefulSetSpec, out *apps.StatefulSetSpec, s conversion.Scope) error {
	if in.Replicas != nil {
		out.Replicas = *in.Replicas
	}
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(metav1.LabelSelector)
		if err := s.Convert(*in, *out, 0); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	if err := s.Convert(&in.Template, &out.Template, 0); err != nil {
		return err
	}
	if in.VolumeClaimTemplates != nil {
		in, out := &in.VolumeClaimTemplates, &out.VolumeClaimTemplates
		*out = make([]api.PersistentVolumeClaim, len(*in))
		for i := range *in {
			if err := s.Convert(&(*in)[i], &(*out)[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.VolumeClaimTemplates = nil
	}
	if err := Convert_v1beta2_StatefulSetUpdateStrategy_To_apps_StatefulSetUpdateStrategy(&in.UpdateStrategy, &out.UpdateStrategy, s); err != nil {
		return err
	}
	if in.RevisionHistoryLimit != nil {
		out.RevisionHistoryLimit = new(int32)
		*out.RevisionHistoryLimit = *in.RevisionHistoryLimit
	} else {
		out.RevisionHistoryLimit = nil
	}
	out.ServiceName = in.ServiceName
	out.PodManagementPolicy = apps.PodManagementPolicyType(in.PodManagementPolicy)
	return nil
}

func Convert_apps_StatefulSetSpec_To_v1beta2_StatefulSetSpec(in *apps.StatefulSetSpec, out *appsv1beta2.StatefulSetSpec, s conversion.Scope) error {
	out.Replicas = new(int32)
	*out.Replicas = in.Replicas
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(metav1.LabelSelector)
		if err := s.Convert(*in, *out, 0); err != nil {
			return err
		}
	} else {
		out.Selector = nil
	}
	if err := s.Convert(&in.Template, &out.Template, 0); err != nil {
		return err
	}
	if in.VolumeClaimTemplates != nil {
		in, out := &in.VolumeClaimTemplates, &out.VolumeClaimTemplates
		*out = make([]v1.PersistentVolumeClaim, len(*in))
		for i := range *in {
			if err := s.Convert(&(*in)[i], &(*out)[i], 0); err != nil {
				return err
			}
		}
	} else {
		out.VolumeClaimTemplates = nil
	}
	if in.RevisionHistoryLimit != nil {
		out.RevisionHistoryLimit = new(int32)
		*out.RevisionHistoryLimit = *in.RevisionHistoryLimit
	} else {
		out.RevisionHistoryLimit = nil
	}
	out.ServiceName = in.ServiceName
	out.PodManagementPolicy = appsv1beta2.PodManagementPolicyType(in.PodManagementPolicy)
	if err := Convert_apps_StatefulSetUpdateStrategy_To_v1beta2_StatefulSetUpdateStrategy(&in.UpdateStrategy, &out.UpdateStrategy, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1beta2_StatefulSetUpdateStrategy_To_apps_StatefulSetUpdateStrategy(in *appsv1beta2.StatefulSetUpdateStrategy, out *apps.StatefulSetUpdateStrategy, s conversion.Scope) error {
	out.Type = apps.StatefulSetUpdateStrategyType(in.Type)
	if in.RollingUpdate != nil {
		out.RollingUpdate = new(apps.RollingUpdateStatefulSetStrategy)
		out.RollingUpdate.Partition = *in.RollingUpdate.Partition
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func Convert_apps_StatefulSetUpdateStrategy_To_v1beta2_StatefulSetUpdateStrategy(in *apps.StatefulSetUpdateStrategy, out *appsv1beta2.StatefulSetUpdateStrategy, s conversion.Scope) error {
	out.Type = appsv1beta2.StatefulSetUpdateStrategyType(in.Type)
	if in.RollingUpdate != nil {
		out.RollingUpdate = new(appsv1beta2.RollingUpdateStatefulSetStrategy)
		out.RollingUpdate.Partition = new(int32)
		*out.RollingUpdate.Partition = in.RollingUpdate.Partition
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func Convert_v1beta2_StatefulSetStatus_To_apps_StatefulSetStatus(in *appsv1beta2.StatefulSetStatus, out *apps.StatefulSetStatus, s conversion.Scope) error {
	out.ObservedGeneration = new(int64)
	*out.ObservedGeneration = in.ObservedGeneration
	out.Replicas = in.Replicas
	out.ReadyReplicas = in.ReadyReplicas
	out.CurrentReplicas = in.CurrentReplicas
	out.UpdatedReplicas = in.UpdatedReplicas
	out.CurrentRevision = in.CurrentRevision
	out.UpdateRevision = in.UpdateRevision
	if in.CollisionCount != nil {
		out.CollisionCount = new(int64)
		*out.CollisionCount = *in.CollisionCount
	}
	return nil
}

func Convert_apps_StatefulSetStatus_To_v1beta2_StatefulSetStatus(in *apps.StatefulSetStatus, out *appsv1beta2.StatefulSetStatus, s conversion.Scope) error {
	if in.ObservedGeneration != nil {
		out.ObservedGeneration = *in.ObservedGeneration
	}
	out.Replicas = in.Replicas
	out.ReadyReplicas = in.ReadyReplicas
	out.CurrentReplicas = in.CurrentReplicas
	out.UpdatedReplicas = in.UpdatedReplicas
	out.CurrentRevision = in.CurrentRevision
	out.UpdateRevision = in.UpdateRevision
	if in.CollisionCount != nil {
		out.CollisionCount = new(int64)
		*out.CollisionCount = *in.CollisionCount
	}
	return nil
}

func Convert_extensions_ScaleStatus_To_v1beta2_ScaleStatus(in *extensions.ScaleStatus, out *appsv1beta2.ScaleStatus, s conversion.Scope) error {
	out.Replicas = int32(in.Replicas)

	out.Selector = nil
	out.TargetSelector = ""
	if in.Selector != nil {
		if in.Selector.MatchExpressions == nil || len(in.Selector.MatchExpressions) == 0 {
			out.Selector = in.Selector.MatchLabels
		}

		selector, err := metav1.LabelSelectorAsSelector(in.Selector)
		if err != nil {
			return fmt.Errorf("invalid label selector: %v", err)
		}
		out.TargetSelector = selector.String()
	}
	return nil
}

func Convert_v1beta2_ScaleStatus_To_extensions_ScaleStatus(in *appsv1beta2.ScaleStatus, out *extensions.ScaleStatus, s conversion.Scope) error {
	out.Replicas = in.Replicas

	// Normally when 2 fields map to the same internal value we favor the old field, since
	// old clients can't be expected to know about new fields but clients that know about the
	// new field can be expected to know about the old field (though that's not quite true, due
	// to kubectl apply). However, these fields are readonly, so any non-nil value should work.
	if in.TargetSelector != "" {
		labelSelector, err := metav1.ParseToLabelSelector(in.TargetSelector)
		if err != nil {
			out.Selector = nil
			return fmt.Errorf("failed to parse target selector: %v", err)
		}
		out.Selector = labelSelector
	} else if in.Selector != nil {
		out.Selector = new(metav1.LabelSelector)
		selector := make(map[string]string)
		for key, val := range in.Selector {
			selector[key] = val
		}
		out.Selector.MatchLabels = selector
	} else {
		out.Selector = nil
	}
	return nil
}

func Convert_v1beta2_DeploymentSpec_To_extensions_DeploymentSpec(in *appsv1beta2.DeploymentSpec, out *extensions.DeploymentSpec, s conversion.Scope) error {
	if in.Replicas != nil {
		out.Replicas = *in.Replicas
	}
	out.Selector = in.Selector
	if err := k8s_api_v1.Convert_v1_PodTemplateSpec_To_api_PodTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	if err := Convert_v1beta2_DeploymentStrategy_To_extensions_DeploymentStrategy(&in.Strategy, &out.Strategy, s); err != nil {
		return err
	}
	out.RevisionHistoryLimit = in.RevisionHistoryLimit
	out.MinReadySeconds = in.MinReadySeconds
	out.Paused = in.Paused
	if in.ProgressDeadlineSeconds != nil {
		out.ProgressDeadlineSeconds = new(int32)
		*out.ProgressDeadlineSeconds = *in.ProgressDeadlineSeconds
	}
	return nil
}

func Convert_extensions_DeploymentSpec_To_v1beta2_DeploymentSpec(in *extensions.DeploymentSpec, out *appsv1beta2.DeploymentSpec, s conversion.Scope) error {
	out.Replicas = &in.Replicas
	out.Selector = in.Selector
	if err := k8s_api_v1.Convert_api_PodTemplateSpec_To_v1_PodTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	if err := Convert_extensions_DeploymentStrategy_To_v1beta2_DeploymentStrategy(&in.Strategy, &out.Strategy, s); err != nil {
		return err
	}
	if in.RevisionHistoryLimit != nil {
		out.RevisionHistoryLimit = new(int32)
		*out.RevisionHistoryLimit = int32(*in.RevisionHistoryLimit)
	}
	out.MinReadySeconds = int32(in.MinReadySeconds)
	out.Paused = in.Paused
	if in.ProgressDeadlineSeconds != nil {
		out.ProgressDeadlineSeconds = new(int32)
		*out.ProgressDeadlineSeconds = *in.ProgressDeadlineSeconds
	}
	return nil
}

func Convert_extensions_DeploymentStrategy_To_v1beta2_DeploymentStrategy(in *extensions.DeploymentStrategy, out *appsv1beta2.DeploymentStrategy, s conversion.Scope) error {
	out.Type = appsv1beta2.DeploymentStrategyType(in.Type)
	if in.RollingUpdate != nil {
		out.RollingUpdate = new(appsv1beta2.RollingUpdateDeployment)
		if err := Convert_extensions_RollingUpdateDeployment_To_v1beta2_RollingUpdateDeployment(in.RollingUpdate, out.RollingUpdate, s); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func Convert_v1beta2_DeploymentStrategy_To_extensions_DeploymentStrategy(in *appsv1beta2.DeploymentStrategy, out *extensions.DeploymentStrategy, s conversion.Scope) error {
	out.Type = extensions.DeploymentStrategyType(in.Type)
	if in.RollingUpdate != nil {
		out.RollingUpdate = new(extensions.RollingUpdateDeployment)
		if err := Convert_v1beta2_RollingUpdateDeployment_To_extensions_RollingUpdateDeployment(in.RollingUpdate, out.RollingUpdate, s); err != nil {
			return err
		}
	} else {
		out.RollingUpdate = nil
	}
	return nil
}

func Convert_v1beta2_RollingUpdateDeployment_To_extensions_RollingUpdateDeployment(in *appsv1beta2.RollingUpdateDeployment, out *extensions.RollingUpdateDeployment, s conversion.Scope) error {
	if err := s.Convert(in.MaxUnavailable, &out.MaxUnavailable, 0); err != nil {
		return err
	}
	if err := s.Convert(in.MaxSurge, &out.MaxSurge, 0); err != nil {
		return err
	}
	return nil
}

func Convert_extensions_RollingUpdateDeployment_To_v1beta2_RollingUpdateDeployment(in *extensions.RollingUpdateDeployment, out *appsv1beta2.RollingUpdateDeployment, s conversion.Scope) error {
	if out.MaxUnavailable == nil {
		out.MaxUnavailable = &intstr.IntOrString{}
	}
	if err := s.Convert(&in.MaxUnavailable, out.MaxUnavailable, 0); err != nil {
		return err
	}
	if out.MaxSurge == nil {
		out.MaxSurge = &intstr.IntOrString{}
	}
	if err := s.Convert(&in.MaxSurge, out.MaxSurge, 0); err != nil {
		return err
	}
	return nil
}

func Convert_extensions_ReplicaSetSpec_To_v1beta2_ReplicaSetSpec(in *extensions.ReplicaSetSpec, out *appsv1beta2.ReplicaSetSpec, s conversion.Scope) error {
	out.Replicas = new(int32)
	*out.Replicas = int32(in.Replicas)
	out.MinReadySeconds = in.MinReadySeconds
	out.Selector = in.Selector
	if err := k8s_api_v1.Convert_api_PodTemplateSpec_To_v1_PodTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1beta2_Deployment_To_extensions_Deployment(in *appsv1beta2.Deployment, out *extensions.Deployment, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1beta2_DeploymentSpec_To_extensions_DeploymentSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}

	// Copy annotation to deprecated rollbackTo field for roundtrip
	// TODO: remove this conversion after we delete extensions/v1beta1 and apps/v1beta1 Deployment
	if revision, _ := in.Annotations[appsv1beta2.DeprecatedRollbackTo]; revision != "" {
		if revision64, err := strconv.ParseInt(revision, 10, 64); err != nil {
			return fmt.Errorf("failed to parse annotation[%s]=%s as int64: %v", appsv1beta2.DeprecatedRollbackTo, revision, err)
		} else {
			out.Spec.RollbackTo = new(extensions.RollbackConfig)
			out.Spec.RollbackTo.Revision = revision64
		}
		delete(out.Annotations, appsv1beta2.DeprecatedRollbackTo)
	} else {
		out.Spec.RollbackTo = nil
	}

	if err := Convert_v1beta2_DeploymentStatus_To_extensions_DeploymentStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

func Convert_v1beta2_ReplicaSetSpec_To_extensions_ReplicaSetSpec(in *appsv1beta2.ReplicaSetSpec, out *extensions.ReplicaSetSpec, s conversion.Scope) error {
	if in.Replicas != nil {
		out.Replicas = *in.Replicas
	}
	out.MinReadySeconds = in.MinReadySeconds
	out.Selector = in.Selector
	if err := k8s_api_v1.Convert_v1_PodTemplateSpec_To_api_PodTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	return nil
}

func Convert_extensions_Deployment_To_v1beta2_Deployment(in *extensions.Deployment, out *appsv1beta2.Deployment, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_extensions_DeploymentSpec_To_v1beta2_DeploymentSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}

	// Copy deprecated rollbackTo field to annotation for roundtrip
	// TODO: remove this conversion after we delete extensions/v1beta1 and apps/v1beta1 Deployment
	if in.Spec.RollbackTo != nil {
		if out.Annotations == nil {
			out.Annotations = make(map[string]string)
		}
		out.Annotations[appsv1beta2.DeprecatedRollbackTo] = strconv.FormatInt(in.Spec.RollbackTo.Revision, 10)
	} else {
		delete(out.Annotations, appsv1beta2.DeprecatedRollbackTo)
	}

	if err := Convert_extensions_DeploymentStatus_To_v1beta2_DeploymentStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}
