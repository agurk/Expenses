#!/usr/bin/python

import sqlite3
import re

class MonthGraph:

    def __init__(self, date):
        self.date = date
        self.CanvasMaxX = 4000
        self.CanvasMaxY = 3000
        self.Padding = 100
        self.XIncrement = (self.CanvasMaxX - self.Padding) / 31 
        self.AmountMaximum = 0
        self.MaxX = 0

    def Graph(self):
        amount = self.CumulativeSpend()
        yFactor = (self.CanvasMaxY) / self.AmountMaximum
        svg = '<svg height="100%" width="100%" viewBox="0 0 {0} {1}">'.format(self.CanvasMaxX, self.CanvasMaxY + self.Padding)
        prevX = self.Padding
        prevY = (self.AmountMaximum * yFactor)
        for i in range (0, self.MaxX+1):
            nextX = prevX + self.XIncrement
            nextY = ((self.AmountMaximum - int(abs(amount[i]))) * yFactor)
            svg += '<line x1="{0}" y1="{1}" x2="{2}" y2="{3}" style="stroke:rgb(0,0,0);stroke-width:20" />'.format(str(prevX), str(prevY), str(nextX), str(nextY))
            prevX = nextX
            prevY = nextY
        svg = self.Axis(svg)
        return svg + '</svg>'

    def Axis(self, svg):
        svg += '<line x1="{1}" y1="{0}" x2="{1}" y2="0" style="stroke:rgb(0,0,0);stroke-width:10" />'.format(self.CanvasMaxY + self.Padding, self.Padding)
        svg += '<line x1="0" y1="{0}" x2="{1}" y2="{0}" style="stroke:rgb(0,0,0);stroke-width:10" />'.format(str(self.CanvasMaxY), self.CanvasMaxX)
        for i in range (1, 32):
            xPos = (i * self.XIncrement)+self.Padding
            yPos = self.CanvasMaxY + (self.Padding / 1.5)
            svg += '<line x1="{0}" y1="{1}" x2="{0}" y2="{2}" style="stroke:rgb(0,0,0);stroke-width:10" />'.format(xPos, yPos, self.CanvasMaxY)
        return svg
        

    def CumulativeSpend(self):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = 'select sum(amount), date from expenses e, classifications c, classificationdef cd where e.eid = c.eid and c.cid = cd.cid and cd.isexpense and strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\') group by date order by date'.format(self.date)
        cursor = conn.execute(query)
        amounts = [0] * 32
        for row in cursor:
            date = re.match('[0-9]{4}-[0-9]{2}-([0-9]{2})',row[1])
            print date.group(1)
            amounts[int(date.group(1))] = float(row[0])
            self.MaxX = int(date.group(1))
        for i in range(1, self.MaxX + 1):
            amounts[i] = amounts[i] + amounts[i-1]
        self.AmountMaximum = abs(amounts[self.MaxX])
        return amounts

