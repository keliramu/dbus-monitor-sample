package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
)

func main() {
	fmt.Println("Starting dbus sample")
	defer fmt.Println("Finishing dbus sample")

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

	// Call AddMatch so that this process will be notified for
	// the PrepareForShutdown signal. This will allow us to do
	// custom logic when the machine is getting ready to shutdown.
	// err = conn.AddMatchSignal(
	// 	dbus.WithMatchInterface("org.freedesktop.login1.Manager"),
	// 	dbus.WithMatchObjectPath("/org/freedesktop/login1"),
	// 	dbus.WithMatchMember("PrepareForShutdown"),
	// )
	// if err != nil {
	// 	fmt.Printf("error adding match signal: %v\n", err)
	// 	os.Exit(1)
	// }

	addMatch(conn, "org.freedesktop.login1.Manager", "/org/freedesktop/login1", "PrepareForShutdown")
	addMatch(conn, "org.freedesktop.login1.Manager", "/org/freedesktop/login1", "PowerOff")
	addMatch(conn, "org.freedesktop.login1.Manager", "/org/freedesktop/login1", "Reboot")
	addMatch(conn, "org.freedesktop.systemd1.Manager", "/org/freedesktop/systemd1", "JobNew") // expect: string "reboot.target", string "halt.target"

	fmt.Println("Waiting for shutdown signal")

	// AddMatch is already called, but we need to setup a signal
	// handler, which is just a channel.
	shutdownSignal := make(chan *dbus.Signal, 1)
	conn.Signal(shutdownSignal)
	for signal := range shutdownSignal {

		fmt.Printf("~~~ Signal: %+v\n", signal)
		for i, b := range signal.Body {
			fmt.Printf("~~~ Signal.Body[%d]: %+v\n", i, b)
		}
		fmt.Printf("~~~ Signal: Sender[%s] Path[%s] Name[%s]\n", signal.Sender, signal.Path, signal.Name)

		if isShutdown(signal) {
			fmt.Println("~~~ SHUTDOWN DETECTED !!! ~~~")
		}
	}
}

func addMatch(conn *dbus.Conn, intf, path, membr string) {
	fmt.Printf("~~~ add match signal: intf[%s] path[%s] membr[%s]\n", intf, path, membr)
	err := conn.AddMatchSignal(
		dbus.WithMatchInterface(intf),
		dbus.WithMatchObjectPath(dbus.ObjectPath(path)),
		dbus.WithMatchMember(membr),
	)
	if err != nil {
		fmt.Printf("~~~ add match signal error: %v\n", err)
	}
}

func isShutdown(sig *dbus.Signal) bool {
	if sig.Name == "org.freedesktop.systemd1.Manager.JobNew" {
		for _, b := range sig.Body {
			ss, ok := b.(string)
			if !ok {
				continue
			}
			switch strings.ToLower(ss) {
			case "poweroff.target", "halt.target", "reboot.target":
				return true
			}
		}
	}
	return false
}
