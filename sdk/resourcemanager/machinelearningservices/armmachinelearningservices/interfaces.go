//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmachinelearningservices

// ComputeClassification provides polymorphic access to related types.
// Call the interface's GetCompute() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AKS, *AmlCompute, *Compute, *ComputeInstance, *DataFactory, *DataLakeAnalytics, *Databricks, *HDInsight, *Kubernetes,
// - *SynapseSpark, *VirtualMachine
type ComputeClassification interface {
	// GetCompute returns the Compute content of the underlying type.
	GetCompute() *Compute
}

// ComputeSecretsClassification provides polymorphic access to related types.
// Call the interface's GetComputeSecrets() method to access the common type.
// Use a type switch to determine the concrete type.  The possible types are:
// - *AksComputeSecrets, *ComputeSecrets, *DatabricksComputeSecrets, *VirtualMachineSecrets
type ComputeSecretsClassification interface {
	// GetComputeSecrets returns the ComputeSecrets content of the underlying type.
	GetComputeSecrets() *ComputeSecrets
}

