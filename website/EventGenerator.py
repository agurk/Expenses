#!/usr/bin/python

import dbus

class EventGenerator:

    EVENT_TYPE = 'ExpensesTimeEvent'
    SERVICE_OBJECT_NAME = '/ExpensesTime/EventServiceObject'
    DBUS_SERVICE_NAME = 'ExpensesTime.events'
    DBUS_INTERFACE_NAME = 'ExpensesTime.interface'

    def __init__(self): 
        bus = dbus.SessionBus()
        self.es = bus.get_object(self.DBUS_SERVICE_NAME, self.SERVICE_OBJECT_NAME)

    def sendEvent(self, event, args):
        payload={}
        for key in args.keys() :
            payload[key] = args[key]
        if (len(payload)):
            self.es.sendEvent(event, payload)
        else:
            self.es.sendEvent(event)

