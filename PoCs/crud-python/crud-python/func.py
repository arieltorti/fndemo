import io
import json


import sys
sys.path.insert(0, '/function')

import actions
from fdk import response

ROUTER_MAP = {
    (("POST", "crud/create"), actions.create),
    (("GET", "crud/get"), actions.retrieve),
    (("PUT", "crud/update"), actions.update),
    (("DELETE", "crud/delete"), actions.delete),
}

def get_trigger_source(url):
    # Urls are of the form:
    #   protocol://host:port/t/app/trigger...
    trigger_uri = url.split('/t/', 1)[1]  # app/trigger...
    return trigger_uri.split('/', 1)[1] # trigger...


def handler(ctx, data: io.BytesIO=None):
    method = ctx._method
    url = ctx._request_url

    body = None
    if data:
        try:
            body = json.loads(data.getvalue())
        except json.JSONDecodeError:
            pass

    if method is None or url is None:
            # Called using invoke
            return response.Response(
                ctx, response_data=json.dumps({ "error": "Called using invoke, try using the triggers" }),
                headers={"Content-Type": "application/json"}
            )

    trigger_source = get_trigger_source(url)

    for (r_method, r_url), func in ROUTER_MAP:
        if method == r_method and trigger_source == r_url:
            try:
                output = { "status": "ok", "value": func(body) }
            except Exception as e:
                output = { "error": str(e) }

            return response.Response(ctx, response_data=json.dumps(output))
    else:
        return response.Response(
            ctx, response_data=json.dumps(
                {
                    "error": "Method {} not allowed on url {}".format(method, url),
                    "debug": "method {} trigger_source {}".format(method, trigger_source)
                }
            )
    )
