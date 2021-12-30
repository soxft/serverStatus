import React, { useState } from 'react';
import ReactDOM from 'react-dom';
import { Progress } from 'antd'
import 'antd/dist/antd.css';

import ReconnectingWebSocket from 'reconnecting-websocket'


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


  return <>
    {
      Object.keys(serverList).map((client_id, index) => {
        if (!client_id) return <></>;
        let serverInfo = serverList[client_id]
        let itemData = serverInfo['data']

        if (itemData === undefined) {
          return (
            <div key={index}>
              <div>{serverInfo.tag}</div>
              <p>waiting data from server...</p>
            </div>
          );
        } else {
          let cpu_percent = itemData.Cpu.Percent
          let memory_percent = itemData.Memory.Percent

          return (
            <div key={index}>
              <div>{serverInfo.tag}</div>
              <Progress width={80} type="circle" percent={cpu_percent} />
              <Progress width={80} type="circle" percent={memory_percent} />
            </div>
          )
        }
      })
    }
  </>;
}


ReactDOM.render(
  <App />,
  document.getElementById('root')
);

