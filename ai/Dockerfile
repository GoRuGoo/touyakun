FROM python:latest

RUN mkdir /usr/src/ai

WORKDIR /usr/src/ai

RUN pip install --upgrade pip

COPY . .

EXPOSE 8080

CMD ["uvicorn","main:app","--port=8080","--host=ai","--reload"]
