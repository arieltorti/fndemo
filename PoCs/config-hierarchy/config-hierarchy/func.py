import io
import json

from fdk import response


def handler(ctx, data: io.BytesIO=None):
    name = "World"
    input_param = ctx.Config().get("param", "name")

    try:
        body = json.loads(data.getvalue())
        name = body.get(input_param, name)
    except (Exception, ValueError) as ex:
        print(str(ex))

    return response.Response(
        ctx, response_data=json.dumps(
            {
                "message": "Hello {0}".format(name)
            }),
        headers={"Content-Type": "application/json"}
    )
