FROM python:3.8

RUN apt-get update
RUN apt-get install -y libdbus-1-dev
#RUN apt-get install -y python3-dbus

RUN pip install dbus-python
# RUN apt-get install -y python3-dbus

WORKDIR /usr/src/app

# Install Adafruit Python BLE library
RUN git clone https://github.com/adafruit/Adafruit_Python_BluefruitLE.git
RUN cd Adafruit_Python_BluefruitLE; python setup.py install
RUN rm -rf Adafruit_Python_BluefruitLE

RUN apt-get install -y pkg-config libcairo2-dev gcc libgirepository1.0-dev
RUN pip install gobject PyGObject

COPY collect.py .

RUN python3 collect.py

CMD ["python3", "collect.py"]
