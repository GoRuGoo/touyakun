from fastapi import FastAPI
from .routers import medications
from .exception.exception_handler import register_exception

app = FastAPI()
app.include_router(medications.router)
register_exception(app)


@app.get("/")
async def read_root():
    return {"message": "Hello, world! Python"}
