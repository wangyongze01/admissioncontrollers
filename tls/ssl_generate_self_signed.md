# Use CFSSL to generate certificates

More about [CFSSL here]("https://github.com/cloudflare/cfssl")

```

cd kubernetes\admissioncontrollers\introduction

docker run -it --rm -v ${PWD}:/work -w /work debian bash

apt-get update && apt-get install -y curl &&
curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssl_1.5.0_linux_amd64 -o /usr/local/bin/cfssl && \
curl -L https://github.com/cloudflare/cfssl/releases/download/v1.5.0/cfssljson_1.5.0_linux_amd64 -o /usr/local/bin/cfssljson && \
chmod +x /usr/local/bin/cfssl && \
chmod +x /usr/local/bin/cfssljson

#generate ca in ./tls/ca
cfssl gencert -initca ./tls/config/ca-csr.json | cfssljson -bare ./tls/ca/ca

#generate certificate in /tmp
cfssl gencert \
  -ca=./tls/ca/ca.pem \
  -ca-key=./tls/ca/ca-key.pem \
  -config=./tls/config/ca-config.json \
  -profile=default \
  ./tls/config/ca-csr.json | cfssljson -bare ./tls/test-webhook/test-webhook

# -hostname="example-webhook,example-webhook.default.svc.cluster.local,example-webhook.default.svc,localhost,127.0.0.1" \

#make a secret
cat <<EOF > ./tls/test-webhook-tls.yaml
apiVersion: v1
kind: Secret
metadata:
  name: test-webhook-tls
type: Opaque
data:
  tls.crt: $(cat ./tls/test-webhook/test-webhook.pem | base64 | tr -d '\n')
  tls.key: $(cat ./tls/test-webhook/test-webhook-key.pem | base64 | tr -d '\n') 
EOF

#generate CA Bundle + inject into template
ca_pem_b64="$(openssl base64 -A <"./tls/ca/ca.pem")"

sed -e 's@${CA_PEM_B64}@'"$ca_pem_b64"'@g' <"webhook-template.yaml" \
    > webhook.yaml
```