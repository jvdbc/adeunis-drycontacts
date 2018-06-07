exports.handler = (event, context, callback) => {
    if(typeof event === 'undefined' || event === null) {
        console.error('Received wrong event ');
        callback('Received wrong event');
        return;
    }
    console.log('Received event ', JSON.stringify(event) );

    var DEFAULT_BATTERY_VALUE = "23"

    var data = event.data;
    console.log("data", data);
    var trueData;
    
    if ( data.startsWith("30") ) { //a keepalive frame
        trueData = "50"+ DEFAULT_BATTERY_VALUE+ getFrameCounter(data.substr(2, 2)) +"000000000000000000";
    } else if ( data.startsWith("40") ) { // a state frame
        trueData = "1"+ getInputsState(data.substr(data.length-2 , 2)) + DEFAULT_BATTERY_VALUE+ getFrameCounter(data.substr(2, 2)) +"0";
    } else {
        callback('Wrong data type. data = '+ data);
        return;
    }

    var payload = {
        time:event.time, 
        device:event.device, 
        duplicate:event.duplicate, 
        snr:event.snr, 
        rssi:event.rssi, 
        avgSignal:event.avgSignal, 
        station:event.station, 
        data:trueData, 
        lat:event.lat, 
        lng:event.lng, 
        seqNumber:event.seqNumber
    }
    postToSmartConnect(payload);
    callback(null, "OK");
}

function getFrameCounter(statusData){
    var intValue = parseInt("0x"+statusData, 16);
    var counter = (intValue & 0xE0) >>> 5;
    var counterString = counter.toString(16);
    counterString = counterString.length == 1 ? "0"+counterString : counterString;
    console.log("Frame counter = 0x"+counterString);
    return counterString;
}

function getInputsState(inputStateData){
    var data = 0;
    var intValue = parseInt("0x"+inputStateData, 16);
    var tor1State = intValue & 0x01;
    var tor2State = (intValue & 0x04) >>> 2;
    var tor3State = (intValue & 0x10) >>> 4;
    console.log("tor3 tor2 tor1 = "+tor3State+tor2State+tor1State);
    //  7    6    5    4    3    2    1    0
    // tor1 tor2 tor3 tor4 tor5 tor6 tor7 tor8

    //tor1 => tor4; tor2 => tor3; tor3 => tor5
    data = data | (tor1State << 4);
    data = data | (tor2State << 5);
    data = data | (tor3State << 3);

    var stringData = data.toString(16);
    stringData =  stringData.length == 1 ? "0"+stringData : stringData;
    console.log("InputsStateData = 0x"+stringData);
    return stringData;
}



const https = require('https');

function postToSmartConnect(data) {
    console.log('Data to send: ', JSON.stringify(data));

    var post_options = {
      host: process.env.smartconnectcallbackhost, // 'connector-demoenv.devinno.fr'
      path: process.env.smartconnectcallbackIldPath, // '/ild/data'
      method: 'POST',
      headers: {
          'Content-Type': 'application/json'
      }
    };

    // POST HTTPS request definition 
    var post_req = https.request(post_options, function(res) {
    res.setEncoding('utf8');
    res.on('data', function (chunk) {
            console.log('Response: ' + chunk);
        });
    });

    //post to smartconnect
    post_req.write(JSON.stringify(data, null, 2));
    post_req.end();
}