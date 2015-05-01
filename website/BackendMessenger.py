#!/usr/bin/python

import socket

class BackendMessenger:

    def __init__(self): 
        self.TCP_IP = '127.0.0.1'
        self.TCP_PORT = 7070 

    def exit(self):
        self.s.close()

    def ProcessRequest(self, request, args):
        print 'processing request ' + request
        if request=='CLASSIFY':
            self.ClassifyExpenses()
        elif request=='LOAD_RAW':
            self.PullRawData()
        elif request=='CONFIRM_CLASSIFICATION':
            self.ConfirmRequest(args)
        elif request=='CHANGE_CLASSIFICATION':
            self.ChangeClassification(args)
        elif request=='CHANGE_AMOUNT':
            self.ChangeAmount(args)
        elif request=='SAVE_CLASSIFICATION':
            self.SaveClassification(args)
        elif request=='TAG_EXPENSE':
            self.TagExpense(args)
        elif request=='DUPLICATE_EXPENSE':
            self.DuplicateExpense(args)
        elif request=='PROCESS_DOCUMENT':
            self.ProcessDocument(args)
        else:
            print "Unknown Command: " + request
        return 'foo';

    def SendMessage(self, message, args=[]):
        for arg in args:
            print arg
            message = message + '|' + arg
            print message
        try:
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            s.connect((self.TCP_IP, self.TCP_PORT))
            s.send(message)
            s.close()
        except socket.error:
            print 'send failed: '

    def TagExpense(self, args):
        eid = args['eid']
        tag = args['tag']
        self.SendMessage('tag_expense', [eid, tag])

    def SaveClassification(self, args):
        cid = args['cid']
        desc = args['description']
        validfrom = args['validfrom']
        validto = args['validto']
        isexpense = args['isexpense']
        self.SendMessage('save_classification', [cid, desc, validfrom, validto, isexpense])


    def ChangeClassification(self, args):
        eid = args['eid']
        cid = args['cid']
        self.SendMessage('change_classification', [eid, cid])

    def ChangeAmount(self, args):
        eid = args['eid']
        amount = args['amount']
        self.SendMessage('change_amount', [eid, amount])

    def ConfirmRequest(self, args):
        eid = args['eid']
        self.SendMessage('confirm_classification', [eid])

    def PullRawData(self):
        self.SendMessage('load_raw')

    def ClassifyExpenses(self):
        self.SendMessage('classify')

    def DuplicateExpense(self, args):
        eid = args['eid']
        self.SendMessage('duplicate_expense', [eid])

    def ProcessDocument(self, args):
        eid = args['did']
        self.SendMessage('process_document', [eid])

#data = s.recv(BUFFER_SIZE)
        #BUFFER_SIZE = 1024
        #MESSAGE = "lassify"

