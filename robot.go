package main

import (
	"log"

	"github.com/juju/errgo"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

type Robot struct {
	adaptor   gobot.Connection
	motorA_BW *gpio.DirectPinDriver
	motorA_FW *gpio.DirectPinDriver
	motorB_BW *gpio.DirectPinDriver
	motorB_FW *gpio.DirectPinDriver
	onError   func(data interface{})
}

// NewRobot connects to a robot at the given address (ip:port)
func NewRobot(addr string, onError func(data interface{})) (*Robot, error) {
	adaptor := firmata.NewTCPAdaptor(addr)
	if err := adaptor.Connect(); err != nil {
		return nil, errgo.WithCausef(nil, err, "Failed to connect to robot: %s", err.Error())
	}

	motorA_BW := gpio.NewDirectPinDriver(adaptor, "5") // D1 A-IA Backward
	if err := motorA_BW.Start(); err != nil {
		return nil, errgo.WithCausef(nil, err, "Cannot start pin 5")
	}
	motorA_FW := gpio.NewDirectPinDriver(adaptor, "4") // D2 A-IB Forward
	if err := motorA_FW.Start(); err != nil {
		return nil, errgo.WithCausef(nil, err, "Cannot start pin 4")
	}
	motorB_BW := gpio.NewDirectPinDriver(adaptor, "0") // D3 B-IA Backward
	if err := motorB_BW.Start(); err != nil {
		return nil, errgo.WithCausef(nil, err, "Cannot start pin 0")
	}
	motorB_FW := gpio.NewDirectPinDriver(adaptor, "2") // D4 B-IB Forward
	if err := motorB_FW.Start(); err != nil {
		return nil, errgo.WithCausef(nil, err, "Cannot start pin 2")
	}

	go func() {
		for e := range adaptor.Subscribe() {
			log.Printf("Event: %s=%v\n", e.Name, e.Data)
			if e.Name == "Error" {
				onError(e.Data)
				log.Printf("Disconnecting...\n")
				adaptor.Disconnect()
			}
		}
	}()

	return &Robot{
		adaptor:   adaptor,
		motorA_BW: motorA_BW,
		motorA_FW: motorA_FW,
		motorB_BW: motorB_BW,
		motorB_FW: motorB_FW,
		onError:   onError,
	}, nil
}

// MotorStop stops both motors
func (r *Robot) MotorStop() error {
	if err := r.setMotor(false, true, false, true); err != nil {
		return maskAny(err)
	}
	return nil
}

// MotorForward drivers both motors forward
func (r *Robot) MotorForward() error {
	if err := r.setMotor(true, true, true, true); err != nil {
		return maskAny(err)
	}
	return nil
}

// MotorBackward drivers both motors backward
func (r *Robot) MotorBackward() error {
	if err := r.setMotor(true, false, true, false); err != nil {
		return maskAny(err)
	}
	return nil
}

// MotorTurnRight drivers both motors to make a right turn
func (r *Robot) MotorTurnRight() error {
	if err := r.setMotor(true, false, true, true); err != nil {
		return maskAny(err)
	}
	return nil
}

// MotorTurnLeft drivers both motors to make a left turn
func (r *Robot) MotorTurnLeft() error {
	if err := r.setMotor(true, true, true, false); err != nil {
		return maskAny(err)
	}
	return nil
}

func (r *Robot) setMotor(activeA, forwardA, activeB, forwardB bool) error {
	on := func(active, forward, isForward bool) byte {
		if !active {
			return 0
		}
		if forward == isForward {
			return 1
		}
		return 0
	}

	if err := r.motorA_BW.DigitalWrite(on(activeA, forwardA, false)); err != nil {
		return maskAny(err)
	}
	if err := r.motorA_FW.DigitalWrite(on(activeA, forwardA, true)); err != nil {
		return maskAny(err)
	}
	if err := r.motorB_BW.DigitalWrite(on(activeB, forwardB, false)); err != nil {
		return maskAny(err)
	}
	if err := r.motorB_FW.DigitalWrite(on(activeB, forwardB, true)); err != nil {
		return maskAny(err)
	}
	return nil
}
