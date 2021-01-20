package alertsmanagement

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// SmartDetectorAlertRulesClient is the alertsManagement Client
type SmartDetectorAlertRulesClient struct {
	BaseClient
}

// NewSmartDetectorAlertRulesClient creates an instance of the SmartDetectorAlertRulesClient client.
func NewSmartDetectorAlertRulesClient(subscriptionID string) SmartDetectorAlertRulesClient {
	return NewSmartDetectorAlertRulesClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewSmartDetectorAlertRulesClientWithBaseURI creates an instance of the SmartDetectorAlertRulesClient client using a
// custom endpoint.  Use this when interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds,
// Azure stack).
func NewSmartDetectorAlertRulesClientWithBaseURI(baseURI string, subscriptionID string) SmartDetectorAlertRulesClient {
	return SmartDetectorAlertRulesClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// CreateOrUpdate create or update a Smart Detector alert rule.
// Parameters:
// resourceGroupName - the name of the resource group.
// alertRuleName - the name of the alert rule.
// parameters - parameters supplied to the operation.
func (client SmartDetectorAlertRulesClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, alertRuleName string, parameters AlertRule) (result AlertRule, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SmartDetectorAlertRulesClient.CreateOrUpdate")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: parameters,
			Constraints: []validation.Constraint{{Target: "parameters.AlertRuleProperties", Name: validation.Null, Rule: false,
				Chain: []validation.Constraint{{Target: "parameters.AlertRuleProperties.Frequency", Name: validation.Null, Rule: true, Chain: nil},
					{Target: "parameters.AlertRuleProperties.Detector", Name: validation.Null, Rule: true,
						Chain: []validation.Constraint{{Target: "parameters.AlertRuleProperties.Detector.ID", Name: validation.Null, Rule: true, Chain: nil}}},
					{Target: "parameters.AlertRuleProperties.Scope", Name: validation.Null, Rule: true, Chain: nil},
					{Target: "parameters.AlertRuleProperties.ActionGroups", Name: validation.Null, Rule: true,
						Chain: []validation.Constraint{{Target: "parameters.AlertRuleProperties.ActionGroups.GroupIds", Name: validation.Null, Rule: true, Chain: nil}}},
				}}}}}); err != nil {
		return result, validation.NewError("alertsmanagement.SmartDetectorAlertRulesClient", "CreateOrUpdate", err.Error())
	}

	req, err := client.CreateOrUpdatePreparer(ctx, resourceGroupName, alertRuleName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	resp, err := client.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "CreateOrUpdate", resp, "Failure responding to request")
		return
	}

	return
}

