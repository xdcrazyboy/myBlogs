# 安装

# Building Your First Network


## 创建通道配置事务
接下来，我们需要创建通道事务工件。请务必替换`$CHANNEL_NAME`或将`CHANNEL_NAME`设置为可在整个说明中使用的环境变量：
``` shell
# The channel.tx artifact contains the definitions for our sample channel

export CHANNEL_NAME=mychannel  && ../bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
```

接下来，我们将在我们构建的通道上为Org1定义锚点对等体。同样，请务必替换`$CHANNEL_NAME`或为以下命令设置环境变量。终端输出将模仿通道事务工件的输出：
``` shell
$ ../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
```

## 启动网络
> 如果您之前运行过上面的byfn.sh示例，请确保在继续之前已关闭测试网络（请参阅关闭网络）。

我们将利用脚本来启动我们的网络。docker-compose文件引用我们先前下载的image，并使用我们之前生成的`genesis.block` bootstraps the orderer。

我们希望手动完成命令，以便公开每个调用的语法和功能。

首先让我们开始我们的网络：
```shell
docker-compose -f docker-compose-cli.yaml up -d
```
>如果要查看网络的实时日志，请不要提供`-d`标志。如果您让日志流，那么您将需要打开第二个终端来执行CLI调用。

### 环境变量
要使以下针对`peer0.org1.example.com`的CLI命令起作用，我们需要在命令前面加上下面给出的四个环境变量。`peer0.org1.example.com`的这些变量被`baked into`到CLI容器中，因此我们可以在不传递它们的情况下进行操作。但是，如果要将呼叫发送到其他对等方或订货方，则可以通过在启动容器之前编辑`docker-compose-base.yaml`来相应地提供这些值。修改以下四个环境变量以使用不同的对等方和组织。
``` shell
# Environment variables for PEER0

CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
CORE_PEER_ADDRESS=peer0.org1.example.com:7051
CORE_PEER_LOCALMSPID="Org1MSP"
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
```

### 创建和加入频道
回想一下，我们使用上面的`Create a Channel Configuration Transaction`部分中的`configtxgen`工具创建了通道配置事务。您可以使用`configtx.yaml`中传递给configtxgen工具的相同或不同配置文件重复该过程以创建其他通道配置事务。然后，您可以重复本节中定义的过程，以在您的网络中建立其他通道。

我们将使用docker exec命令进入CLI容器:
```
docker exec -it cli bash
```

您不希望对默认对等`peer0.org1.example.com`运行CLI命令，在四个环境变量中替换`peer0`或`org1`的值并运行命令：
``` shell
# Environment variables for PEER0
export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=peer0.org1.example.com:7051
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
```

接下来，我们将作为创建通道请求的一部分，将我们在“创建通道配置事务”部分（我们称之为channel.tx）中创建的生成的通道配置事务工件传递给订购者。

我们使用`-c`标志指定通道名称，使用`-f`标志指定通道配置事务。在这种情况下，它是`channel.tx`，但是您可以使用其他名称装入自己的配置事务。我们将再次在CLI容器中设置`CHANNEL_NAME`环境变量，以便我们不必显式传递此参数。通道名称必须全部小写，长度小于250个字符，并且与正则表达式`[a-z] [a-z0-9 .-] *`匹配。

``` shell
export CHANNEL_NAME=mychannel

# the channel.tx file is mounted in the channel-artifacts directory within your CLI container
# as a result, we pass the full path for the file
# we also pass the path for the orderer ca-cert in order to verify the TLS handshake
# be sure to export or replace the $CHANNEL_NAME variable appropriately

peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```
>注意我们作为此命令的一部分传递的`--cafile`。它是orderer的根证书的本地路径，允许我们验证TLS握手。

此命令返回一个创世块 -  `<channel-ID.block> ` - 我们将用它来加入频道。它包含`channel.tx中`指定的配置信息如果您没有对默认通道名称进行任何修改，那么该命令将返回一个名为`mychannel.block`的原型。

>对于其余的这些手动命令，您将保留在CLI容器中。在定位`peer0.org1.example.com`以外的对等方时，您还必须记住在所有命令前加上相应的环境变量。

现在让我们将`peer0.org1.example.com`加入频道。
```shell
# By default, this joins ``peer0.org1.example.com`` only
# the <channel-ID.block> was returned by the previous command
# if you have not modified the channel name, you will join with mychannel.block
# if you have created a different channel name, then pass in the appropriately named block

 peer channel join -b mychannel.block
```

