# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  config.vm.define "engino" do |engino|
    engino.vm.box = "ubuntu/trusty64"
    engino.vm.network :private_network, ip: "192.168.33.21"
    engino.vm.provision "shell", path: "install.sh"
  end

  config.vm.provider "virtualbox" do |v|
    v.customize ["modifyvm", :id, "--memory", "1024"]
  end

end
