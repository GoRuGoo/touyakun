from app.schemas.medications import Medications, Medication


class GeminiRecognitionService:
    def __init__(self):
        # VertexAIの初期化処理
        pass
    def get_medications(self, image:bytes)->Medications:
        # routerから呼ばれる
        # とりあえずダミーデータを返す
        return Medications(medications=[Medication(name="test", morning=True, afternoon=True, evening=True, dosage=1, duration_days=1), Medication(name="test2", morning=True, afternoon=False, evening=False, dosage=2, duration_days=5)])
    
    def recognize(self):
        # VertexAIのAPIを呼ぶ
        pass

gemini_recognition_service = GeminiRecognitionService()