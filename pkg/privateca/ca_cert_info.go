/*
 * api
 *
 * <br/>https://pca.apigw.ntruss.com/api/v1
 *
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package privateca

type CaCertInfo struct {

	// Ca Cert Id
	CaCertId *float32 `json:"caCertId,omitempty"`

	// 인증서 이름
	CommonName *string `json:"commonName,omitempty"`

	// 인증서 발행자 이름
	IssuerName *string `json:"issuerName,omitempty"`

	// 인증서 유효 기간
	NotBeforeDate *float32 `json:"notBeforeDate,omitempty"`

	// 인증서 만료 기간
	NotAfterDate *float32 `json:"notAfterDate,omitempty"`

	// 공개키 알고리즘
	PublicKeyAlgorithm *string `json:"publicKeyAlgorithm,omitempty"`

	// 서명 알고리즘
	SignatureAlgorithm *string `json:"signatureAlgorithm,omitempty"`

	// 시리얼 번호
	SerialNo *string `json:"serialNo,omitempty"`

	// 국가
	Country *string `json:"country,omitempty"`

	// 구/동
	StateProvince *string `json:"stateProvince,omitempty"`

	// 도시
	Locality *string `json:"locality,omitempty"`

	// 회사
	Organization *string `json:"organization,omitempty"`

	// 조직
	OrganizationUnit *string `json:"organizationUnit,omitempty"`

	// Pem 형식 인증서
	CertPem *string `json:"certPem,omitempty"`

	// Pem 형식 인증서 체인
	ChainPem *string `json:"chainPem,omitempty"`
}