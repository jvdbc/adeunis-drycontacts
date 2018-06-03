exports.handler = (event, context, callback) => {
    if(typeof event === 'undefined' || event === null) {
        console.error('Received wrong event ');
        callback('Received wrong event');
        return;
    }
    console.log('Received event ', JSON.stringify(event) );
    var DEFAULT_BATTERY_VALUE = "0DAC";

    var time =  Math.round(new Date(event.timestamp).getTime()/1000);
    var deveui = event.device_properties.deveui;
    var data = event.payload_cleartext;
    var trueData;

    if ( data.startsWith("30") ) { //a keepalive frame
        trueData = "110A0050000641050104"+ DEFAULT_BATTERY_VALUE +"04";
    } else if ( data.startsWith("40") ) { // a state frame
        trueData = "110A8005000019"+ getInputsState(data.substr(data.length-2 , 2));
    } else {
        callback('Wrong data type. data = '+ data);
        return;
    }

    postToSmartConnect(buildOutPayload(event, time, trueData));
    callback(null, "OK");
}

function getInputsState(inputStateData){
    var data = 0;
    var intValue = parseInt("0x"+inputStateData, 16);
    var tor1State = intValue & 0x01;
    var tor2State = (intValue & 0x04) >>> 2;
    var tor3State = (intValue & 0x10) >>> 4;
    console.log("tor3 tor2 tor1 = "+tor3State+tor2State+tor1State);
    
    //inverted logic
    tor1State = ~tor1State & 1;
    tor2State = ~tor2State & 1;
    tor3State = ~tor3State & 1;
    console.log("Inverted logic tor3 tor2 tor1 = "+tor3State+tor2State+tor1State);

    data = data | tor1State;
    data = data | (tor2State << 1);
    data = data | (tor3State << 2);

    var stringData = data.toString(16);
    stringData =  "000"+stringData;
    console.log("InputsStateData = 0x"+stringData);
    return stringData;
}

function buildOutPayload(event, time, data) {
    var lorawanOutput = {
        id: event.id, 
        device_id:event.device_id, 
        type:event.type, 
        timestamp:event.timestamp, 
        numericTimestamp:time, 
        count:event.count, 
        payload_encrypted: event.payload_encrypted, 
        payload_cleartext: data, 
        device_properties: event.device_properties,
        protocol_data: event.protocol_data
    };
    return lorawanOutput;
}

const https = require('https');

function postToSmartConnect(data) {
    console.log('Data to send: ', JSON.stringify(data));

    var post_options = {
      host: process.env.smartconnectcallbackhost, // 'connector-demoenv.devinno.fr'
      path: process.env.smartconnectcallbackPath, // '/ildlo/message'
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