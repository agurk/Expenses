#!/usr/bin/python

import os
import redis
import urlparse
from werkzeug.wrappers import Request, Response
from werkzeug.routing import Map, Rule
from werkzeug.exceptions import HTTPException, NotFound
from werkzeug.wsgi import SharedDataMiddleware
from werkzeug.utils import redirect
from jinja2 import Environment, FileSystemLoader
import sqlite3
from MonthView import MonthView
from ItemView import ItemView
from Search import Search
from MonthGraph import MonthGraph
from BackendMessenger import BackendMessenger
import time

class Expenses:
    def __init__(self, config):
        self.redis = redis.Redis(config['redis_host'], config['redis_port'])
        template_path = os.path.join(os.path.dirname(__file__), 'templates')
        self.jinja_env = Environment(loader=FileSystemLoader(template_path),
                                 autoescape=True)
        self.url_map = Map([
            Rule('/', endpoint='expenses'),
            Rule('/expenses', endpoint='expenses'),
            Rule('/expense_details', endpoint='expense_details'),
            Rule('/search', endpoint='search'),
            Rule('/backend/<command>', endpoint='backend'),
        ])
        self.BackendMessenger = BackendMessenger()

    def render_template(self, template_name, **context):
        t = self.jinja_env.get_template(template_name)
        return Response(t.render(context), mimetype='text/html')

    def dispatch_request(self, request):
        adapter = self.url_map.bind_to_environ(request.environ)
        try:
            endpoint, values = adapter.match()
            return getattr(self, 'on_' + endpoint)(request, **values)
        except HTTPException, e:
            return e

    def on_backend(self, request, command):
        response = self.BackendMessenger.ProcessRequest(command, request.args)
        return Response(response);

    def on_search(self, request):
        search = Search()
        if 'description' in request.args.keys():
            description = request.args['description']
            similar_ex = search.SimilarExpenses(description)
        else:
            description = ''
            similar_ex = ''
        return self.render_template('search.html', description=description, similar_ex=similar_ex)

    def on_expense_details(self, request):
        idno = request.args['eid']
        expense = ItemView(idno)
        for i in expense.RawStr():
            print i
        return self.render_template('expense.html',rawData=expense.RawStr(), classifications=expense.Classifications(), classification=expense.Classification(), amount=expense.Amount(),eid=idno)

    def on_expenses(self, request):
        if 'date' in request.args.keys():
            mv = MonthView(request.args['date'])
            mg = MonthGraph(request.args['date'])
        else:
            mv = MonthView(time.strftime("%Y-%m-%d"))
            mg = MonthGraph(time.strftime("%Y-%m-%d"))
        return self.render_template('expenses.html', cursor=mv.OverallExpenses(), cursor2=mv.IndividualExpenses(), previous_month=mv.PreviousMonth(), next_month=mv.NextMonth(), total_amount=mv.TotalAmount(), month_name=mv.MonthName(),month_graph=mg.Graph())

    def wsgi_app(self, environ, start_response):
        request = Request(environ)
        response = self.dispatch_request(request)
        return response(environ, start_response)

    def __call__(self, environ, start_response):
        return self.wsgi_app(environ, start_response)

def create_app(redis_host='localhost', redis_port=6379, with_static=True):
    app = Expenses({
        'redis_host':       redis_host,
        'redis_port':       redis_port
    })
    if with_static:
        app.wsgi_app = SharedDataMiddleware(app.wsgi_app, {
            '/static':  os.path.join(os.path.dirname(__file__), 'static')
        })
    return app

if __name__ == '__main__':
    from werkzeug.serving import run_simple
    app = create_app()
    run_simple('127.0.0.1', 5000, app, use_debugger=True, use_reloader=True)
