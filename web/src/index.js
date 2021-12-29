import React, { useState } from 'react';
import ReactDOM from 'react-dom';
import ReconnectingWebSocket from 'reconnecting-websocket'
import { RingProgress } from '@ant-design/charts';


var ws = new ReconnectingWebSocket("ws://127.0.0.1:8282");
ws.onopen = function () {
  ws.send(JSON.stringify({ 'type': 'login' }));
};

const App = () => {

  const [serverList, setServerList] = useState([]);

  ws.onclose = function () {
    console.error("connect down")
  }

  ws.onmessage = function (evt) {
    var resData = JSON.parse(evt.data)
    //console.log(resData)
    var type = resData.type

    if (type === 'login_success') {
      //登录成功
      console.log('Establish connection')
      setServerList(resData.server_list)
    } else if (type === 'server_offline') {
      // 服务器掉线
      let serverListTemp = { ...serverList }
      delete serverListTemp[resData.client_id]
      setServerList(serverListTemp);
    } else if (type === 'server_online') {
      // 服务器上线
      let serverListTemp = { ...serverList }
      serverListTemp[resData.client_id] = resData.data
      setServerList(serverListTemp)
    } else if (type === 'server_info') {
      let serverListTemp = { ...serverList }
      serverListTemp[resData.client_id]['data'] = resData.data
      setServerList(serverListTemp)
    }
  };

  const cpuBaseConfig = {
    height: 100,
    width: 100,
    autoFit: false,
    percent: 1,
    color: ['#5B8FF9', '#E8EDF3'],
  };

  return <>
    {
      Object.keys(serverList).map((client_id, index) => {
        if (!client_id) return <></>;
        let serverInfo = serverList[client_id]
        let itemData = serverInfo['data']
        return (
          <div key={index}>
            <div>{serverInfo.tag}</div>
            <RingProgress key={index} {...cpuBaseConfig} percent={itemData === undefined ? 0.1 : itemData.cpu_percent} />
          </div>
        )
      })
    }
  </>;
}


ReactDOM.render(
  <App />,
  document.getElementById('root')
);

