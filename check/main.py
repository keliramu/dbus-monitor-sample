import dbus
sysbus = dbus.SystemBus()
systemd1 = sysbus.get_object('org.freedesktop.systemd1', '/org/freedesktop/systemd1')
manager = dbus.Interface(systemd1, 'org.freedesktop.systemd1.Manager')
service = sysbus.get_object('org.freedesktop.systemd1', object_path=manager.GetUnit('snapd.service'))
interface = dbus.Interface(service, dbus_interface='org.freedesktop.DBus.Properties')
print(interface.Get('org.freedesktop.systemd1.Unit', 'ActiveState'))