@echo off
set "CODEHEART_OPERATING_KIT_CLI=1"
set "PYTHONPATH=%LOCALAPPDATA%\Codeheart\OperatingKit\lib;%PYTHONPATH%"
python -m codeheart_operating_kit.cli %*