// CreateOrUpdatePreparer prepares the CreateOrUpdate request.
func (client SmartDetectorAlertRulesClient) CreateOrUpdatePreparer(ctx context.Context, resourceGroupName string, alertRuleName string, parameters AlertRule) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"alertRuleName":     autorest.Encode("path", alertRuleName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.alertsManagement/smartDetectorAlertRules/{alertRuleName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateOrUpdateSender sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (client SmartDetectorAlertRulesClient) CreateOrUpdateSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client SmartDetectorAlertRulesClient) CreateOrUpdateResponder(resp *http.Response) (result AlertRule, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete delete an existing Smart Detector alert rule.
// Parameters:
// resourceGroupName - the name of the resource group.
// alertRuleName - the name of the alert rule.
func (client SmartDetectorAlertRulesClient) Delete(ctx context.Context, resourceGroupName string, alertRuleName string) (result autorest.Response, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SmartDetectorAlertRulesClient.Delete")
		defer func() {
			sc := -1
			if result.Response != nil {
				sc = result.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("alertsmanagement.SmartDetectorAlertRulesClient", "Delete", err.Error())
	}

	req, err := client.DeletePreparer(ctx, resourceGroupName, alertRuleName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "Delete", nil, "Failure preparing request")
		return
	}

	resp, err := client.DeleteSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "Delete", resp, "Failure sending request")
		return
	}

	result, err = client.DeleteResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "Delete", resp, "Failure responding to request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client SmartDetectorAlertRulesClient) DeletePreparer(ctx context.Context, resourceGroupName string, alertRuleName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"alertRuleName":     autorest.Encode("path", alertRuleName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.alertsManagement/smartDetectorAlertRules/{alertRuleName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client SmartDetectorAlertRulesClient) DeleteSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client SmartDetectorAlertRulesClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get get a specific Smart Detector alert rule.
// Parameters:
// resourceGroupName - the name of the resource group.
// alertRuleName - the name of the alert rule.
// expandDetector - indicates if Smart Detector should be expanded.
func (client SmartDetectorAlertRulesClient) Get(ctx context.Context, resourceGroupName string, alertRuleName string, expandDetector *bool) (result AlertRule, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SmartDetectorAlertRulesClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("alertsmanagement.SmartDetectorAlertRulesClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, alertRuleName, expandDetector)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client SmartDetectorAlertRulesClient) GetPreparer(ctx context.Context, resourceGroupName string, alertRuleName string, expandDetector *bool) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"alertRuleName":     autorest.Encode("path", alertRuleName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if expandDetector != nil {
		queryParameters["expandDetector"] = autorest.Encode("query", *expandDetector)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.alertsManagement/smartDetectorAlertRules/{alertRuleName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client SmartDetectorAlertRulesClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client SmartDetectorAlertRulesClient) GetResponder(resp *http.Response) (result AlertRule, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List list all the existing Smart Detector alert rules within the subscription.
// Parameters:
// expandDetector - indicates if Smart Detector should be expanded.
func (client SmartDetectorAlertRulesClient) List(ctx context.Context, expandDetector *bool) (result AlertRulesListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SmartDetectorAlertRulesClient.List")
		defer func() {
			sc := -1
			if result.arl.Response.Response != nil {
				sc = result.arl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("alertsmanagement.SmartDetectorAlertRulesClient", "List", err.Error())
	}

	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, expandDetector)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.arl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "List", resp, "Failure sending request")
		return
	}

	result.arl, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "List", resp, "Failure responding to request")
		return
	}
	if result.arl.hasNextLink() && result.arl.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListPreparer prepares the List request.
func (client SmartDetectorAlertRulesClient) ListPreparer(ctx context.Context, expandDetector *bool) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if expandDetector != nil {
		queryParameters["expandDetector"] = autorest.Encode("query", *expandDetector)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/providers/microsoft.alertsManagement/smartDetectorAlertRules", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client SmartDetectorAlertRulesClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client SmartDetectorAlertRulesClient) ListResponder(resp *http.Response) (result AlertRulesList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client SmartDetectorAlertRulesClient) listNextResults(ctx context.Context, lastResults AlertRulesList) (result AlertRulesList, err error) {
	req, err := lastResults.alertRulesListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "listNextResults", resp, "Failure responding to next results request")
		return
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client SmartDetectorAlertRulesClient) ListComplete(ctx context.Context, expandDetector *bool) (result AlertRulesListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SmartDetectorAlertRulesClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, expandDetector)
	return
}

// ListByResourceGroup list all the existing Smart Detector alert rules within the subscription and resource group.
// Parameters:
// resourceGroupName - the name of the resource group.
// expandDetector - indicates if Smart Detector should be expanded.
func (client SmartDetectorAlertRulesClient) ListByResourceGroup(ctx context.Context, resourceGroupName string, expandDetector *bool) (result AlertRulesListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SmartDetectorAlertRulesClient.ListByResourceGroup")
		defer func() {
			sc := -1
			if result.arl.Response.Response != nil {
				sc = result.arl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("alertsmanagement.SmartDetectorAlertRulesClient", "ListByResourceGroup", err.Error())
	}

	result.fn = client.listByResourceGroupNextResults
	req, err := client.ListByResourceGroupPreparer(ctx, resourceGroupName, expandDetector)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "ListByResourceGroup", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.arl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "ListByResourceGroup", resp, "Failure sending request")
		return
	}

	result.arl, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "ListByResourceGroup", resp, "Failure responding to request")
		return
	}
	if result.arl.hasNextLink() && result.arl.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListByResourceGroupPreparer prepares the ListByResourceGroup request.
