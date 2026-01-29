#!/bin/sh
set -e

AUDIO_DIR="internal/quote/data/audio"

if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

ZIP_SOURCE="${VOICE_ZIP_URL:?VOICE_ZIP_URL is not set. Create a .env file with VOICE_ZIP_URL=<url or path>}"

if [ -d "$AUDIO_DIR" ]; then
    echo "Audio directory already exists at $AUDIO_DIR, skipping download."
    exit 0
fi

if [ -f "$ZIP_SOURCE" ]; then
    echo "Extracting from local file: $ZIP_SOURCE"
    mkdir -p internal/quote/data
    unzip -qo "$ZIP_SOURCE" -d /tmp/voice
    mv /tmp/voice/voice "$AUDIO_DIR"
    rm -rf /tmp/voice
else
    echo "Downloading voice files..."
    curl -fSL -o /tmp/voice.zip "$ZIP_SOURCE"
    echo "Extracting..."
    mkdir -p internal/quote/data
    unzip -qo /tmp/voice.zip -d /tmp/voice
    mv /tmp/voice/voice "$AUDIO_DIR"
    rm -rf /tmp/voice.zip /tmp/voice
fi

echo "Done. Audio files extracted to $AUDIO_DIR"
