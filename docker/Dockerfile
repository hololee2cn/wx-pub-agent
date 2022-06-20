FROM python:2

WORKDIR ./
ADD mybin ./
RUN apt-get -q update && apt-get -qy install netcat
RUN chmod +x ./mybin
EXPOSE 80
EXPOSE 50051
CMD ["./mybin","-h"]