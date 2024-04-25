FROM python:latest

RUN mkdir /usr/src/ai

WORKDIR /usr/src/ai

RUN pip install --upgrade pip

RUN pip install FastAPI uvicorn

COPY . .

EXPOSE 8080

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080", "--reload"]
