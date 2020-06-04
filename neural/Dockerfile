FROM python:3.7

COPY . /usr/src/neural
WORKDIR /usr/src/neural
RUN pip3 install -r requirements.txt
EXPOSE 80

CMD ["gunicorn", "-b", "0.0.0.0:80", "app:app"]