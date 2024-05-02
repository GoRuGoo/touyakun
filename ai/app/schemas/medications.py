from typing import List
from pydantic import BaseModel


class Medication(BaseModel):
    name: str
    isMorning: bool
    isAfternoon: bool
    isEvening: bool
    duration: int
    amount: int


class Medications(BaseModel):
    medications: List[Medication]
