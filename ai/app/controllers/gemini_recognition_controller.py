from fastapi import UploadFile, HTTPException, status
from app.services.gemini_recognition_service import gemini_recognition_service
from app.schemas.medications import Medications
acceptable_content_types = ["image/jpeg", "image/png"] # jpegとpngのみ受け付ける


class GeminiRecognitionController:
    def __init__(self) -> None:
        pass

    def get_medications(self, file:UploadFile) -> Medications:
        # routerから呼ばれる
        # バリデーション処理
        if file.size > 10485760: # 10MB
            raise HTTPException(status_code=status.HTTP_406_NOT_ACCEPTABLE, detail=f"File size is too large. Max size is 10MB.") 
        if file.content_type not in acceptable_content_types:
            raise HTTPException(status_code=status.HTTP_406_NOT_ACCEPTABLE, detail=f"Content type is not supported. Supported content type is image/jpeg.")
        return gemini_recognition_service.get_medications(file.read())
    

gemini_recognition_controller = GeminiRecognitionController()