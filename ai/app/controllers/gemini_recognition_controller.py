from fastapi import UploadFile, HTTPException, status
from app.services.gemini_recognition_service import gemini_recognition_service
from app.schemas.medications import Medications


class GeminiRecognitionController:
    ACCEPTABLE_CONTENT_TYPES = ["image/jpeg", "image/png"] # jpegとpngのみ受け付ける

    def __init__(self) -> None:
        pass
    
    @staticmethod
    def get_medications(file: UploadFile) -> Medications:
        # routerから呼ばれる
        # バリデーション処理
        if file.size > 10485760: # 10MB
            raise HTTPException(status_code=status.HTTP_406_NOT_ACCEPTABLE, detail=f"File size is too large. Max size is 10MB.") 
        if file.content_type not in GeminiRecognitionController.ACCEPTABLE_CONTENT_TYPES:
            raise HTTPException(status_code=status.HTTP_406_NOT_ACCEPTABLE, detail=f"Content type is not supported. Supported content type is image/jpeg and image/png.")
        return gemini_recognition_service.get_medications(file.read())
