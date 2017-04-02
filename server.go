package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	robot *Robot
	srv   http.Server
}

// NewServer prepares a new server
func NewServer(robot *Robot, port int) *Server {
	s := &Server{
		robot: robot,
		srv: http.Server{
			Addr: fmt.Sprintf("0.0.0.0:%d", port),
		},
	}
	s.srv.Handler = s.createHandler()
	return s
}

// ListenAndServe starts listening on the desired port.
func (s *Server) ListenAndServe() error {
	err := s.srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return maskAny(err)
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown() error {
	return maskAny(s.srv.Shutdown(context.Background()))
}

func (s *Server) createHandler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/crossdomain.xml", s.crossDomainHandler)
	mux.HandleFunc("/poll", s.pollHandler)
	mux.HandleFunc("/reset_all", s.resetAllHandler)
	mux.HandleFunc("/motor/", s.motorHandler)

	return mux
}

const policy = `<cross-domain-policy>
  <allow-access-from-domain="*" to-port="8089"/>
</cross-domain-policy>\0`

func (s *Server) crossDomainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("crossDomain")
	w.Write([]byte(policy))
}

func (s *Server) pollHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) resetAllHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("reset_all")
	s.robot.MotorStop()
}

func (s *Server) motorHandler(w http.ResponseWriter, r *http.Request) {
	cmd := strings.Split(strings.TrimPrefix(r.URL.Path, "/motor/"), "/")
	if len(cmd) >= 1 {
		action := cmd[0]
		switch action {
		case "stop":
			log.Println("Stop")
			s.robot.MotorStop()
		case "vooruit":
			log.Println("Vooruit")
			s.robot.MotorForward()
		case "achteruit":
			log.Println("Achteruit")
			s.robot.MotorBackward()
		case "links":
			log.Println("Links")
			s.robot.MotorTurnLeft()
		case "rechts":
			log.Println("Rechts")
			s.robot.MotorTurnRight()
		}
	}
}
