import serial

def sendATOCommand(ser):
    print('Send ATO')
    ser.write(b'ATO\r\n')
    line = ser.readline()
    print(line)
    if line == b'O\r\n' :
        return True
    else :
        print('ERROR')
        return False    

def sendWriteCommand(ser, register, value):
    print('Send Write command')
    ser.write(b'ATS' + register + b'=' + value + b'\r\n')
    line = ser.readline()
    print(line)
    if line == b'O\r\n' :
        return True
    else :
        print('ERROR')
        return False

def sendStoreCommand(ser):
    print('Send AT&W')
    ser.write(b'AT&W\r\n')
    line = ser.readline()
    print(line)
    if line == b'O\r\n' :
        return True
    else :
        print('ERROR')
        return False  

def initSequence(ser):
    print('Send Init sequence')
    ser.write(b'+++')
    line = ser.readline()
    print(line)
    if line != b'CONNECTING...\r\n':
        if sendATOCommand(ser) == True :
            initSequence(ser)
            return True
        else :
            print('ERROR')
            return False       
    else :
        line = ser.readline()
        print(line)
        if line == b'CM\r\n' :
            print('Init sequence OK')
            return True
        else :
            print('ERROR')
            return False
    return

def configurationSequence(ser):
    print('Start configuration sequence')
    if sendWriteCommand(ser, b'300', b'254') == False:
        return False
    if sendWriteCommand(ser, b'301', b'255') == False:
        return False
    if sendWriteCommand(ser, b'302', b'46') == False:
        return False
    if sendWriteCommand(ser, b'303', b'46') == False:
        return False   
    if sendWriteCommand(ser, b'304', b'46') == False:
        return False
    if sendWriteCommand(ser, b'305', b'0') == False:
        return False
    if sendWriteCommand(ser, b'307', b'0') == False:
        return False
    if sendStoreCommand(ser) == False:
        return False
    if sendATOCommand(ser) == False:
        return False
    else :
        return True
    
ser = serial.Serial(port='COM6',
                    baudrate=115200,
                    parity=serial.PARITY_NONE,
                    stopbits=1,
                    bytesize=8,
                    timeout=1
                    )
if ser.is_open :
    if initSequence(ser) == True :
        if configurationSequence(ser) == True :
            print('Configuraton OK')
        else :
            print('Configuration KO')
    else :
        print('Unable to configure')
    
    ser.close()
print('End')
