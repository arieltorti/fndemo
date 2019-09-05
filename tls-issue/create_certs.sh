mkdir -p .certificates
openssl genrsa -out .certificates/server.key 2048
openssl rsa -in .certificates/server.key -out .certificates/server.key

# We need to create a certificate for each host (in theory we could create a certificate with many alternative names).
openssl req -sha256 -new -key .certificates/server.key -out .certificates/server.csr -subj '/CN=caddy'
openssl x509 -req -sha256 -days 365 -in .certificates/server.csr -signkey .certificates/server.key -out .certificates/server.crt

openssl req -sha256 -new -key .certificates/server.key -out .certificates/fn_api.csr -subj '/CN=fn_api'
openssl x509 -req -sha256 -days 365 -in .certificates/fn_api.csr -signkey .certificates/server.key -out .certificates/fn_api.crt

openssl req -sha256 -new -key .certificates/server.key -out .certificates/fn_runner.csr -subj '/CN=fn_runner'
openssl x509 -req -sha256 -days 365 -in .certificates/fn_runner.csr -signkey .certificates/server.key -out .certificates/fn_runner.crt

openssl req -sha256 -new -key .certificates/server.key -out .certificates/fn_lb.csr -subj '/CN=fn_lb'
openssl x509 -req -sha256 -days 365 -in .certificates/fn_lb.csr -signkey .certificates/server.key -out .certificates/fn_lb.crt

cat .certificates/server.crt .certificates/server.key > .certificates/cert.pem
