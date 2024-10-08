# 1. Generate CA's private key and self-sign certificate
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj="/C=VN/ST=Ho Chi Minh/L=Linh Dong/O=Taddster/OU=Technology/CN=*.taddster.tech/emailAddress=taddster.tech@gmail.com"
# 2. Geneate web server's private key and certificate signing request (CSR)
openssl req  -newkey rsa:4096 -nodes -keyout server-key.pem  -out server-req.pem -subj="/C=VN/ST=Ho Chi Minh/L=Linh Dong/O=PC Book/OU=Computer/CN=*.pcbook.com/emailAddress=pcbook@gmail.com"

# 3. Use CA's private key to sign web web server's CSR and get back the signed certificate.
openssl x509 -req -in server-req.pem -days 60  -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in server-cert.pem -noout -text

# 4. Geneate web client's private key and certificate signing request (CSR)
openssl req  -newkey rsa:4096 -nodes -keyout client-key.pem  -out client-req.pem -subj="/C=VN/ST=Thua Thien Hue/L=Thanh pho Hue/O=PC Client/OU=client/CN=*.pcclient.com/emailAddress=pcclient@gmail.com"

# 3. Use CA's private key to sign web client's CSR and get back the signed certificate.
openssl x509 -req -in client-req.pem -days 60  -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.cnf

echo "Client's signed certificate"
openssl x509 -in client-cert.pem -noout -text
