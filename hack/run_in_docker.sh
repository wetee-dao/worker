# get shell path
SOURCE="$0"
while [ -h "$SOURCE"  ]; do
    DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
    SOURCE="$(readlink "$SOURCE")"
    [[ $SOURCE != /*  ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE"  )" && pwd  )"
cd $DIR/../

sudo chmod 777 /etc/rancher/k3s/k3s.yaml
sudo mkdir /opt/wetee-worker
sudo chmod 777 /opt/wetee-worker

echo $KUBE_CONFIG_PATH

rm -f  bin/manager
cd bin && export KUBECONFIG=/etc/kube/config &&  \
    ego-go build -o manager ../cmd/main.go

cd $DIR/../
docker build -t wetee/worker:dev .

docker run -e KUBECONFIG=/etc/kube/config -v /etc/rancher:/etc/rancher --device /dev/sgx/enclave --device /dev/sgx/provision wetee/worker:2024-03-05-21_45