from fastapi import APIRouter, UploadFile

router = APIRouter()


@router.post("/medications/", tags=["medications"])
async def recognize_medications(file: UploadFile):
    return {"filename": file.filename}
    