func (client SmartDetectorAlertRulesClient) ListByResourceGroupPreparer(ctx context.Context, resourceGroupName string, expandDetector *bool) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if expandDetector != nil {
		queryParameters["expandDetector"] = autorest.Encode("query", *expandDetector)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.alertsManagement/smartDetectorAlertRules", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByResourceGroupSender sends the ListByResourceGroup request. The method will close the
// http.Response Body if it receives an error.
func (client SmartDetectorAlertRulesClient) ListByResourceGroupSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListByResourceGroupResponder handles the response to the ListByResourceGroup request. The method always
// closes the http.Response Body.
func (client SmartDetectorAlertRulesClient) ListByResourceGroupResponder(resp *http.Response) (result AlertRulesList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByResourceGroupNextResults retrieves the next set of results, if any.
func (client SmartDetectorAlertRulesClient) listByResourceGroupNextResults(ctx context.Context, lastResults AlertRulesList) (result AlertRulesList, err error) {
	req, err := lastResults.alertRulesListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "listByResourceGroupNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByResourceGroupSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "listByResourceGroupNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByResourceGroupResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "listByResourceGroupNextResults", resp, "Failure responding to next results request")
		return
	}
	return
}

// ListByResourceGroupComplete enumerates all values, automatically crossing page boundaries as required.
func (client SmartDetectorAlertRulesClient) ListByResourceGroupComplete(ctx context.Context, resourceGroupName string, expandDetector *bool) (result AlertRulesListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SmartDetectorAlertRulesClient.ListByResourceGroup")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.ListByResourceGroup(ctx, resourceGroupName, expandDetector)
	return
}

// Patch patch a specific Smart Detector alert rule.
// Parameters:
// resourceGroupName - the name of the resource group.
// alertRuleName - the name of the alert rule.
// parameters - parameters supplied to the operation.
func (client SmartDetectorAlertRulesClient) Patch(ctx context.Context, resourceGroupName string, alertRuleName string, parameters AlertRulePatchObject) (result AlertRule, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/SmartDetectorAlertRulesClient.Patch")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.SubscriptionID,
			Constraints: []validation.Constraint{{Target: "client.SubscriptionID", Name: validation.MinLength, Rule: 1, Chain: nil}}}}); err != nil {
		return result, validation.NewError("alertsmanagement.SmartDetectorAlertRulesClient", "Patch", err.Error())
	}

	req, err := client.PatchPreparer(ctx, resourceGroupName, alertRuleName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "Patch", nil, "Failure preparing request")
		return
	}

	resp, err := client.PatchSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "Patch", resp, "Failure sending request")
		return
	}

	result, err = client.PatchResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "alertsmanagement.SmartDetectorAlertRulesClient", "Patch", resp, "Failure responding to request")
		return
	}

	return
}

// PatchPreparer prepares the Patch request.
func (client SmartDetectorAlertRulesClient) PatchPreparer(ctx context.Context, resourceGroupName string, alertRuleName string, parameters AlertRulePatchObject) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"alertRuleName":     autorest.Encode("path", alertRuleName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2019-06-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	parameters.ID = nil
	parameters.Type = nil
	parameters.Name = nil
	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.alertsManagement/smartDetectorAlertRules/{alertRuleName}", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// PatchSender sends the Patch request. The method will close the
// http.Response Body if it receives an error.
func (client SmartDetectorAlertRulesClient) PatchSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// PatchResponder handles the response to the Patch request. The method always
// closes the http.Response Body.
func (client SmartDetectorAlertRulesClient) PatchResponder(resp *http.Response) (result AlertRule, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
