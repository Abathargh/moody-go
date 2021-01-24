FROM python:3.8-slim-buster

COPY . /usr/src/neural
WORKDIR /usr/src/neural
RUN arch=$(uname -m); \
    if [ "$arch" = "aarch64" ] || [ "$arch" = "armv7" ]; then \
        apt-get -y update && apt-get -y install wget libgomp1;
        if [ "$arch" = "armv7" ]; then \
            wget https://raw.githubusercontent.com/abathargh/python3-arm-wheels/master/wheels/numpy/python3.8/numpy-1.20.0.dev0%2Be6b8b19-cp38-cp38-linux_armv7l.whl && \
            pip3 install numpy-1.20.0.dev0%2Be6b8b19-cp38-cp38-linux_armv7l.whl; \
        fi \ 
        wget https://raw.githubusercontent.com/abathargh/python3-arm-wheels/master/wheels/scikit-learn/python3.8/scikit_learn-1.0.dev0-cp38-cp38-linux_aarch64.whl && \
        pip3 install scipy scikit_learn-1.0.dev0-cp38-cp38-linux_aarch64.whl && apt-get -y remove wget && \
        apt-get -y autoremove; \
    fi
RUN pip3 install -r requirements.txt
EXPOSE 80

CMD ["gunicorn", "-b", "0.0.0.0:80", "app:app"]