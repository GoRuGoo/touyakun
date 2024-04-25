FROM python:3.12.3

RUN mkdir /usr/src/ai

WORKDIR /usr/src/ai

COPY . .

RUN pip install --upgrade pip

# Poetryのインストール
RUN curl -sSL https://install.python-poetry.org | python3 -

ENV PATH /root/.local/bin:$PATH

# Poetryが仮想環境を作らないようにする
RUN poetry config virtualenvs.create false

RUN poetry install --no-root


EXPOSE 8080

CMD ["poetry" ,"run","uvicorn", "main:app", "--host", "ai", "--port", "8080", "--reload"]
