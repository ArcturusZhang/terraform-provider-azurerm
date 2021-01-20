package servicefabricmesh

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
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// ServiceClient is the service Fabric Mesh Management Client
type ServiceClient struct {
	BaseClient
}

// NewServiceClient creates an instance of the ServiceClient client.
func NewServiceClient(subscriptionID string) ServiceClient {
	return NewServiceClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewServiceClientWithBaseURI creates an instance of the ServiceClient client using a custom endpoint.  Use this when
// interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewServiceClientWithBaseURI(baseURI string, subscriptionID string) ServiceClient {
	return ServiceClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get gets the information about the service resource with the given name. The information include the description and
// other properties of the service.
// Parameters:
// resourceGroupName - azure resource group name
// applicationResourceName - the identity of the application.
// serviceResourceName - the identity of the service.
func (client ServiceClient) Get(ctx context.Context, resourceGroupName string, applicationResourceName string, serviceResourceName string) (result ServiceResourceDescription, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	req, err := client.GetPreparer(ctx, resourceGroupName, applicationResourceName, serviceResourceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicefabricmesh.ServiceClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "servicefabricmesh.ServiceClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicefabricmesh.ServiceClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client ServiceClient) GetPreparer(ctx context.Context, resourceGroupName string, applicationResourceName string, serviceResourceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"applicationResourceName": applicationResourceName,
		"resourceGroupName":       autorest.Encode("path", resourceGroupName),
		"serviceResourceName":     serviceResourceName,
		"subscriptionId":          autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabricMesh/applications/{applicationResourceName}/services/{serviceResourceName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ServiceClient) GetResponder(resp *http.Response) (result ServiceResourceDescription, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List gets the information about all services of an application resource. The information include the description and
// other properties of the Service.
// Parameters:
// resourceGroupName - azure resource group name
// applicationResourceName - the identity of the application.
func (client ServiceClient) List(ctx context.Context, resourceGroupName string, applicationResourceName string) (result ServiceResourceDescriptionListPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceClient.List")
		defer func() {
			sc := -1
			if result.srdl.Response.Response != nil {
				sc = result.srdl.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, resourceGroupName, applicationResourceName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicefabricmesh.ServiceClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.srdl.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "servicefabricmesh.ServiceClient", "List", resp, "Failure sending request")
		return
	}

	result.srdl, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicefabricmesh.ServiceClient", "List", resp, "Failure responding to request")
		return
	}
	if result.srdl.hasNextLink() && result.srdl.IsEmpty() {
		err = result.NextWithContext(ctx)
	}

	return
}

// ListPreparer prepares the List request.
func (client ServiceClient) ListPreparer(ctx context.Context, resourceGroupName string, applicationResourceName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"applicationResourceName": applicationResourceName,
		"resourceGroupName":       autorest.Encode("path", resourceGroupName),
		"subscriptionId":          autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2018-09-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabricMesh/applications/{applicationResourceName}/services", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client ServiceClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client ServiceClient) ListResponder(resp *http.Response) (result ServiceResourceDescriptionList, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client ServiceClient) listNextResults(ctx context.Context, lastResults ServiceResourceDescriptionList) (result ServiceResourceDescriptionList, err error) {
	req, err := lastResults.serviceResourceDescriptionListPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "servicefabricmesh.ServiceClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "servicefabricmesh.ServiceClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "servicefabricmesh.ServiceClient", "listNextResults", resp, "Failure responding to next results request")
		return
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client ServiceClient) ListComplete(ctx context.Context, resourceGroupName string, applicationResourceName string) (result ServiceResourceDescriptionListIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/ServiceClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, resourceGroupName, applicationResourceName)
	return
}
