package main

import (
	"fmt"
	"os"
	//"runtime"
	//"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
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
	sendMQTT(c, publishUri+"/usedmem", vmStat.UsedPercent)

	// disk - start from "/" mount point for Linux
	// might have to change for Windows!!
	// don't have a Window to test this out, if detect OS == windows
	// then use "\" instead of "/"
	diskStat, err := disk.Usage("/")
	dealwithErr(err)
	sendMQTT(c, publishUri+"/usedrootfs", diskStat.UsedPercent)

	// cpu - get CPU number of cores and speed
	//cpuStat, err := cpu.Info()
	//dealwithErr(err)
	percentage, err := cpu.Percent(0, true)
	dealwithErr(err)
	sendMQTT(c, publishUri+"/cpuPercentage", percentage)

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	dealwithErr(err)
	sendMQTT(c, publishUri+"/uptime", hostStat.Uptime)

	// get interfaces MAC/hardware address
	interfStat, err := net.Interfaces()
	dealwithErr(err)
	// TODO: iterate interfaces
	sendMQTT(c, publishUri+"/interfaces", interfStat[0].Addrs[0].Addr)

	/*time.Sleep(3 * time.Second)

	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe("fhem/wetter/aussentemp"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}*/

	c.Disconnect(250)
}

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func sendMQTT(client MQTT.Client, uri string, message interface{}) {
	var token MQTT.Token
	switch m := message.(type) {
	case string:
		token = client.Publish(uri, 0, false, m)
	case float64, float32, []float64, []float32:
		token = client.Publish(uri, 0, false, fmt.Sprintf("%f", m))
	case uint, uint64, uint32:
		token = client.Publish(uri, 0, false, fmt.Sprintf("%d", m))
	}
	token.Wait()

}
