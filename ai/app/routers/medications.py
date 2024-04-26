from fastapi import APIRouter

router = APIRouter()


@router.get("/medications/", tags=["medications"])
async def read_router():
    return {"message": "Hello, world! Medications"}
