#!/usr/bin/python

import sqlite3
import re
import config
import math
import datetime
from FXValues import FXValues
from datetime import date, timedelta, datetime
from calendar import monthrange
import calendar
import time
from MonthView import MonthView

class Line:
    def __init__(self, points, colour, stroke=20, extrapolate=False, exStr="20, 20"):
        self.points = points
        self.colour = colour
        self.stroke = stroke
        self.extrapolate = extrapolate
        self.extrapolateStroke = exStr

class Area:
    def __init__(self, e0, e1, colour):
        self.E0 = e0
        self.E1 = e1
        self.Colour = colour

class MonthGraph:

    def __init__(self, date, period, ccy):
        self.date = date
        if period == 'year':
            self.period = period
            self.periodDays = 366
            self.lookbackYears = 2
        else:
            self.period = 'month'
            self.periodDays = 31
            self.lookbackMonths = 12
        self.ccy = str(ccy)
        self.fxValues = FXValues()
        self.CanvasMaxX = 4750
        self.CanvasMaxY = 2500
        self.Padding = 100
        self.XIncrement = (self.CanvasMaxX - self.Padding) / self.periodDays 
        self.MinYIncrement = 100
        self.AmountMaximum = 2000
        self.MaxX = 0
        self.IncNonDayBox = True
        self.NonDayBoxStyle = ' style="fill:rgb(20, 20, 20)" fill-opacity="0.3"'
        self.AxisStyle = 'style="stroke:rgb(0,0,0);stroke-width:10"'
        self.Lines = []
        self.Areas = []

    def Graph(self):
        result = self.AverageSpend()
        self.AddLine(result['cumulative'], 'rgb(165, 165, 165)', stroke=4)
        if self.period == 'year':
            prevYear = time.strptime(self.date, "%Y-%m-%d")
            mv = MonthView(self.date, self.period)
            self.AddLine(self.CumulativeSpend(mv.add_months(prevYear, -12)), 'rgb(229,84,84)')
            self.AddLine(self.CumulativeSpend(mv.add_months(prevYear, -24)), 'rgb(239,157,157)')
            self.AddLine(self.CumulativeSpend(self.date), 'rgb(165,0,0)')
        else:
            self.AddLine(self.CumulativeSpend(self.date), 'rgb(165,0,0)', extrapolate=True)
            self.AddSD(result['cumulative'], result['sd'])
        return (self._buildSVG())

    def AddSD(self, average, sd):
        means = [self.CanvasMaxY] * (self.periodDays + 1)
        sdUp = [self.CanvasMaxY] * (self.periodDays + 1)
        sdDown = [self.CanvasMaxY] * (self.periodDays + 1)
        twosdUp = [self.CanvasMaxY] * (self.periodDays + 1)
        twosdDown = [self.CanvasMaxY] * (self.periodDays + 1)
        yFactor = float(self.CanvasMaxY) / float(self.AmountMaximum)
        for i in range (1, (self.periodDays + 1)):
            means[i] = ((self.AmountMaximum - int(abs(average[i]))) * yFactor)
        for i in range (1, (self.periodDays + 1)):
            sdUp[i] = means[i] - sd[i] * yFactor
            twosdUp[i] = means[i] - sd[i]*2 * yFactor
            sdDown[i] = means[i] + sd[i] * yFactor
            twosdDown[i] = means[i] + sd[i]*2 * yFactor
            if sdDown[i] > self.CanvasMaxY:
                sdDown[i] = self.CanvasMaxY
            if twosdDown[i] > self.CanvasMaxY:
                twosdDown[i] = self.CanvasMaxY
        self.AddArea(twosdUp, twosdDown, 'rgb(240, 240, 240)')
        self.AddArea(sdUp, sdDown, 'rgb(225, 225, 225)')

    def AddLine(self, amount, colour, stroke=20, extrapolate=False):
        self.Lines.append(Line(amount, colour, stroke, extrapolate))
        for amt in amount:
            if abs(amt) > self.AmountMaximum:
                self.AmountMaximum = abs(amt)

    def AddArea(self, e0, e1, colour):
        self.Areas.append(Area(e0, e1, colour))

    def _buildSVG(self):
        svg = self.SVGHead()
        for area in self.Areas:
            svg += self.Area(area.E0, area.E1, area.Colour)
        for line in self.Lines:
            svg += self._buildLine(line.points, line.colour, line.stroke, line.extrapolate, line.extrapolateStroke)
        if self.IncNonDayBox:
            svg += self.NonMonthDayBox()
        svg += self.Axis()
        return (svg + str(self.SVGEnd()))

    def _buildLine(self, amount, colour, stroke, extrapolate, extraStroke):
        yFactor = float(self.CanvasMaxY) / float(self.AmountMaximum)
        maxRange = len(amount)
        points = [0 for x in range(maxRange)]
        points[0] = (self.CanvasMaxY)
        for i in range (1, maxRange):
            points[i] = ((self.AmountMaximum - int(abs(amount[i]))) * yFactor)
        xPos = self.Padding
        line = '<polyline points="'
        for yPos in points:
            line += " " + str(xPos) + " " + str(yPos)
            xPos += self.XIncrement
        line += '" stroke="{0}" stroke-width="{1}" stroke-linecap="square" fill="none" stroke-linejoin="round"/>'.format(colour, stroke)
        if extrapolate:
            line += self.ExtrapolatedLine(points, colour, stroke, extraStroke)
        return line

    def ExtrapolatedLine(self, points, color, stroke, extraStroke):
        todaysDate = date.today()
        line = ''
        if (self._periodDiff(todaysDate) == 0):
            if self.period == 'year':
                day = self._getExpenseDay(todaysDate)
            else:
                day = todaysDate.day
            x1 = self.MaxX * self.XIncrement + self.Padding
            x2 = (day - self.MaxX) * self.XIncrement + x1
            line = '<line stroke="{0}" stroke-width="{1}" stroke-dasharray="{5}"  x1="{2}" y1="{3}" x2="{4}" y2="{3}" />'.format(color, stroke, x1, points[self.MaxX], x2, extraStroke)
        elif ( self._periodDiff(todaysDate) < 0 ):
            graphDate = datetime.strptime(self.date, '%Y-%m-%d')
            lastDay = self._daysInPeriod(graphDate)
            if (lastDay > graphDate.day ):
                x1 = self.MaxX * self.XIncrement + self.Padding
                x2 = (lastDay - self.MaxX) * self.XIncrement + x1
                line = '<line stroke="{0}" stroke-width="{1}" x1="{2}" y1="{3}" x2="{4}" y2="{3}" />'.format(color, stroke, x1, points[self.MaxX], x2)
        return line

    def Area(self, pointsUp, pointsDown, color):
        xPos = self.Padding
        polygon = '<polygon points="'
        for yPos in pointsUp:
            polygon += '{0}, {1} '.format(str(xPos), str(yPos))
            xPos += self.XIncrement
        for i in reversed(range(0, (self.periodDays + 1))):
            yPos = pointsDown[i]
            xPos -= self.XIncrement
            polygon += '{0}, {1} '.format(str(xPos), str(yPos))
        polygon += '" fill="{0}" stroke-width="0" />'.format(color)
        return polygon

    def SVGHead(self):
        return '<svg viewBox="{0} {1} {2} {3}">'.format(self.Padding*-2, self.Padding*-1, self.CanvasMaxX+3*self.Padding, self.CanvasMaxY + 2*self.Padding)

    def SVGEnd(self):
        return '</svg>'

    def Axis(self):
        svg = '<line x1="{1}" y1="{0}" x2="{1}" y2="0" {2}/>'.format(self.CanvasMaxY + self.Padding, self.Padding, self.AxisStyle)
        svg += '<line x1="0" y1="{0}" x2="{1}" y2="{0}" {2} />'.format(str(self.CanvasMaxY), self.CanvasMaxX, self.AxisStyle)
        if self.period == 'year':
            months = self._months()
            for i in range(1, 13):
                xPos = (self._monthYearDay(i) * self.XIncrement)+self.Padding
                yPos = self.CanvasMaxY + (self.Padding / 2.5)
                svg += '<line x1="{0}" y1="{1}" x2="{0}" y2="{2}" {3} />'.format(xPos, yPos, self.CanvasMaxY, self.AxisStyle)
                svg += '<text x="{0}" y={1} font-size="80" text-anchor="middle">{2}</text>'.format(xPos, yPos+60, months[i])
        else:
            for i in range (1, (self.periodDays + 1)):
                xPos = (i * self.XIncrement)+self.Padding
                yPos = self.CanvasMaxY + (self.Padding / 2.5)
                svg += '<line x1="{0}" y1="{1}" x2="{0}" y2="{2}" {3} />'.format(xPos, yPos, self.CanvasMaxY, self.AxisStyle)
                svg += '<text x="{0}" y={1} font-size="80" text-anchor="middle">{2}</text>'.format(xPos, yPos+60, i)
        increment = int( math.pow(10, round(math.log10(self.AmountMaximum) -1 )) )
        if (increment  < self.MinYIncrement):
            increment = self.MinYIncrement
        amount = 0
        yFactor = float(self.CanvasMaxY) / float(self.AmountMaximum)
        while amount <= self.AmountMaximum:
            amount += increment
            yPos =  self.CanvasMaxY - (amount * yFactor)
            xPos =  (1/3 * self.Padding)
            svg += '<line x1="{0}" y1="{2}" x2="{1}" y2="{2}" {3} />'.format(xPos, self.Padding, yPos, self.AxisStyle)
            svg += '<text x="{0}" y={1} font-size="80" text-anchor="end" dominant-baseline="middle">{2}</text>'.format(xPos, yPos, amount)
        return svg

    def NonMonthDayBox(self):
        height = self.CanvasMaxY + self.Padding
        width = self.XIncrement * (self.periodDays- self._daysInPeriod(self.date))
        xPos = self.CanvasMaxX - width
        yPos = 0 - self.Padding
        svg = '<rect x="{0}" y="{1}" width="{2}" height="{3}" {4}/>'.format(xPos, yPos, width, height, self.NonDayBoxStyle)
        return svg

    def _monthAverageQuery(self):
        return 'select sum (e.amount), strftime(\'%d\', e.date) day, strftime(\'%m\', e.date) month, strftime(\'%Y\', e.date) year, ccy from expenses e, classifications c, classificationdef cd where date(e.date) < date(\'{0}\',\'start of month\') and date(e.date) > date(\'{0}\',\'start of month\',\'-12 months\') and e.eid = c.eid and c.cid = cd.cid and cd.isexpense group by day, month, ccy'.format(self.date)

    def _yearAverageQuery(self):
        return 'select sum (e.amount), strftime(\'%d\', e.date) day, strftime(\'%m\', e.date) month, strftime(\'%Y\', e.date) year, ccy from expenses e, classifications c, classificationdef cd where date(e.date) < date(\'{0}\',\'start of year\') and date(e.date) >= date(\'{0}\',\'start of year\',\'-{1} years\') and e.eid = c.eid and c.cid = cd.cid and cd.isexpense group by day, month, ccy'.format(self.date, self.lookbackYears)
       
    def AverageSpend(self): 
        conn = sqlite3.connect(config.SQLITE_DB, uri=True)
        conn.text_factory = str 
        if self.period == 'year':
            query = self._yearAverageQuery()
            lookbackPeriod = self.lookbackYears
        else:
            query = self._monthAverageQuery()
            lookbackPeriod = self.lookbackMonths
        cursor = conn.execute(query)
        averageSpend = [[0 for x in range(self.periodDays + 1)] for x in range(lookbackPeriod + 1)]
        totalSpend = [0 for x in range(self.periodDays + 1)]
        for row in cursor:
            amount = float(row[0])
            date = row[3] +'-'+ row[2] +'-'+ row[1]
            year = int(row[3])
            ccy = row[4]
            if self.period == 'year':
                key = abs(self._periodDiff(date))
                day = self._getExpenseDay(date)
            else:
                #month
                key = int(row[2])
                day = int(row[1])
            averageSpend[key][day] += self.fxValues.FXAmount(amount, ccy, self.ccy, date)
        cumulativeAmount = [0 for x in range((self.periodDays + 1))]
        diff = [0 for x in range((self.periodDays + 1))]
        for day in range(1, (self.periodDays + 1)):
            for year in range(1, lookbackPeriod + 1):
                cumulativeAmount[day] += abs(averageSpend[year][day])
                averageSpend[year][day] += averageSpend[year][day-1]
            cumulativeAmount[day] = cumulativeAmount[day] / lookbackPeriod 
            cumulativeAmount[day] += cumulativeAmount[day -1]
            averageDiff = 0
            for year in range(1, lookbackPeriod + 1):
                diff[day] += math.pow((abs(averageSpend[year][day]) - cumulativeAmount[day]),2)
            diff[day] = math.sqrt( diff[day] / lookbackPeriod )
            if abs(cumulativeAmount[day]) > self.AmountMaximum:
                self.AmountMaximum = abs(cumulativeAmount[day])
        conn.close()
        return {'cumulative': cumulativeAmount, 'sd': diff}

    def _monthCumulativeQuery(self, date):
        return 'select amount, ccy, date from expenses e, classifications c, classificationdef cd where e.eid = c.eid and c.cid = cd.cid and cd.isexpense and strftime(date) >= date(\'{0}\',\'start of month\') and strftime(date) < date(\'{0}\',\'start of month\',\'+1 month\')'.format(date)

    def _yearCumulativieQuery(self, date):
        return 'select amount, ccy, date from expenses e, classifications c, classificationdef cd where e.eid = c.eid and c.cid = cd.cid and cd.isexpense and strftime(date) >= date(\'{0}\',\'start of year\') and strftime(date) < date(\'{0}\',\'start of year\',\'+1 year\')'.format(date)

    def CumulativeSpend(self, date):
        conn = sqlite3.connect(config.SQLITE_DB, uri=True)
        conn.text_factory = str 
        if self.period == 'year':
            query = self._yearCumulativieQuery(date)
        else:
            query = self._monthCumulativeQuery(date)
        cursor = conn.execute(query)
        amounts = [0] * (self.periodDays + 1)
        localMaxX = 0
        for row in cursor:
            day = self._getExpenseDay(row[2])
            amounts[int(day)] += self.fxValues.FXAmount(row[0],row[1],self.ccy,self.date)
            if int(day) > localMaxX:
                localMaxX = int(day)
        if (localMaxX > self.MaxX):
            self.MaxX = localMaxX
        for i in range(1, localMaxX + 1):
            amounts[i] = amounts[i] + amounts[i-1]
        for i in reversed(range (localMaxX + 1, self.periodDays + 1)):
            amounts.pop(i)
        conn.close()
        return amounts

    def _averageMonthlySpend(self, cursor):
        averageSpend = [[0 for x in range(32)] for x in range(13)]
        totalSpend = [0 for x in range(32)]
        for row in cursor:
            amount = float(row[0])
            day = int(row[1])
            month = int(row[2])
            year = int(row[3])
            ccy = row[4]
            averageSpend[month][day] += self.fxValues.FXAmount(amount, ccy, self.ccy, str(year) +'-'+ str(month) +'-'+ str(day))
        cumulativeAmount = [0 for x in range((self.periodDays + 1))]
        diff = [0 for x in range((self.periodDays + 1))]
        for day in range(1, (self.periodDays + 1)):
            for month in range(1, 13):
                cumulativeAmount[day] += abs(averageSpend[month][day])
                averageSpend[month][day] += averageSpend[month][day-1]
            cumulativeAmount[day] = cumulativeAmount[day] / 12
            cumulativeAmount[day] += cumulativeAmount[day -1]
            averageDiff = 0
            for month in range(1, 13):
                diff[day] += math.pow((abs(averageSpend[month][day]) - cumulativeAmount[day]),2)
            diff[day] = math.sqrt( diff[day] / 12 )
            if abs(cumulativeAmount[day]) > self.AmountMaximum:
                self.AmountMaximum = abs(cumulativeAmount[day])
        return {'cumulative': cumulativeAmount, 'sd': diff}

    def _daysInPeriod(self, date):
        graphDate = datetime.strptime(self.date, '%Y-%m-%d')
        if self.period == 'year':
           lastDay = 365
           if calendar.isleap(graphDate.year):
               lastDay += 1
           return lastDay
        else:
            return monthrange(graphDate.year, graphDate.month)[1]

    def _periodDiff(self, comparedDate):
        if type (comparedDate) is str:
            comparedDate = datetime.strptime(comparedDate, '%Y-%m-%d')
        graphDate = datetime.strptime(self.date, '%Y-%m-%d')
        if self.period == 'year':
            return graphDate.year - comparedDate.year
        else:
            yearDiff = graphDate.year - comparedDate.year
            monthDiff = graphDate.month - comparedDate.month
            return yearDiff * 12 + monthDiff

    def _getExpenseDay(self, date):
        if type(date) is str:
            date = datetime.strptime(date, "%Y-%m-%d")
        if self.period == 'year':
            return date.timetuple().tm_yday
        else:
            return date.day

    def _months(self):
        return ['None', 'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

    def _monthYearDay(self, month):
        graphDate = datetime.strptime(self.date, '%Y-%m-%d')
        newDate = graphDate.replace(month = month, day = 1)
        return self._getExpenseDay(newDate)
