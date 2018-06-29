
@set path=%1
@set action=%2

@set BL=bootloader-v1.0.bin
@set FW=firmware-v1.0.bin
@set MT=metadata-v1.0.bin
@set FTM=ftm.bin
@set FTM_M=metadata-ftm.bin

@set BL_ADDR=0x08000000
@set MT_ADDR=0x08008000
@set FW_ADDR=0x08010000

@if {"%action%"} == {"b"} (
    @echo %path%\%BL% %BL_ADDR%
    ST-LINK_CLI.exe -c SWD UR LPM -CmpFile %path%\%BL% %BL_ADDR%
) 

@if {"%action%"} == {"m"} (
    @echo %path%\%MT% %MT_ADDR%
    ST-LINK_CLI.exe -c SWD UR LPM -CmpFile %path%\%MT% %MT_ADDR%
) 

@if {"%action%"} == {"ftm-m"} (
    @echo %path%\%FTM_M% %MT_ADDR% 
    ST-LINK_CLI.exe -c SWD UR LPM -CmpFile %path%\%FTM_M% %MT_ADDR%
) 

@if {"%action%"} == {"f"} (
    @echo %path%\%FW% %FW_ADDR% 
    ST-LINK_CLI.exe -c SWD UR LPM -CmpFile %path%\%FW% %FW_ADDR%
) 

@if {"%action%"} == {"ftm-f"} (
    @echo %path%\%FTM% %FW_ADDR% 
    ST-LINK_CLI.exe -c SWD UR LPM -CmpFile %path%\%FTM% %FW_ADDR%
) 

@if {"%action%"} == {"all"} (
    ST-LINK_CLI.exe -c SWD UR LPM -CmpFile %path%\%BL% %BL_ADDR%
	ST-LINK_CLI.exe -c SWD UR LPM -CmpFile %path%\%MT% %MT_ADDR%
	ST-LINK_CLI.exe -c SWD UR LPM -CmpFile %path%\%FW% %FW_ADDR%
) 

@if {"%action%"} == {"ftm-all"} (
    ST-LINK_CLI.exe -c SWD UR LPM -CmpFile %path%\%BL% %BL_ADDR%
	ST-LINK_CLI.exe -c SWD UR LPM -CmpFile %path%\%FTM_M% %MT_ADDR%
	ST-LINK_CLI.exe -c SWD UR LPM -CmpFile %path%\%FTM% %FW_ADDR%
) 

