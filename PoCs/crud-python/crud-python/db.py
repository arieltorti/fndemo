class Singleton(object):
    _instance = None
    def __new__(class_, *args, **kwargs):
        if not isinstance(class_._instance, class_):
            class_._instance = object.__new__(class_, *args, **kwargs)
        return class_._instance

class MemDatabase(Singleton):
    def __init__(self, data={}):
        self._data = data
        self.index = 0

    def create(self, user):
        self._data[self.index] = user
        self.index += 1

    def delete(self, id):
        del self._data[id]

    def update(self, id, user):
        self._data[id].update(user)

    def get(self, id):
        return self._data[id]

    def list(self):
        return self._data

Db = MemDatabase()
