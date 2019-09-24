from db import Db

def _extract_user(body):
    try:
        user = {
            "name": body["name"],
            "last_name": body["last_name"]
        }
        return user
    except (ValueError, KeyError):
        raise ValueError("name or last_name not provided")

def create(body):
    user = _extract_user(body)

    Db.create(user)
