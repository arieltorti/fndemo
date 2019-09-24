from db import Db

def _extract_id(body):
    if body is None:
        return None
    return body.get("id", None)

def retrieve(body):
    id = _extract_id(body)

    if id is None:
        return Db.list()
    try:
        return Db.get(id)
    except KeyError:
        raise Exception("Invalid id")
