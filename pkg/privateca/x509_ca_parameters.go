/*
 * api
 *
 * <br/>https://pca.apigw.ntruss.com/api/v1
 *
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package privateca

type X509CaParameters struct {

	// 인증서 이름 혹은 도메인. 최대 64자 입력 가능
	CommonName *string `json:"commonName"`

	// 인증서의 다른 이름
	AltName *string `json:"altName,omitempty"`

	// 인증서에 해당하는 IP, \",\" 구분자 리스트
	Ip *string `json:"ip,omitempty"`

	// 국가
	Country *string `json:"country,omitempty"`

	// 회사
	Organization *string `json:"organization,omitempty"`

	// 조직
	OrganizationUnit *string `json:"organizationUnit,omitempty"`

	// 도시
	Locality *string `json:"locality,omitempty"`

	// 구/동
	StateProvince *string `json:"stateProvince,omitempty"`
}