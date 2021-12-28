import React, { useEffect, useState } from 'react';
import ReactDOM from 'react-dom';

const App = () => {

  const [serverInfo, serServerInfo] = useState(null);

  const establish = () => {
    console.log("Try to connect to server")

    var ws = new WebSocket("ws://10.11.11.11:8282");

    ws.onopen = function () {
      ws.send(JSON.stringify({ 'type': 'login' }));
    };

    ws.onmessage = function (evt) {
      var resData = JSON.parse(evt.data)
      console.log(resData)

      var type = resData.type
      if (type === 'login_success') {
        console.log('Establish connection')
        serServerInfo(JSON.stringify(resData.server_list))
      }
    };

    ws.onclose = function () {
      console.log('Lost connection, try to reconnect in 1 second')
      setTimeout(() => establish(), 1000)
    };
  }

  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(() => establish(), [])

  return (
    <p>{serverInfo}</p>
  );
}


ReactDOM.render(
  <App />,
  document.getElementById('root')
);

