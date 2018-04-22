package main

import (
	"fmt"
	"os"
	//"runtime"
	//"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	//	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	//	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	//	"github.com/shirou/gopsutil/net"
)

//define a function for the default message handler
var publishHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	hostname, _ := os.Hostname()
	clientId := "simplemqtt_" + hostname
	publishUri := "simplemqtt/monitoring/" + hostname

	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("tcp://raspi1.fritz.box:1883")
	opts.SetClientID(clientId)
	opts.SetDefaultPublishHandler(publishHandler)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	/*if token := c.Subscribe("fhem/wetter/aussentemp", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}*/

	//////
	//runtimeOS := runtime.GOOS
	// memory
	vmStat, err := mem.VirtualMemory()
	dealwithErr(err)

	text := fmt.Sprintf("%f", vmStat.UsedPercent)
	token := c.Publish(publishUri+"/usedmem", 0, false, text)
	token.Wait()

	// disk - start from "/" mount point for Linux
	// might have to change for Windows!!
	// don't have a Window to test this out, if detect OS == windows
	// then use "\" instead of "/"
	diskStat, err := disk.Usage("/")
	dealwithErr(err)
	text = fmt.Sprintf("%f", diskStat.UsedPercent)
	token = c.Publish(publishUri+"/usedrootfs", 0, false, text)
	token.Wait()

	/*	// cpu - get CPU number of cores and speed
		cpuStat, err := cpu.Info()
		dealwithErr(err)
		percentage, err := cpu.Percent(0, true)
		dealwithErr(err)

		// host or machine kernel, uptime, platform Info
		hostStat, err := host.Info()
		dealwithErr(err)

		// get interfaces MAC/hardware address
		interfStat, err := net.Interfaces()
		dealwithErr(err)

		//////
		for i := 0; i < 5; i++ {
			text := fmt.Sprintf("this is msg #%d!", i)
			token := c.Publish(publishUri, 0, false, text)
			token.Wait()
			fmt.Println("Published to " + publishUri)
		}
	*/
	/*time.Sleep(3 * time.Second)

	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe("fhem/wetter/aussentemp"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}*/

	c.Disconnect(250)
}
