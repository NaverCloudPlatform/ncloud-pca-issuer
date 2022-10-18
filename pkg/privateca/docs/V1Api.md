# \V1Api

All URIs are relative to *https://pca.apigw.ntruss.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CaCaTagActivatePost**](V1Api.md#CaCaTagActivatePost) | **Post** /ca/{caTag}/activate | 
[**CaCaTagCertGet**](V1Api.md#CaCaTagCertGet) | **Get** /ca/{caTag}/cert | 
[**CaCaTagCertPost**](V1Api.md#CaCaTagCertPost) | **Post** /ca/{caTag}/cert | 
[**CaCaTagCertSerialNoGet**](V1Api.md#CaCaTagCertSerialNoGet) | **Get** /ca/{caTag}/cert/{serialNo} | 
[**CaCaTagCertSerialNoRevokePost**](V1Api.md#CaCaTagCertSerialNoRevokePost) | **Post** /ca/{caTag}/cert/{serialNo}/revoke | 
[**CaCaTagCertSignPost**](V1Api.md#CaCaTagCertSignPost) | **Post** /ca/{caTag}/cert/sign | 
[**CaCaTagChainGet**](V1Api.md#CaCaTagChainGet) | **Get** /ca/{caTag}/chain | 
[**CaCaTagCrlConfigGet**](V1Api.md#CaCaTagCrlConfigGet) | **Get** /ca/{caTag}/crl/config | 
[**CaCaTagCrlConfigPut**](V1Api.md#CaCaTagCrlConfigPut) | **Put** /ca/{caTag}/crl/config | 
[**CaCaTagCrlGet**](V1Api.md#CaCaTagCrlGet) | **Get** /ca/{caTag}/crl | 
[**CaCaTagDelete**](V1Api.md#CaCaTagDelete) | **Delete** /ca/{caTag} | 
[**CaCaTagGet**](V1Api.md#CaCaTagGet) | **Get** /ca/{caTag} | 
[**CaCaTagPut**](V1Api.md#CaCaTagPut) | **Put** /ca/{caTag} | 
[**CaCaTagSubCsrGet**](V1Api.md#CaCaTagSubCsrGet) | **Get** /ca/{caTag}/sub/csr | 
[**CaCaTagSubSignPost**](V1Api.md#CaCaTagSubSignPost) | **Post** /ca/{caTag}/sub/sign | 
[**CaCaTagTrimPost**](V1Api.md#CaCaTagTrimPost) | **Post** /ca/{caTag}/trim | 
[**CaCaTagUrlsDelete**](V1Api.md#CaCaTagUrlsDelete) | **Delete** /ca/{caTag}/urls | 
[**CaCaTagUrlsPut**](V1Api.md#CaCaTagUrlsPut) | **Put** /ca/{caTag}/urls | 
[**CaCaTagUsersGet**](V1Api.md#CaCaTagUsersGet) | **Get** /ca/{caTag}/users | 
[**CaCaTagUsersIdNoDelete**](V1Api.md#CaCaTagUsersIdNoDelete) | **Delete** /ca/{caTag}/users/{idNo} | 
[**CaCaTagUsersPost**](V1Api.md#CaCaTagUsersPost) | **Post** /ca/{caTag}/users | 
[**CaGet**](V1Api.md#CaGet) | **Get** /ca | 
[**CaPost**](V1Api.md#CaPost) | **Post** /ca | 


# **CaCaTagActivatePost**
> ActivateSubCaResponse CaCaTagActivatePost(activateSubCa, caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**activateSubCa** | **[\*ActivateSubCa](ActivateSubCa.md)** |  | **caTag** | **string** | caTag | 

### Return type

*[**ActivateSubCaResponse**](ActivateSubCaResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagCertGet**
> ListIssuedEndCertsResponse CaCaTagCertGet(caTag, optional)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | 
 **optional** | **map[string]interface{}** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a map[string]interface{}.

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | **pageNo** | **string** |  | 

### Return type

*[**ListIssuedEndCertsResponse**](ListIssuedEndCertsResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagCertPost**
> IssueEndCertResponse CaCaTagCertPost(createEndCert, caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**createEndCert** | **[\*CreateEndCert](CreateEndCert.md)** |  | **caTag** | **string** | caTag | 

### Return type

*[**IssueEndCertResponse**](IssueEndCertResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagCertSerialNoGet**
> GetEndCertInfoResponse CaCaTagCertSerialNoGet(caTag, serialNo)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | **serialNo** | **string** | serialNo | 

### Return type

*[**GetEndCertInfoResponse**](GetEndCertInfoResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagCertSerialNoRevokePost**
> BaseResponse CaCaTagCertSerialNoRevokePost(caTag, serialNo)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | **serialNo** | **string** | serialNo | 

### Return type

*[**BaseResponse**](BaseResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagCertSignPost**
> CaCaTagCertSignPost(signCsr, caTag, optional)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**signCsr** | **[\*SignCsr](SignCsr.md)** |  | **caTag** | **string** | caTag | 
 **optional** | **map[string]interface{}** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a map[string]interface{}.

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**signCsr** | **[]\*[SignCsr](SignCsr.md)** |  | **caTag** | **string** | caTag | **period** | **string** |  | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagChainGet**
> GetChainResponse CaCaTagChainGet(caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | 

### Return type

*[**GetChainResponse**](GetChainResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagCrlConfigGet**
> CaCaTagCrlConfigGet(caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagCrlConfigPut**
> CaCaTagCrlConfigPut(updateCrlConfig, caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**updateCrlConfig** | **[\*UpdateCrlConfig](UpdateCrlConfig.md)** |  | **caTag** | **string** | caTag | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagCrlGet**
> GetCrlResponse CaCaTagCrlGet(caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | 

### Return type

*[**GetCrlResponse**](GetCrlResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagDelete**
> BaseResponse CaCaTagDelete(caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | 

### Return type

*[**BaseResponse**](BaseResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagGet**
> GetCaInfoResponse CaCaTagGet(caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | 

### Return type

*[**GetCaInfoResponse**](GetCaInfoResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagPut**
> UpdateCaResponse CaCaTagPut(updateCa, caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**updateCa** | **[\*UpdateCa](UpdateCa.md)** |  | **caTag** | **string** | caTag | 

### Return type

*[**UpdateCaResponse**](UpdateCaResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagSubCsrGet**
> CaCaTagSubCsrGet(caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagSubSignPost**
> CaCaTagSubSignPost(signCsr, caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**signCsr** | **[\*SignCsr](SignCsr.md)** |  | **caTag** | **string** | caTag | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagTrimPost**
> CaCaTagTrimPost(caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagUrlsDelete**
> CaCaTagUrlsDelete(caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagUrlsPut**
> CaCaTagUrlsPut(modifyOcspUrl, caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**modifyOcspUrl** | **[\*ModifyOcspUrl](ModifyOcspUrl.md)** |  | **caTag** | **string** | caTag | 

### Return type

 (empty response body)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagUsersGet**
> ListCaUsersResponse CaCaTagUsersGet(caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | 

### Return type

*[**ListCaUsersResponse**](ListCaUsersResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagUsersIdNoDelete**
> BaseResponse CaCaTagUsersIdNoDelete(caTag, idNo)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**caTag** | **string** | caTag | **idNo** | **string** | idNo | 

### Return type

*[**BaseResponse**](BaseResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaCaTagUsersPost**
> BaseResponse CaCaTagUsersPost(addCaUser, caTag)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**addCaUser** | **[\*AddCaUser](AddCaUser.md)** |  | **caTag** | **string** | caTag | 

### Return type

*[**BaseResponse**](BaseResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaGet**
> ListCaResponse CaGet(optional)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **optional** | **map[string]interface{}** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a map[string]interface{}.

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**pageNo** | **int32** |  | 

### Return type

*[**ListCaResponse**](ListCaResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CaPost**
> CreateCaResponse CaPost(createCa, caType, optional)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**createCa** | **[\*CreateCa](CreateCa.md)** |  | **caType** | **string** |  | 
 **optional** | **map[string]interface{}** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a map[string]interface{}.

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**createCa** | **[]\*[CreateCa](CreateCa.md)** |  | **caType** | **string** |  | **issuerId** | **int64** |  | 

### Return type

*[**CreateCaResponse**](CreateCaResponse.md)

### Authorization

[x-ncp-iam](../README.md#x-ncp-iam)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

