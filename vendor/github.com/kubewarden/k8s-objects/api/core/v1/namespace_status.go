// Code generated by go-swagger; DO NOT EDIT.

package v1

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// NamespaceStatus NamespaceStatus is information about the current status of a Namespace.
//
// swagger:model NamespaceStatus
type NamespaceStatus struct {

	// Represents the latest available observations of a namespace's current state.
	Conditions []*NamespaceCondition `json:"conditions,omitempty"`

	// Phase is the current lifecycle phase of the namespace. More info: https://kubernetes.io/docs/tasks/administer-cluster/namespaces/
	Phase string `json:"phase,omitempty"`
}
