from db import Db

def _extract_id(body):
    try:
        return body["id"]
    except (ValueError, KeyError):
        raise ValueError("id not provided")


def delete(body):
    id = _extract_id(body)

    Db.delete(id)
