#!/usr/bin/python3

from flask import Flask, request, make_response, render_template

from MonthView import MonthView
from MonthGraph import MonthGraph
from Document import Document
from MetaData import MetaData
from Analysis import Analysis
from EventGenerator import EventGenerator
from Expense import Expense
from OverallExpenses import OverallExpenses
from flask import send_file
from ResourceHandler import ResourceHandler, ResourceType

import time

app = Flask(__name__)

def _getParam(key, default=''):
    if key in request.args.keys():
        return request.args[key]
    return default

@app.route('/')
@app.route('/expenses')
def main():
    oe = OverallExpenses()
    date = _getParam('date', time.strftime("%Y-%m-%d"))
    period = _getParam('period', 'month')
    ccy = _getParam('ccy', 'GBP')
    mv = MonthView(date, period)
    ex = Expense(period)
    overall = oe.OverallExpenses(date, period, ccy)
    mg = MonthGraph(date, period, ccy)
    return render_template('monthview.html', overall_expenses=overall, expenses=ex.Expenses(date, ''), previous_period=mv.PreviousPeriod(period), previous_year=mv.PreviousYear(), next_period=mv.NextPeriod(period), month_name=mv.MonthName(),month_graph=mg.Graph(), this_month=mv.ThisMonth(), period=period, ccy=ccy)

@app.route('/analysis')
def on_analysis():
    dateFrom = _getParam('from', '2011')
    dateTo = _getParam('to', '2018')
    ccy = _getParam('ccy', 'GBP')
    analysis = Analysis(dateFrom, dateTo, ccy)
    results = analysis.YearlySpend()
    yearTotals = analysis.yearTotals()
    return (render_template('analysis.html', yearly_spend = results['salary'], reimbursements=results['reimbursements'], expenses=results['expenses'], withholding=results['withholding'], year_totals=yearTotals, date_from=dateFrom, date_to=dateTo, date_range = analysis.DateRange(), ccy=ccy))

@app.route('/documents')
def on_documents():
    doc = Document()
    return (render_template('documents.html', documents=doc.Documents()))

@app.route('/document_fragment')
def on_document_fragment():
    did = request.args['did']
    doc = Document()
    return render_template('document_fragment.html', document=doc.Document(did))


@app.route('/document')
@app.route('/receipt')
def on_receipt():
    did = request.args['did']
    doc = Document()
    return render_template('receipt.html', document=doc.Document(did), item_id=did, item_type='did')

@app.route('/document_all_expense_fragments')
def on_document_expense_fragment():
    did = _getParam('did') 
    doc = Document()
    return render_template('document_all_expense_fragments.html', document=doc.Document(did))

@app.route('/expense')
def on_edit_expense():
    eid = _getParam('eid')
    # period shouldn't matter for this, so setting to month
    ex = Expense('month')
    md = MetaData()
    if eid == 'NEW':
        did = _getParam('did')
        return render_template('expense.html', expense=ex.NewExpense(did=did), classifications=md.Classifications(), accounts=md.Accounts())
    else:
        return render_template('expense.html', expense=ex.Expense(eid), classifications=md.Classifications(eid), item_id=eid, item_type='eid', accounts=md.Accounts())

@app.route('/detailed_expenses')
def on_detailed_expenses():
    date = _getParam('date', time.strftime("%Y-%m-%d")) 
    allExes = _getParam('all', '')
    ccy = _getParam('ccy')
    period = _getParam('period', 'month')
    ex = Expense(period)
    return render_template('detailedexpenses_fragment.html', expenses=ex.Expenses(date, allExes, ccy))

@app.route('/config')
def on_config():
    md = MetaData()
    return render_template('config.html', classifications=md.AllClassifications(), accountloaders=md.AccountLoaders(), accounts=md.Accounts());

@app.route('/expense_summary')
def on_expense_summary():
    eid = request.args['eid']
    period = _getParam('period', 'month')
    ex = Expense(period)
    return render_template('expense_fragment.html', expense=ex.Expense(eid))

@app.route('/expense_details')
def on_expense_details():
    eid = request.args['eid']
    period = _getParam('period', 'month')
    ex = Expense(period)
    md = MetaData()
    return render_template('expense_details_fragment.html',expense=ex.Expense(eid), classifications=md.Classifications(eid))

@app.route('/search')
def on_search():
    # assuming period doesn't matter here so setting to month
    ex = Expense('month')
    classification = _getParam('classification')
    if 'description' in request.args.keys():
        description = request.args['description']
        similar_ex = ex.Search(description, classification)
    else:
        description = ''
        similar_ex = ''
    return render_template('search.html', description=description, similar_ex=similar_ex)

@app.route('/image/<path>/<filename>')
@app.route('/image/<filename>')
@app.route('/pdf/<filename>')
@app.route('/pdf/<path>/<filename>')
def serveImage(filename, path=''):
    if '/pdf/' in request.path:
        resType = ResourceType.pdf
    elif '/image/' in request.path:
        resType = ResourceType.image
    rh = ResourceHandler()
    return send_file(rh.Resource(resType, filename, path))


@app.route('/backend/<command>', methods=['GET', 'POST'])
def generateEvent(command):
    extraArgs = {}
    if command == 'MERGE_EXPENSE' or command == 'MERGE_EXPENSE_COMMISSION':
        extraArgs[request.cookies.get('pinned_type') + '_merged'] = request.cookies.get('pinned_id')
    _generateEvent(command, request.args, extraArgs)
    return '200';

def _generateEvent(command, args, extraArgs={}):
    eg = EventGenerator()
    eg.sendEvent(command, args, extraArgs)

@app.route('/pinned', methods=['GET', 'POST', 'DELETE', 'PUT'])
def pinned():
    if request.method == 'POST':
        pt = request.args['pinned_type']
        pid = request.args['pinned_id']
        pin = pt and pid
        resp = make_response(render_template('pinned.html', pin=pin))
        resp.set_cookie('pinned_type', pt)
        resp.set_cookie('pinned_id', pid)
        return resp
    elif request.method == 'DELETE':
        resp = make_response(render_template('pinned.html'))
        resp.set_cookie('pinned_type', '')
        resp.set_cookie('pinned_id', '')
        return resp
    elif request.method == 'PUT':
        eventArgs={}
        eventArgs[request.args['pinned_type']] = request.args['pinned_id']
        eventArgs[request.cookies.get('pinned_type')] = request.cookies.get('pinned_id')
        _generateEvent('PIN_ITEM', eventArgs)
        resp = make_response(render_template('pinned.html'))
        resp.set_cookie('pinned_type', '')
        resp.set_cookie('pinned_id', '')
        return resp
    else:
        pin = request.cookies.get('pinned_type') and request.cookies.get('pinned_id')
        return render_template('pinned.html', pin=pin)

if __name__ == '__main__':
    app.run(debug=True, use_reloader=True)

