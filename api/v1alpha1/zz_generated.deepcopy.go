//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023 Red Hat Inc..

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppConfig) DeepCopyInto(out *AppConfig) {
	*out = *in
	if in.ConfigMaps != nil {
		in, out := &in.ConfigMaps, &out.ConfigMaps
		*out = make([]ObjectKeyRef, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppConfig.
func (in *AppConfig) DeepCopy() *AppConfig {
	if in == nil {
		return nil
	}
	out := new(AppConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Application) DeepCopyInto(out *Application) {
	*out = *in
	if in.AppConfig != nil {
		in, out := &in.AppConfig, &out.AppConfig
		*out = new(AppConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.ExtraFiles != nil {
		in, out := &in.ExtraFiles, &out.ExtraFiles
		*out = new(ExtraFiles)
		(*in).DeepCopyInto(*out)
	}
	if in.ExtraEnvs != nil {
		in, out := &in.ExtraEnvs, &out.ExtraEnvs
		*out = new(ExtraEnvs)
		(*in).DeepCopyInto(*out)
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	if in.Image != nil {
		in, out := &in.Image, &out.Image
		*out = new(string)
		**out = **in
	}
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Application.
func (in *Application) DeepCopy() *Application {
	if in == nil {
		return nil
	}
	out := new(Application)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Backstage) DeepCopyInto(out *Backstage) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Backstage.
func (in *Backstage) DeepCopy() *Backstage {
	if in == nil {
		return nil
	}
	out := new(Backstage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Backstage) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackstageList) DeepCopyInto(out *BackstageList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Backstage, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackstageList.
func (in *BackstageList) DeepCopy() *BackstageList {
	if in == nil {
		return nil
	}
	out := new(BackstageList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BackstageList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackstageSpec) DeepCopyInto(out *BackstageSpec) {
	*out = *in
	if in.Application != nil {
		in, out := &in.Application, &out.Application
		*out = new(Application)
		(*in).DeepCopyInto(*out)
	}
	out.RawRuntimeConfig = in.RawRuntimeConfig
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackstageSpec.
func (in *BackstageSpec) DeepCopy() *BackstageSpec {
	if in == nil {
		return nil
	}
	out := new(BackstageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackstageStatus) DeepCopyInto(out *BackstageStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackstageStatus.
func (in *BackstageStatus) DeepCopy() *BackstageStatus {
	if in == nil {
		return nil
	}
	out := new(BackstageStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Env) DeepCopyInto(out *Env) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Env.
func (in *Env) DeepCopy() *Env {
	if in == nil {
		return nil
	}
	out := new(Env)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtraEnvs) DeepCopyInto(out *ExtraEnvs) {
	*out = *in
	if in.ConfigMaps != nil {
		in, out := &in.ConfigMaps, &out.ConfigMaps
		*out = make([]ObjectKeyRef, len(*in))
		copy(*out, *in)
	}
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]ObjectKeyRef, len(*in))
		copy(*out, *in)
	}
	if in.Envs != nil {
		in, out := &in.Envs, &out.Envs
		*out = make([]Env, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtraEnvs.
func (in *ExtraEnvs) DeepCopy() *ExtraEnvs {
	if in == nil {
		return nil
	}
	out := new(ExtraEnvs)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExtraFiles) DeepCopyInto(out *ExtraFiles) {
	*out = *in
	if in.ConfigMaps != nil {
		in, out := &in.ConfigMaps, &out.ConfigMaps
		*out = make([]ObjectKeyRef, len(*in))
		copy(*out, *in)
	}
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]ObjectKeyRef, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExtraFiles.
func (in *ExtraFiles) DeepCopy() *ExtraFiles {
	if in == nil {
		return nil
	}
	out := new(ExtraFiles)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectKeyRef) DeepCopyInto(out *ObjectKeyRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectKeyRef.
func (in *ObjectKeyRef) DeepCopy() *ObjectKeyRef {
	if in == nil {
		return nil
	}
	out := new(ObjectKeyRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RuntimeConfig) DeepCopyInto(out *RuntimeConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RuntimeConfig.
func (in *RuntimeConfig) DeepCopy() *RuntimeConfig {
	if in == nil {
		return nil
	}
	out := new(RuntimeConfig)
	in.DeepCopyInto(out)
	return out
}
