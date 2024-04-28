from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse
from app.exception.errors import RecognitionException


def register_exception(app: FastAPI):
    @app.exception_handler(RecognitionException)
    async def handle_recognition_exception(request: Request, exc: RecognitionException):
        # TODO: ENVIRONMENTがdevならエラー内容を返す、それ以外はコードのみ返す切り分け処理の実装
        return JSONResponse(status_code=500, content={"error": exc.msg})
