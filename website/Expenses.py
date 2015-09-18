#!/usr/bin/python

from flask import Flask, request, make_response, render_template

from MonthView import MonthView
from MonthGraph import MonthGraph
from Document import Document
from MetaData import MetaData
from EventGenerator import EventGenerator
from Expense import Expense
from OverallExpenses import OverallExpenses
from flask import send_file
from ImageHandler import ImageHandler

import time

app = Flask(__name__)

@app.route('/')
@app.route('/expenses')
def main():
    if 'date' in request.args.keys():
        date = request.args['date']
    else:
        date = time.strftime("%Y-%m-%d")
    if 'ccy' in request.args.keys():
        mg = MonthGraph(date, request.args['ccy'])
    else:
        mg = MonthGraph(date)
    mv = MonthView(date)
    ex = Expense()
    oe = OverallExpenses()
    overall = oe.OverallExpenses(date)
    return render_template('monthview.html', cursor=overall, expenses=ex.Expenses(date, ''), previous_month=mv.PreviousMonth(), next_month=mv.NextMonth(), total_amount=oe.TotalAmount(overall), month_name=mv.MonthName(),month_graph=mg.Graph(), this_month=mv.ThisMonth())

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
    did = _getFromArgs(request.args, 'did', '') 
    doc = Document()
    return render_template('document_all_expense_fragments.html', document=doc.Document(did))

@app.route('/expense')
def on_edit_expense():
    eid = _getFromArgs(request.args, 'eid', '') 
    ex = Expense()
    md = MetaData()
    return render_template('expense.html', expense=ex.Expense(eid), classifications=md.Classifications(eid), item_id=eid, item_type='eid')

@app.route('/detailed_expenses')
def on_detailed_expenses():
    date = _getFromArgs(request.args, 'date', time.strftime("%Y-%m-%d")) 
    allExes = _getFromArgs(request.args, 'all', 'false')
    ccy = _getFromArgs(request.args, 'ccy', '')
    ex = Expense()
    return render_template('detailedexpenses_fragment.html', expenses=ex.Expenses(date, allExes, ccy))

@app.route('/config')
def on_config():
    md = MetaData()
    return render_template('config.html', classifications=md.AllClassifications(), accountloaders=md.AccountLoaders());

@app.route('/expense_summary')
def on_expense_summary():
    eid = request.args['eid']
    ex = Expense()
    return render_template('expense_fragment.html', expense=ex.Expense(eid))

@app.route('/expense_details')
def on_expense_details():
    eid = request.args['eid']
    ex = Expense()
    md = MetaData()
    return render_template('expense_details_fragment.html',expense=ex.Expense(eid), classifications=md.Classifications(eid))

@app.route('/search')
def on_search():
    ex = Expense()
    if 'description' in request.args.keys():
        description = request.args['description']
        similar_ex = ex.Search(description)
    else:
        description = ''
        similar_ex = ''
    return render_template('search.html', description=description, similar_ex=similar_ex)

@app.route('/image/<path>/<filename>')
@app.route('/image/<filename>')
def serveImage(filename, path=''):
    img = ImageHandler()
    return send_file(img.Image(filename, path))

@app.route('/backend/<command>', methods=['GET', 'POST'])
def generateEvent(command):
    extraArgs = {}
    if command == 'MERGE_EXPENSE' or command == 'MERGE_EXPENSE_COMMISSION':
        extraArgs[request.cookies.get('pinned_type') + '_merged'] = request.cookies.get('pinned_id')
    _generateEvent(command, request.args, extraArgs)
    return '200';

def _getFromArgs(args, value, default):
    if value in request.args.keys():
        return request.args[value]
    else:
        return default

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

