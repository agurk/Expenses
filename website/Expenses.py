#!/usr/bin/python

from flask import Flask, request
from flask import render_template

from MonthView import MonthView
from MonthGraph import MonthGraph
from Document import Document
from ConfigView import Config
from ItemView import ItemView
from ReadOnlyData import ReadOnlyData
from Search import Search
from EventGenerator import EventGenerator

import time

app = Flask(__name__)

@app.route('/')
@app.route('/expenses')
def main():
    if 'date' in request.args.keys():
        mv = MonthView(request.args['date'])
        mg = MonthGraph(request.args['date'])
    else:
        mv = MonthView(time.strftime("%Y-%m-%d"))
        mg = MonthGraph(time.strftime("%Y-%m-%d"))
    return render_template('monthview.html', cursor=mv.OverallExpenses(), cursor2=mv.IndividualExpenses(), previous_month=mv.PreviousMonth(), next_month=mv.NextMonth(), total_amount=mv.TotalAmount(), month_name=mv.MonthName(),month_graph=mg.Graph(), this_month=mv.ThisMonth())

@app.route('/receipt')
def on_receipt():
    did = request.args['did']
    doc = Document()
    return render_template('receipt.html', document=doc.Document(did))

@app.route('/expense')
def on_edit_expense():
    eid = ''
    if 'eid' in request.args.keys():
        eid = request.args['eid']
    ex = Expense()
    return render_template('expenseview.html', expense=ex.Expense(eid))

@app.route('/detailed_expenses_all')
def on_detailed_expenses_all():
     if 'date' in request.args.keys():
         mv = MonthView(request.args['date'])
         mg = MonthGraph(request.args['date'])
     else:
         mv = MonthView(time.strftime("%Y-%m-%d"))
         mg = MonthGraph(time.strftime("%Y-%m-%d"))
     return render_template('detailedexpenses.html', cursor2=mv.IndividualExpensesAll())

@app.route('/detailed_expenses')
def on_detailed_expenses():
     if 'date' in request.args.keys():
         mv = MonthView(request.args['date'])
         mg = MonthGraph(request.args['date'])
     else:
         mv = MonthView(time.strftime("%Y-%m-%d"))
         mg = MonthGraph(time.strftime("%Y-%m-%d"))
     return render_template('detailedexpenses.html', cursor2=mv.IndividualExpenses())

@app.route('/config')
def on_config():
     config = Config()
     return render_template('config.html', classifications=config.AllClassifications());

@app.route('/expense_summary')
def on_expense_summary():
    eid = request.args['eid']
    rod = ReadOnlyData()
    return render_template('expense.html', row=rod.Expense(eid))

@app.route('/expense_details')
def on_expense_details():
    idno = request.args['eid']
    expense = ItemView(idno)
    return render_template('expense_details.html',rawData=expense.RawStr(), classifications=expense.Classifications(), classification=expense.Classification(), amount=expense.Amount(),eid=idno)

@app.route('/search')
def on_search():
    search = Search()
    if 'description' in request.args.keys():
        description = request.args['description']
        similar_ex = search.SimilarExpenses(description)
    else:
        description = ''
        similar_ex = ''
    return render_template('search.html', description=description, similar_ex=similar_ex)

@app.route('/backend/<command>', methods=['GET', 'POST'])
def on_backend(command):
    eg = EventGenerator()
    eg.sendEvent(command, request.args)
    return '200';

if __name__ == '__main__':
    app.run(debug=True, use_reloader=True)

