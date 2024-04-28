class RecognitionException(Exception):
    def __init__(self, msg: str):
        self.msg = msg


class MedicationsValidationException(Exception):
    def __init__(self, msg: str):
        self.msg = msg
