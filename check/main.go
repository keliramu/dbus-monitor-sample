package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

func main() {
	fmt.Println("Starting service status check via dbus...")
	defer fmt.Println("DONE, exiting.")

	// Get a handle on the system bus. There are two types
	// of buses: system and session. The system bus is for
	// handling system-wide operations (like in this case,
	// shutdown). The session bus is a per-user bus.
	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Printf("error getting system bus: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// obtain systemd dbus object
	intf := "org.freedesktop.systemd1"
	path := "/org/freedesktop/systemd1"

	dbusObj, ok := conn.Object(intf, dbus.ObjectPath(path)).(*dbus.Object)
	if !ok {
		panic("Failed obtaining systemd dbus object")
	}

	// check service: snapd.service state
	service := "snapd.service"
	var unitPath dbus.ObjectPath

	// first need to obtain service unit path
	err = dbusObj.Call("org.freedesktop.systemd1.Manager.GetUnit", 0, service).Store(&unitPath)
	if err != nil {
		panic(err)
	}

	fmt.Println("Service[", service, "] path[", unitPath, "]")

	// now obtain service unit dbus object
	dbusObj, ok = conn.Object(intf, unitPath).(*dbus.Object)
	if !ok {
		panic("Failed obtaining service unit dbus object")
	}

	// get all info on dbus object
	node, err := introspect.Call(dbusObj)
	if err != nil {
		panic(err)
	}
	data, _ := json.MarshalIndent(node, "", "    ")
	os.Stdout.Write(data)

	// now retrieve service unit dbus object ptoperty
	prop := "org.freedesktop.systemd1.Unit.ActiveState"
	rez, err := dbusObj.GetProperty(prop)
	if err != nil {
		panic(fmt.Sprintf("Failed obtaining dbus object property[%s] err[%s]", prop, err))
	}
	fmt.Println("property[", prop, "] val:", rez)
}
