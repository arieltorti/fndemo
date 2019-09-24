from db import Db

def _extract_id_and_user(body):
    try:
        name = body.get("name", None)
        last_name = body.get("last_name", None)
        id = body["id"]

        user = {}
        if name:
            user["name"] = name
        if last_name:
            user["last_name"] = last_name

        return id, user
    except (ValueError, KeyError):
        raise ValueError("id not provided")


def update(body):
    id, user = _extract_id_and_user(body)

    try:
        Db.update(id, user)
    except KeyError:
        raise Exception("Invalid id")
