#!/usr/bin/python

import sqlite3
import re
import config
import math
from FXValues import FXValues

class MonthGraph:

    def __init__(self, date, ccy='GBP'):
        self.date = date
        self.ccy = str(ccy)
        self.fxValues = FXValues()
        self.CanvasMaxX = 4000
        self.CanvasMaxY = 2500
        self.Padding = 100
        self.XIncrement = (self.CanvasMaxX - self.Padding) / 31 
        self.AmountMaximum = 2000
        self.MaxX = 0

    def Graph(self):
        amount = self.CumulativeSpend()
        result = self.AverageSpend()
        svg = self.SVGHead()
        svg = self.BuildSD(svg, result['cumulative'], result['sd'])
        yFactor = float(self.CanvasMaxY) / float(self.AmountMaximum)
        points = [0 for x in range(self.MaxX+1)]
        points[0] = (self.CanvasMaxY)
        for i in range (1, self.MaxX+1):
            points[i] = ((self.AmountMaximum - int(abs(amount[i]))) * yFactor)
        svg += self.Line(points, 'rgb(165,0,0)')
        svg += self.Axis()
        return (svg + str(self.SVGEnd()))

    def BuildSD(self, svg, average, sd):
        means = [self.CanvasMaxY] * 32
        sdUp = [self.CanvasMaxY] * 32
        sdDown = [self.CanvasMaxY] * 32
        twosdUp = [self.CanvasMaxY] * 32
        twosdDown = [self.CanvasMaxY] * 32
        yFactor = float(self.CanvasMaxY) / float(self.AmountMaximum)
        for i in range (1, 32):
            means[i] = ((self.AmountMaximum - int(abs(average[i]))) * yFactor)
        for i in range (1, 32):
            sdUp[i] = means[i] - sd[i] * yFactor
            twosdUp[i] = means[i] - sd[i]*2 * yFactor
            sdDown[i] = means[i] + sd[i] * yFactor
            twosdDown[i] = means[i] + sd[i]*2 * yFactor
            if sdDown[i] > self.CanvasMaxY:
                sdDown[i] = self.CanvasMaxY
            if twosdDown[i] > self.CanvasMaxY:
                twosdDown[i] = self.CanvasMaxY
        svg += self.Area(twosdUp, twosdDown, 'rgb(240, 240, 240)')
        svg += self.Area(sdUp, sdDown, 'rgb(225, 225, 225)')
        svg += self.Line(means, 'rgb(165, 165, 165)', 4)
        return svg

    def Line(self, points, color, stroke=20):
        xPos = self.Padding
        line = '<polyline points="'
        for yPos in points:
            line += " " + str(xPos) + " " + str(yPos)
            xPos += self.XIncrement
        line += '" stroke="{0}" stroke-width="{1}" stroke-linecap="square" fill="none" stroke-linejoin="round"/>'.format(color, stroke)
        return line

    def Area(self, pointsUp, pointsDown, color):
        xPos = self.Padding
        polygon = '<polygon points="'
        for yPos in pointsUp:
            polygon += '{0}, {1} '.format(str(xPos), str(yPos))
            xPos += self.XIncrement
        for i in reversed(range(0, 32)):
            yPos = pointsDown[i]
            xPos -= self.XIncrement
            polygon += '{0}, {1} '.format(str(xPos), str(yPos))
        polygon += '" fill="{0}" stroke-width="0" />'.format(color)
        return polygon

    def SVGHead(self):
        return '<svg height="100%" width="100%" viewBox="{0} {1} {2} {3}">'.format(self.Padding*-2, self.Padding*-1, self.CanvasMaxX+3*self.Padding, self.CanvasMaxY + 2*self.Padding)

    def SVGEnd(self):
        return '</svg>'

    def Axis(self):
        svg = '<line x1="{1}" y1="{0}" x2="{1}" y2="0" style="stroke:rgb(0,0,0);stroke-width:10" />'.format(self.CanvasMaxY + self.Padding, self.Padding)
        svg += '<line x1="0" y1="{0}" x2="{1}" y2="{0}" style="stroke:rgb(0,0,0);stroke-width:10" />'.format(str(self.CanvasMaxY), self.CanvasMaxX)
        for i in range (1, 32):
            xPos = (i * self.XIncrement)+self.Padding
            yPos = self.CanvasMaxY + (self.Padding / 2.5)
            svg += '<line x1="{0}" y1="{1}" x2="{0}" y2="{2}" style="stroke:rgb(0,0,0);stroke-width:10" />'.format(xPos, yPos, self.CanvasMaxY)
            svg += '<text x="{0}" y={1} font-size="80" text-anchor="middle">{2}</text>'.format(xPos, yPos+60, i)
        amount = 0
        yFactor = float(self.CanvasMaxY) / float(self.AmountMaximum)
        while amount <= self.AmountMaximum:
            amount += 100
            yPos =  self.CanvasMaxY - (amount * yFactor)
            xPos =  (1/3 * self.Padding)
            svg += '<line x1="{0}" y1="{2}" x2="{1}" y2="{2}" style="stroke:rgb(0,0,0);stroke-width:10" />'.format(xPos, self.Padding, yPos)
            svg += '<text x="{0}" y={1} font-size="80" text-anchor="end" dominant-baseline="middle">{2}</text>'.format(xPos, yPos, amount)
        return svg
       
    def AverageSpend(self): 
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
        query = 'select sum (e.amount), strftime(\'%d\', e.date) day, strftime(\'%m\', e.date) month, strftime(\'%Y\', e.date) year, ccy from expenses e, classifications c, classificationdef cd where date(e.date) < date(\'{0}\',\'start of month\') and date(e.date) > date(\'{0}\',\'start of month\',\'-12 months\') and e.eid = c.eid and c.cid = cd.cid and cd.isexpense group by day, month, ccy'.format(self.date)
        cursor = conn.execute(query)
        averageSpend = [[0 for x in range(32)] for x in range(13)]
        totalSpend = [0 for x in range(32)]
        for row in cursor:
            amount = float(row[0])
            day = int(row[1])
            month = int(row[2])
            year = int(row[3])
            ccy = row[4]
            averageSpend[month][day] += self.fxValues.FXAmount(amount, ccy, self.ccy, str(year) +'-'+ str(month) +'-'+ str(day))
        cumulativeAmount = [0 for x in range(32)]
        diff = [0 for x in range(32)]
        for day in range(1, 32):
            for month in range(1, 13):
                cumulativeAmount[day] += abs(averageSpend[month][day])
                averageSpend[month][day] -= averageSpend[month][day-1]
            cumulativeAmount[day] = cumulativeAmount[day] / 12
            cumulativeAmount[day] += cumulativeAmount[day -1]
            for month in range(1, 13):
                diff[day] += math.pow((abs(averageSpend[month][day]) - cumulativeAmount[day]),2)
            diff[day] = math.sqrt( diff[day] / cumulativeAmount[day] )
            if abs(cumulativeAmount[day]) > self.AmountMaximum:
                self.AmountMaximum = abs(cumulativeAmount[day])
        return {'cumulative': cumulativeAmount, 'sd': diff}

    def CumulativeSpend(self):
        conn = sqlite3.connect(config.SQLITE_DB)
        conn.text_factory = str 
#        query = 'select sum(amount), date from expenses e, classifications c, classificationdef cd where e.eid = c.eid and c.cid = cd.cid and cd.isexpense and strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\') group by date order by date'.format(self.date)
        query = 'select amount, ccy, date from expenses e, classifications c, classificationdef cd where e.eid = c.eid and c.cid = cd.cid and cd.isexpense and strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\')'.format(self.date)
        cursor = conn.execute(query)
        amounts = [0] * 32
        for row in cursor:
            date = re.match('[0-9]{4}-[0-9]{2}-([0-9]{2})',row[2])
            amounts[int(date.group(1))] += self.fxValues.FXAmount(row[0],row[1],self.ccy,self.date)
            if int(date.group(1)) > self.MaxX:
                self.MaxX = int(date.group(1))
        for i in range(1, self.MaxX + 1):
            amounts[i] = amounts[i] + amounts[i-1]
            if abs(amounts[i]) > self.AmountMaximum:
                self.AmountMaximum = abs(amounts[i])
        return amounts

