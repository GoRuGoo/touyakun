from fastapi import APIRouter, UploadFile
from ..services.gemini_recognition_service import gemini_recognition_service

router = APIRouter()


@router.post("/medications/", tags=["medications"])
async def recognize_medications(file: UploadFile):
    return {"filename": file.filename}
    