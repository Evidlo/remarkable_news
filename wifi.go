package main

import (
	"fmt"
	"os"
	"strings"
	// "io/ioutil"

	"github.com/godbus/dbus"
)


func wait_online(x chan int) {
	// wait for wpa_supplicant to emit message that we have just connected to a network

	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// dbus filter for getting WiFi state change from wpa_supplicant
	// see `dbus-monitor --system`
	var rules = []string{
		"interface='org.freedesktop.DBus.Properties',path='/org/freedesktop/network1/link/_34',member='PropertiesChanged'",
	}

	// begin monitoring dbus events
	var flag uint = 0
	call := conn.BusObject().Call("org.freedesktop.DBus.Monitoring.BecomeMonitor", 0, rules, flag)
	if call.Err != nil {
		fmt.Fprintln(os.Stderr, "Failed to become monitor:", call.Err)
		os.Exit(1)
	}

	c := make(chan *dbus.Message, 10)
	conn.Eavesdrop(c)
	for v := range c {

		// FIXME - this is a terrible hack.  I don't know Go

		// this is v.Body:

		// string "org.freedesktop.network1.Link"
		// array [
		// 	dict entry(
		// 		string "OperationalState"
		// 		variant             string "routable"
		// 	)
		// ]
		// array [
		// ]

		// we want to determine if "routable" is present or not

		str := fmt.Sprintf("%v", v.Body)
		if strings.Contains(str, "routable") {
			x <- 0
		}
	}
}

