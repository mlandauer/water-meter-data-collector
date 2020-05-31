FROM python:3.8

RUN apt-get update && apt-get install -y libdbus-1-dev libgirepository1.0-dev
RUN pip install dbus-python PyGObject

WORKDIR /usr/src/app

# Install Adafruit Python BLE library
RUN git clone https://github.com/donatieng/Adafruit_Python_BluefruitLE.git
RUN cd Adafruit_Python_BluefruitLE; git checkout af46b05cbcfd82110c8bbd08ed3d483de128fed1
RUN cd Adafruit_Python_BluefruitLE; python setup.py install
RUN rm -rf Adafruit_Python_BluefruitLE

COPY collect.py .

CMD ["python3", "collect.py"]
