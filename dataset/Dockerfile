FROM python:3.8-slim-buster

COPY . /usr/src/neural
WORKDIR /usr/src/neural
RUN arch=$(uname -m); \
    if [ "$arch" = "aarch64" ] || [ "$arch" = "armv7l" ]; then \
        apt-get -y update && apt-get -y install wget libgomp1; \
        if [ "$arch" = "armv7l" ]; then \
            wget https://raw.githubusercontent.com/abathargh/python3-arm-wheels/master/wheels/numpy/python3.8/numpy-1.20.0.dev0%2Be6b8b19-cp38-cp38-linux_armv7l.whl && \
            wget https://raw.githubusercontent.com/abathargh/python3-arm-wheels/master/wheels/scikit-learn/python3.8/scikit_learn-0.24.dev0-cp38-cp38-linux_armv7l.whl && \
            wget https://raw.githubusercontent.com/abathargh/python3-arm-wheels/master/wheels/scipy/python3.8/scipy-1.6.0.dev0%2Be4e3b89-cp38-cp38-linux_armv7l.whl && \
            wget https://raw.githubusercontent.com/abathargh/python3-arm-wheels/master/wheels/pandas/python3.8/pandas-1.2.0.dev0%2B600.ga975a754c-cp38-cp38-linux_armv7l.whl && \
            pip3 install $(ls | grep .whl$); \
        else \
            wget https://raw.githubusercontent.com/abathargh/python3-arm-wheels/master/wheels/scikit-learn/python3.8/scikit_learn-1.0.dev0-cp38-cp38-linux_aarch64.whl && \
            pip3 install scipy scikit_learn-1.0.dev0-cp38-cp38-linux_aarch64.whl; \
        fi; \
        apt-get -y remove wget && apt-get -y autoremove; rm -f *.whl; \
    fi
RUN pip3 install -r requirements.txt && sh conf.sh

EXPOSE 80
CMD ["gunicorn", "-b", ":80", "app:app"]