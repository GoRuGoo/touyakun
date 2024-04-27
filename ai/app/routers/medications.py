from fastapi import APIRouter, UploadFile
from app.services.gemini_recognition_service import gemini_recognition_service
from app.schemas.medications import Medications

router = APIRouter()


@router.post("/medications/", tags=["medications"], response_model=Medications)
async def recognize_medications(file: UploadFile):
    return gemini_recognition_service.get_medications(file.read())
    