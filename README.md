# smartcontract
some experiments in hyperledger
## build the network
### Pre-requisite
#### docker
Docker is a container,and the component in hyperledger fabric is simulated by running specific docker.
##### install(by script)
1、download script
$ curl -fsSL get.docker.com -o get-docker.sh

$ ls get*

get-docker.sh

2、install by executing script

$ sudo sh get-docker.sh

3、make user XX can execute docker

$ sudo usermod -aG docker XX

4、run the "hello-world" image 

$ sudo docker run hello-world

If you see " Hello from Docker! This message shows that your installation appears to be working",Congratulation,you succeed!
##### uninstall
sudo apt-get remove --auto-remove docker
##### some common commands about docker

|sss|s|

