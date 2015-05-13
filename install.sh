sudo apt-get update
sudo apt-get -y install unzip

curl -L  https://github.com/coreos/etcd/releases/download/v2.1.0-alpha.0/etcd-v2.1.0-alpha.0-linux-amd64.tar.gz -o /tmp/etcd-v2.1.0-alpha.0-linux-amd64.tar.gz 2>/dev/null 
tar xzvf /tmp/etcd-v2.1.0-alpha.0-linux-amd64.tar.gz -C /opt
cd /opt/etcd-v2.1.0-alpha.0-linux-amd64
./etcd &

curl -L https://dl.bintray.com/mitchellh/consul/0.5.0_linux_amd64.zip -o /tmp/0.5.0_linux_amd64.zip 2>/dev/null
cd /opt
unzip /tmp/0.5.0_linux_amd64.zip 
/opt/consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul > /var/log/consul-agent.log &

sudo apt-get install -y redis-server

curl -L https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz -o /tmp/go1.4.2.linux-amd64.tar.gz 2>/dev/null 
tar zxvf /tmp/go1.4.2.linux-amd64.tar.gz -C /opt

echo export GOPATH=/home/vagrant >> ~/.bashrc
echo export GOROOT=/opt/go >> ~/.bashrc
echo export PATH=$PATH:$GOROOT/bin:$GOPATH/bin >> ~/.bashrc

mkdir ~/src
cp -R /vagrant/* src/engino/
