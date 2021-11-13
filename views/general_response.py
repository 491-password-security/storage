class SuccessResponse:
    def __init__(self, body):
        self.status = 200
        self.msg = "Success"
        self.body = body


class FailureResponse:
    def __init__(self, msg):
        self.status = 400
        self.msg = msg
