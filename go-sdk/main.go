package main

import(
  "fmt"
  "strconv"
  "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
  "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
  "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
  "github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

const (
  channelID      = "mychannel"
  orgName        = "Org1"
  orgAdmin       = "Admin"
  ordererOrgName = "OrdererOrg"
)

const (
  fileConfig = "config_test.yaml"
  create_Channel = false
)

// Initial B values for ExampleCC
const (
  ExampleCCInitB    = "200"
  ExampleCCUpgradeB = "400"
  keyExp            = "key-%s-%s"
)


var (
  ccID = "mycc"
)

// ExampleCC query and transaction arguments
var defaultQueryArgs = [][]byte{[]byte("query"), []byte("b")}
var defaultTxArgs = [][]byte{[]byte("move"), []byte("a"), []byte("b"), []byte("1")}

// ExampleCC init and upgrade args
var initArgs = [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte(ExampleCCInitB)}
var resetArgs = [][]byte{[]byte("a"), []byte("100"), []byte("b"), []byte(ExampleCCInitB)}


//TODO Args
func queryCC(client *channel.Client, targetEndpoints ...string) []byte {
  response, err := client.Query(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: ExampleCCDefaultQueryArgs()},
    channel.WithRetry(retry.DefaultChannelOpts),
    channel.WithTargetEndpoints(targetEndpoints...),
  )
  if err != nil {
    fmt.Println(err)
  }
  return response.Payload
}

func executeCC(client *channel.Client) {
  _, err := client.Execute(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: ExampleCCDefaultTxArgs()},
    channel.WithRetry(retry.DefaultChannelOpts))
  if err != nil {
    fmt.Println(err)
  }
}

// ExampleCCDefaultQueryArgs returns example cc query args
func ExampleCCDefaultQueryArgs() [][]byte {
  return defaultQueryArgs
}

// ExampleCCQueryArgs returns example cc query args
func ExampleCCQueryArgs(key string) [][]byte {
  return [][]byte{[]byte("query"), []byte(key)}
}

// ExampleCCTxArgs returns example cc query args
func ExampleCCTxArgs(from, to, val string) [][]byte {
  return [][]byte{[]byte("move"), []byte(from), []byte(to), []byte(val)}
}

// ExampleCCDefaultTxArgs returns example cc move funds args
func ExampleCCDefaultTxArgs() [][]byte {
  return defaultTxArgs
}



//ExampleCCTxSetArgs sets the given key value in examplecc
func ExampleCCTxSetArgs(key, value string) [][]byte {
  return [][]byte{[]byte("set"), []byte(key), []byte(value)}
}

//ExampleCCInitArgs returns example cc initialization args
func ExampleCCInitArgs() [][]byte {
  return initArgs
}

func main() {

  sdk, err := fabsdk.New(config.FromFile(fileConfig))
  if err != nil {
    fmt.Println(err)
  }
  defer sdk.Close()


  //prepare channel client context using client context
  clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser("User1"), fabsdk.WithOrg(orgName))
  // Channel client is used to query and execute transactions (Org1 is default org)
  client, err := channel.New(clientChannelContext)
  if err != nil {
   fmt.Println(err)
  }

  // Query
  existingValue := queryCC(client)
  valueInt, err := strconv.Atoi(string(existingValue))
  if err != nil {
   fmt.Println(err)
  }else{
    fmt.Println(valueInt)
  }

  //execute transaction

  executeCC(client)


}