您可以根据需要通过对我们在上面的“环境变量”部分中使用的四个环境变量进行适当更改来使其他对等方加入通道。

我们将加入`peer0.org2.example.com`，而不是加入每个对等体，以便我们可以正确更新频道中的锚点对等体定义。由于我们将覆盖CLI容器中的默认环境变量，因此完整命令如下：
```shell
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp CORE_PEER_ADDRESS=peer0.org2.example.com:7051 CORE_PEER_LOCALMSPID="Org2MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt peer channel join -b mychannel.block
```

或者，您可以选择单独设置这些环境变量，而不是传入整个字符串。一旦设置完毕，您只需再次发出`peer channel join`命令，CLI容器将代表`peer0.org2.example.com`。

## 更新锚点Peers
以下命令是通道更新，它们将传播到通道的定义。实质上，我们在通道的创世块之上添加了额外的配置信息。请注意，我们不是修改genesis块，而是简单地将增量添加到将定义锚点对等的链中。

更新通道定义以将Org1的锚点对等体定义为`peer0.org1.example.com`：
```shell
peer channel update -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/Org1MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

```
现在更新通道定义以将Org2的锚点对等体定义为`peer0.org2.example.com`。与Org2对等体的`peer channel join`命令相同，我们需要在此调用前加上适当的环境变量。

```shell
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp CORE_PEER_ADDRESS=peer0.org2.example.com:7051 CORE_PEER_LOCALMSPID="Org2MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt peer channel update -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/Org2MSPanchors.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

```
## 安装和实例化Chaincode
>我们将使用一个简单的现有链码。要了解如何编写自己的链代码，请参阅Chaincode for Developers教程。

应用程序通过链代码与区块链分类帐进行交互。因此，我们需要在每个将执行和支持我们的事务的对等体上安装链代码，然后在通道上实例化链代码。

首先，将示例Go，Node.js或Java链代码安装到Org1中的peer0节点上。这些命令将指定的源代码味道放在我们的对等文件系统上。
>您只能为每个链代码名称和版本安装一个版本的源代码。源代码存在于对等体的文件系统中，在链代码名称和版本的上下文中;它与语言无关。类似地，实例化的链代码容器将反映对等体上安装的任何语言。

**Java**
``` shell
# make note of the -l flag to indicate "java" chaincode
# for java chaincode -p takes the absolute path to the java chaincode
peer chaincode install -n mycc -v 1.0 -l java -p /opt/gopath/src/github.com/chaincode/chaincode_example02/java/

```

**Golang**
``` shell
# this installs the Go chaincode. For go chaincode -p takes the relative path from $GOPATH/src
peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/chaincode_example02/go/
```

**Node.js**
``` shell
# this installs the Node.js chaincode
# make note of the -l flag to indicate "node" chaincode
# for node chaincode -p takes the absolute path to the node.js chaincode
peer chaincode install -n mycc -v 1.0 -l node -p /opt/gopath/src/github.com/chaincode/chaincode_example02/node/
```
当我们在频道上实例化链代码时，认可政策将被设置为要求来自Org1和Org2中的对等方的认可。因此，我们还需要在Org2中的对等端上安装链代码。

修改以下四个环境变量，以便在Org2中针对peer0发出install命令：
``` shell
# Environment variables for PEER0 in Org2

CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
CORE_PEER_ADDRESS=peer0.org2.example.com:7051
CORE_PEER_LOCALMSPID="Org2MSP"
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
```

现在将示例Go，Node.js或Java链代码安装到Org2中的peer0上。这些命令将指定的源代码味道放在我们的对等文件系统上。

**Java**
```shell
# make note of the -l flag to indicate "java" chaincode
# for java chaincode -p takes the absolute path to the java chaincode
peer chaincode install -n mycc -v 1.0 -l java -p /opt/gopath/src/github.com/chaincode/chaincode_example02/java/
```

接下来，在通道上实例化链码。这将初始化通道上的链代码，设置链代码的认可策略，并为目标对等方启动链代码容器。记下`-P`参数。这是我们的政策，我们在此政策中指定针对要验证的此链码的交易所需的认可级别。

在下面的命令中，您会注意到我们将策略指定为`-P“AND（'Org1MSP.peer'，'Org2MSP.peer'）”`。这意味着我们需要来自属于Org1 AND Org2的对等方的“认可”（即两个认可）。如果我们将语法更改为`OR`，那么我们只需要一个认可。

