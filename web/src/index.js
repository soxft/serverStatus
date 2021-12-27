import React from 'react';
import ReactDOM from 'react-dom';

var ws = new WebSocket("ws://127.0.0.1:8282");
ws.onopen = function () {
  // Web Socket 已连接上，使用 send() 方法发送数据
  ws.send(JSON.stringify({ 'type': 'login' }));
};

ws.onmessage = function (evt) {
  var resData = JSON.parse(evt.data)
  console.log(resData)

  var type = resData.type
  if (type === 'login_success') {
    ws.send(JSON.stringify({ 'type': 'get_server_list' }))
  } else if (type === 'server_list') {
    resData.data.map((item, index) => {
      ws.send(JSON.stringify({ 'type': 'get_server_base_info', 'client_id': item.client_id }))
      return index
    })
  }

};

ws.onclose = function () {
  console.log('closed')
};

const App = () => {
  return (
    <h1>HelloWorld</h1>
  );
}


ReactDOM.render(
  <App />,
  document.getElementById('root')
);

