# Smart Contract
some experiments with hyperledger1.4.1 in Ubuntu18.04
# Build the Network
## Pre-requisite
### Docker
Docker is a container,and the component in hyperledger fabric is simulated by running specific docker.
#### Install(by script)
1、download script 
```
$ curl -fsSL get.docker.com -o get-docker.sh
$ ls get*
get-docker.sh 
```
2、install by executing script
```
$ sudo sh get-docker.sh
```
3、make user XX can execute docker
```
$ sudo usermod -aG docker XX
```
4、run the "hello-world" image 
```
$ sudo docker run hello-world
```
If you see " Hello from Docker! This message shows that your installation appears to be working",Congratulation,you succeed!
#### Uninstall
sudo apt-get remove --auto-remove docker
#### Some Common Commands About Docker
```
sudo docker ps -a(check the running dokcer)
sudo docker stop XX(stop the running docker XX)
sudo docker rm XX(rm the running docker XX)
sudo docker image ls(check the images)
sudo docker image rm XX(rm the dokcer image XX)
sudo docker log XX(check the log of the running docker XX)
sudo docker cp local path  docker ID:docker path(copy a file from local machine to docker)
sudo docker cp '/pbc-0.5.14'  be55b27495f3:/home(copy a file from local machine to docker)
sudo docker run -it  a1e3874f338b /bin/bash(into an image of docker)
sudo docker exec -it 775c7c9ee1e1 /bin/bash(into a running docker)
sudo docker commit runningdocker dockeriamge(pack a running docker to docker image)
```
### Docker-compose
Compose is a tool for defining and running multi-container Docker applications. With Compose, you use a YAML file to configure your application’s services. Then, with a single command, you create and start all the services from your configuration. 
#### Install
1、download docker-compose to /usr/local/bin/docker-compose
```
$ sudo curl -L https://github.com/docker/compose/releases/download/1.21.2/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
```
2、allow normal user to execute compose
```
$ sudo chmod +x /usr/local/bin/docker-compose
```
3、verify if succeed
```
$ docker-compose --version
```
### Go(optional)
1、download the tar of golang

2、untar the tar to /usr/local
```
$ sudo tar -C /usr/local -xzf go1.10.3.linux-amd64.tar.gz
```
3、create go directory
```
$ mkdir $HOME/go
```
4、configure the environment variables
```
$ vi ~/.bashrc
```
add following in the file
```
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
To make it effective
```
$ source ~/.bashrc
```
5、To check go
```
go version
```
### Maven
<https://blog.csdn.net/badder2/article/details/89672612> 
## Blockchain-java-sdk
tips：I can not find the way to deploy java chaincode with java sdk of this version.<https://jira.hyperledger.org/browse/FABJ-220>

You can download the blockchain-java-sdk from
https://github.com/IBM/blockchain-application-using-fabric-java-sdk
then follow the step in readme.md to get your first network of blockchain
### My Experience
There are three important folders in the project.

network:to build the network.You can check the docker-compose.yml to see the detail about the docker in the network.

network_resources:something about cryptography,for example,key

java:In the official direct,they coperate java file using command line.If you prefer IDE,you can open java file with Idea and modify the java file according to your needs.In /java/src/main/java/org/example/config/Config.java,you can find the location of chaincode.
# My Experiment1-bgn
I use go chaincode.
I modify the project https://github.com/sachaservan/bgn according to my needs.The bgn is based on pbc.So I need to install pbc first.
## The Detail of Installation
The official direct is  https://godoc.org/github.com/Nik-U/pbc.

In my machine,Ubuntu18.04,I do the following steps.

1、install GMP
```
sudo apt-get install libgmp-dev
```
2、install the dependencies of pbc
```
sudo apt-get install build-essential flex bison
```
3、download PBC library source code

I download PBC library source code in https://crypto.stanford.edu/pbc/download.html

4、compile and install
```
./configure
make
sudo make install
sudo ldconfig(After installing, you may need to rebuild the search path for libraries)
```
5、I download the PBC Go Wrapper(https://github.com/Nik-U/pbc)
## Run Bgn in Chaincode
There is a way for you to import go packages into the chaincode

you can put your packages under src/vendor folder
src/

______chaincode.go

______vendor/

__________[other packages]

But this does not solve my problem completely,the bgn need the support of C language.

I recompile the image--hyperledger/fabric-ccenv, add the library I need.Because the chaincode is compiled in a container which is run base on fabric-ccenv image.
### The Problems I come across
1、pt-get install vim：
E:Unable to locate package vim

We can type "apt-get update"

2、create a new user and delete a user in Ubuntu18.04
https://blog.csdn.net/qq_33867131/article/details/86490083
