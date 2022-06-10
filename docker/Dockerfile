FROM python:2

WORKDIR ./
COPY conf ./conf
ADD mybin ./
RUN chmod +x ./mybin
EXPOSE 80
ENTRYPOINT [ "./mybin" ]