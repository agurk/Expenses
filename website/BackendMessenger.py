#!/usr/bin/python

import socket

class BackendMessenger:

    def __init__(self): 
        self.TCP_IP = '127.0.0.1'
        self.TCP_PORT = 7070 

    def exit(self):
        self.s.close()

    def ProcessRequest(self, request):
        print 'processing request ' + request
        if request=='CLASSIFY':
            self.ClassifyExpenses()
        elif request=='LOAD_RAW':
            self.PullRawData()
        else:
            print "Unknown Command: " + request
        return 'foo';

    def SendMessage(self, message):        
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            s.connect((self.TCP_IP, self.TCP_PORT))
            s.send(message)
            s.close()
        except socket.error:
            print 'send failed: '

    def PullRawData(self):
        self.SendMessage('load_raw')

    def ClassifyExpenses(self):
        self.SendMessage('classify')

#data = s.recv(BUFFER_SIZE)
        #BUFFER_SIZE = 1024
        #MESSAGE = "lassify"

