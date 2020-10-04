FROM python:3.8-slim-buster

COPY . /usr/src/neural
WORKDIR /usr/src/neural
RUN arch=$(uname -m); \
    if [ "$arch" = "aarch64" ] || [ "$arch" = "armv7l" ]; then \
        apt-get -y update && apt-get -y install wget; \
        wget https://github.com/Abathargh/python3-arm-wheels/releases/download/0.0.1/python3-"$arch".tar; \
        tar -xf python3-"$arch".tar; \
        for lib in numpy scipy scikit_learn pandas; do \
        	pip3 install python3-"$arch"/"$lib"*; \
        done; \
        rm python3-"$arch".tar; \
        rm -r python3-"$arch"; \
        apt-get -y remove wget; \
        apt-get -y autoremove; \
    fi
RUN pip3 install -r requirements.txt
EXPOSE 80

CMD ["gunicorn", "-b", "0.0.0.0:80", "app:app"]