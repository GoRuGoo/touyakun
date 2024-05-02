from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from .routers import medications
from .exception.exception_handler import register_exception

app = FastAPI()
app.include_router(medications.router)
register_exception(app)

origins = [
    "http://api:8080",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
async def read_root():
    return {"message": "Hello, world! Python"}
