package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

const (
	defaultIP = "192.168.4.1"
)

var (
	robotIP    string // IP address of robot
	robotPort  int    // Port on robot
	robotName  string // Hostname of robot
	serverPort int    // Port on which to listen

	projectVersion = "dev"
	projectBuild   = "dev"
)

func init() {
	flag.StringVar(&robotIP, "ip", "", "IP address of the robot")
	flag.IntVar(&robotPort, "port", 3030, "Port of the robot")
	flag.StringVar(&robotName, "name", "", "MDNS Name of the robot")
	flag.IntVar(&serverPort, "server-port", 8089, "Port on which our server will listen")
}

func main() {
	flag.Parse()

	log.Printf("Starting Auto-Bridge version %s, build %s\n", projectVersion, projectBuild)

	attempt := 0
	for {
		if attempt > 0 {
			log.Println("Waiting a bit...")
			time.Sleep(time.Second * 5)
		}
		attempt++
		ip := robotIP
		if robotName == "" && ip == "" {
			ip = defaultIP
		} else if ip == "" {
			var err error
			ip, err = findRobotIP(robotName)
			if err != nil {
				log.Printf("Kan robot '%s' niet vinden\n", robotName)
				continue
			}
		}

		shutdownServer := make(chan struct{})
		addr := fmt.Sprintf("%s:%d", ip, robotPort)
		log.Printf("Connecting to %s\n", addr)
		robot, err := NewRobot(addr, func(data interface{}) {
			log.Printf("Received error (%v), restarting...\n", data)
			shutdownServer <- struct{}{}
		})
		if err != nil {
			log.Printf("Cannot connect to robot at '%s': %#v\n", addr, err)
			continue
		}

		// Now start the server
		s := NewServer(robot, serverPort)
		log.Printf("Listening on port %d\n", serverPort)
		go func() {
			if err := s.ListenAndServe(); err != nil {
				log.Fatalf("Kan de server niet starten: %#v\n", err)
			}
		}()
		<-shutdownServer
		log.Printf("Shutting down server...\n")
		if err := s.Shutdown(); err != nil {
			log.Printf("Failed to shutdown server: %#v\n", err)
		}
	}
}
