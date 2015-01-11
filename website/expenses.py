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

class Expenses:
    def __init__(self, config):
        self.redis = redis.Redis(config['redis_host'], config['redis_port'])
        template_path = os.path.join(os.path.dirname(__file__), 'templates')
        self.jinja_env = Environment(loader=FileSystemLoader(template_path),
                                 autoescape=True)
        self.url_map = Map([
            Rule('/', endpoint='expenses'),
            #Rule('/<short_id>', endpoint='follow_short_link'),
            #Rule('/<short_id>+', endpoint='short_link_details')
        ])


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
#        return Response('Hello World!')

    def get_cursor(self):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str
        cursor = conn.execute ('select count (*), classificationdef.name, sum(amount) from expenses, classifications, classificationdef where strftime(date) >= strftime(\'2015\') and expenses.eid = classifications.eid and classifications.cid = classificationdef.cid group by classifications.cid;')
        return cursor

    def get_cursor2(self):
        conn = sqlite3.connect('../expenses.db')
        conn.text_factory = str
        cursor = conn.execute ('select date, description, amount, classificationdef.name from expenses, classifications, classificationdef where strftime(date) >= strftime(\'2015\') and expenses.eid = classifications.eid and classifications.cid = classificationdef.cid order by date desc;')
        return cursor

    def on_expenses(self, request):
        error = None
        url = ''
#        if request.method == 'POST':
#            url = request.form['url']
#            if not is_valid_url(url):
#                error = 'Please enter a valid URL'
#            else:
#                short_id = self.insert_url(url)
#                return redirect('/%s+' % short_id)
        return self.render_template('expenses.html', error=error, url=url, cursor=self.get_cursor(), cursor2=self.get_cursor2())

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
