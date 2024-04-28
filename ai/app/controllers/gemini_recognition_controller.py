from fastapi import UploadFile
from app.services.gemini_recognition_service import GeminiRecognitionService
from app.schemas.medications import Medications
from app.exception.errors import MedicationsValidationException


class GeminiRecognitionController:
    ACCEPTABLE_CONTENT_TYPES = ["image/jpeg", "image/png"]  # jpegとpngのみ受け付ける

    @staticmethod
    async def get_medications(file: UploadFile) -> Medications:
        # routerから呼ばれる
        # バリデーション処理
        if file.size > 10485760:  # 10MB
            raise MedicationsValidationException(msg="File size is too large. Max size is 10MB.")
        if file.content_type not in GeminiRecognitionController.ACCEPTABLE_CONTENT_TYPES:
            raise MedicationsValidationException(
                msg="Content type is not supported. Supported content type is image/jpeg and image/png."
            )
        return await GeminiRecognitionService.get_medications(await file.read())
