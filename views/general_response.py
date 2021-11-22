from types import ClassMethodDescriptorType


class SuccessResponse:
    def __init__(self, body):
        self.success = True
        self.message = "Success"
        self.data = body
        self.code = 0


class FailureResponse:
    def __init__(self, msg):
        self.success = False
        self.message = msg
        self.data = {}
        self.code = -1

class Response:
    def __init__(self, success, message, data, code):
        self.success = True
        self.message = message
        self.data = data
        self.code = code

class AuthErrorResponse:
    def __init__(self, msg):
        self.success = False
        self.message = msg
        self.data = {}
        self.code = -2