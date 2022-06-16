FROM python:2

WORKDIR ./
ADD mybin ./
RUN chmod +x ./mybin
EXPOSE 80
EXPOSE 50051
CMD ["/mybin","-h"]