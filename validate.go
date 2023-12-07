package main

import (
	// "encoding/json"
	// "fmt"
	// "strings"

	// onelog "github.com/francoispqt/onelog"
	// corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	"encoding/json"
	"fmt"
	"strconv"

	onelog "github.com/francoispqt/onelog"
	appsv1 "github.com/kubewarden/k8s-objects/api/apps/v1"
	metav1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"

	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	// kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

// TODO: Remove
type Sql struct {
	APIVersion string             `json:"apiVersion,omitempty"`
	Kind       string             `json:"kind,omitempty"`
	Metadata   *metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec       *SqlSpec           `json:"spec,omitempty"`
}

// TODO: Remove
type SqlSpec struct {
	ID         string            `json:"id"`
	Parameters SqlSpecParameters `json:"parameters"`
}

// TODO: Remove
type SqlSpecParameters struct {
	Version string `json:"version"`
	Size    string `json:"size"`
}

func validate(payload []byte) ([]byte, error) {
	// Create a ValidationRequest instance from the incoming payload
	validationRequest := kubewarden_protocol.ValidationRequest{}
	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	// Create a Settings instance from the ValidationRequest object
	settings, err := NewSettingsFromValidationReq(&validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(400))
	}

	// Access the **raw** JSON that describes the object
	deploymentJson := validationRequest.Request.Object

	// Try to create a Deployment instance using the RAW JSON we got from the
	// ValidationRequest.
	deployment := &appsv1.Deployment{}
	if err := json.Unmarshal([]byte(deploymentJson), deployment); err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(
				fmt.Sprintf("Cannot decode Deployment object: %s", err.Error())),
			kubewarden.Code(400))
	}

	logger.DebugWithFields("validating Deployment object", func(e onelog.Entry) {
		e.String("name", deployment.Metadata.Name)
		e.String("namespace", deployment.Metadata.Namespace)
	})

	if !settings.IsReplicasGreaterThanAllowed(int(deployment.Spec.Replicas)) {
		logger.InfoWithFields("rejecting Deployment object", func(e onelog.Entry) {
			e.String("name", deployment.Metadata.Name)
			e.String("replicas greater than", strconv.Itoa(settings.ReplicasGreaterThan))
		})

		message := fmt.Sprintf("The '%s' Deployment is on the deny list. The spec.replicas must be greater than %d.", deployment.Metadata.Name, deployment.Spec.Replicas)
		return kubewarden.RejectRequest(kubewarden.Message(message), kubewarden.NoCode)
	}

	return kubewarden.AcceptRequest()
}