**Java**
>请注意，Java链代码实例化可能需要一些时间，因为它编译链代码并使用java环境下载docker容器。
```shell
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc -l java -v 1.0 -c '{"Args":["init","a", "100", "b","200"]}' -P "AND ('Org1MSP.peer','Org2MSP.peer')"

```
有关策略实施的更多详细信息，请参阅认可政策文档。

如果您希望其他对等方与分类帐进行交互，则需要将它们连接到通道，并将链代码源的相同名称，版本和语言安装到相应的对等文件系统上。一旦他们尝试与特定的链代码进行交互，就会为每个对等体启动一个链代码容器。再次，要认识到Node.js图像的编译速度会慢一些。

一旦在通道上实例化了链代码，我们就可以放弃`l`标志。我们只需传递频道标识符和链码的名称。

### Query
让我们查询a的值，以确保链代码被正确实例化并填充状态DB。查询语法如下：
```shell
# be sure to set the -C and -n flags appropriately

peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","a"]}'
```

### Invoke
现在让我们从a到b移动10。此事务将剪切新块并更新状态DB。调用的语法如下：
```shell
# be sure to set the -C and -n flags appropriately

peer chaincode invoke -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n mycc --peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["invoke","a","b","10"]}'
```
### Query
让我们确认我们之前的调用是否正确执行。我们使用值100初始化了键值a，并在之前的调用中删除了10。因此，对a的查询应返回90.查询的语法如下。
```shell
# be sure to set the -C and -n flags appropriately

peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","a"]}'
```
我们应该看到以下内容：
```shell
Query Result: 90
```
随意重新开始并操纵键值对和后续调用。

### Install
现在我们将在Org2中的第三个对等体peer1上安装链代码。修改以下四个环境变量，以便在Org2中针对peer1发出install命令：
```shell
# Environment variables for PEER1 in Org2

CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
CORE_PEER_ADDRESS=peer1.org2.example.com:7051
CORE_PEER_LOCALMSPID="Org2MSP"
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/ca.crt
```

现在将示例Go，Node.js或Java链代码安装到Org2中的peer1上。这些命令将指定的源代码味道放在我们的对等文件系统上。

**Java**
```shell
# make note of the -l flag to indicate "java" chaincode
# for java chaincode -p takes the absolute path to the java chaincode
peer chaincode install -n mycc -v 1.0 -l java -p /opt/gopath/src/github.com/chaincode/chaincode_example02/java/
```

### Query
让我们确认我们可以在Org2中向Peer1发出查询。我们使用值100初始化了键值a，并在之前的调用中删除了10。因此，对a的查询仍应返回90。

Org2中的peer1必须首先加入通道才能响应查询。可以通过发出以下命令来连接通道：
```shell
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp CORE_PEER_ADDRESS=peer1.org2.example.com:7051 CORE_PEER_LOCALMSPID="Org2MSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/ca.crt peer channel join -b mychannel.block

```
返回join命令后，可以发出查询。查询的语法如下。
```shell
# be sure to set the -C and -n flags appropriately

peer chaincode query -C $CHANNEL_NAME -n mycc -c '{"Args":["query","a"]}'
```
我们应该看到以下内容：
```
Query Result: 90
```
随意重新开始并操纵键值对和后续调用。

### 这个命令背后发生了什么？
>这些步骤描述了`script.sh`由'./byfn.sh up'运行的场景。使用`./byfn.sh`清理网络并确保此命令处于活动状态。然后使用相同的docker-compose提示再次启动您的网络.

1. 脚本 -  `script.sh`  - 在CLI容器中运行。该脚本根据提供的通道名称驱动`createChannel`命令，并使用channel.tx文件进行通道配置。
2. `createChannel`的输出是一个创世块 -  `<your_channel_name> .block`  - 存储在对等体的文件系统中，包含channel.tx指定的通道配置。
3. 对所有四个对等体执行`joinChannel`命令，该命令将先前生成的生成块作为输入。此命令指示对等方加入`<your_channel_name>`并创建以`<your_channel_name> .block`开头的链。
4. 现在我们有一个由四个同行和两个组织组成的频道。这是我们的`TwoOrgsChannel`个人资料。
5. `peer0.org1.example.com`和`peer1.org1.example.com`属于Org1;
   `peer0.org2.example.com`和`peer1.org2.example.com`属于Org2
