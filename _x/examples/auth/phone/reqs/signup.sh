#!/bin/zsh

http -v \
  POST http://localhost:7450/api/auth/phone/signup \
  phone_number="+11234567890"
