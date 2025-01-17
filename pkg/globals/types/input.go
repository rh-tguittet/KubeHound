package types

import (
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

type PodType *corev1.Pod
type NodeType *corev1.Node
type ContainerType *corev1.Container
type VolumeMountType *corev1.VolumeMount
type RoleType *rbacv1.Role
type RoleBindingType *rbacv1.RoleBinding
type ClusterRoleType *rbacv1.ClusterRole
type ClusterRoleBindingType *rbacv1.ClusterRoleBinding
type EndpointType *discoveryv1.EndpointSlice

// Openshift specific
type RouteType *routev1.Route

type InputType interface {
	PodType | NodeType | ContainerType | VolumeMountType | RoleType | RoleBindingType | ClusterRoleType | ClusterRoleBindingType | EndpointType | RouteType
}

// Openshift specific types for ListInputType
type openshiftListInputType interface {
	routev1.RouteList
}

type ListInputType interface {
	corev1.PodList | corev1.NodeList | rbacv1.RoleList | rbacv1.RoleBindingList | rbacv1.ClusterRoleList | rbacv1.ClusterRoleBindingList | discoveryv1.EndpointSliceList | openshiftListInputType
}
