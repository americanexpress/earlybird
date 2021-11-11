### Converting PDF and rtf files
The standalone Earlybird binary cannot convert PDF or rtf files to binary. To enable scanning of these file types, install the following dependencies:
```bash
brew install poppler-utils wv unrtf tidy
go get github.com/JalfResi/justext
```