6. 这些关系是通过`crypto-config.yaml`定义的，MSP路径是在我们的docker compose中指定的。
7. 然后更新Org1MSP（`peer0.org1.example.com`）和Org2MSP（`peer0.org2.example.com`）的锚点对等体。我们通过将`Org1MSPanchors.tx`和`Org2MSPanchors.tx`工件与我们的通道名称一起传递给订购服务来完成此操作。
8. 链码 -  `chaincode_example02`  - 安装在`peer0.org1.example.com`和`peer0.org2.example.com`上。
9. 然后在`mychannel上`“实例化”链代码。实例化将链代码添加到通道，启动目标对等方的容器，并初始化与链代码关联的键值对。该示例的初始值是`[“a”，“100”“b”，“200”]`。这个“实例化”导致一个名为`dev-peer0.org2.example.com-mycc-1.0`的容器启动。
10. 实例化也传递了背书政策的论据。该策略定义为`-P“AND（'Org1MSP.peer'，'Org2MSP.peer'）”`，表示任何事务必须由与Org1和Org2相关联的对等方签署。
11. 向`peer0.org2.example.com`发出针对“a”值的查询。在实例化链代码时，启动了名为`dev-peer0.org2.example.com-mycc-1.0`的Org2 peer0的容器。返回查询的结果。没有发生写入操作，因此对“a”的查询仍将返回值“100”。
12. 将调用发送到`peer0.org1.example.com`和`peer0.org2.example.com`以将“10”从“a”移动到“b”。
13. 查询将发送到`peer0.org2.example.com`以获取“a”的值。返回值90，正确反映上一个事务，在此期间，键“a”的值被修改为10。
14. chaincode  -  `chaincode_example02`  - 安装在`peer1.org2.example.com`上。
15. 查询将发送到`peer1.org2.example.com`以获取“a”的值。这将启动名为`dev-peer1.org2.example.com-mycc-1.0`的第三个链代码容器。返回值90，正确反映上一个事务，在此期间，键“a”的值被修改为10。

### 这表明了什么？
Chaincode必须安装在对等体上，以便它能够成功地对分类帐执行读/写操作。此外，在针对该链代码执行init或传统事务（读/写）之前，不会为对等体启动链代码容器（例如，查询“a”的值）。该事务导致容器启动。此外，通道中的所有对等体都保持分类帐的精确副本，其包括用于以块的形式存储不可变的有序记录的区块链，以及用于维护当前状态的快照的状态数据库。这包括那些没有安装链代码的对等体（如上例中的peer1.org1.example.com）。最后，链代码在安装后可以访问（如上例中的peer1.org2.example.com），因为它已经被实例化了。

### 我如何查看这些交易？

检查CLI Docker容器的日志。
``` 
docker logs -f cli
```
### 如何查看链码日志？
检查各个链代码容器，以查看针对每个容器执行的单独事务。以下是每个容器的组合输出：
```s
$ docker logs dev-peer0.org2.example.com-mycc-1.0
04:30:45.947 [BCCSP_FACTORY] DEBU : Initialize BCCSP [SW]
ex02 Init
Aval = 100, Bval = 200

$ docker logs dev-peer0.org1.example.com-mycc-1.0
04:31:10.569 [BCCSP_FACTORY] DEBU : Initialize BCCSP [SW]
ex02 Invoke
Query Response:{"Name":"a","Amount":"100"}
ex02 Invoke
Aval = 90, Bval = 210

$ docker logs dev-peer1.org2.example.com-mycc-1.0
04:31:30.420 [BCCSP_FACTORY] DEBU : Initialize BCCSP [SW]
ex02 Invoke
Query Response:{"Name":"a","Amount":"90"}
```

