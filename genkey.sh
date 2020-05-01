keysize=${1:-1024}
keyfile=${2:-devkey}

ssh-keygen -t rsa -P "" -b $keysize -m PEM -f $keyfile
openssl rsa -in $keyfile -pubout -outform PEM -out $keyfile.pub
