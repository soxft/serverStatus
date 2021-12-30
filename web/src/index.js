import React, { useState } from 'react';
import ReactDOM from 'react-dom';
import { Helmet } from "react-helmet";
import bTS from './class/byteToSize';
import timeTransform from './class/timeTransform';
import timestampToDate from './class/timestampToDate';
import config from './config';
import { RingProgress } from '@ant-design/charts';

import { Card, Col, Row, Spin, Popover, Divider, Button, Typography } from 'antd';
import { LoadingOutlined } from '@ant-design/icons';
import 'antd/dist/antd.css';

import ReconnectingWebSocket from 'reconnecting-websocket'

const { Text, Link } = Typography;

var ws = new ReconnectingWebSocket(config.server);
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
    <Helmet>
      <title>{config.title} - Powered By ServerStatus</title>
    </Helmet>
    <Card bordered={false}>
      <Row gutter={[16, 16]}>
        {
          Object.keys(serverList).map((client_id, index) => {
            if (!client_id) return <></>;
            let serverInfo = serverList[client_id]
            let itemData = serverInfo['data']

            if (itemData === undefined) {
              return (
                <Col key={index}>
                  <Spin indicator={<LoadingOutlined style={{ fontSize: 24 }} spin />}>
                    <Card title={serverInfo.tag}>
                      <p>waiting response from the server</p>
                    </Card>
                  </Spin>
                </Col>
              );
            } else {
              let cpu_percent = itemData.Cpu.Percent / 100
              let memory_percent = itemData.Memory.Percent / 100
              let swap_percent = itemData.Swap.Percent / 100

              let RingProgressBaseConfig = {
                height: 100,
                width: 100,
                autoFit: false,
                color: ['#5B8FF9', '#E8EDF3'],
              }

              return (
                <Col
                  key={index}
                >
                  <Card
                    title={serverInfo.tag}
                    extra={(
                      <Popover
                        placement="leftBottom"
                        trigger="click"
                        content={(<>
                          <p>主机名: {itemData.Host.HostName}</p>
                          <p>开机时间: {timestampToDate(itemData.Host.BootTime)}</p>
                          <Divider />
                          <p>内核架构: {itemData.Host.KernelArch}</p>
                          <p>内核版本: {itemData.Host.KernelVersion}</p>
                          <p>系统类型: {itemData.Host.Os}</p>
                          <Divider />
                          <p>平台: {itemData.Host.Platform}</p>
                          <p>平台家族: {itemData.Host.PlatformFamily}</p>
                          <p>平台版本: {itemData.Host.PlatformVersion}</p>
                          <Divider />
                          <p>虚拟化角色: {itemData.Host.VirtualizationRole}</p>
                          <p>虚拟化类型: {itemData.Host.VirtualizationSystem}</p>
                          <Divider />
                          <p>数据更新时间: {timestampToDate(itemData.Time)}</p>
                        </>
                        )}
                      >
                        <Button type="link" size="small">Details</Button>
                      </Popover>
                    )}
                  >

                    <Row
                      gutter={[24, 16]}
                      justify="space-between"
                    >

                      <Col
                        span={8}
                        align='middle'
                      >
                        <Popover
                          placement="bottom"
                          content={(
                            <>
                              <p>{itemData.Cpu.ModalName}</p>
                              <p>{itemData.Cpu.PhysicalCores}个物理核心</p>
                              <p>{itemData.Cpu.LogicalCores}个逻辑核心</p>
                            </>
                          )}
                        >
                          <RingProgress
                            {...RingProgressBaseConfig}
                            percent={cpu_percent}
                          />
                          <div style={{ height: '5px' }}></div>
                          CPU占用率
                        </Popover>
                      </Col>

                      <Col
                        span={8}
                        align='middle'
                      >
                        <Popover
                          placement="bottom"
                          content={bTS(itemData.Memory.Used) + "/" + bTS(itemData.Memory.Total)}
                        >
                          <RingProgress
                            {...RingProgressBaseConfig}
                            percent={memory_percent}
                          />
                          <div style={{ height: '5px' }}></div>
                          内存占用率
                        </Popover>
                      </Col>

                      <Col
                        span={8}
                        align='middle'
                      >
                        <Popover
                          placement="bottom"
                          content={bTS(itemData.Swap.Used) + "/" + bTS(itemData.Swap.Total)}
                        >
                          <RingProgress
                            {...RingProgressBaseConfig}
                            percent={swap_percent}
                          />
                          <div style={{ height: '5px' }}></div>
                          swap占用率
                        </Popover>
                      </Col>
                    </Row>
                    <Divider />
                    <p>进程数: {itemData.Host.Procs}</p>
                    <p>在线时间: {timeTransform(itemData.Host.UpTime)}</p>
                    <p>负载:{itemData.Load.M1}, {itemData.Load.M5}, {itemData.Load.M15}</p>
                  </Card>
                </Col>
              )
            }
          })
        }
      </Row >
    </Card >
    <Divider />
    <Text style={{ color: "#D3D3D3" }}>&emsp; CopyRight 2021 <Link style={{ color: "#D3D3D3" }} href="https://xsot.cn" target="_blank">xcsoft</Link> All Rights Reserved.</Text>
  </>;
}


ReactDOM.render(
  <App />,
  document.getElementById('root')
);

