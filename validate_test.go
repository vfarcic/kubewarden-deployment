package main

import (
	"encoding/json"
	"testing"

	appsv1 "github.com/kubewarden/k8s-objects/api/apps/v1"
	metav1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
	kubewarden_testing "github.com/kubewarden/policy-sdk-go/testing"
)

var deploymentObject = appsv1.Deployment{
	Metadata: &metav1.ObjectMeta{
		Name:      "my-deployment",
		Namespace: "production",
	},
	Spec: &appsv1.DeploymentSpec{
		Replicas: 3,
	},
}

func TestEmptyReplicasLeadsToApproval(t *testing.T) {
	settings := Settings{}
	deployment := appsv1.Deployment{
		Metadata: &metav1.ObjectMeta{
			Name:      "my-deployment",
			Namespace: "production",
		},
		Spec: &appsv1.DeploymentSpec{},
	}

	payload, err := kubewarden_testing.BuildValidationRequest(&deployment, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if !response.Accepted {
		t.Errorf("Unexpected rejection: msg %s - code %d", *response.Message, response.Code)
	}
}

func TestApproval(t *testing.T) {
	settings := Settings{ReplicasGreaterThan: 2}

	payload, err := kubewarden_testing.BuildValidationRequest(&deploymentObject, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if !response.Accepted {
		t.Error("Unexpected rejection")
	}
}

func TestApproveFixtureDenied(t *testing.T) {
	settings := Settings{
		ReplicasGreaterThan: 2,
	}

	payload, err := kubewarden_testing.BuildValidationRequestFromFixture(
		"test_data/deployment.json",
		&settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted {
		t.Error("Unexpected rejection")
	}
}

func TestRejectionBecauseReplicasIsLessThan(t *testing.T) {
	settings := Settings{
		ReplicasGreaterThan: 2,
	}
	deployment := deploymentObject
	deployment.Spec.Replicas = 2
	payload, err := kubewarden_testing.BuildValidationRequest(&deployment, &settings)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err := json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted {
		t.Error("Unexpected approval")
	}

	expected_message := "The 'my-deployment' Deployment is on the deny list. The spec.replicas must be greater than 2."

	if response.Message == nil {
		t.Errorf("expected response to have a message")
	}

	if *response.Message != expected_message {
		t.Errorf("Got '%s' instead of '%s'", *response.Message, expected_message)
	}
}
