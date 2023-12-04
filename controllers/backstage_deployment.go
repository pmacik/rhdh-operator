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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	_defaultBackstageInitContainerName = "install-dynamic-plugins"
	_defaultBackstageMainContainerName = "backstage-backend"
	_containersWorkingDir              = "/opt/app-root/src"
)

//var (
//	DefaultBackstageDeployment = fmt.Sprintf(`
//apiVersion: apps/v1
//kind: Deployment
//metadata:
// name: backstage
//spec:
// replicas: 1
// selector:
//   matchLabels:
//     janus-idp.io/app:  # placeholder for 'backstage-<cr-name>'
// template:
//   metadata:
//     labels:
//       janus-idp.io/app:  # placeholder for 'backstage-<cr-name>'
//   spec:
//#      serviceAccountName: default
//
//     volumes:
//       - ephemeral:
//           volumeClaimTemplate:
//             spec:
//               accessModes:
//               - ReadWriteOnce
//               resources:
//                 requests:
//                   storage: 1Gi
//         name: dynamic-plugins-root
//       - name: dynamic-plugins-npmrc
//         secret:
//           defaultMode: 420
//           optional: true
//           secretName: dynamic-plugins-npmrc
//
//     initContainers:
//       - command:
//         - ./install-dynamic-plugins.sh
//         - /dynamic-plugins-root
//         env:
//         - name: NPM_CONFIG_USERCONFIG
//           value: %[3]s/.npmrc.dynamic-plugins
//         image: 'quay.io/janus-idp/backstage-showcase:next'
//         imagePullPolicy: IfNotPresent
//         name: %[1]s
//         volumeMounts:
//         - mountPath: /dynamic-plugins-root
//           name: dynamic-plugins-root
//         - mountPath: %[3]s/.npmrc.dynamic-plugins
//           name: dynamic-plugins-npmrc
//           readOnly: true
//           subPath: .npmrc
//         workingDir: %[3]s
//
//     containers:
//       - name: %[2]s
//         image: quay.io/janus-idp/backstage-showcase:next
//         imagePullPolicy: IfNotPresent
//         args:
//           - "--config"
//           - "dynamic-plugins-root/app-config.dynamic-plugins.yaml"
//         readinessProbe:
//           failureThreshold: 3
//           httpGet:
//             path: /healthcheck
//             port: 7007
//             scheme: HTTP
//           initialDelaySeconds: 30
//           periodSeconds: 10
//           successThreshold: 2
//           timeoutSeconds: 2
//         livenessProbe:
//           failureThreshold: 3
//           httpGet:
//             path: /healthcheck
//             port: 7007
//             scheme: HTTP
//           initialDelaySeconds: 60
//           periodSeconds: 10
//           successThreshold: 1
//           timeoutSeconds: 2
//         ports:
//           - name: http
//             containerPort: 7007
//         env:
//           - name: APP_CONFIG_backend_listen_port
//             value: "7007"
//         envFrom:
//           - secretRef:
//               name: postgres-secrets
//#            - secretRef:
//#                name: backstage-secrets
//         volumeMounts:
//           - mountPath: %[3]s/dynamic-plugins-root
//             name: dynamic-plugins-root
//`, _defaultBackstageInitContainerName, _defaultBackstageMainContainerName, _containersWorkingDir)
//)

// ContainerVisitor is called with each container
type ContainerVisitor func(container *v1.Container)

// visitContainers invokes the visitor function for every container in the given pod template spec
func visitContainers(podTemplateSpec *v1.PodTemplateSpec, visitor ContainerVisitor) {
	for i := range podTemplateSpec.Spec.InitContainers {
		visitor(&podTemplateSpec.Spec.InitContainers[i])
	}
	for i := range podTemplateSpec.Spec.Containers {
		visitor(&podTemplateSpec.Spec.Containers[i])
	}
	for i := range podTemplateSpec.Spec.EphemeralContainers {
		visitor((*v1.Container)(&podTemplateSpec.Spec.EphemeralContainers[i].EphemeralContainerCommon))
	}
}

