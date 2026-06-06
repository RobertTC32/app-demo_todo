FROM alpine:3.21.2
RUN mkdir -p /app
COPY ./temp/app/main.exe /app/main.exe
WORKDIR /app
EXPOSE 80/tcp
RUN chmod +x main.exe
ENTRYPOINT [ "./main.exe" ]
