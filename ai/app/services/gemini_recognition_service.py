from vertexai.preview.generative_models import GenerativeModel, Image
from app.schemas.medications import Medications, Medication
import json
from typing import AsyncIterable, Union
from app.exception.errors import RecognitionException
from urllib import request
import os


class GeminiRecognitionService:
    __model = GenerativeModel("gemini-pro-vision")  # 未来の自分が治す

    @staticmethod
    async def get_medications(image: bytes) -> Union[
        "Medications",
        AsyncIterable["Medications"],
    ]:
        # routerから呼ばれる
        img = Image.from_bytes(image)
        json_dict = await GeminiRecognitionService.recognize(img)
        # Medicationsオブジェクトに変換して返す
        medications = []
        for med in json_dict["medications"]:
            medications.append(Medication(**med))
        return Medications(medications=medications)

    @staticmethod
    async def get_medications_url(messageId: str):
        url = "https://api-data.line.me/v2/bot/message/" + messageId + "/content"
        headers = {"Authorization": "Bearer " + os.environ["CHANNEL_TOKEN"]}
        get_req = request.Request(url, headers=headers)
        with request.urlopen(get_req) as response:
            img = Image.from_bytes(response.read())
        json_dict = await GeminiRecognitionService.recognize(img)
        # Medicationsオブジェクトに変換して返す
        medications = []
        for med in json_dict["medications"]:
            medications.append(Medication(**med))
        return Medications(medications=medications)

    @staticmethod
    async def recognize(img: Image):
        # VertexAIのAPIを呼ぶ
        response = await GeminiRecognitionService.__model.generate_content_async(
            [
                """この画像から、薬の名前、朝昼夜のうちいつ何錠の薬を何日間にわたって飲めばよいかを認識し、json形式で返してください。
        以下は返却するjson形式の例です。
        {
            "medications": [
                {
                    "name": "Medication 1",
                    "isMorning": false,
                    "isAfternoon": true,
                    "isEvening": true,
                    "duration": 14,
                    "amount": 1
                },
                {
                    "name": "Medication 2",
                    "isMorning": true,
                    "isAfternoon": false,
                    "isEvening": false,
                    "duration": 7,
                    "amount": 2
                }
            ]
        }
        """,
                img,
            ]
        )
        # ここでresponseをパースしてMedicationsオブジェクトに変換して返す
        # エラー処理
        try:
            res = response.text.replace(" ```json\n", "").replace("\n```", "")
            json_dict = json.loads(res)
        except Exception as e:
            raise RecognitionException(
                msg="Failed to parse response from VertexAI. (Data returned from VertexAI API was invalid.)"
            ) from e
        return json_dict
