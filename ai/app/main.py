from fastapi import FastAPI
from .routers import medications

app = FastAPI()
app.include_router(medications.router)


@app.get("/")
async def read_root():
    return {"message": "Hello, world! Python"}
