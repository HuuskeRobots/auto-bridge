package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

const (
	defaultIP = "192.168.4.1"
)

var (
	robotIP   string // IP address of robot
	robotPort int    // Port on robot
	robotName string // Hostname of robot

	robot     *gobot.Robot
	motorA_BW *gpio.DirectPinDriver
	motorA_FW *gpio.DirectPinDriver
	motorB_BW *gpio.DirectPinDriver
	motorB_FW *gpio.DirectPinDriver
)

func init() {
	flag.StringVar(&robotIP, "ip", "", "IP address of the robot")
	flag.IntVar(&robotPort, "port", 3030, "IP address of the robot")
	flag.StringVar(&robotName, "name", "", "MDNS Name of the robot")
}

func main() {
	flag.Parse()

	if robotName == "" && robotIP == "" {
		robotIP = defaultIP
	} else if robotIP == "" {
		var err error
		robotIP, err = findRobotIP(robotName)
		if err != nil {
			log.Fatalf("Kan robot '%s' niet vinden\n", robotName)
		}
	}

	addr := fmt.Sprintf("%s:%d", robotIP, robotPort)
	log.Printf("Connecting to %s\n", addr)
	firmataAdaptor := firmata.NewTCPAdaptor(addr)
	n := func(name string, d gobot.Device) gobot.Device {
		d.SetName(name)
		return d
	}

	go func() {
		for e := range firmataAdaptor.Subscribe() {
			log.Printf("Event: %s=%v\n", e.Name, e.Data)
		}
	}()

	work := func() {
		initHttpServer()
		http.ListenAndServe("0.0.0.0:8089", nil)
		/*gobot.Every(1*time.Second, func() {
			log.Println("Toggle")
			led.Toggle()
		})*/
	}

	motorA_BW = gpio.NewDirectPinDriver(firmataAdaptor, "5") // D1 A-IA Backward
	motorA_FW = gpio.NewDirectPinDriver(firmataAdaptor, "4") // D2 A-IB Forward
	motorB_BW = gpio.NewDirectPinDriver(firmataAdaptor, "0") // D3 B-IA Backward
	motorB_FW = gpio.NewDirectPinDriver(firmataAdaptor, "2") // D4 B-IB Forward

	robot = gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{
			n("m-a-ia", motorA_BW), // D1
			n("m-a-ib", motorA_FW), // D2
			n("m-b-ia", motorB_BW), // D3
			n("m-b-ib", motorB_FW), // D4
		},
		work,
	)

	robot.Start()
}

func initHttpServer() {
	http.HandleFunc("/crossdomain.xml", crossDomainHandler)
	http.HandleFunc("/poll", pollHandler)
	http.HandleFunc("/reset_all", resetAllHandler)
	http.HandleFunc("/motor/", motorHandler)
}

const policy = `<cross-domain-policy>
  <allow-access-from-domain="*" to-port="8089"/>
</cross-domain-policy>\0`

func crossDomainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("crossDomain")
	w.Write([]byte(policy))
}

func pollHandler(w http.ResponseWriter, r *http.Request) {
	var result []string
	/*	for _, d := range *robot.Devices() {
		if pin, ok := d.(*gpio.DirectPinDriver); ok {
			pin.DigitalRead
			val, _ := pin.DigitalRead()
			result = append(result, fmt.Sprintf("%s %d", pin.Name(), val))
		}
	}*/
	w.Write([]byte(strings.Join(result, "\n")))
}

func resetAllHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("reset_all")
	motorA_BW.Off()
	motorA_FW.Off()
	motorB_BW.Off()
	motorB_FW.Off()
}

func motorHandler(w http.ResponseWriter, r *http.Request) {
	cmd := strings.Split(strings.TrimPrefix(r.URL.Path, "/motor/"), "/")
	if len(cmd) >= 1 {
		action := cmd[0]
		switch action {
		case "stop":
			log.Println("Stop")
			motorA_BW.Off()
			motorA_FW.Off()
			motorB_BW.Off()
			motorB_FW.Off()
		case "vooruit":
			log.Println("Vooruit")
			motorA_BW.Off()
			motorA_FW.On()
			motorB_BW.Off()
			motorB_FW.On()
		case "achteruit":
			log.Println("Achteruit")
			motorA_BW.On()
			motorA_FW.Off()
			motorB_BW.On()
			motorB_FW.Off()
		case "links":
			log.Println("Links")
			motorA_BW.On()
			motorA_FW.Off()
			motorB_BW.Off()
			motorB_FW.On()
		case "rechts":
			log.Println("Rechts")
			motorA_BW.Off()
			motorA_FW.On()
			motorB_BW.On()
			motorB_FW.Off()
		}
	}
}
