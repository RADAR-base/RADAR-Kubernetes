FROM docker.io/spark:python3
USER root
RUN pip install --no-cache-dir pyiceberg
RUN mkdir -p /opt/spark/scripts
COPY output-restructure.py /opt/spark/scripts
USER spark
WORKDIR /opt/spark/scripts