func (r *BackstageReconciler) applyBackstageDeployment(ctx context.Context, backstage bs.Backstage, ns string) error {

	//lg := log.FromContext(ctx)

	deployment := &appsv1.Deployment{}
	err := r.readConfigMapOrDefault(ctx, backstage.Spec.RawRuntimeConfig.BackstageConfigName, "deployment.yaml", ns, deployment)
	if err != nil {
		return fmt.Errorf("failed to read config: %s", err)
	}

	foundDeployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: ns}, foundDeployment)
	if err != nil {
		if errors.IsNotFound(err) {

			setBackstageAppLabel(&deployment.Spec.Template.ObjectMeta.Labels, backstage)
			setBackstageAppLabel(&deployment.Spec.Selector.MatchLabels, backstage)
			r.labels(&deployment.ObjectMeta, backstage)

			if r.OwnsRuntime {
				if err = controllerutil.SetControllerReference(&backstage, deployment, r.Scheme); err != nil {
					return fmt.Errorf("failed to set owner reference: %s", err)
				}
			}

			err = r.addVolumes(ctx, backstage, ns, deployment)
			if err != nil {
				return fmt.Errorf("failed to add volumes to Backstage deployment, reason: %s", err)
			}

			err = r.addVolumeMounts(ctx, backstage, ns, deployment)
			if err != nil {
				return fmt.Errorf("failed to add volume mounts to Backstage deployment, reason: %s", err)
			}

			err = r.addContainerArgs(ctx, backstage, ns, deployment)
			if err != nil {
				return fmt.Errorf("failed to add container args to Backstage deployment, reason: %s", err)
			}

			err = r.addEnvVars(ctx, backstage, ns, deployment)
			if err != nil {
				return fmt.Errorf("failed to add env vars to Backstage deployment, reason: %s", err)
			}

			if backstage.Spec.Backstage != nil {
				deployment.Spec.Replicas = backstage.Spec.Backstage.Replicas

				if backstage.Spec.Backstage.Image != nil {
					visitContainers(&deployment.Spec.Template, func(container *v1.Container) {
						container.Image = *backstage.Spec.Backstage.Image
					})
				}

				if backstage.Spec.Backstage.ImagePullSecret != nil {
					deployment.Spec.Template.Spec.ImagePullSecrets = append(deployment.Spec.Template.Spec.ImagePullSecrets, v1.LocalObjectReference{
						Name: *backstage.Spec.Backstage.ImagePullSecret,
					})
				}
			}

			err = r.Create(ctx, deployment)
			if err != nil {
				return fmt.Errorf("failed to create backstage deployment, reason: %s", err)
			}

		} else {
			return fmt.Errorf("failed to get backstage deployment.yaml, reason: %s", err)
		}
	} else {
		//lg.Info("CR update is ignored for the time")
		err = r.Update(ctx, foundDeployment)
		if err != nil {
			return fmt.Errorf("failed to update backstage deplyment, reason: %s", err)
		}
	}
	return nil
}

func (r *BackstageReconciler) addVolumes(ctx context.Context, backstage bs.Backstage, ns string, deployment *appsv1.Deployment) error {
	dpConfVol, err := r.getDynamicPluginsConfVolume(ctx, backstage, ns)
	if err != nil {
		return err
	}
	if dpConfVol != nil {
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, *dpConfVol)
	}

	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, r.appConfigsToVolumes(backstage)...)
	deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, r.extraConfigsToVolumes(backstage)...)
	return nil
}

func (r *BackstageReconciler) addVolumeMounts(ctx context.Context, backstage bs.Backstage, ns string, deployment *appsv1.Deployment) error {
	err := r.addDynamicPluginsConfVolumeMount(ctx, backstage, ns, deployment)
	if err != nil {
		return err
	}
	err = r.addAppConfigsVolumeMounts(ctx, backstage, ns, deployment)
	if err != nil {
		return err
	}
	return r.addExtraConfigsVolumeMounts(ctx, backstage, ns, deployment)
}

func (r *BackstageReconciler) addContainerArgs(ctx context.Context, backstage bs.Backstage, ns string, deployment *appsv1.Deployment) error {
	return r.addAppConfigsContainerArgs(ctx, backstage, ns, deployment)
}

func (r *BackstageReconciler) addEnvVars(ctx context.Context, backstage bs.Backstage, ns string, deployment *appsv1.Deployment) error {
	err := r.addBackendAuthEnvVar(ctx, backstage, ns, deployment)
	if err != nil {
		return err
	}
	if backstage.Spec.Backstage == nil {
		return nil
	}

	for _, env := range backstage.Spec.Backstage.Env {
		visitContainers(&deployment.Spec.Template, func(container *v1.Container) {
			container.Env = append(container.Env, v1.EnvVar{
				Name:  env.Name,
				Value: env.Value,
			})
		})
	}

	for _, envFrom := range backstage.Spec.Backstage.EnvFrom {
		var (
			name   string
			cmSrc  *v1.ConfigMapEnvSource
			secSrc *v1.SecretEnvSource
		)
		switch {
		case envFrom.ConfigMapRef != nil:
			name = envFrom.ConfigMapRef.Name
			cmSrc = &v1.ConfigMapEnvSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: name,
				},
			}
		case envFrom.SecretRef != nil:
			name = envFrom.SecretRef.Name
			secSrc = &v1.SecretEnvSource{
				LocalObjectReference: v1.LocalObjectReference{
					Name: name,
				},
			}
		}
		visitContainers(&deployment.Spec.Template, func(container *v1.Container) {
			container.EnvFrom = append(container.EnvFrom, v1.EnvFromSource{
				ConfigMapRef: cmSrc,
				SecretRef:    secSrc,
			})
		})
	}

	return nil
}
