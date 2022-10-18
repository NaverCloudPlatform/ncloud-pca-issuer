/*
 * api
 *
 * <br/>https://pca.apigw.ntruss.com/api/v1
 *
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package privateca

type IssueEndCertResponseData struct {

	// 인증서
	Certificate *string `json:"certificate,omitempty"`

	// 발행자 인증서
	Issuer *string `json:"issuer,omitempty"`

	// Ca 인증서 체인
	CaChain *string `json:"caChain,omitempty"`

	// 인증서 개인키
	PrivateKey *string `json:"privateKey,omitempty"`

	// 인증서 시리얼 번호
	SerialNo *string `json:"serialNo,omitempty"`

	// Ocsp Responder 인증서
	OcspResponder *string `json:"ocspResponder,omitempty"`
}