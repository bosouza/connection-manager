package main

import (
	"fmt"

	"github.com/souza-bruno/connection-manager/pkg/connmngr"
	"github.com/souza-bruno/connection-manager/pkg/connmngr/cnnchan"
)

func main() {
	cm := connmngr.CreateConnManager()

	myDeviceRecvChan := make(chan string)
	myDeviceSendChan := make(chan string)

	myDeviceCnnFactory := cnnchan.CreateChannelCnnFactory(myDeviceRecvChan, myDeviceSendChan)

	go func() {
		v := <-myDeviceRecvChan
		fmt.Printf("myDevice received %q, sending it back to sender\n", v)
		myDeviceSendChan <- v
	}()

	cm.AddClient("myDevice", myDeviceCnnFactory)

	cnnToMyDevice, err := cm.ConnectTo("myDevice")
	if err != nil {
		fmt.Println("error connecting to myDevice")
	}

	fmt.Println("sending potato to myDevice")
	cnnToMyDevice.Send("potato")
	fmt.Println("reading what myDevice sent back")
	response, err := cnnToMyDevice.Receive()
	if err != nil {
		fmt.Println("error receiving from myDevice")
	}
	fmt.Printf("received %q back from myDevice\n", response)
}
