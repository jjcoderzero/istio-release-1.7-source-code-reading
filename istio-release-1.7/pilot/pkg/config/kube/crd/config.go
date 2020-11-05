package crd

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// IstioKind is the generic Kubernetes API object wrapper
type IstioKind struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	Spec               map[string]interface{} `json:"spec"`
}

// GetSpec from a wrapper
func (in *IstioKind) GetSpec() map[string]interface{} {
	return in.Spec
}

// GetObjectMeta from a wrapper
func (in *IstioKind) GetObjectMeta() meta_v1.ObjectMeta {
	return in.ObjectMeta
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IstioKind) DeepCopyInto(out *IstioKind) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IstioKind.
func (in *IstioKind) DeepCopy() *IstioKind {
	if in == nil {
		return nil
	}
	out := new(IstioKind)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IstioKind) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}

	return nil
}

// IstioObject is a k8s wrapper interface for config objects
type IstioObject interface {
	runtime.Object
	GetSpec() map[string]interface{}
	GetObjectMeta() meta_v1.ObjectMeta
}
