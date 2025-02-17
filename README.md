# Ncloud Private CA Issuer for cert-manager
이 레포지토리는 네이버 클라우드 [Private CA](https://www.ncloud.com/product/security/privateCA)를 사용하는 cert-manager의 [external Issuer](https://cert-manager.io/docs/configuration/external/#known-external-issuers)를 포함하고 있습니다.

## 시작하기

### 사전준비

#### Private CA 이용신청
[Private CA](https://www.ncloud.com/product/security/privateCA) 이용 신청하기 버튼을 눌러 신청합니다.

#### Private CA 서비스로 사설 CA 생성
[사설 CA 생성](https://guide.ncloud-docs.com/docs/privateca-use#%EC%82%AC%EC%84%A4-ca-%EC%83%9D%EC%84%B1) 가이드를 참고하여 CA를 생성합니다.

#### cert-manager
클러스터에 cert-manager가 설치되지 않았다면, cert-manager [official documentation](https://cert-manager.io/docs/installation/kubernetes/)를 참고하여 설치합니다.

#### Ncloud Private CA Issuer for cert-manager 설치
`cert-manager` Namespace에 설치한다고 가정합니다.

[Realeses](https://github.com/NaverCloudPlatform/ncloud-pca-issuer/releases)를 확인하여, 최신버전을 확인합니다.
버전에 맞춰 아래 명령을 수행합니다.

```shell
kubectl apply -f https://github.com/NaverCloudPlatform/ncloud-pca-issuer/releases/download/v0.1.0/ncloud-pca-issuer-v0.1.0.yaml
```


#### Custom 배포(개발자용)

`config/rbac/role.yaml`와 `config/rbac/role_binding.yaml` 파일의 ClusterRole 과 ClusterRolebinding을 확인합니다. 
기본적으로 cet-manager 네임스페이스에 위치한  `ncloud-pca-issuer` 서비스어카운트가 필요한 모든 권한을 가지고 있습니다. 
필요에 따라 수정해서 사용 하십시오.


```shell
kubectl create serviceaccount -n cert-manager ncloud-pca-issuer

kubectl apply -f config/rbac/role.yaml
kubectl apply -f config/rbac/role_binding.yaml
```

Install the Ncloud PCA Issuer CRDs in `config/crd` 디렉토리의 Ncloud PCA Issuer CRD들을 설치합니다. 이 매니페스트는 kustomization을 사용합니다(`-k`옵션을 추가하세요)


```shell
kubectl apply -k config/crd
```

##### Controller 이미지 빌드와 배포

**참고**: Ncloud 공식 이미지를 사용하면 이 단계를 건너뛸 수 있습니다.

이미지를 빌드하기 위해서
[kubebuilder installed](https://book.kubebuilder.io/quick-start.html#installation)가 완료되어야 합니다.

docker 이미지 빌드:

```shell
make docker-build
```

Docker 이미지 푸시하고 kind에 load해 테스트:

```shell
make docker-push || kind load docker-image nks-release.kr.ncr.ntruss.com/cert-manager-ncloud-pca-issuer:latest
```

#### controller 배포하기

Issuer controller 배포:

```shell
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ncloud-pca-issuer
  namespace: cert-manager
  labels:
    app: ncloud-pca-issuer
spec:
  selector:
    matchLabels:
      app: ncloud-pca-issuer
  replicas: 1
  template:
    metadata:
      labels:
        app: ncloud-pca-issuer
    spec:
      serviceAccountName: controller-manager
      containers:
      # update the image to your registry if you built and pushed your own image.
      - image: nks-release.kr.ncr.ntruss.com/cert-manager-ncloud-pca-issuer:latest 
        imagePullPolicy: IfNotPresent
        name: ncloud-pca-issuer
        resources:
          limits:
            cpu: 100m
            memory: 30Mi
          requests:
            cpu: 100m
            memory: 20Mi
      terminationGracePeriodSeconds: 10
EOF
```

기본적으로 Ncloud PCA Issuer controller는 `cert-manager` 네임스페이스에 배포됩니다.

```shell
NAME                                                    READY   STATUS    RESTARTS   AGE
cert-manager-b4d6fd99b-jsr8p                            1/1     Running   0          38h
cert-manager-cainjector-74bfccdfdf-q62qx                1/1     Running   0          38h
cert-manager-webhook-65b766b5f8-m2cpb                   1/1     Running   0          38h
ncloud-pca-issuer-controller-manager-6b777b58b9-vfm5b   2/2     Running   0          37h
```

### SubAccount 인증 설정
Private CA 권한이 있는 [SubAccount](https://guide.ncloud-docs.com/docs/management-management-4-1)를 생성하고, AccessKey/SecretKey를 준비합니다.
AccessKey, SecretKey로 Ncloud Private CA Issuer가 사용할 secret을 생성합니다.

```shell
kubectl create secret generic ncloud-secret --from-literal=NCLOUD_ACCESS_KEY=ACCESSKEYIDACCESSKEY --from-literal=NCLOUD_SECRET_KEY=SECRETACCESSKEYSECRETACCESSKEYSECRETACCE -n cert-manager
```

### Issuer 설정

cert-manager는 Ncloud Private CA를 사용하여 `NcloudPCAIssuer` (네임스페이스 내부) 와 `NcloudPCAClusterIssuer` (클러스터 전체)를 등록 할 수 있습니다.

아래의 샘플 설정을 참고하여주십시오.

```yaml
# ncloudpcaissuer-sample.yaml
apiVersion: privateca-issuer.ncloud.com/v1alpha1
kind: NcloudPCAIssuer
metadata:
  name: issuer-sample
  namespace: cert-manager
spec:
  # Private CA에 생성된 Root CA TAG
  caTag: 12345678-abcdefg
  # 민간클라우드 : https://pca.apigw.ntruss.com, 공공클라우드 : https://privateca.apigw.gov-ntruss.com
  ncloudApiGw: https://pca.apigw.ntruss.com
  # AccessKey, SecretKey가 저장된 Secret
  secretRef:
    name: ncloud-secret
```

```shell
kubectl apply -f ncloudpcaissuer-sample.yaml
```

또는

```yaml
# ncloudpcaclusterissuer-sample.yaml
apiVersion: privateca-issuer.ncloud.com/v1alpha1
kind: NcloudPCAClusterIssuer
metadata:
  name: clusterissuer-sample
spec:
  # Private CA에 생성된 Root CA TAG
  caTag: 12345678-abcdefg
  # 민간클라우드 : https://pca.apigw.ntruss.com, 공공클라우드 : https://privateca.apigw.gov-ntruss.com
  ncloudApiGw: https://pca.apigw.ntruss.com
  # AccessKey, SecretKey가 저장된 Secret
  secretRef:
    name: ncloud-secret
    namespace: cert-manager
```

```shell
kubectl apply -f ncloudpcaclusterissuer-sample.yaml
```

### 첫 인증서 생성
인증서를 생성하기 위해서는 위 단계에서 `NcludPCAIssuer`나 `NcloudPCACluserIssuer`가 생성되어 있어야 합니다.

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: demo-certificate
  namespace: cert-manager
spec:
  # 인증서가 저장될 Secret Name
  secretName: demo-cert-tls
  # Common Name
  commonName: cert-manager.io.demo
  # DNS SAN
  dnsNames:
    - cert-manager.io
    - jetstack.io
  # 인증서 duration
  duration: 24h0m0s
  # 인증서가 만료되기 8시간 전에 갱신
  renewBefore: 8h0m0s
  # Issuer나 ClusterIssuer가 먼저 설정되어 있어야 합니다.
  issuerRef:
    group: privateca-issuer.ncloud.com
    kind: NcloudPCAIssuer #또는 NcloudPCAClusterIssuer
    name: issuer-sample
```

```shell
kubectl apply -f demo-certificate.yaml
```

곧 클러스터내에 인증서가 생성되고 사용 가능합니다.

```shell
kubectl get certificates,secret
NAME                                           READY   SECRET         AGE
certificate.cert-manager.io/demo-certificate   True    demo-cert-tls  1m

NAME                                     TYPE                                  DATA   AGE
secret/demo-cert-tls                     kubernetes.io/tls                     3      1m
```
