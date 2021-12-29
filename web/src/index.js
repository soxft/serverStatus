import React, { useState } from 'react';
import ReactDOM from 'react-dom';

var ws;

const establish = () => {
  console.log("Try to connect to server")
  ws = new WebSocket("ws://127.0.0.1:8282");
}

establish()

const App = () => {

  const [serverInfo, setServerInfo] = useState([]);


  ws.onopen = function () {
    ws.send(JSON.stringify({ 'type': 'login' }));
  };

  ws.onclose = function () {
    console.log('Lost connection, reconnect in 1 second')
    setTimeout(() => establish(), 1000)
  };

  ws.onmessage = function (evt) {
    var resData = JSON.parse(evt.data)
    console.log(resData)

    var type = resData.type
    if (type === 'login_success') {
      console.log('Establish connection')
      setServerInfo(resData.server_list)
    } else if (type === 'server_offline') {
      var server_list_new = [];
      serverInfo.forEach((item) => {
        if (item.client_id !== resData.client_id) server_list_new.push(item)
      })
      setServerInfo(server_list_new);
    } else if (type === 'server_online') {
      setServerInfo(serverInfo => [...serverInfo, resData['data']]);
    }
  };



  return (
    <p>{JSON.stringify(serverInfo)}</p>
  );
}


ReactDOM.render(
  <App />,
  document.getElementById('root')
);

