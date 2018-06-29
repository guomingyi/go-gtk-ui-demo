#!/bin/bash

path=$1
ACTION=$2

BL=bootloader-v1.0.bin
FW=firmware-v1.0.bin
MT=metadata-v1.0.bin
FTM=ftm.bin
FTM_M=metadata-ftm.bin

BL_ADDR=0x08000000
MT_ADDR=0x08008000
FW_ADDR=0x08010000

function device_reset()
{
  st-flash reset
}

if [ "y$ACTION" = "yb" ] ;then
  echo -e "st-flash write $path/$BL $BL_ADDR. "
  st-flash write $path/$BL $BL_ADDR 2>&1 |tee $BL-log.txt
  ret=cat $BL-log.txt
  echo -e "$ret"
  device_reset
  exit 0
fi

if [ "y$ACTION" = "yf" ]; then
  echo -e " st-flash write  $path/$FW $FW_ADDR. "
  st-flash write  $path/$FW $FW_ADDR 2>&1 |tee $FW-log.txt
  ret=cat $FW-log.txt
  echo -e "$ret"
  device_reset
  exit 0
fi

if [ "y$ACTION" = "ym" ]; then
  echo -e " st-flash write  $path/$MT $MT_ADDR. "
  st-flash write  $path/$MT $MT_ADDR 2>&1 |tee $MT-log.txt
  ret=cat $MT-log.txt
  echo -e "$ret"
  device_reset
  exit 0
fi

if [ "y$ACTION" = "yftm-m" ]; then
  echo -e " st-flash write  $path/$FTM_M $MT_ADDR. "
  st-flash write  $path/$FTM_M $MT_ADDR 2>&1 |tee $FTM_M-log.txt
  ret=cat $FTM_M-log.txt
  echo -e "$ret"
  device_reset
  exit 0
fi

if [ "y$ACTION" = "yftm-f" ]; then
  echo -e " st-flash write  $path/$FTM $FW_ADDR. "
  st-flash write  $path/$FTM $FW_ADDR 2>&1 |tee $FTM-log.txt
  ret=cat $FTM-log.txt
  echo -e "$ret"
  device_reset
  exit 0
fi

if [ "y$ACTION" = "yreset-stlink" ]; then
  killall st-flash
  killall st-util
  exit 0
fi



