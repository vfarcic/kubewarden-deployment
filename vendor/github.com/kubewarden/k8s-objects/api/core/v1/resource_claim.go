// Code generated by go-swagger; DO NOT EDIT.

package v1

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

// ResourceClaim ResourceClaim references one entry in PodSpec.ResourceClaims.
//
// swagger:model ResourceClaim
type ResourceClaim struct {

	// Name must match the name of one entry in pod.spec.resourceClaims of the Pod where this field is used. It makes that resource available inside a container.
	// Required: true
	Name *string `json:"name"`
}
