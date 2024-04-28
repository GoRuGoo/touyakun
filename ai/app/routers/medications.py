from fastapi import APIRouter, UploadFile, Form
from app.controllers.gemini_recognition_controller import GeminiRecognitionController
from app.schemas.medications import Medications

router = APIRouter()


@router.post("/medications/", tags=["medications"], response_model=Medications)
async def recognize_medications(file: UploadFile):
    return await GeminiRecognitionController.get_medications(file)
