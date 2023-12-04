//
// Copyright (c) 2023 Red Hat, Inc.
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

package controller

import (
	"context"
	"fmt"

	bs "janus-idp.io/backstage-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
)

func (r *BackstageReconciler) extraConfigsToVolumes(backstage bs.Backstage) (result []v1.Volume) {
	if backstage.Spec.Backstage == nil || backstage.Spec.Backstage.ExtraConfig == nil {
		return nil
	}
	for _, extraConfig := range backstage.Spec.Backstage.ExtraConfig.Items {
		var volumeSource v1.VolumeSource
		var name string
		switch {
		case extraConfig.ConfigMapRef != nil:
			name = extraConfig.ConfigMapRef.Name
			volumeSource.ConfigMap = &v1.ConfigMapVolumeSource{
				DefaultMode:          pointer.Int32(420),
				LocalObjectReference: v1.LocalObjectReference{Name: name},
			}
		case extraConfig.SecretRef != nil:
			name = extraConfig.SecretRef.Name
			volumeSource.Secret = &v1.SecretVolumeSource{
				DefaultMode: pointer.Int32(420),
				SecretName:  name,
			}
		}
		result = append(result,
			v1.Volume{
				Name:         name,
				VolumeSource: volumeSource,
			},
		)
	}

	return result
}

func (r *BackstageReconciler) addExtraConfigsVolumeMounts(ctx context.Context, backstage bs.Backstage, ns string, deployment *appsv1.Deployment) error {
	if backstage.Spec.Backstage == nil || backstage.Spec.Backstage.ExtraConfig == nil {
		return nil
	}

	appConfigFilenamesList, err := r.extractExtraConfigFileNames(ctx, backstage, ns)
	if err != nil {
		return err
	}

	for i, c := range deployment.Spec.Template.Spec.Containers {
		if c.Name == _defaultBackstageMainContainerName {
			for _, appConfigFilenames := range appConfigFilenamesList {
				for _, f := range appConfigFilenames.files {
					deployment.Spec.Template.Spec.Containers[i].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[i].VolumeMounts,
						v1.VolumeMount{
							Name:      appConfigFilenames.ref,
							MountPath: fmt.Sprintf("%s/%s", backstage.Spec.Backstage.ExtraConfig.MountPath, f),
							SubPath:   f,
						})
				}
			}
			break
		}
	}
	return nil
}

// extractExtraConfigFileNames returns a mapping of extra-config object name and the list of files in it.
// We intentionally do not return a Map, to preserve the iteration order of the ExtraConfigs in the Custom Resource,
// even though we can't guarantee the iteration order of the files listed inside each ConfigMap or Secret.
func (r *BackstageReconciler) extractExtraConfigFileNames(ctx context.Context, backstage bs.Backstage, ns string) (result []appConfigData, err error) {
	if backstage.Spec.Backstage == nil || backstage.Spec.Backstage.ExtraConfig == nil {
		return nil, nil
	}

	for _, appConfig := range backstage.Spec.Backstage.ExtraConfig.Items {
		var files []string
		var name string
		switch {
		case appConfig.ConfigMapRef != nil:
			name = appConfig.ConfigMapRef.Name
			cm := v1.ConfigMap{}
			if err = r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, &cm); err != nil {
				return nil, err
			}
			for filename := range cm.Data {
				// Bear in mind that iteration order over this map is not guaranteed by Go
				files = append(files, filename)
			}
			for filename := range cm.BinaryData {
				// Bear in mind that iteration order over this map is not guaranteed by Go
				files = append(files, filename)
			}
		case appConfig.SecretRef != nil:
			name = appConfig.SecretRef.Name
			sec := v1.Secret{}
			if err = r.Get(ctx, types.NamespacedName{Name: name, Namespace: ns}, &sec); err != nil {
				return nil, err
			}
			for filename := range sec.Data {
				// Bear in mind that iteration order over this map is not guaranteed by Go
				files = append(files, filename)
			}
		}
		result = append(result, appConfigData{
			ref:   name,
			files: files,
		})
	}
	return result, nil
}
