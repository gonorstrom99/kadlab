@REM @echo off
@REM setlocal enabledelayedexpansion

@REM set START_PORT=8001
@REM set NUM_CONTAINERS=50

@REM for /l %%i in (1,1,%NUM_CONTAINERS%) do (
@REM     set /a PORT=%START_PORT% + %%i - 1
@REM     set ID=container%%i
@REM     echo Starting !ID! on port !PORT!
@REM     start cmd /k docker run -it --network host -e ID=!ID! -e ADDRESS=127.0.0.1:!PORT! kadlab
@REM )

@REM endlocal



@echo off
setlocal enabledelayedexpansion

rem Starting port number
set START_PORT=8001

rem Number of containers to create
set NUM_CONTAINERS=9

rem Loop to create the containers
for /l %%i in (1,1,%NUM_CONTAINERS%) do (
    set /a PORT=!START_PORT! + %%i - 1
    set NAME=container%%i

    echo Starting !ID! on port !PORT!
    
    @REM docker run -d -it --network host -e ADDRESS=127.0.0.1:!PORT! kadlab
    docker run -d -it --name !NAME! --network host -e ID=!NAME! -e ADDRESS=127.0.0.1:!PORT! kadlab
    @REM timeout /t 1 /nobreak >nul

)

endlocal
