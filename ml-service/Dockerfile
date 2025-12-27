#FROM python:3.11-slim
#
#WORKDIR /app
#
## Установка системных зависимостей
#RUN apt-get update && apt-get install -y \
#    gcc \
#    g++ \
#    libgl1-mesa-glx \
#    libglib2.0-0 \
#    libsm6 \
#    libxext6 \
#    libxrender-dev \
#    libgomp1 \
#    && rm -rf /var/lib/apt/lists/*
#
## Копирование файлов зависимостей
#COPY requirements.txt .
#
## Установка Python зависимостей
#RUN pip install --no-cache-dir --upgrade pip && \
#    pip install --no-cache-dir -r requirements.txt
#
## Копирование исходного кода
#COPY . .
#
## Копирование модели (если она находится локально)
## COPY best.pt /app/best.pt
#
## Создание директории для временных файлов
#RUN mkdir -p /tmp
#
## Установка переменных окружения
#ENV PYTHONUNBUFFERED=1
#ENV PYTHONPATH=/app
#
## Запуск сервиса
#CMD ["python", "-u", "main.py"]

FROM python:3.10-slim

WORKDIR /app

# Создаем необходимые директории
RUN mkdir -p model logs

RUN pip install --upgrade pip

# Обновляем и устанавливаем системные зависимости
RUN apt-get update && apt-get install -y \
    gcc \
    libjpeg-dev \
    libxslt-dev \
    libpq-dev \
    ffmpeg \
    libsm6 \
    libxext6 \
    libglib2.0-0 \
    libxrender-dev \
    libgomp1 \
    git \
    build-essential \
    libmariadb-dev \
    libmariadb-dev-compat \
    gettext \
    cron \
    openssh-client \
    flake8 \
    locales \
    curl \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Копируем и устанавливаем зависимости Python
COPY requirements.txt ./

# Устанавливаем PyTorch (CPU версия)
RUN #pip install torch torchvision --index-url https://download.pytorch.org/whl/cpu
RUN pip install --no-cache-dir torch torchvision --extra-index-url https://download.pytorch.org/whl/cpu
RUN #pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cpu

# Устанавливаем остальные зависимости
RUN pip install -r requirements.txt --no-cache-dir

# Копируем исходный код
COPY . ./

# Создаем non-root пользователя для безопасности
RUN useradd -m -u 1000 appuser && \
    chown -R appuser:appuser /app/

USER appuser

# Устанавливаем переменные окружения
ENV PYTHONUNBUFFERED=1
ENV PYTHONPATH=/app
ENV NATS_URL=nats://nats:4222
ENV MODEL_PATH=/app/model/best.pt

# Открываем порт если нужно (для отладки)
EXPOSE 8000

# Запускаем сервис
CMD ["python", "-u", "main.py"]