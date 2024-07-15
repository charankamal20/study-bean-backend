#!/bin/bash
export PATH=$PATH:/usr/local/go/bin

cd /src/study-bean-backend
git pull origin main
pmgo kill study-bean
pmgo start study-bean-backend study-bean