### 了解Docker Compose拓扑
BYFN示例为我们提供了两种Docker Compose文件，这两种文件都是从`docker-compose-base.yaml`（位于`base `文件夹中）扩展而来的。我们的第一个版本`docker-compose-cli.yaml`为我们提供了一个CLI容器，以及一个订购者，四个同行。我们将此文件用于此页面上的所有说明。
>本节的其余部分介绍了为SDK设计的docker-compose文件。有关运行这些测试的详细信息，请参阅[Node SDK repo](https://github.com/hyperledger/fabric-sdk-node)。

第二种风格`docker-compose-e2e.yaml`构建为使用Node.js SDK运行端到端测试。除了使用SDK之外，它的主要区别在于Fabric-ca服务器还有容器。因此，我们可以向组织CA发送REST调用以进行用户注册和注册。

如果你想在没有先运行`byfn.sh`脚本的情况下使用`docker-compose-e2e.yaml`，那么我们需要进行四次略微修改。我们需要指向组织CA的私钥。您可以在`crypto-config`文件夹中找到这些值。例如，要找到Org1的私钥，我们将遵循此路径 -  `crypto-config / peerOrganizations / org1.example.com / ca /​​`。私钥是一个长哈希值，后跟_sk。Org2的路径是 -  `crypto-config / peerOrganizations / org2.example.com / ca /`​​。

在`docker-compose-e2e.yaml`中更新ca0和ca1的`FABRIC_CA_SERVER_TLS_KEYFILE`变量。您还需要编辑命令中提供的路径以启动ca服务器。您为每个CA容器提供两次相同的私钥。

### 使用CouchDB
状态数据库可以从默认（goleveldb）切换到CouchDB。CouchDB提供了相同的链代码功能，但是，根据链代码数据被建模为JSON，还可以根据状态数据库数据内容执行丰富而复杂的查询。

要使用CouchDB而不是默认数据库（goleveldb），请按照前面概述的相同步骤生成工件，除非在启动网络时也通过`docker-compose-couch.yaml`：
```
docker-compose -f docker-compose-cli.yaml -f docker-compose-couch.yaml up -d
```
`chaincode_example02`现在应该使用下面的CouchDB。
>Note：如果您选择将fabric-couchdb容器端口映射到主机端口，请确保您了解安全隐患。在开发环境中映射端口使CouchDB REST API可用，并允许通过CouchDB Web界面（Fauxton）可视化数据库。生产环境可能会避免实施端口映射，以限制对CouchDB容器的外部访问。

您可以使用上面列出的步骤对CouchDB状态数据库使用**chaincode_example02**链代码，但是为了运用CouchDB查询功能，您需要使用具有建模为JSON的数据的链代码（例如marbles02）。您可以在`fabric/examples/chaincode/go`目录中找到**marbles02**链代码。

我们将按照相同的流程创建和加入频道，如上面的“创建和加入频道”部分所述。将对等方加入频道后，请使用以下步骤与**marbles02**链代码进行交互：
- 在`peer0.org1.example.com`上安装并实例化链代码：
```shell
# be sure to modify the $CHANNEL_NAME variable accordingly for the instantiate command

peer chaincode install -n marbles -v 1.0 -p github.com/chaincode/marbles02/go
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n marbles -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org0MSP.peer','Org1MSP.peer')"
```
- 制作一些marbles并移动它们：
```
# be sure to modify the $CHANNEL_NAME variable accordingly

peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n marbles -c '{"Args":["initMarble","marble1","blue","35","tom"]}'
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n marbles -c '{"Args":["initMarble","marble2","red","50","tom"]}'
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n marbles -c '{"Args":["initMarble","marble3","blue","70","tom"]}'
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n marbles -c '{"Args":["transferMarble","marble2","jerry"]}'
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n marbles -c '{"Args":["transferMarblesBasedOnColor","blue","jerry"]}'
peer chaincode invoke -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C $CHANNEL_NAME -n marbles -c '{"Args":["delete","marble1"]}'
```

- 如果您选择在docker-compose中映射CouchDB端口，现在可以通过打开浏览器并导航到以下URL，通过CouchDB Web界面（Fauxton）查看状态数据库：`http://localhost:5984/_utils`

您应该看到名为mychannel（或您的唯一通道名称）的数据库及其中的文档。
>对于以下命令，请确保适当更新$ CHANNEL_NAME变量。
您可以从CLI运行常规查询（例如，阅读marble2）：
```
peer chaincode query -C $CHANNEL_NAME -n marbles -c '{"Args":["readMarble","marble2"]}'
```
输出应显示marble2的详细信息：
```
Query Result: {"color":"red","docType":"marble","name":"marble2","owner":"jerry","size":50}

```
您可以检索特定marble的历史记录 - 例如marble1：
```
peer chaincode query -C $CHANNEL_NAME -n marbles -c '{"Args":["getHistoryForMarble","marble1"]}'
```
输出应显示marble1上的事务：
```
Query Result: [{"TxId":"1c3d3caf124c89f91a4c0f353723ac736c58155325f02890adebaa15e16e6464", "Value":{"docType":"marble","name":"marble1","color":"blue","size":35,"owner":"tom"}},{"TxId":"755d55c281889eaeebf405586f9e25d71d36eb3d35420af833a20a2f53a3eefd", "Value":{"docType":"marble","name":"marble1","color":"blue","size":35,"owner":"jerry"}},{"TxId":"819451032d813dde6247f85e56a89262555e04f14788ee33e28b232eef36d98f", "Value":}]

```
您还可以对数据内容执行丰富的查询，例如按所有者jerry查询marble字段：
```
peer chaincode query -C $CHANNEL_NAME -n marbles -c '{"Args":["queryMarblesByOwner","jerry"]}'

```
输出应该显示杰里拥有的两个marble：
```
Query Result: [{"Key":"marble2", "Record":{"color":"red","docType":"marble","name":"marble2","owner":"jerry","size":50}},{"Key":"marble3", "Record":{"color":"blue","docType":"marble","name":"marble3","owner":"jerry","size":70}}]

```

## 为什么选择CouchDB
CouchDB是一种NoSQL解决方案。它是一个面向文档的数据库，其中文档字段存储为键值映射。字段可以是简单的键值对，列表或映射。除了LevelDB支持的键控/复合键/键范围查询之外，CouchDB还支持完全数据丰富的查询功能，例如针对整个区块链数据的非键查询，因为其数据内容以JSON格式存储，完全可查询。因此，CouchDB可以满足LevelDB不支持的许多用例的链代码，审计和报告要求。

CouchDB还可以增强区块链中的合规性和数据保护的安全性。因为它能够通过过滤和屏蔽事务中的各个属性来实现字段级安全性，并且只在需要时授权只读权限。

此外，CouchDB属于CAP定理的AP类型（可用性和分区容差）。它使用具有**最终一致性**的主 - 主复制模型。可以在CouchDB文档的[ Eventual Consistency ](http://docs.couchdb.org/en/latest/intro/consistency.html)页面上找到更多信息。但是，在每个结构对等体下，没有数据库副本，对数据库的写入保证一致且持久（不是最终一致性）。

CouchDB是Fabric的第一个外部可插拔状态数据库，可以而且应该有其他外部数据库选项。例如，IBM为其区块链启用了关系数据库。并且CP类型（一致性和分区容差）数据库也可能需要，以便在没有应用程序级别保证的情况下实现数据一致性。

## 关于数据持久性的注记
如果在对等容器或CouchDB容器上需要数据持久性，则可以选择将docker-host中的目录安装到容器中的相关目录中。例如，您可以在`docker-compose-base.yaml`文件的peer容器规范中添加以下两行：
```
volumes:
 - /var/hyperledger/peer0:/var/hyperledger/production
```
对于CouchDB容器，您可以在CouchDB容器规范中添加以下两行：
```
volumes:
 - /var/hyperledger/couchdb0:/opt/couchdb/data
```

## 故障排除
- 始终保持网络新鲜。使用以下命令删除工件，加密，容器和链代码图像：
    ```
    ./byfn.sh down
    ```
    >如果不删除旧容器和图像，将会看到错误。
- 如果您看到Docker错误，请首先检查您的docker版本[先决条件](https://hyperledger-fabric.readthedocs.io/en/master/prereqs.html)，然后尝试重新启动Docker进程。Docker的问题通常无法立即识别。例如，您可能会看到因无法访问容器中安装的加密材料而导致的错误。

    如果他们坚持删除你的图像并从头开始：
    ```
    docker rm -f $(docker ps -aq)
    docker rmi -f $(docker images -q)
    ```
- 如果在创建，实例化，调用或查询命令上看到错误，请确保已正确更新通道名称和链代码名称。提供的示例命令中有占位符值。
- 如果您看到以下错误：
    ```
    Error: Error endorsing chaincode: rpc error: code = 2 desc = Error installing chaincode code mycc:1.0(chaincode /var/hyperledger/production/chaincodes/mycc.1.0 exits)
    ```
    您可能从之前的运行中获得了链代码图像（例如`dev-peer1.org2.example.com-mycc-1.0`或`dev-peer0.org1.example.com-mycc-1.0`）。删除它们然后再试一次。
    ```
    docker rmi -f $(docker images | grep peer[0-9]-peer[0-9] | awk '{print $3}')

    ```
    
- 如果您看到以下错误：
    ```
    Error connecting: rpc error: code = 14 desc = grpc: RPC failed fast due to transport failure
    Error: rpc error: code = 14 desc = grpc: RPC failed fast due to transport failure
    ```
    确保您正在使用已被重新标记为“最新”的“1.0.0”图像运行您的网络。