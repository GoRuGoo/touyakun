from fastapi import APIRouter, UploadFile, Form
from app.controllers.gemini_recognition_controller import gemini_recognition_controller
from app.schemas.medications import Medications

router = APIRouter()


@router.post("/medications/", tags=["medications"], response_model=Medications)
async def recognize_medications(file: UploadFile):
    return gemini_recognition_controller.get_medications(file)
    