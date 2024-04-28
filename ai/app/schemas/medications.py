from typing import List
from pydantic import BaseModel


class Medication(BaseModel):
    name: str
    morning: bool
    afternoon: bool
    evening: bool
    dosage: int
    duration_days: int


class Medications(BaseModel):
    medications: List[Medication]
