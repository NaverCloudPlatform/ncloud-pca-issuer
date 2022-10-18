# X509EndCertParameters

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CommonName** | ***string** | 인증서 이름 혹은 도메인. 최대 64자 입력 가능 | [default to null]
**AltName** | ***string** | 인증서의 다른 이름 | [optional] [default to null]
**Ip** | ***string** | 인증서에 해당하는 IP, \&quot;,\&quot; 구분자 리스트 | [optional] [default to null]
**Uri** | ***string** | 인증서에 해당하는 URI, \&quot;,\&quot; 구분자 리스트 | [optional] [default to null]
**Other** | ***string** | 기타 인증서 정보, {OID};UTF8:{value} 포맷, \&quot;,\&quot; 구분자 리스트 | [optional] [default to null]
**Country** | ***string** | 국가 | [optional] [default to null]
**Organization** | ***string** | 회사 | [optional] [default to null]
**OrganizationUnit** | ***string** | 조직 | [optional] [default to null]
**Locality** | ***string** | 도시 | [optional] [default to null]
**StateProvince** | ***string** | 구/동 | [optional] [default to null]
**StreetAddress** | ***string** | 도로명 | [optional] [default to null]
**PostalCode** | ***string** | 우편번호 | [optional] [default to null]
**UseCnAsSan** | ***bool** | CSR 설정 값 사용 여부 | [optional] [default to null]
**KeyUsage** | **[]\*string** | 키 용도 | [optional] [default to null]
**ExtendedKeyUsage** | **[]\*string** | 확장 키 용도 | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


