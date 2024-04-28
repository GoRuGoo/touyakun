from vertexai.preview.generative_models import GenerativeModel, Image
from app.schemas.medications import Medications, Medication
import json
from typing import AsyncIterable, Union
from fastapi import HTTPException


class GeminiRecognitionService:
    def __init__(self):
        # VertexAIの初期化処理
        self.model = GenerativeModel("gemini-pro-vision")

    async def get_medications(self, image: bytes) -> Union[
        "Medications",
        AsyncIterable["Medications"],
    ]:
        # routerから呼ばれる
        img = Image.from_bytes(image)
        json_dict = await self.recognize(img)
        # Medicationsオブジェクトに変換して返す
        medications = []
        for med in json_dict["medications"]:
            medications.append(Medication(**med))
        return Medications(medications=medications)

    async def recognize(self, img: Image):
        # VertexAIのAPIを呼ぶ
        response = await self.model.generate_content_async(
            [
                r"""この画像から、薬の名前、朝昼夜のうちいつ何錠の薬を何日間にわたって飲めばよいかを認識し、json形式で返してください。認識に失敗したときは、認識に失敗しました。というエラーメッセージを返してください。
        以下は認識に成功したときに返却するjson形式の例です。
        {
            "medications": [
                {
                    "name": "Medication 1",
                    "morning": false,
                    "afternoon": true,
                    "evening": true,
                    "dosage": 1,
                    "duration_days": 5
                },
                {
                    "name": "Medication 2",
                    "morning": true,
                    "afternoon": false,
                    "evening": true,
                    "dosage": 2,
                    "duration_days": 20
                }
            ]
        }
        """,
                img,
            ]
        )
        # ここでresponseをパースしてMedicationsオブジェクトに変換して返す
        # エラー処理
        if "認識に失敗しました。" in response.text:
            raise HTTPException(
                status_code=500,
                detail="Error: Failed to recognize medications from image. "
                "(Data returned from VertexAI API was empty. Could be due to invalid image or other reasons.)",
            )
        try:
            res = response.text.replace(" ```json\n", "").replace("\n```", "")
            json_dict = json.loads(res)
        except Exception as e:
            raise HTTPException(
                status_code=500,
                detail="Error: Failed to parse response from VertexAI. "
                "(Data returned from VertexAI API was invalid.)",
            ) from e
        return json_dict


gemini_recognition_service = GeminiRecognitionService()
