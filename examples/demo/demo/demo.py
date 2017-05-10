class Demo(object):

    def __init__(self, stra, uc):
        self._stra = stra
        self._uc = uc
        print(stra, uc)

    def order(self, event):
        print(event)
        print()
        print("version:7")
        return []

    def timing(self, event):
        print(event)
        return [
            {'type': "order", 'params': {'code': '000001'}}
        ]
