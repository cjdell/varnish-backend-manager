#!/bin/bash

mkdir /opt/varnish-backend-manager
cp ../../bin/varnish-backend-manager /opt/varnish-backend-manager
cp ../start.sh /opt/varnish-backend-manager
cp ../varnish-restart.sh /opt/varnish-backend-manager

cp varnish-backend-manager.service /etc/systemd/system/
systemctl enable varnish-backend-manager
