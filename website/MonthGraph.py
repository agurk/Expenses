#!/usr/bin/python

import sqlite3
import re

class MonthGraph:

    def __init__(self, date):
        self.date = date
        self.CanvasMaxX = 4000
        self.CanvasMaxY = 2500
        self.Padding = 100
        self.XIncrement = (self.CanvasMaxX - self.Padding) / 31 
        self.AmountMaximum = 2000
        self.MaxX = 0

    def Graph(self):
        amount = self.CumulativeSpend()
        yFactor = float(self.CanvasMaxY) / float(self.AmountMaximum)
        svg = self.SVGHead()
        svg += self.Axis()
        points = [0] * 32
        points[0] = (self.CanvasMaxY)
        average = self.AverageSpend()
        for i in range (1, 32):
            points[i] = ((self.AmountMaximum - int(abs(average[i]))) * yFactor)
        svg += self.Line(points, 'rgb(165, 165, 165)')
        points = [0 for x in range(self.MaxX+1)]
        points[0] = (self.CanvasMaxY)
        for i in range (1, self.MaxX+1):
            points[i] = ((self.AmountMaximum - int(abs(amount[i]))) * yFactor)
        svg += self.Line(points, 'rgb(165,0,0)')
        return (svg + str(self.SVGEnd()))

    def Line(self, points, color):
        xPos = self.Padding
        line = '<polyline points="'
        for yPos in points:
            line += " " + str(xPos) + " " + str(yPos)
            xPos += self.XIncrement
        line += '" stroke="{0}" stroke-width="20" stroke-linecap="square" fill="none" stroke-linejoin="round"/>'.format(color)
        return line

    def SVGHead(self):
        return '<svg height="100%" width="100%" viewBox="0 0 {0} {1}">'.format(self.CanvasMaxX, self.CanvasMaxY + self.Padding)

    def SVGEnd(self):
        return '</svg>'

    def Axis(self):
        svg = '<line x1="{1}" y1="{0}" x2="{1}" y2="0" style="stroke:rgb(0,0,0);stroke-width:10" />'.format(self.CanvasMaxY + self.Padding, self.Padding)
        svg += '<line x1="0" y1="{0}" x2="{1}" y2="{0}" style="stroke:rgb(0,0,0);stroke-width:10" />'.format(str(self.CanvasMaxY), self.CanvasMaxX)
        for i in range (1, 32):
            xPos = (i * self.XIncrement)+self.Padding
            yPos = self.CanvasMaxY + (self.Padding / 1.5)
            svg += '<line x1="{0}" y1="{1}" x2="{0}" y2="{2}" style="stroke:rgb(0,0,0);stroke-width:10" />'.format(xPos, yPos, self.CanvasMaxY)
        amount = 0
        yFactor = float(self.CanvasMaxY) / float(self.AmountMaximum)
        while amount <= self.AmountMaximum:
            amount += 100
            yPos =  self.CanvasMaxY - (amount * yFactor)
            xPos =  (1/3 * self.Padding)
            svg += '<line x1="{0}" y1="{2}" x2="{1}" y2="{2}" style="stroke:rgb(0,0,0);stroke-width:10" />'.format(xPos, self.Padding, yPos)
        return svg
       
    def AverageSpend(self): 
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = 'select sum (e.amount)/12, strftime(\'%d\', e.date) day from expenses e, classifications c, classificationdef cd where date(e.date) < date(\'{0}\',\'start of month\',\'-1 month\') and date(e.date) > date(\'{0}\',\'start of month\',\'-12 months\') and e.eid = c.eid and c.cid = cd.cid and cd.isexpense group by day'.format(self.date)
        cursor = conn.execute(query)
        averageSpend = [0] * 32
        for row in cursor:
            averageSpend[int(row[1])] = float(row[0])
        for i in range (1, 32):
            averageSpend[i] += averageSpend[i-1]
        if abs(averageSpend[31]) > self.AmountMaximum:
            self.AmountMaximum = abs(averageSpend[31])
        return averageSpend

    def CumulativeSpend(self):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str 
        query = 'select sum(amount), date from expenses e, classifications c, classificationdef cd where e.eid = c.eid and c.cid = cd.cid and cd.isexpense and strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\') group by date order by date'.format(self.date)
        cursor = conn.execute(query)
        amounts = [0] * 32
        for row in cursor:
            date = re.match('[0-9]{4}-[0-9]{2}-([0-9]{2})',row[1])
            amounts[int(date.group(1))] = float(row[0])
            self.MaxX = int(date.group(1))
        for i in range(1, self.MaxX + 1):
            amounts[i] = amounts[i] + amounts[i-1]
        if abs(amounts[self.MaxX]) > self.AmountMaximum:
            self.AmountMaximum = abs(amounts[self.MaxX])
        return amounts

