exports.handler = (event, context, callback) => {
    if(typeof event === 'undefined' || event === null) {
        console.error('Received wrong event ', JSON.stringify(event) );
        callback('Received wrong event');
        return;
    }
    console.log('Received event ', JSON.stringify(event) );
    
    if(typeof event.payload_cleartext === 'undefined' || event.payload_cleartext === null ) {
        console.error('Received wrong event ', JSON.stringify(event) );
        callback('Received wrong event');
        return;
    }

    var time = event.timestamp;
    var deveui = event.device_properties.deveui.toUpperCase();
    var data = event.payload_cleartext;

    if ( data.startsWith("47") || data.startsWith("48") ) { //data frame or alert frame
        var lightningState = 0;
        var batteryState = 0;
        var ipTorData = [];

        var intData = parseInt(data,16);

        lightningState = ( intData >>> 2 ) & 1;
        batteryState   = ( intData >>> 3 ) & 1;

        ipTorData.push({
            id : 1,
            label : "DÃ©faut parafoudre",
            state : (lightningState == 1) ? true : false,
            enabled: true
        });
        ipTorData.push({
            id : 2,
            label : "DÃ©faut Batterie",
            state : (batteryState == 1) ? true : false,
            enabled: true
        });
        
        var dataTosend = {
            id : deveui,
            timestamp : time,
            values : ipTorData
        }
        postToSmartConnect(dataTosend);
        
    } else {
        callback('Received wrong data type. Data = '+data);
        return;
    }

    callback(null, "OK");
}

const https = require('https');

function postToSmartConnect(data) {
    console.log('Data to send: ', JSON.stringify(data));

    var post_options = {
      host: process.env.smartconnectcallbackhost, // 'connector-demoenv.devinno.fr'
      path: process.env.smartconnectcallbackpath, // '/inode/data'
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