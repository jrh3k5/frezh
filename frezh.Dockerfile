FROM golang:1.23.1

# Install tesseract
RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/
RUN apt-get install -y -qq \
  tesseract-ocr-eng

COPY go.mod .
COPY go.sum .
# Only download updates if modules files have changed
RUN go mod download

WORKDIR /app
COPY . .
CMD ["go", "run", "main.go"]